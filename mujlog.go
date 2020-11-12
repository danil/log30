package mujlog

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"
)

// Mujlog (Multiline JSON Log) is an formatter and writer.
type Mujlog struct {
	Output    io.Writer                     // destination for output
	Flag      int                           // log properties
	Fields    map[string]string             // additional fields
	Functions map[string]func() interface{} // dynamically calculated fields
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
	clean := strings.TrimSpace(full)
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

	if full == "" {
		short = "_EMPTY_"
	} else if clean == "" {
		short = "_BLANK_"
	} else {
		short = strings.SplitN(clean, "\n", 2)[0]
	}

	if mjl.Fields["host"] == "" {
		m["short_message"] = short
	} else {
		m["short_message"] = mjl.Fields["host"] + " " + short
	}

	if short != full {
		m["full_message"] = full
	}

	if file != "" {
		m["_file"] = file
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}

func GELF() Mujlog {
	return Mujlog{
		Flag: log.Llongfile,
		Fields: map[string]string{
			"version": "1.1",
		},
		Functions: map[string]func() interface{}{
			"timestamp": func() interface{} { return time.Now().Unix() },
		},
	}
}
