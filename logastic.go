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
	runeP     = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func logastic(
	flg int,
	kv,
	kv2 map[string]interface{}, // kv2 is a temporary key-value map in addition to the permanent kv key-value map
	fns map[string]func() interface{},
	trunc int,
	keys [4]string,
	key uint8,
	marks [3][]byte,
	replace [][]byte,
	original ...byte,
) ([]byte, error) {
	if kv2 == nil {
		kv2 = make(map[string]interface{})
	}

	for k, v := range kv {
		if _, ok := kv2[k]; ok {
			continue
		}
		kv2[k] = v
	}

	for k, fn := range fns {
		kv2[k] = fn()
	}

	if v, ok := kv2[keys[Original]]; ok {
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

	if kv2[keys[Excerpt]] != nil {
		switch v := kv2[keys[Excerpt]].(type) {
		case string:
			kv2[keys[Excerpt]] = v
		case []byte:
			kv2[keys[Excerpt]] = string(v)
		case []rune:
			kv2[keys[Excerpt]] = string(v)
		default:
			kv2[keys[Excerpt]] = v
		}
	} else {
		if tail == len(original) {
			excerpt = append(excerpt, marks[emptyMark]...)
		} else {
			i := tail
			beg := true

			for {
				r, n := utf8.DecodeRune(original[i:])
				if n == 0 {
					break
				}

				// Rids of off all leading space, as defined by Unicode.
				// Fast path for ASCII: look for the first ASCII non-space byte or
				// if we run into a non-ASCII byte, fall back to the slower unicode-aware method
				if beg {
					c := original[i]
					if c < utf8.RuneSelf && asciiSpace[c] == 1 || unicode.IsSpace(r) {
						i++
						tail++
						continue
					} else {
						beg = false
					}
				}

				if i-tail >= trunc {
					break
				}

				p := *runeP.Get().(*[]byte)
				p = p[:0]
				defer runeP.Put(&p)

				p = append(p, make([]byte, utf8.RuneLen(r))...)
				utf8.EncodeRune(p, r)
				excerpt = append(excerpt, p...)

				i += n
			}

			truncate := len(original[tail:]) > len(excerpt)

			// Rids of off all trailing white space,
			// as defined by Unicode.
			// Look for the first ASCII non-space byte from the end.
			i = len(excerpt)
			for ; i > 0; i-- {
				c := excerpt[i-1]
				if c >= utf8.RuneSelf {
					excerpt = bytes.TrimFunc(excerpt[0:i], unicode.IsSpace)
					break
				}
				if asciiSpace[c] == 0 {
					excerpt = excerpt[:i]
					break
				}
			}

			if len(excerpt) == 0 {
				excerpt = append(excerpt, marks[blankMark]...)
			}

			if kv2[keys[Host]] != nil {
				excerpt = append(excerpt[:0], append([]byte(fmt.Sprint(kv2[keys[Host]])), append([]byte(" "), excerpt...)...)...)
			}

			if len(excerpt) != 0 && truncate {
				excerpt = append(excerpt, marks[truncMark]...)
			}

			for i := 0; i < len(replace); i += 2 {
				excerpt = bytes.Replace(excerpt, replace[i], replace[i+1], -1)
			}
		}
	}

	if bytes.Equal(original, excerpt) {
		if key != Excerpt {
			key = Original
		}

		if key == Original {
			delete(kv2, keys[Excerpt])
		} else {
			delete(kv2, keys[Original])
		}

		if kv2[keys[key]] == nil {
			kv2[keys[key]] = string(original)
		}
	} else {
		kv2[keys[Original]] = string(original)

		if kv2[keys[Excerpt]] == nil && len(excerpt) != 0 {
			kv2[keys[Excerpt]] = string(excerpt)
		}
	}

	if file != 0 {
		kv2[keys[File]] = string(original[:file])
	}

	p, err := json.Marshal(kv2)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
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
