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

func (mjl Log) Write(p []byte) (int, error) {
	msg, err := message(mjl, p)
	if err != nil {
		return 0, err
	}

	return mjl.Output.Write(msg)
}

func message(mjl Log, p []byte) ([]byte, error) {
	m := make(map[string]interface{})

	for k, v := range mjl.Fields {
		m[k] = v
	}

	for k, fn := range mjl.Functions {
		m[k] = fn()
	}

	full := string(p)
	tail := full
	var file string

	switch mjl.Flag {
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

	if mjl.Fields[mjl.Short] != nil {
		short = fmt.Sprint(mjl.Fields[mjl.Short])
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

				if n == mjl.Truncate {
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

	if mjl.Fields["host"] == nil {
		m[mjl.Short] = short
	} else {
		m[mjl.Short] = fmt.Sprintf("%s %s", mjl.Fields["host"], short)
	}

	if short != full {
		m[mjl.Full] = full
	}

	if file != "" {
		m[mjl.File] = file
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}
