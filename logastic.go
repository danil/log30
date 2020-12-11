package logastic

import (
	"bytes"
	"encoding/json"
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
	Trail
	File
	Host
)

const (
	truncMark = iota
	emptyMark
	blankMark
)

// Logger is a JSON logger/writer.
type Logger struct {
	Output  io.Writer                        // destination for output
	Flag    int                              // log properties
	KV      map[string]json.Marshaler        // key-values
	Funcs   map[string]func() json.Marshaler // dynamically calculated key-values
	Trunc   int                              // maximum length of the message excerpt after which the message excerpt is truncated
	Keys    [4]string                        // 0 = original message; 1 = message excerpt; 2 = message trail; 3 = file path;
	Key     uint8                            // default/sticky message key: all except 1 = original message; 1 = message excerpt;
	Marks   [3][]byte                        // 0 = truncate; 1 = empty; 2 = blank;
	Replace [][2][]byte                      // pairs of byte slices to replace in the message excerpt
}

func (l Logger) Write(p []byte) (int, error) {
	j, err := logastic(l.Flag, l.KV, nil, l.Funcs, l.Trunc, l.Keys, l.Key, l.Marks, l.Replace, p)
	if err != nil {
		return 0, err
	}
	return l.Output.Write(j)
}

// Log is a JSON writer with additional key-value map.
type Log struct {
	log Logger
	kv  map[string]json.Marshaler
}

func (l Logger) With(kv map[string]json.Marshaler) Log {
	return Log{log: l, kv: kv}
}

func (l Log) Write(p []byte) (int, error) {
	j, err := logastic(l.log.Flag, l.log.KV, l.kv, l.log.Funcs, l.log.Trunc, l.log.Keys, l.log.Key, l.log.Marks, l.log.Replace, p)
	if err != nil {
		return 0, err
	}
	return l.log.Output.Write(j)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	excerptP = sync.Pool{New: func() interface{} { return new([]byte) }}
	kvP      = sync.Pool{New: func() interface{} { m := make(map[string]json.Marshaler); return &m }}
)

func logastic(
	flg int,
	permKV, // permKV is a permanent key-value map
	optKV map[string]json.Marshaler, // optKV is a optional key-value map in addition to the permanent key-value map
	fns map[string]func() json.Marshaler,
	trunc int,
	keys [4]string,
	key uint8,
	marks [3][]byte,
	replace [][2][]byte,
	original []byte,
) ([]byte, error) {
	tmpKV := *kvP.Get().(*map[string]json.Marshaler)
	for k := range tmpKV {
		delete(tmpKV, k)
	}
	defer kvP.Put(&tmpKV)

	for k, v := range optKV {
		tmpKV[k] = v
	}

	for k, v := range permKV {
		if _, ok := tmpKV[k]; ok {
			continue
		}
		tmpKV[k] = v
	}

	for k, fn := range fns {
		if _, ok := tmpKV[k]; ok {
			continue
		}
		tmpKV[k] = fn()
	}

	var tail, file int

	if len(original) != 0 {
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
	}

	excerpt := *excerptP.Get().(*[]byte)
	excerpt = excerpt[:0]
	defer excerptP.Put(&excerpt)

	end := tail

	if tmpKV[keys[Excerpt]] == nil {
		if tail == len(original) && tmpKV[keys[Original]] == nil {
			excerpt = append(excerpt[:0], marks[emptyMark]...)
		} else if tail != len(original) {
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
			for _, rep := range replace {
				for offset := 0; ; {
					if len(rep[0]) == 0 || bytes.Equal(rep[0], rep[1]) {
						continue replace
					}

					idx := bytes.Index(excerpt[offset:], rep[0])
					if idx == -1 {
						continue replace
					}

					offset += idx

					excerpt = append(excerpt[:offset], append(rep[1], excerpt[offset+len(rep[0]):]...)...)

					offset += len(rep[1])
				}
			}

			if end-tail == 0 {
				excerpt = append(excerpt, marks[blankMark]...)
			}

			if end-tail != 0 && truncate {
				excerpt = append(excerpt, marks[truncMark]...)
			}
		}
	}

	if bytes.Equal(original, excerpt) && original != nil {
		if key == Excerpt {
			tmpKV[keys[Excerpt]] = Bytes(original)

		} else {
			if tmpKV[keys[Original]] == nil {
				tmpKV[keys[Original]] = Bytes(original)
			} else if len(original) != 0 {
				tmpKV[keys[Trail]] = Bytes(original)
			}
		}

	} else if !bytes.Equal(original, excerpt) {
		if tmpKV[keys[Original]] == nil {
			tmpKV[keys[Original]] = Bytes(original)
		} else if tmpKV[keys[Original]] != nil && len(original) != 0 {
			tmpKV[keys[Trail]] = Bytes(original)
		}

		if tmpKV[keys[Excerpt]] == nil && len(excerpt) != 0 {
			tmpKV[keys[Excerpt]] = Bytes(excerpt)
		}
	}

	if file != 0 {
		tmpKV[keys[File]] = Bytes(original[:file])
	}

	p, err := json.Marshal(tmpKV)
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

func GELF() Logger {
	return Logger{
		KV: map[string]json.Marshaler{
			"version": String("1.1"),
		},
		Funcs: map[string]func() json.Marshaler{
			"timestamp": func() json.Marshaler { return Int64(time.Now().Unix()) },
		},
		Trunc:   120,
		Keys:    [4]string{"full_message", "short_message", "_trail", "_file"},
		Key:     Excerpt,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
	}
}
