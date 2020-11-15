package mujlog

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode"
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

func mujlog(l Log, p []byte) ([]byte, error) {
	m := make(map[string]interface{})

	for k, v := range l.Fields {
		m[k] = v
	}

	for k, fn := range l.Functions {
		m[k] = fn()
	}

	full := string(p)
	tail := full
	var file string

	switch l.Flag {
	case log.Lshortfile, log.Llongfile:
		a := strings.SplitN(full, ": ", 2)
		if len(a) == 1 {
			file = strings.TrimRight(a[0], ":")
			tail = ""
		} else {
			file = a[0]
			tail = a[1]
		}
	}

	var short string

	if l.Fields[l.Short] != nil {
		short = fmt.Sprint(l.Fields[l.Short])
	} else {
		if tail == "" {
			short = "_EMPTY_"
		} else {
			var n int
			beg := true

			for i, r := range tail {
				if beg && unicode.IsSpace(r) {
					continue
				} else {
					beg = false
				}

				n++

				if n == l.Truncate {
					short = tail[:i]
					break
				}
			}

			if short == "" {
				short = tail
			}

			var trunc bool
			if len(tail) != len(short) {
				trunc = true
			}

			short = strings.TrimSpace(short)

			if short == "" {
				short = "_BLANK_"
			}

			if short != "" && trunc {
				short = short + "…"
			}

			short = strings.Replace(short, "\n", " ", -1)
		}
	}

	if l.Fields["host"] == nil {
		m[l.Short] = short
	} else {
		m[l.Short] = fmt.Sprintf("%s %s", l.Fields["host"], short)
	}

	if short != full {
		m[l.Full] = full
	}

	if file != "" {
		m[l.File] = file
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}
