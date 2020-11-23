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

// Log is a Multiline JSON Log and formatter and writer.
type Log struct {
	Output    io.Writer                     // destination for output
	Flag      int                           // log properties
	Fields    map[string]interface{}        // additional fields
	Functions map[string]func() interface{} // dynamically calculated fields
	Short     string
	Full      string
	File      string
	Truncate  int
}

func GELF() Log {
	return Log{
		Flag: log.Llongfile,
		Fields: map[string]interface{}{
			"version": "1.1",
		},
		Functions: map[string]func() interface{}{
			"timestamp": func() interface{} { return time.Now().Unix() },
		},
		Short:    "short_message",
		Full:     "full_message",
		File:     "_file",
		Truncate: 120,
	}
}

func (muj Log) Write(p []byte) (int, error) {
	return muj.Log(p, make(map[string]interface{}))
}

func (muj Log) Log(p []byte, kv map[string]interface{}) (int, error) {
	if kv == nil {
		kv = make(map[string]interface{})
	}

	msg, err := mujlog(muj, p, kv)
	if err != nil {
		return 0, err
	}

	return muj.Output.Write(msg)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	shortP = sync.Pool{New: func() interface{} { return new([]byte) }}
	runeP  = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func mujlog(muj Log, full []byte, kv map[string]interface{}) ([]byte, error) {
	for k, v := range muj.Fields {
		if _, ok := kv[k]; ok {
			continue
		}
		kv[k] = v
	}

	for k, fn := range muj.Functions {
		kv[k] = fn()
	}

	var tail, file int

	switch muj.Flag {
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

	if muj.Fields[muj.Short] != nil {
		switch v := muj.Fields[muj.Short].(type) {
		case string:
			kv[muj.Short] = v
		case []byte:
			kv[muj.Short] = string(v)
		case []rune:
			kv[muj.Short] = string(v)
		default:
			kv[muj.Short] = v
		}
	} else {
		if tail == len(full) {
			short = append(short, []byte("_EMPTY_")...)
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

				if i-tail >= muj.Truncate-1 {
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

			var trunc bool
			if len(full[tail:]) > len(short) {
				trunc = true
			}

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
				short = append(short, []byte("_BLANK_")...)
			}

			if len(short) != 0 && trunc {
				short = append(short, []byte("…")...)
			}

			short = bytes.Replace(short, []byte("\n"), []byte(" "), -1)
		}
	}

	if _, ok := kv[muj.Short]; !ok {
		if muj.Fields["host"] == nil {
			kv[muj.Short] = string(short)
		} else {
			kv[muj.Short] = fmt.Sprintf("%s %s", muj.Fields["host"], short)
		}
	}

	if !bytes.Equal(short, full) {
		kv[muj.Full] = string(full)
	}

	if file != 0 {
		kv[muj.File] = string(full[:file])
	}

	p, err := json.Marshal(kv)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}
