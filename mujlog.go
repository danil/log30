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
	Keys      [3]string                     // key names: 0 = message; 1 = short message; 2 = file;
	Truncate  int                           // maximum length of the short message
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
		Keys:     [3]string{"full_message", "short_message", "_file"},
		Truncate: 120,
	}
}

func (muj Log) Write(p []byte) (int, error) {
	return muj.Log(p, nil)
}

func (muj Log) Log(p []byte, kv map[string]interface{}) (int, error) {
	if kv == nil {
		kv = make(map[string]interface{})
	}
	j, err := mujlog(muj, p, kv)
	if err != nil {
		return 0, err
	}
	return muj.Output.Write(j)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	msgP   = sync.Pool{New: func() interface{} { return new([]byte) }}
	shortP = sync.Pool{New: func() interface{} { return new([]byte) }}
	runeP  = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func mujlog(muj Log, msg []byte, kv map[string]interface{}) ([]byte, error) {
	for k, v := range muj.Fields {
		if _, ok := kv[k]; ok {
			continue
		}
		kv[k] = v
	}

	for k, fn := range muj.Functions {
		kv[k] = fn()
	}

	if v, ok := kv[muj.Keys[0]]; ok {
		p := *msgP.Get().(*[]byte)
		p = p[:0]
		defer msgP.Put(&p)

		if v != nil {
			p = append(p, []byte(fmt.Sprint(v))...)
		}

		if msg == nil {
			msg = p
		} else {
			msg = append(p, msg...)
		}
	}

	var tail, file int

	switch muj.Flag {
	case log.Lshortfile, log.Llongfile:
		i := bytes.Index(msg, []byte(": "))
		if i == -1 {
			file = len(msg) - 1
			tail = file + 1
		} else {
			file = i
			tail = i + 2
		}
	}

	short := *shortP.Get().(*[]byte)
	short = short[:0]
	defer shortP.Put(&short)

	if muj.Fields[muj.Keys[1]] != nil {
		switch v := muj.Fields[muj.Keys[1]].(type) {
		case string:
			kv[muj.Keys[1]] = v
		case []byte:
			kv[muj.Keys[1]] = string(v)
		case []rune:
			kv[muj.Keys[1]] = string(v)
		default:
			kv[muj.Keys[1]] = v
		}
	} else {
		if tail == len(msg) {
			short = append(short, []byte("_EMPTY_")...)
		} else {
			i := tail
			beg := true

			for {
				r, n := utf8.DecodeRune(msg[i:])
				if n == 0 {
					break
				}

				// Rids of off all leading space, as defined by Unicode.
				// Fast path for ASCII: look for the first ASCII non-space byte or
				// if we run into a non-ASCII byte, fall back to the slower unicode-aware method
				if beg {
					c := msg[i]
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
			if len(msg[tail:]) > len(short) {
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

	if _, ok := kv[muj.Keys[1]]; !ok {
		if muj.Fields["host"] == nil {
			kv[muj.Keys[1]] = string(short)
		} else {
			kv[muj.Keys[1]] = fmt.Sprintf("%s %s", muj.Fields["host"], short)
		}
	}

	if bytes.Equal(short, msg) {
		delete(kv, muj.Keys[0])
	} else {
		kv[muj.Keys[0]] = string(msg)
	}

	if file != 0 {
		kv[muj.Keys[2]] = string(msg[:file])
	}

	p, err := json.Marshal(kv)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}
