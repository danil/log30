package mujlog

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
	FullKey = iota
	ShortKey
	FileKey
	HostKey
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
	Trunc   int                           // maximum length of the short message after which the short message is truncated
	Keys    [4]string                     // 0 = full message key; 1 = short message key; 2 = file key; 3 = host key;
	Key     uint8                         // sticky message key: all except 1 = full message; 1 = short message;
	Marks   [3][]byte                     // 0 = truncate mark; 1 = empty mark; 2 = blank mark;
	Replace [][]byte                      // pairs of byte slices to replace in a short message
}

func (muj Log) Write(p []byte) (int, error) {
	j, err := mujlog(p, muj.Flag, muj.KV, nil, muj.Funcs, muj.Trunc, muj.Keys, muj.Key, muj.Marks, muj.Replace)
	if err != nil {
		return 0, err
	}
	return muj.Output.Write(j)
}

func (muj Log) Log(p []byte, kv map[string]interface{}) (int, error) {
	j, err := mujlog(p, 0, muj.KV, kv, muj.Funcs, muj.Trunc, muj.Keys, muj.Key, muj.Marks, muj.Replace)
	if err != nil {
		return 0, err
	}
	return muj.Output.Write(j)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	fullP  = sync.Pool{New: func() interface{} { return new([]byte) }}
	shortP = sync.Pool{New: func() interface{} { return new([]byte) }}
	runeP  = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func mujlog(
	full []byte,
	flg int,
	kv,
	kv2 map[string]interface{}, // kv2 is a temporary key-value map in addition to the permanent kv key-value map
	fns map[string]func() interface{},
	trunc int,
	keys [4]string,
	key uint8,
	marks [3][]byte,
	replace [][]byte,
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

	if v, ok := kv2[keys[FullKey]]; ok {
		p := *fullP.Get().(*[]byte)
		p = p[:0]
		defer fullP.Put(&p)

		if v != nil {
			p = append(p, []byte(fmt.Sprint(v))...)
		}

		if full == nil {
			full = p
		} else {
			full = append(p, full...)
		}
	}

	var tail, file int

	switch flg {
	case log.Lshortfile, log.Llongfile:
		i := bytes.Index(full, []byte(": "))
		if i == -1 {
			file = len(full) - 1
			tail = file + 1
		} else {
			file = i
			tail = i + 2
		}
	}

	short := *shortP.Get().(*[]byte)
	short = short[:0]
	defer shortP.Put(&short)

	if kv2[keys[ShortKey]] != nil {
		switch v := kv2[keys[ShortKey]].(type) {
		case string:
			kv2[keys[ShortKey]] = v
		case []byte:
			kv2[keys[ShortKey]] = string(v)
		case []rune:
			kv2[keys[ShortKey]] = string(v)
		default:
			kv2[keys[ShortKey]] = v
		}
	} else {
		if tail == len(full) {
			short = append(short, marks[emptyMark]...)
		} else {
			i := tail
			beg := true

			for {
				r, n := utf8.DecodeRune(full[i:])
				if n == 0 {
					break
				}

				// Rids of off all leading space, as defined by Unicode.
				// Fast path for ASCII: look for the first ASCII non-space byte or
				// if we run into a non-ASCII byte, fall back to the slower unicode-aware method
				if beg {
					c := full[i]
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
				short = append(short, p...)

				i += n
			}

			truncate := len(full[tail:]) > len(short)

			// Rids of off all trailing white space,
			// as defined by Unicode.
			// Look for the first ASCII non-space byte from the end.
			i = len(short)
			for ; i > 0; i-- {
				c := short[i-1]
				if c >= utf8.RuneSelf {
					short = bytes.TrimFunc(short[0:i], unicode.IsSpace)
					break
				}
				if asciiSpace[c] == 0 {
					short = short[:i]
					break
				}
			}

			if len(short) == 0 {
				short = append(short, marks[blankMark]...)
			}

			if kv2[keys[HostKey]] != nil {
				short = append(short[:0], append([]byte(fmt.Sprint(kv2[keys[HostKey]])), append([]byte(" "), short...)...)...)
			}

			if len(short) != 0 && truncate {
				short = append(short, marks[truncMark]...)
			}

			for i := 0; i < len(replace); i += 2 {
				short = bytes.Replace(short, replace[i], replace[i+1], -1)
			}
		}
	}

	if bytes.Equal(full, short) {
		if key != ShortKey {
			key = FullKey
		}

		if key == FullKey {
			delete(kv2, keys[ShortKey])
		} else {
			delete(kv2, keys[FullKey])
		}

		if kv2[keys[key]] == nil {
			kv2[keys[key]] = string(full)
		}
	} else {
		kv2[keys[FullKey]] = string(full)

		if kv2[keys[ShortKey]] == nil && len(short) != 0 {
			kv2[keys[ShortKey]] = string(short)
		}
	}

	if file != 0 {
		kv2[keys[FileKey]] = string(full[:file])
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
		Key:     ShortKey,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][]byte{[]byte("\n"), []byte(" ")},
	}
}
