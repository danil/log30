// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log0 is a JSON logging.
package log0

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
)

type Logger interface {
	io.Writer
	// Get returns copy of the logger with an additional key-values.
	// Copy of the original key-values should have a lower priority
	// than the priority of the newer key-values.
	Get(...KV) Logger
	// Put puts the logger into the sync pool.
	Put()
}

// KV is a key-value pair.
type KV interface {
	encoding.TextMarshaler
	json.Marshaler
}

// KVS is a key-value pair with severity.
type KVS interface {
	encoding.TextMarshaler
	json.Marshaler
	fmt.Stringer
}

const (
	Original = iota
	Excerpt
	Trail
	File
)

const (
	Trunc = iota
	Empty
	Blank
)

// Log is a JSON logger/writer.
type Log struct {
	Output   io.Writer                                // Output is a destination for output.
	Flag     int                                      // Flag is a log properties.
	KV       []KV                                     // KV is a key-values.
	Severity func(severity string) (output io.Writer) // Severity function receives severity level and returns a output writer for a severity level.
	Keys     [4]encoding.TextMarshaler                // Keys: 0 = original message; 1 = message excerpt; 2 = message trail; 3 = file path.
	Key      uint8                                    // Key is a default/sticky message key: all except 1 = original message; 1 = message excerpt.
	Trunc    int                                      // Trunc is a maximum length of an excerpt, after which it is truncated.
	Marks    [3][]byte                                // Marks: 0 = truncate; 1 = empty; 2 = blank.
	Replace  [][2][]byte                              // Replace ia a pairs of byte slices to replace in the message excerpt.
}

var logPool = sync.Pool{New: func() interface{} { return new(Log) }}

// Get returns copy of the logger with additional key-values.
// If first key-value pair implements the KVS interface and the Severity field
// of the Log is not null then calls the function from Severity field
// with the severity level as argument which obtained from KVS interface.
// Then the function from Severity field returns writer for output of the logger.
// Copy of the original key-values has the priority lower
// than the priority of the newer key-values.
func (l *Log) Get(kv ...KV) Logger {
	l0 := logPool.Get().(*Log)
	l0.Output = l.Output
	l0.Flag = l.Flag
	l0.KV = append(l0.KV[:0], append(l.KV, kv...)...)
	l0.Severity = l.Severity
	l0.Keys = l.Keys
	l0.Key = l.Key
	l0.Trunc = l.Trunc
	l0.Marks = l.Marks
	l0.Replace = append(l0.Replace[:0], l.Replace...)

	if l0.Severity != nil && len(kv) > 0 {
		s, ok := kv[0].(KVS)
		if ok {
			out := l0.Severity(s.String())
			if out != nil {
				l0.Output = out
			}
		}
	}

	return l0
}

// Put puts a log into sync pool.
func (l *Log) Put() { logPool.Put(l) }

