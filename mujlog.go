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

func (l Log) Write(p []byte) (int, error) {
	msg, err := mujlog(l, p)
	if err != nil {
		return 0, err
	}

	return l.Output.Write(msg)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	shortp = sync.Pool{New: func() interface{} { return new([]byte) }}
	runep  = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func mujlog(l Log, full []byte) ([]byte, error) {
	m := make(map[string]interface{})

	for k, v := range l.Fields {
		m[k] = v
	}

	for k, fn := range l.Functions {
		m[k] = fn()
	}

	var tail, file int

	switch l.Flag {
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

	short := *shortp.Get().(*[]byte)
	short = short[:0]
	defer shortp.Put(&short)

	if l.Fields[l.Short] != nil {
		switch v := l.Fields[l.Short].(type) {
		case string:
			m[l.Short] = v
		case []byte:
			m[l.Short] = string(v)
		case []rune:
			m[l.Short] = string(v)
		default:
			m[l.Short] = v
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

				if i-tail >= l.Truncate-1 {
					break
				}

				p := *runep.Get().(*[]byte)
				p = p[:0]
				defer runep.Put(&p)

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

	if _, ok := m[l.Short]; !ok {
		if l.Fields["host"] == nil {
			m[l.Short] = string(short)
		} else {
			m[l.Short] = fmt.Sprintf("%s %s", l.Fields["host"], short)
		}
	}

	if !bytes.Equal(short, full) {
		m[l.Full] = string(full)
	}

	if file != 0 {
		m[l.File] = string(full[:file])
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}
