// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log0_test

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/danil/log0"
	"github.com/kinbiko/jsonassert"
)

var WriteTestCases = []struct {
	name      string
	line      int
	log       log0.Logger
	input     []byte
	kv        []log0.KV
	expected  string
	benchmark bool
}{
	{
		name:     "nil",
		line:     line(),
		log:      dummy(),
		input:    nil,
		expected: `{}`,
	},
	{
		name: `nil message with "foo" key with "bar" value`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("foo", "bar")},
		},
		input: nil,
		expected: `{
	    "foo":"bar"
		}`,
	},
	{
		name:  "empty",
		line:  line(),
		log:   dummy(),
		input: []byte{},
		expected: `{
	    "message":"",
			"excerpt":"_EMPTY_"
		}`,
	},
	{
		name: "blank",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Original,
			Trunc:  120,
			Marks:  [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		},
		input: []byte(" "),
		expected: `{
	    "message":" ",
			"excerpt":"_BLANK_"
		}`,
	},
	{
		name: `"string" key with "foo" value and "string" key with "bar" value`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("string", "foo")},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
			Trunc:  120,
		},
		input: []byte("Hello, World!"),
		kv:    []log0.KV{log0.Strings("string", "bar")},
		expected: `{
			"message":"Hello, World!",
		  "string": "bar"
		}`,
		benchmark: true,
	},
	{
		name:  "kv is nil",
		line:  line(),
		log:   dummy(),
		input: []byte("Hello, World!"),
		kv:    nil,
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: `bytes appends to the "message" key with "string value"`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			KV:      []log0.KV{log0.Strings("message", "string value")},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("Hello,\nWorld!"),
		expected: `{
			"message":"string value",
			"excerpt":"Hello, World!",
			"trail":"Hello,\nWorld!"
		}`,
	},
	{
		name:  `bytes appends to the "message" key with "string value"`,
		line:  line(),
		log:   dummy(),
		input: []byte("Hello,\nWorld!"),
		kv:    []log0.KV{log0.Strings("message", "string value")},
		expected: `{
			"message":"string value",
			"excerpt":"Hello, World!",
			"trail":"Hello,\nWorld!"
		}`,
	},
	{
		name: `bytes is nil and "message" key with "string value"`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("message", "string value")},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
			Trunc:  120,
		},
		expected: `{
			"message":"string value"
		}`,
	},
	{
		name: `input is nil and "message" key with "string value"`,
		line: line(),
		log:  dummy(),
		kv:   []log0.KV{log0.Strings("message", "string value")},
		expected: `{
			"message":"string value"
		}`,
	},
	{
		name:  `bytes appends to the integer key "message"`,
		line:  line(),
		log:   dummy(),
		input: []byte("Hello, World!\n"),
		kv:    []log0.KV{log0.StringInt("message", 1)},
		expected: `{
			"message":1,
			"excerpt":"Hello, World!",
			"trail":"Hello, World!\n"
		}`,
	},
	{
		name:  `bytes appends to the float 32 bit key "message"`,
		line:  line(),
		log:   dummy(),
		input: []byte("Hello,\nWorld!"),
		kv:    []log0.KV{log0.StringFloat32("message", 4.2)},
		expected: `{
			"message":4.2,
			"excerpt":"Hello, World!",
			"trail":"Hello,\nWorld!"
		}`,
	},
	{
		name:  `bytes appends to the float 64 bit key "message"`,
		line:  line(),
		log:   dummy(),
		input: []byte("Hello,\nWorld!"),
		kv:    []log0.KV{log0.StringFloat64("message", 4.2)},
		expected: `{
			"message":4.2,
			"excerpt":"Hello, World!",
			"trail":"Hello,\nWorld!"
		}`,
	},
	{
		name:  `bytes appends to the boolean key "message"`,
		line:  line(),
		log:   dummy(),
		input: []byte("Hello,\nWorld!"),
		kv:    []log0.KV{log0.StringBool("message", true)},
		expected: `{
			"message":true,
			"excerpt":"Hello, World!",
			"trail":"Hello,\nWorld!"
		}`,
	},
	{
		name:  `bytes will appends to the nil key "message"`,
		line:  line(),
		log:   dummy(),
		input: []byte("Hello, World!"),
		kv:    []log0.KV{log0.StringReflect("message", nil)},
		expected: `{
			"message":null,
			"trail":"Hello, World!"
		}`,
	},
	{
		name: `default key is original and bytes is nil and "message" key is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
			Key:    log0.Original,
			Trunc:  120,
		},
		kv: []log0.KV{log0.Strings("message", "foo")},
		expected: `{
			"message":"foo"
		}`,
	},
	{
		name: `default key is original and bytes is nil and "message" key is present and with replace`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		kv: []log0.KV{log0.Strings("message", "foo\n")},
		expected: `{
			"message":"foo\n"
		}`,
	},
	{
		name: `default key is original and bytes is present and "message" key is present`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo"),
		kv:    []log0.KV{log0.Strings("message", "bar")},
		expected: `{
			"message":"bar",
			"trail":"foo"
		}`,
	},
	{
		name: `default key is original and bytes is present and "message" key is present and with replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Trunc:  120,
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:    log0.Original,
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar")},
		expected: `{
			"message":"bar",
			"excerpt":"foo",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is original and bytes is present and "message" key is present and with replace input bytes and key`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar\n")},
		expected: `{
			"message":"bar\n",
			"excerpt":"foo",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is original and bytes is nil and "excerpt" key is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Original,
			Trunc:  120,
		},
		kv: []log0.KV{log0.Strings("excerpt", "foo")},
		expected: `{
			"excerpt":"foo"
		}`,
	},
	{
		name: `default key is original and bytes is nil and "excerpt" key is present and with replace`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		kv: []log0.KV{log0.Strings("excerpt", "foo\n")},
		expected: `{
			"excerpt":"foo\n"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" key is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Original,
			Trunc:  120,
		},
		input: []byte("foo"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" key is present and with replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" key is present and with replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" key is present and with replace input bytes and rey`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar\n")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar\n"
		}`,
	},
	{
		name: `default key is original and bytes is nil and "excerpt" and "message" keys is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Original,
			Trunc:  120,
		},
		kv: []log0.KV{log0.Strings("message", "foo"), log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is original and bytes is nil and "excerpt" and "message" keys is present and replace keys`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		kv: []log0.KV{log0.Strings("message", "foo\n"), log0.Strings("excerpt", "bar\n")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar\n"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" and "message" keys is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:    log0.Original,
			Trunc:  120,
		},
		input: []byte("foo"),
		kv:    []log0.KV{log0.Strings("message", "bar"), log0.Strings("excerpt", "xyz")},
		expected: `{
			"message":"bar",
			"excerpt":"xyz",
			"trail":"foo"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" and "message" keys is present and replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv: []log0.KV{
			log0.Strings("message", "bar"),
			log0.Strings("excerpt", "xyz"),
		},
		expected: `{
			"message":"bar",
			"excerpt":"xyz",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is original and bytes is present and "excerpt" and "message" keys is present and replace input bytes and keys`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Original,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar\n"), log0.Strings("excerpt", "xyz\n")},
		expected: `{
			"message":"bar\n",
			"excerpt":"xyz\n",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is nil and "message" key is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
			Key:    log0.Excerpt,
			Trunc:  120,
		},
		kv: []log0.KV{log0.Strings("message", "foo")},
		expected: `{
			"message":"foo"
		}`,
	},
	{
		name: `default key is excerpt and bytes is nil and "message" key is present and with replace`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		kv: []log0.KV{log0.Strings("message", "foo\n")},
		expected: `{
			"message":"foo\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "message" key is present`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo"),
		kv:    []log0.KV{log0.Strings("message", "bar")},
		expected: `{
			"message":"bar",
			"excerpt":"foo"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "message" key is present and with replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:    log0.Excerpt,
			Trunc:  120,
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar")},
		expected: `{
			"message":"bar",
			"excerpt":"foo",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "message" key is present and with replace input bytes and key`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar\n")},
		expected: `{
			"message":"bar\n",
			"excerpt":"foo",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is nil and "excerpt" key is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Excerpt,
			Trunc:  120,
		},
		kv: []log0.KV{log0.Strings("excerpt", "foo")},
		expected: `{
			"excerpt":"foo"
		}`,
	},
	{
		name: `default key is excerpt and bytes is nil and "excerpt" key is present and with replace`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		kv: []log0.KV{log0.Strings("excerpt", "foo\n")},
		expected: `{
			"excerpt":"foo\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" key is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Excerpt,
			Trunc:  120,
		},
		input: []byte("foo"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" key is present and with replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" key is present and with replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" key is present and with replace input bytes and rey`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("excerpt", "bar\n")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is nil and "excerpt" and "message" keys is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:    log0.Excerpt,
			Trunc:  120,
		},
		kv: []log0.KV{log0.Strings("message", "foo"), log0.Strings("excerpt", "bar")},
		expected: `{
			"message":"foo",
			"excerpt":"bar"
		}`,
	},
	{
		name: `default key is excerpt and bytes is nil and "excerpt" and "message" keys is present and replace keys`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		kv: []log0.KV{log0.Strings("message", "foo\n"), log0.Strings("excerpt", "bar\n")},
		expected: `{
			"message":"foo\n",
			"excerpt":"bar\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" and "message" keys is present`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:    log0.Excerpt,
			Trunc:  120,
		},
		input: []byte("foo"),
		kv:    []log0.KV{log0.Strings("message", "bar"), log0.Strings("excerpt", "xyz")},
		expected: `{
			"message":"bar",
			"excerpt":"xyz",
			"trail":"foo"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" and "message" keys is present and replace input bytes`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar"), log0.Strings("excerpt", "xyz")},
		expected: `{
			"message":"bar",
			"excerpt":"xyz",
			"trail":"foo\n"
		}`,
	},
	{
		name: `default key is excerpt and bytes is present and "excerpt" and "message" keys is present and replace input bytes and keys`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail")},
			Key:     log0.Excerpt,
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: []byte("foo\n"),
		kv:    []log0.KV{log0.Strings("message", "bar\n"), log0.Strings("excerpt", "xyz\n")},
		expected: `{
			"message":"bar\n",
			"excerpt":"xyz\n",
			"trail":"foo\n"
		}`,
	},
	{
		name: `bytes is nil and bytes "message" key with json`,
		line: line(),
		log:  dummy(),
		kv:   []log0.KV{log0.StringBytes("message", []byte(`{"foo":"bar"}`))},
		expected: `{
			"message":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		name: `bytes is nil and raw "message" key with json`,
		line: line(),
		log:  dummy(),
		kv:   []log0.KV{log0.StringRaw("message", []byte(`{"foo":"bar"}`))},
		expected: `{
			"message":{"foo":"bar"}
		}`,
	},
	{
		name: "bytes is nil and flag is long file",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		kv: []log0.KV{log0.Strings("foo", "bar")},
		expected: `{
			"foo":"bar"
		}`,
	},
	{
		name: "bytes is one char and flag is long file",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		input: []byte("a"),
		expected: `{
			"message":"a"
		}`,
	},
	{
		name: "bytes is two chars and flag is long file",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail"), log0.String("file")},
		},
		input: []byte("ab"),
		expected: `{
			"message":"ab",
			"file":"a"
		}`,
	},
	{
		name: "bytes is three chars and flag is long file",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail"), log0.String("file")},
		},
		input: []byte("abc"),
		expected: `{
			"message":"abc",
			"file":"ab"
		}`,
	},
	{
		name: "permanent kv overwritten by the additional kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("foo", "bar")},
		},
		kv: []log0.KV{log0.Strings("foo", "baz")},
		expected: `{
			"foo":"baz"
		}`,
	},
	{
		name: "permanent kv and first additional kv overwritten by the second additional kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("foo", "bar")},
		},
		kv: []log0.KV{
			log0.Strings("foo", "baz"),
			log0.Strings("foo", "xyz"),
		},
		expected: `{
			"foo":"xyz"
		}`,
	},
}

func TestWrite(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range WriteTestCases {
		tc := tc
		t.Run(fmt.Sprintf("io writer %s %d", tc.name, tc.line), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			l0, ok := tc.log.(*log0.Log)
			if !ok {
				t.Fatal("unexpected logger type")
			}

			buf, ok := l0.Output.(*bytes.Buffer)
			if !ok {
				t.Fatal("unexpected output type")
			}

			*buf = bytes.Buffer{}

			l := tc.log.Get(tc.kv...)
			defer l.Put()

			_, err := l.Write(tc.input)
			if err != nil {
				t.Fatalf("unexpected write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}

var FprintWriteTestCases = []struct {
	name      string
	line      int
	log       log0.Logger
	input     interface{}
	expected  string
	benchmark bool
}{
	{
		name: "readme example 1",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Marks:   [3][]byte{[]byte("…")},
			Trunc:   12,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		line:  line(),
		input: "Hello,\nWorld!",
		expected: `{
			"message":"Hello,\nWorld!",
			"excerpt":"Hello, World…"
		}`,
	},
	{
		name: "readme example 2",
		line: line(),
		log: func() log0.Logger {
			l1 := log0.GELF()
			l1.Output = &bytes.Buffer{}
			l := l1.Get(
				log0.Strings("version", "1.1"),
				log0.StringFunc("timestamp", func() log0.KV {
					t := time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC)
					return log0.Int64(t.Unix())
				}),
			)
			return l
		}(),
		input: "Hello,\nGELF!",
		expected: `{
			"version":"1.1",
			"short_message":"Hello, GELF!",
			"full_message":"Hello,\nGELF!",
			"timestamp":1602785340
		}`,
	},
	{
		name: "readme example 3.1",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		line:  line(),
		input: 3.21,
		expected: `{
			"message":"3.21"
		}`,
	},
	{
		name: "readme example 3.2",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		line:  line(),
		input: 123,
		expected: `{
			"message":"123"
		}`,
	},
	{
		name:  "string",
		line:  line(),
		log:   dummy(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name:  "integer type appears in the messages excerpt as a string",
		line:  line(),
		log:   dummy(),
		input: 123,
		expected: `{
			"message":"123"
		}`,
	},
	{
		name:  "float type appears in the messages excerpt as a string",
		line:  line(),
		log:   dummy(),
		input: 3.21,
		expected: `{
			"message":"3.21"
		}`,
	},
	{
		name:  "nil message",
		line:  line(),
		log:   dummy(),
		input: nil,
		expected: `{
			"message":"<nil>"
		}`,
	},
	// FIXME: not working(
	// {
	// 	name: "empty message",
	// 	line: line(),
	// 	log: &log0.Log{
	// 		Output: &bytes.Buffer{},
	// 		Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
	// 		Key:    log0.Original,
	// 		Trunc:  120,
	// 		Marks:  [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
	// 		Replace:  [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
	// 	},
	// 	input: "",
	// 	expected: `{
	// 		"message":"",
	// 		"excerpt":"_EMPTY_"
	// 	}`,
	// },
	{
		name:  "blank message",
		line:  line(),
		log:   dummy(),
		input: " ",
		expected: `{
	    "message":" ",
			"excerpt":"_BLANK_"
		}`,
	},
	{
		name:  "single quotes",
		line:  line(),
		log:   dummy(),
		input: "foo 'bar'",
		expected: `{
			"message":"foo 'bar'"
		}`,
	},
	{
		name:  "double quotes",
		line:  line(),
		log:   dummy(),
		input: `foo "bar"`,
		expected: `{
			"message":"foo \"bar\""
		}`,
	},
	{
		name:  `leading/trailing "spaces"`,
		line:  line(),
		log:   dummy(),
		input: " \n\tHello, World! \t\n",
		expected: `{
			"message":" \n\tHello, World! \t\n",
			"excerpt":"Hello, World!"
		}`,
	},
	{
		name:  "JSON string",
		line:  line(),
		log:   dummy(),
		input: `{"foo":"bar"}`,
		expected: `{
			"message":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		name: `"string" key with "foo" value`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("string", "foo")},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
		  "string": "foo"
		}`,
	},
	{
		name: `"integer" key with 123 value`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringInt("integer", 123)},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
		  "integer": 123
		}`,
	},
	{
		name: `"float" key with 3.21 value`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringFloat32("float", 3.21)},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
		  "float": 3.21
		}`,
	},
	{
		name:  "fmt.Fprint prints nil as <nil>",
		line:  line(),
		log:   dummy(),
		input: nil,
		expected: `{
			"message":"<nil>"
		}`,
	},
	{
		name:  "multiline string",
		line:  line(),
		log:   dummy(),
		input: "Hello,\nWorld\n!",
		expected: `{
			"message":"Hello,\nWorld\n!",
			"excerpt":"Hello, World !"
		}`,
	},
	{
		name:  "long string",
		line:  line(),
		log:   dummy(),
		input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"message":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"excerpt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliq…"
		}`,
	},
	{
		name:  "multiline long string with leading spaces",
		line:  line(),
		log:   dummy(),
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"excerpt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliq…"
		}`,
	},
	{
		name:  "multiline long string with leading spaces and multibyte character",
		line:  line(),
		log:   dummy(),
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
		expected: `{
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
			"excerpt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna Ää…"
		}`,
		benchmark: true,
	},
	{
		name: "zero maximum length",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
			Trunc:  0,
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: "without message key names",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{},
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"":"Hello, World!"
		}`,
	},
	{
		name: "only original message key name",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: "explicit byte slice as message excerpt key",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringBytes("excerpt", []byte("Explicit byte slice"))},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:  120,
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":"Explicit byte slice"
		}`,
	},
	{
		name: "explicit string as message excerpt key",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.Strings("excerpt", "Explicit string")},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:  120,
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":"Explicit string"
		}`,
	},
	{
		name: "explicit integer as message excerpt key",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringInt("excerpt", 42)},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:  120,
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":42
		}`,
	},
	{
		name: "explicit float as message excerpt key",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringFloat32("excerpt", 4.2)},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:  120,
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":4.2
		}`,
	},
	{
		name: "explicit boolean as message excerpt key",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringBool("excerpt", true)},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:  120,
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":true
		}`,
	},
	{
		name: "explicit rune slice as messages excerpt key",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV:     []log0.KV{log0.StringRunes("excerpt", []rune("Explicit rune slice"))},
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:  120,
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":"Explicit rune slice"
		}`,
	},
	{
		name: `dynamic "time" key`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV: []log0.KV{
				log0.StringFunc("time", func() log0.KV {
					t := time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC)
					return log0.String(t.String())
				}),
			},
			Keys: [4]encoding.TextMarshaler{log0.String("message")},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
			"time":"2020-10-15 18:09:00 +0000 UTC"
		}`,
	},
	{
		name: `"standard flag" do not respects file path`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.LstdFlags,
			Keys:   [4]encoding.TextMarshaler{log0.String("message")},
		},
		input: "path/to/file1:23: Hello, World!",
		expected: `{
			"message":"path/to/file1:23: Hello, World!"
		}`,
	},
	{
		name: `"long file" flag respects file path`,
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail"), log0.String("file")},
			Trunc:  120,
		},
		input: "path/to/file1:23: Hello, World!",
		expected: `{
			"message":"path/to/file1:23: Hello, World!",
			"excerpt":"Hello, World!",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name: "replace newline character by whitespace character",
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: "Hello,\nWorld!",
		expected: `{
			"message":"Hello,\nWorld!",
			"excerpt":"Hello, World!"
		}`,
	},
	{
		name: "remove exclamation marks",
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("!")}},
		},
		input: "Hello, World!!!",
		expected: `{
			"message":"Hello, World!!!",
			"excerpt":"Hello, World"
		}`,
	},
	{
		name: `replace word "World" by world "Work"`,
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt")},
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("World"), []byte("Work")}},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
			"excerpt":"Hello, Work!"
		}`,
	},
	{
		name: "ignore pointless replace",
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message")},
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{[]byte("!"), []byte("!")}},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: "ignore empty replace",
		line: line(),
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Keys:    [4]encoding.TextMarshaler{log0.String("message")},
			Trunc:   120,
			Replace: [][2][]byte{[2][]byte{}},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: "file path with empty message",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail"), log0.String("file")},
			Trunc:  120,
			Marks:  [3][]byte{[]byte("…"), []byte("_EMPTY_")},
		},
		input: "path/to/file1:23:",
		expected: `{
			"message":"path/to/file1:23:",
			"excerpt":"_EMPTY_",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name: "file path with blank message",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Flag:   log.Llongfile,
			Keys:   [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail"), log0.String("file")},
			Trunc:  120,
			Marks:  [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		},
		input: "path/to/file4:56:  ",
		expected: `{
			"message":"path/to/file4:56:  ",
			"excerpt":"_BLANK_",
			"file":"path/to/file4:56"
		}`,
	},
	{
		name: "GELF",
		line: line(),
		log: func() log0.Logger {
			l1 := log0.GELF()
			l1.Output = &bytes.Buffer{}
			l := l1.Get(
				log0.Strings("version", "1.1"),
				log0.Strings("host", "example.tld"),
				log0.StringFunc("timestamp", func() log0.KV {
					t := time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC)
					return log0.Int64(t.Unix())
				}),
			)
			return l
		}(),
		input: "Hello, GELF!",
		expected: `{
			"version":"1.1",
			"short_message":"Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340
		}`,
	},
	{
		name: "GELF with file path",
		line: line(),
		log: func() log0.Logger {
			l1 := log0.GELF()
			l1.Output = &bytes.Buffer{}
			l1.Flag = log.Llongfile
			l := l1.Get(
				log0.Strings("version", "1.1"),
				log0.Strings("host", "example.tld"),
				log0.StringFunc("timestamp", func() log0.KV {
					t := time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC)
					return log0.Int64(t.Unix())
				}),
			)
			return l
		}(),
		input: "path/to/file7:89: Hello, GELF!",
		expected: `{
			"version":"1.1",
			"short_message":"Hello, GELF!",
			"full_message":"path/to/file7:89: Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340,
			"_file":"path/to/file7:89"
		}`,
	},
}

func TestFprintWrite(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range FprintWriteTestCases {
		tc := tc
		t.Run(fmt.Sprintf("io writer via fmt fprint %s %d", tc.name, tc.line), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			l, ok := tc.log.(*log0.Log)
			if !ok {
				t.Fatal("unexpected logger type")
			}

			var buf bytes.Buffer
			l.Output = &buf

			_, err := fmt.Fprint(l, tc.input)
			if err != nil {
				t.Fatalf("unexpected write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}

func BenchmarkLog0(b *testing.B) {
	for _, tc := range WriteTestCases {
		if !tc.benchmark {
			continue
		}
		b.Run(fmt.Sprintf("io.Writer %d", tc.line), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := tc.log.Get(tc.kv...)
				_, err := l.Write(tc.input)
				if err != nil {
					fmt.Println(err)
				}
				l.Put()
			}
		})
	}

	for _, tc := range FprintWriteTestCases {
		if !tc.benchmark {
			continue
		}
		b.Run(fmt.Sprintf("fmt.Fprint io.Writer %d", tc.line), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := fmt.Fprint(tc.log, tc.input)
				if err != nil {
					fmt.Println(err)
				}
			}
		})
	}
}

var dummy = func() log0.Logger {
	return &log0.Log{
		Output:  &bytes.Buffer{},
		Keys:    [4]encoding.TextMarshaler{log0.String("message"), log0.String("excerpt"), log0.String("trail"), log0.String("file")},
		Key:     log0.Original,
		Trunc:   120,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
	}
}

func TestLogWriteTrailingNewLine(t *testing.T) {
	var buf bytes.Buffer

	l := &log0.Log{Output: &buf}

	_, err := l.Write([]byte("Hello, Wrold!"))
	if err != nil {
		t.Fatalf("unexpected write error: %s", err)
	}

	if buf.Bytes()[len(buf.Bytes())-1] != '\n' {
		t.Errorf("trailing new line expected but not present: %q", buf.String())
	}
}

var TruncateTestCases = []struct {
	name      string
	line      int
	log       log0.Logger
	input     []byte
	expected  []byte
	benchmark bool
}{
	{
		name: "do nothing",
		log: &log0.Log{
			Output: &bytes.Buffer{},
		},
		line:     line(),
		input:    []byte("Hello,\nWorld!"),
		expected: []byte("Hello,\nWorld!"),
	},
	{
		name: "truncate last character",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Trunc:  12,
		},
		line:     line(),
		input:    []byte("Hello, World!"),
		expected: []byte("Hello, World"),
	},
	{
		name: "truncate last character and places ellipsis instead",
		log: &log0.Log{
			Output: &bytes.Buffer{},
			Trunc:  12,
			Marks:  [3][]byte{[]byte("…")},
		},
		line:     line(),
		input:    []byte("Hello, World!"),
		expected: []byte("Hello, World…"),
	},
	{
		name: "replace new lines by spaces",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		line:     line(),
		input:    []byte("Hello\n,\nWorld\n!"),
		expected: []byte("Hello , World !"),
	},
	{
		name: "replace new lines by empty string",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte("")}},
		},
		line:     line(),
		input:    []byte("Hello\n,\nWorld\n!"),
		expected: []byte("Hello,World!"),
	},
	{
		name: "remove new lines",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Replace: [][2][]byte{[2][]byte{[]byte("\n")}},
		},
		line:     line(),
		input:    []byte("Hello\n,\nWorld\n!"),
		expected: []byte("Hello,World!"),
	},
	{
		name: "replace three characters by one",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Replace: [][2][]byte{[2][]byte{[]byte("foo"), []byte("f")}, [2][]byte{[]byte("bar"), []byte("b")}},
		},
		line:     line(),
		input:    []byte("foobar"),
		expected: []byte("fb"),
	},
	{
		name: "replace one characters by three",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Replace: [][2][]byte{[2][]byte{[]byte("f"), []byte("foo")}, [2][]byte{[]byte("b"), []byte("bar")}},
		},
		line:     line(),
		input:    []byte("fb"),
		expected: []byte("foobar"),
	},
	{
		name: "remove three characters",
		log: &log0.Log{
			Output:  &bytes.Buffer{},
			Replace: [][2][]byte{[2][]byte{[]byte("foo")}, [2][]byte{[]byte("bar")}},
		},
		line:     line(),
		input:    []byte("foobar foobar"),
		expected: []byte(" "),
	},
}

func TestTruncate(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range TruncateTestCases {
		tc := tc
		t.Run(fmt.Sprintf("truncate %s %d", tc.name, tc.line), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			l0, ok := tc.log.(*log0.Log)
			if !ok {
				t.Fatal("unexpected logger type")
			}

			n := len(tc.input) + 10*10
			for _, m := range l0.Marks {
				if n < len(m) {
					n = len(m)
				}
			}

			excerpt := make([]byte, n)

			n, err := l0.Truncate(excerpt, tc.input)
			if err != nil {
				t.Fatalf("unexpected write error: %s", err)
			}

			excerpt = excerpt[:n]

			if !bytes.Equal(excerpt, tc.expected) {
				t.Errorf("unexpected excerpt, expected: %q, received: %q %s", tc.expected, excerpt, linkToExample)
			}
		})
	}
}

var NewTestCases = []struct {
	name      string
	line      int
	log       log0.Logger
	kv        []log0.KV
	expected  string
	benchmark bool
}{
	{
		name: "one kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV: []log0.KV{
				log0.Strings("foo", "bar"),
			},
		},
		expected: `{
			"foo":"bar"
		}`,
	},
	{
		name: "two kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV: []log0.KV{
				log0.Strings("foo", "bar"),
				log0.Strings("baz", "xyz"),
			},
		},
		expected: `{
			"foo":"bar",
			"baz":"xyz"
		}`,
	},
	{
		name: "one additional kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
		},
		kv: []log0.KV{
			log0.Strings("baz", "xyz"),
		},
		expected: `{
			"baz":"xyz"
		}`,
	},
	{
		name: "two additional kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
		},
		kv: []log0.KV{
			log0.Strings("foo", "bar"),
			log0.Strings("baz", "xyz"),
		},
		expected: `{
			"foo":"bar",
			"baz":"xyz"
		}`,
	},
	{
		name: "one kv with additional one kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV: []log0.KV{
				log0.Strings("foo", "bar"),
			},
		},
		kv: []log0.KV{
			log0.Strings("baz", "xyz"),
		},
		expected: `{
			"foo":"bar",
			"baz":"xyz"
		}`,
	},
	{
		name: "two kv with two additional kv",
		line: line(),
		log: &log0.Log{
			Output: &bytes.Buffer{},
			KV: []log0.KV{
				log0.Strings("foo", "bar"),
				log0.Strings("abc", "dfg"),
			},
		},
		kv: []log0.KV{
			log0.Strings("baz", "xyz"),
			log0.Strings("hjk", "lmn"),
		},
		expected: `{
			"foo":"bar",
			"abc":"dfg",
			"baz":"xyz",
			"hjk":"lmn"
		}`,
	},
}

func TestNew(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range NewTestCases {
		tc := tc
		t.Run(fmt.Sprintf("with %s %d", tc.name, tc.line), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			l0, ok := tc.log.(*log0.Log)
			if !ok {
				t.Fatal("unexpected logger type")
			}

			buf, ok := l0.Output.(*bytes.Buffer)
			if !ok {
				t.Fatal("unexpected output type")
			}

			*buf = bytes.Buffer{}

			l := tc.log.Get(tc.kv...)
			defer l.Put()

			_, err := l.Write(nil)
			if err != nil {
				t.Fatalf("unexpected write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}

var SeverityTestCases = []struct {
	name       string
	line       int
	log        log0.Logger
	severities []string
	kv         []log0.KV
	expected   string
	benchmark  bool
}{
	{
		name: "just severity 7",
		line: line(),
		log: &log0.Log{
			Output:   &bytes.Buffer{},
			Severity: func(severity string) io.Writer { return nil },
		},
		severities: []string{"7"},
		expected: `{
			"severity":"7"
		}`,
	},
	{
		name: "just severity 7 next to the key-value pair",
		line: line(),
		log: &log0.Log{
			Output:   &bytes.Buffer{},
			Severity: func(severity string) io.Writer { return nil },
			KV: []log0.KV{
				log0.Strings("foo", "bar"),
			},
		},
		severities: []string{"7"},
		expected: `{
			"foo":"bar",
			"severity":"7"
		}`,
	},
	{
		name: "two consecutive severities call 7 and 6",
		line: line(),
		log: &log0.Log{
			Output:   &bytes.Buffer{},
			Severity: func(severity string) io.Writer { return nil },
		},
		severities: []string{"7", "6"},
		expected: `{
			"severity":"6"
		}`,
	},
	{
		name: "severity 7 overwrites severity 42 from key-value pair",
		line: line(),
		log: &log0.Log{
			Output:   &bytes.Buffer{},
			Severity: func(severity string) io.Writer { return nil },
			KV: []log0.KV{
				log0.Strings("severity", "42"),
			},
		},
		severities: []string{"7"},
		expected: `{
			"severity":"7"
		}`,
	},
}

func TestSeverity(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range SeverityTestCases {
		tc := tc
		t.Run(fmt.Sprintf("with %s %d", tc.name, tc.line), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			l0, ok := tc.log.(*log0.Log)
			if !ok {
				t.Fatal("unexpected logger type")
			}

			buf, ok := l0.Output.(*bytes.Buffer)
			if !ok {
				t.Fatal("unexpected output type")
			}

			*buf = bytes.Buffer{}

			l := tc.log.Get(tc.kv...)
			defer l.Put()

			for _, sev := range tc.severities {
				l = l.Get(log0.StringSeverity("severity", sev))
			}

			_, err := l.Write(nil)
			if err != nil {
				t.Fatalf("unexpected write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}