// Write implements io.Writer. Do nothing if log does not have output.
func (l *Log) Write(src []byte) (int, error) {
	if l.Output == nil {
		return 0, nil
	}
	j, err := l.json(src)
	if err != nil {
		return 0, err
	}
	return l.Output.Write(j)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var (
	mapPool     = sync.Pool{New: func() interface{} { m := make(map[string]json.Marshaler); return &m }}
	excerptPool = sync.Pool{New: func() interface{} { return new([]byte) }}
)

func (l Log) json(src []byte) ([]byte, error) {
	tmpKV := *mapPool.Get().(*map[string]json.Marshaler)
	for k := range tmpKV {
		delete(tmpKV, k)
	}
	defer mapPool.Put(&tmpKV)

	for _, kv := range l.KV {
		p, err := kv.MarshalText()
		if err != nil {
			return nil, err
		}
		tmpKV[string(p)] = kv
	}

	var tail, file int

	if len(src) != 0 {
		switch l.Flag {
		case log.Lshortfile, log.Llongfile:
			i := bytes.Index(src, []byte(": "))
			if i == -1 {
				file = len(src) - 1
				tail = file + 1
			} else {
				file = i
				tail = i + 2
			}
		}
	}

	var originalKey string

	if l.Keys[Original] == nil {
		originalKey = ""
	} else {
		p, err := l.Keys[Original].MarshalText()
		if err != nil {
			return nil, err
		}
		originalKey = string(p)
	}

	var excerptKey string

	if l.Keys[Excerpt] == nil {
		excerptKey = ""
	} else {
		p, err := l.Keys[Excerpt].MarshalText()
		if err != nil {
			return nil, err
		}
		excerptKey = string(p)
	}

	excerpt := *excerptPool.Get().(*[]byte)
	excerpt = excerpt[:0]
	defer excerptPool.Put(&excerpt)

	if tmpKV[excerptKey] == nil {
		if src != nil && tail == len(src) && tmpKV[originalKey] == nil {
			excerpt = append(excerpt, l.Marks[Empty]...)

		} else if tail != len(src) {
			n := len(src) + len(l.Marks[Trunc])
			for _, m := range l.Marks {
				if n < len(m) {
					n = len(m)
				}
			}

			excerpt = append(excerpt, make([]byte, n)...)
			n, err := l.Truncate(excerpt, src[tail:])
			if err != nil {
				return nil, err
			}

			excerpt = excerpt[:n]
		}
	}

	var trailKey string

	if l.Keys[Trail] == nil {
		trailKey = ""
	} else {
		p, err := l.Keys[Trail].MarshalText()
		if err != nil {
			return nil, err
		}
		trailKey = string(p)
	}

	if bytes.Equal(src, excerpt) && src != nil {
		if l.Key == Excerpt {
			tmpKV[excerptKey] = Bytes(src)

		} else {
			if tmpKV[originalKey] == nil {
				tmpKV[originalKey] = Bytes(src)
			} else if len(src) != 0 {
				tmpKV[trailKey] = Bytes(src)
			}
		}

	} else if !bytes.Equal(src, excerpt) {
		if tmpKV[originalKey] == nil {
			tmpKV[originalKey] = Bytes(src)
		} else if tmpKV[originalKey] != nil && len(src) != 0 {
			tmpKV[trailKey] = Bytes(src)
		}

		if tmpKV[excerptKey] == nil && len(excerpt) != 0 {
			tmpKV[excerptKey] = Bytes(excerpt)
		}
	}

	var fileKey string

	if l.Keys[File] == nil {
		fileKey = ""
	} else {
		p, err := l.Keys[File].MarshalText()
		if err != nil {
			return nil, err
		}
		fileKey = string(p)
	}

	if file != 0 {
		tmpKV[fileKey] = Bytes(src[:file])
	}

	p, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(tmpKV)
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

// Truncate writes excerpt of the src to the dst and returns number of the written bytes
// and error if occurre.
func (l Log) Truncate(dst, src []byte) (int, error) {
	var start, end int
	begin := true

	for {
		r, n := utf8.DecodeRune(src[end:])
		if n == 0 {
			break
		}

		// Rids of off all leading space, as defined by Unicode.
		if begin {
			c := src[end]

			// Fast path for ASCII: look for the first ASCII non-space byte or
			// if we run into a non-ASCII byte, fall back
			// to the slower unicode-aware method
			if c < utf8.RuneSelf && asciiSpace[c] == 1 {
				start++
				end++

				continue
			} else if unicode.IsSpace(r) {
				start += n
				end += n

				continue
			} else {
				begin = false
			}
		}

		if end-start >= len(src) || (l.Trunc > 0 && end-start >= l.Trunc) {
			break
		}

		end += n
	}

	truncate := end-start < len(src[start:])

	// Rids of off all trailing white space,
	// as defined by Unicode.
	// Look for the first ASCII non-space byte from the end.
	for ; end > start; end-- {
		c := src[end-1]
		if c >= utf8.RuneSelf {
			end = lastIndexFunc(src[:end], unicode.IsSpace, false)
			if end >= 0 && src[end] >= utf8.RuneSelf {
				_, wid := utf8.DecodeRune(src[end:])
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

	n := copy(dst, src[start:end])

replc:
	for _, r := range l.Replace {
		for offset := 0; offset < n; {
			if len(r[0]) == 0 || bytes.Equal(r[0], r[1]) {
				continue replc
			}

			idx := bytes.Index(dst[offset:n], r[0])
			if idx == -1 {
				continue replc
			}

			offset += idx

			copy(dst, append(dst[:offset], append(r[1], dst[offset+len(r[0]):]...)...))

			offset += len(r[1])
			n += len(r[1]) - len(r[0])
		}
	}

	if end-start == 0 {
		n += copy(dst[n:], l.Marks[Blank])
	}

	if end-start != 0 && truncate {
		n += copy(dst[n:], l.Marks[Trunc])
	}

	return n, nil
}

// GELF returns a GELF formater <https://docs.graylog.org/en/latest/pages/gelf.html>.
func GELF() *Log {
	return &Log{
		// GELF spec version – "1.1"; Must be set by client library.
		// <https://docs.graylog.org/en/latest/pages/gelf.html#gelf-payload-specification>,
		// <https://github.com/graylog-labs/gelf-rb/issues/41#issuecomment-198266505>.
		KV: []KV{
			Strings("version", "1.1"),
			StringFunc("timestamp", func() KV { return Int64(time.Now().Unix()) }),
		},
		Trunc: 120,
		Keys: [4]encoding.TextMarshaler{
			String("full_message"),
			String("short_message"),
			String("_trail"),
			String("_file"),
		},
		Key:     Excerpt,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
	}
}
