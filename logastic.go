package logastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	Original = iota
	Excerpt
	File
	Host
)

const (
	truncMark = iota
	emptyMark
	blankMark
)

// Log is a Multiline JSON Log and formatter and writer.
type Log struct {
	Output  io.Writer                     // destination for output
	Flag    int                           // log properties
	KV      map[string]interface{}        // key-values
	Funcs   map[string]func() interface{} // dynamically calculated key-values
	Trunc   int                           // maximum length of the message excerpt after which the message excerpt is truncated
	Keys    [4]string                     // 0 = original message; 1 = message excerpt; 2 = file; 3 = host;
	Key     uint8                         // sticky message key: all except 1 = original message; 1 = message excerpt;
	Marks   [3][]byte                     // 0 = truncate; 1 = empty; 2 = blank;
	Replace [][]byte                      // pairs of byte slices to replace in the message excerpt
}

func (l Log) Write(p []byte) (int, error) {
	j, err := logastic(l.Flag, l.KV, nil, l.Funcs, l.Trunc, l.Keys, l.Key, l.Marks, l.Replace, p...)
	if err != nil {
		return 0, err
	}
	return l.Output.Write(j)
}

func (l Log) Log(kv map[string]interface{}, p ...byte) (int, error) {
	j, err := logastic(0, l.KV, kv, l.Funcs, l.Trunc, l.Keys, l.Key, l.Marks, l.Replace, p...)
	if err != nil {
		return 0, err
	}
	return l.Output.Write(j)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	originalP = sync.Pool{New: func() interface{} { return new([]byte) }}
	excerptP  = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func logastic(
	flg int,
	permKV,
	tempKV map[string]interface{}, // tempKV is a temporary key-value map in addition to the permanent kv key-value map
	fns map[string]func() interface{},
	trunc int,
	keys [4]string,
	key uint8,
	marks [3][]byte,
	replace [][]byte,
	original ...byte,
) ([]byte, error) {
	if tempKV == nil {
		tempKV = make(map[string]interface{})
	}

	for k, v := range permKV {
		if _, ok := tempKV[k]; ok {
			continue
		}
		tempKV[k] = v
	}

	for k, fn := range fns {
		if _, ok := tempKV[k]; ok {
			continue
		}
		tempKV[k] = fn()
	}

	if v, ok := tempKV[keys[Original]]; ok {
		p := *originalP.Get().(*[]byte)
		p = p[:0]
		defer originalP.Put(&p)

		if v != nil {
			p = append(p, []byte(fmt.Sprint(v))...)
		}

		if original == nil {
			original = p
		} else {
			original = append(p, original...)
		}
	}

	var tail, file int

	switch flg {
	case log.Lshortfile, log.Llongfile:
		i := bytes.Index(original, []byte(": "))
		if i == -1 {
			file = len(original) - 1
			tail = file + 1
		} else {
			file = i
			tail = i + 2
		}
	}

	excerpt := *excerptP.Get().(*[]byte)
	excerpt = excerpt[:0]
	defer excerptP.Put(&excerpt)

	end := tail

	if tempKV[keys[Excerpt]] != nil {
		switch v := tempKV[keys[Excerpt]].(type) {
		case string:
			tempKV[keys[Excerpt]] = v
		case []byte:
			tempKV[keys[Excerpt]] = string(v)
		case []rune:
			tempKV[keys[Excerpt]] = string(v)
		default:
			tempKV[keys[Excerpt]] = v
		}
	} else {
		if tail == len(original) {
			excerpt = append(excerpt[:0], marks[emptyMark]...)
		} else {
			beg := true

			for {
				r, n := utf8.DecodeRune(original[end:])
				if n == 0 {
					break
				}

				// Rids of off all leading space, as defined by Unicode.
				if beg {
					c := original[end]

					// Fast path for ASCII: look for the first ASCII non-space byte or
					// if we run into a non-ASCII byte, fall back
					// to the slower unicode-aware method
					if c < utf8.RuneSelf && asciiSpace[c] == 1 {
						tail++
						end++

						continue
					} else if unicode.IsSpace(r) {
						tail += n
						end += n

						continue
					} else {
						beg = false
					}
				}

				if end-tail >= trunc {
					break
				}

				end += n
			}

			truncate := end-tail < len(original[tail:])

			// Rids of off all trailing white space,
			// as defined by Unicode.
			// Look for the first ASCII non-space byte from the end.
			for ; end > tail; end-- {
				c := original[end-1]
				if c >= utf8.RuneSelf {
					end = lastIndexFunc(original[:end], unicode.IsSpace, false)
					if end >= 0 && original[end] >= utf8.RuneSelf {
						_, wid := utf8.DecodeRune(original[end:])
						end += wid
					} else {
						end++
					}
					break
				}
				if asciiSpace[c] == 0 {
					break
				}
			}

			excerpt = append(excerpt[:0], original[tail:end]...)

		replace:
			for n := 0; n < len(replace); n += 2 {
				for i := 0; ; {
					j := bytes.Index(original[tail+i:end], replace[n])
					if j == -1 {
						continue replace
					}
					j += i
					i = j + len(replace[n+1])
					excerpt = append(excerpt[:j], append(replace[n+1], excerpt[i:]...)...)
				}
			}

			if end-tail == 0 {
				excerpt = append(excerpt, marks[blankMark]...)
			}

			if tempKV[keys[Host]] != nil {
				excerpt = append(excerpt[:0], append([]byte(fmt.Sprint(tempKV[keys[Host]])), append([]byte(" "), excerpt...)...)...)
			}

			if end-tail != 0 && truncate {
				excerpt = append(excerpt, marks[truncMark]...)
			}
		}
	}

	if bytes.Equal(original, excerpt) {
		if key != Excerpt {
			key = Original
		}

		if key == Original {
			delete(tempKV, keys[Excerpt])
		} else {
			delete(tempKV, keys[Original])
		}

		if tempKV[keys[key]] == nil {
			tempKV[keys[key]] = string(original)
		}
	} else {
		tempKV[keys[Original]] = string(original)

		if tempKV[keys[Excerpt]] == nil && len(excerpt) != 0 {
			tempKV[keys[Excerpt]] = string(excerpt)
		}
	}

	if file != 0 {
		tempKV[keys[File]] = string(original[:file])
	}

	p, err := json.Marshal(tempKV)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}

// lastIndexFunc is the same as bytes.LastIndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
// lastIndexFunc copied from the bytes package.
func lastIndexFunc(s []byte, f func(r rune) bool, truth bool) int {
	for i := len(s); i > 0; {
		r, size := rune(s[i-1]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeLastRune(s[0:i])
		}
		i -= size
		if f(r) == truth {
			return i
		}
	}
	return -1
}

func GELF() Log {
	return Log{
		KV: map[string]interface{}{
			"version": "1.1",
		},
		Funcs: map[string]func() interface{}{
			"timestamp": func() interface{} { return time.Now().Unix() },
		},
		Trunc:   120,
		Keys:    [4]string{"full_message", "short_message", "_file", "host"},
		Key:     Excerpt,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][]byte{[]byte("\n"), []byte(" ")},
	}
}
