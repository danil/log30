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

const gelf = "GELF"

// Mujlog (Multiline JSON Log) is a formatter and writer.
type Mujlog struct {
	Output    io.Writer                     // destination for output
	Flag      int                           // log properties
	Fields    map[string]interface{}        // additional fields
	Functions map[string]func() interface{} // dynamically calculated fields
	Metadata  map[string]string
}

func Metadata() map[string]string {
	return map[string]string{
		"short": "shortMessage",
		"full":  "fullMessage",
		"file":  "file",
	}
}

func GELF() Mujlog {
	m := map[string]string{
		"format": gelf,
		"short":  "short_message",
		"full":   "full_message",
		"file":   "_file",
	}

	return Mujlog{
		Flag: log.Llongfile,
		Fields: map[string]interface{}{
			"version": "1.1",
		},
		Functions: map[string]func() interface{}{
			"timestamp": func() interface{} { return time.Now().Unix() },
		},
		Metadata: m,
	}
}

func (mjl Mujlog) Write(p []byte) (int, error) {
	msg, err := message(mjl, p)
	if err != nil {
		return 0, err
	}

	return mjl.Output.Write(msg)
}

func message(mjl Mujlog, p []byte) ([]byte, error) {
	m := make(map[string]interface{})

	for k, v := range mjl.Fields {
		m[k] = v
	}

	for k, fn := range mjl.Functions {
		m[k] = fn()
	}

	full := string(p)

	var clean string
	ir := 0
	for i, r := range full {
		if unicode.IsSpace(r) {
			continue
		}
		ir++
		if ir > 1024 {
			clean = full[:i]
			break
		}
	}
	if clean == "" {
		clean = full
	}
	clean = strings.TrimSpace(clean)

	file := ""
	switch mjl.Flag {
	case log.Lshortfile, log.Llongfile:
		a := strings.SplitN(clean, ": ", 2)
		if len(a) == 1 {
			file = strings.TrimRight(a[0], ":")
			clean = ""
		} else {
			file = a[0]
			clean = a[1]
		}
	}

	var short string

	if mjl.Fields[mjl.Metadata["short"]] != nil {
		short = fmt.Sprint(mjl.Fields[mjl.Metadata["short"]])
	} else {
		if full == "" {
			short = "_EMPTY_"
		} else if clean == "" {
			short = "_BLANK_"
		} else {
			ir := 0
			for i, _ := range clean {
				ir++
				if ir > 119 {
					short = clean[:i] + "…"
					break
				}
			}
			if short == "" {
				short = clean
			}
			short = strings.Replace(short, "\n", " ", -1)
		}
	}

	if mjl.Fields["host"] == nil {
		m[mjl.Metadata["short"]] = short
	} else {
		m[mjl.Metadata["short"]] = fmt.Sprintf("%s %s", mjl.Fields["host"], short)
	}

	if short != full {
		m[mjl.Metadata["full"]] = full
	}

	if file != "" {
		m[mjl.Metadata["file"]] = file
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}
