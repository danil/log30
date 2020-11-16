package mujlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func mujlog(l Log, full []byte) ([]byte, error) {
	m := make(map[string]interface{})

	for k, v := range l.Fields {
		m[k] = v
	}

	for k, fn := range l.Functions {
		m[k] = fn()
	}

	tail := make([]byte, len(full))
	copy(tail, full)

	var file []byte

	switch l.Flag {
	case log.Lshortfile, log.Llongfile:
		a := bytes.SplitN(full, []byte(": "), 2)
		if len(a) == 1 {
			file = bytes.TrimRight(a[0], ":")
			tail = tail[:0]
		} else {
			file = a[0]
			tail = a[1]
		}
	}

	var short []byte

	if l.Fields[l.Short] != nil {
		switch v := l.Fields[l.Short].(type) {
		case []byte:
			short = v
		case string:
			short = []byte(v)
		default:
			short = []byte(fmt.Sprint(v))
		}
	} else {
		if bytes.Equal(tail, []byte{}) {
			short = []byte("_EMPTY_")
		} else {
			tail = trimSpaceLeft(tail)

			buf := bytes.NewBuffer(tail)
			i := 0

			for {
				_, n, err := buf.ReadRune()
				if err == io.EOF {
					break
				} else if err != nil {
					return []byte{}, err
				}
				if i >= l.Truncate-1 {
					break
				}
				i += n
			}

			short = tail[:i]

			if bytes.Equal(short, []byte{}) {
				short = tail
			}

			var trunc bool
			if len(tail) != len(short) {
				trunc = true
			}

			short = trimSpaceRight(short)

			if bytes.Equal(short, []byte{}) {
				short = []byte("_BLANK_")
			}

			if !bytes.Equal(short, []byte{}) && trunc {
				short = append(short, []byte("…")...)
			}

			short = bytes.Replace(short, []byte("\n"), []byte(" "), -1)
		}
	}

	if l.Fields["host"] == nil {
		m[l.Short] = string(short)
	} else {
		m[l.Short] = fmt.Sprintf("%s %s", l.Fields["host"], short)
	}

	if !bytes.Equal(short, full) {
		m[l.Full] = string(full)
	}

	if !bytes.Equal(file, []byte{}) {
		m[l.File] = string(file)
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return append(p, '\n'), nil
}

// trimSpaceLeft returns a subslice of s by slicing off all leading space,
// as defined by Unicode.
// trimSpaceLeft function is partially copied from the bytes package.
func trimSpaceLeft(s []byte) []byte {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return bytes.TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[start:stop] starts and ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	if start == len(s) {
		// Special case to preserve previous TrimLeftFunc behavior,
		// returning nil instead of empty slice if all spaces.
		return nil
	}
	return s[start:]
}

// trimSpaceRight returns a subslice of s by slicing off all trailing white space,
// as defined by Unicode.
// trimSpaceRight function is partially copied from the bytes package.
func trimSpaceRight(s []byte) []byte {
	// Now look for the first ASCII non-space byte from the end
	stop := len(s)
	for ; stop > 0; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return bytes.TrimFunc(s[0:stop], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[start:stop] starts and ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	if stop == 0 {
		// Special case to preserve previous TrimLeftFunc behavior,
		// returning nil instead of empty slice if all spaces.
		return nil
	}
	return s[:stop]
}
