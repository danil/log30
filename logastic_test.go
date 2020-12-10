package logastic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/danil/logastic"
	"github.com/kinbiko/jsonassert"
)

var (
	dummy = logastic.Log{
		Trunc:   120,
		Keys:    [4]string{"message", "excerpt", "file", "host"},
		Key:     logastic.Original,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
	}

	gelf = func() logastic.Log {
		l := logastic.GELF()
		l.Funcs = map[string]func() json.Marshaler{"timestamp": func() json.Marshaler {
			return logastic.Int64(time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix())
		}}
		return l
	}()
)

func TestWriteTrailingNewLine(t *testing.T) {
	var buf bytes.Buffer

	mjl := logastic.Log{Output: &buf}

	_, err := mjl.Write([]byte("Hello, Wrold!"))
	if err != nil {
		t.Fatalf("write error: %s", err)
	}

	if buf.Bytes()[len(buf.Bytes())-1] != '\n' {
		t.Errorf("trailing new line expected but not present: %q", buf.String())
	}
}

var WriteTestCases = []struct {
	name      string
	line      int
	log       logastic.Log
	input     interface{}
	expected  string
	benchmark bool
}{
	{
		name: "readme example 1",
		log: logastic.Log{
			Trunc:   12,
			Keys:    [4]string{"message", "excerpt"},
			Marks:   [3][]byte{[]byte("…")},
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
		name:  "readme example 2",
		line:  line(),
		log:   gelf,
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
		log: logastic.Log{
			Keys: [4]string{"message"},
		},
		line:  line(),
		input: 3.21,
		expected: `{
			"message":"3.21"
		}`,
	},
	{
		name: "readme example 3.2",
		log: logastic.Log{
			Keys: [4]string{"message"},
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
		log:   dummy,
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name:  "integer type appears in the messages excerpt as a string",
		line:  line(),
		log:   dummy,
		input: 123,
		expected: `{
			"message":"123"
		}`,
	},
	{
		name:  "float type appears in the messages excerpt as a string",
		line:  line(),
		log:   dummy,
		input: 3.21,
		expected: `{
			"message":"3.21"
		}`,
	},
	{
		name:  "empty message",
		line:  line(),
		log:   dummy,
		input: "",
		expected: `{
	    "message":"",
			"excerpt":"_EMPTY_"
		}`,
	},
	{
		name:  "blank message",
		line:  line(),
		log:   dummy,
		input: " ",
		expected: `{
	    "message":" ",
			"excerpt":"_BLANK_"
		}`,
	},
	{
		name:  "single quotes",
		line:  line(),
		log:   dummy,
		input: "foo 'bar'",
		expected: `{
			"message":"foo 'bar'"
		}`,
	},
	{
		name:  "double quotes",
		line:  line(),
		log:   dummy,
		input: `foo "bar"`,
		expected: `{
			"message":"foo \"bar\""
		}`,
	},
	{
		name:  `leading/trailing "spaces"`,
		line:  line(),
		log:   dummy,
		input: " \n\tHello, World! \t\n",
		expected: `{
			"message":" \n\tHello, World! \t\n",
			"excerpt":"Hello, World!"
		}`,
	},
	{
		name:  "JSON string",
		line:  line(),
		log:   dummy,
		input: `{"foo":"bar"}`,
		expected: `{
			"message":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		name: `"string" key with "foo" value`,
		line: line(),
		log: logastic.Log{
			KV:   map[string]json.Marshaler{"string": logastic.String("foo")},
			Keys: [4]string{"message"},
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
		log: logastic.Log{
			KV:   map[string]json.Marshaler{"integer": logastic.Int(123)},
			Keys: [4]string{"message"},
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
		log: logastic.Log{
			KV:   map[string]json.Marshaler{"float": logastic.Float32(3.21)},
			Keys: [4]string{"message"},
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
		log:   dummy,
		input: nil,
		expected: `{
			"message":"<nil>"
		}`,
	},
	{
		name:  "multiline string",
		line:  line(),
		log:   dummy,
		input: "Hello,\nWorld\n!",
		expected: `{
			"message":"Hello,\nWorld\n!",
			"excerpt":"Hello, World !"
		}`,
	},
	{
		name:  "long string",
		line:  line(),
		log:   dummy,
		input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"message":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"excerpt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliq…"
		}`,
	},
	{
		name:  "multiline long string with leading spaces",
		line:  line(),
		log:   dummy,
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"excerpt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliq…"
		}`,
	},
	{
		name:  "multiline long string with leading spaces and multibyte character",
		line:  line(),
		log:   dummy,
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
		expected: `{
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
			"excerpt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna Ää…"
		}`,
		benchmark: true,
	},
	{
		name: "zero maximum length",
		log: logastic.Log{
			Keys:  [4]string{"message"},
			Trunc: 0,
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: "without message key names",
		log: logastic.Log{
			Keys: [4]string{},
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"":"Hello, World!"
		}`,
	},
	{
		name: "only original message key name",
		log: logastic.Log{
			Keys: [4]string{"message"},
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: `explicit byte slice as message excerpt key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"excerpt": logastic.Bytes([]byte("Explicit byte slice"))},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":"Explicit byte slice"
		}`,
	},
	{
		name: `explicit string as message excerpt key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"excerpt": logastic.String("Explicit string")},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":"Explicit string"
		}`,
	},
	{
		name: `explicit integer as message excerpt key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"excerpt": logastic.Int(42)},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":42
		}`,
	},
	{
		name: `explicit float as message excerpt key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"excerpt": logastic.Float32(4.2)},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":4.2
		}`,
	},
	{
		name: `explicit boolean as message excerpt key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"excerpt": logastic.Bool(true)},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"excerpt":true
		}`,
	},
	{
		name: `explicit rune slice as messages excerpt key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"excerpt": logastic.Runes([]rune("Explicit rune slice"))},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
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
		log: logastic.Log{
			Funcs: map[string]func() json.Marshaler{"time": func() json.Marshaler {
				return logastic.String(time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).String())
			}},
			Keys: [4]string{"message"},
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
		log: logastic.Log{
			Flag: log.LstdFlags,
			Keys: [4]string{"message"},
		},
		input: "path/to/file1:23: Hello, World!",
		expected: `{
			"message":"path/to/file1:23: Hello, World!"
		}`,
	},
	{
		name: `"long file" flag respects file path`,
		line: line(),
		log: logastic.Log{
			Flag:  log.Llongfile,
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt", "file"},
		},
		input: "path/to/file1:23: Hello, World!",
		expected: `{
			"message":"path/to/file1:23: Hello, World!",
			"excerpt":"Hello, World!",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name: `replace newline character by whitespace character`,
		line: line(),
		log: logastic.Log{
			Trunc:   120,
			Keys:    [4]string{"message", "excerpt"},
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		input: "Hello,\nWorld!",
		expected: `{
			"message":"Hello,\nWorld!",
			"excerpt":"Hello, World!"
		}`,
	},
	{
		name: `remove exclamation marks`,
		line: line(),
		log: logastic.Log{
			Trunc:   120,
			Keys:    [4]string{"message", "excerpt"},
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
		log: logastic.Log{
			Trunc:   120,
			Keys:    [4]string{"message", "excerpt"},
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
		log: logastic.Log{
			Trunc:   120,
			Keys:    [4]string{"message"},
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
		log: logastic.Log{
			Trunc:   120,
			Keys:    [4]string{"message"},
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
		log: logastic.Log{
			Flag:  log.Llongfile,
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt", "file"},
			Marks: [3][]byte{[]byte("…"), []byte("_EMPTY_")},
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
		log: logastic.Log{
			Flag:  log.Llongfile,
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt", "file"},
			Marks: [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		},
		input: "path/to/file4:56:  ",
		expected: `{
			"message":"path/to/file4:56:  ",
			"excerpt":"_BLANK_",
			"file":"path/to/file4:56"
		}`,
	},
	{
		name: `"magic" host key`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"host": logastic.String("example.tld")},
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt", "file", "host"},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
			"excerpt":"example.tld Hello, World!",
			"host":"example.tld"
		}`,
	},
	{
		name: "GELF",
		line: line(),
		log: func() logastic.Log {
			l := logastic.GELF()
			l.Funcs = map[string]func() json.Marshaler{"timestamp": func() json.Marshaler {
				return logastic.Int64(time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix())
			}}
			l.KV = map[string]json.Marshaler{"version": logastic.String("1.1"), "host": logastic.String("example.tld")}
			return l
		}(),
		input: "Hello, GELF!",
		expected: `{
			"version":"1.1",
			"short_message":"example.tld Hello, GELF!",
			"full_message":"Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340
		}`,
	},
	{
		name: "GELF with file path",
		line: line(),
		log: func() logastic.Log {
			l := logastic.GELF()
			l.Flag = log.Llongfile
			l.Funcs = map[string]func() json.Marshaler{"timestamp": func() json.Marshaler {
				return logastic.Int64(time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix())
			}}
			l.KV = map[string]json.Marshaler{"version": logastic.String("1.1"), "host": logastic.String("example.tld")}
			return l
		}(),
		input: "path/to/file7:89: Hello, GELF!",
		expected: `{
			"version":"1.1",
			"short_message":"example.tld Hello, GELF!",
			"full_message":"path/to/file7:89: Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340,
			"_file":"path/to/file7:89"
		}`,
	},
}

func TestWrite(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range WriteTestCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			var buf bytes.Buffer

			tc.log.Output = &buf

			_, err := fmt.Fprint(tc.log, tc.input)
			if err != nil {
				t.Fatalf("write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}

func BenchmarkLogastic(b *testing.B) {
	for _, tc := range WriteTestCases {
		if !tc.benchmark {
			continue
		}
		b.Run(strconv.Itoa(tc.line), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer

				tc.log.Output = &buf

				_, err := fmt.Fprint(tc.log, tc.input)
				if err != nil {
					fmt.Println(err)
				}
			}
		})
	}
}

var LogTestCases = []struct {
	name     string
	line     int
	log      logastic.Log
	bytes    []byte
	kv       map[string]json.Marshaler
	expected string
}{
	{
		name:  "nil",
		line:  line(),
		bytes: nil,
		log:   dummy,
		expected: `{
	    "message":null,
			"excerpt":"_EMPTY_"
		}`,
	},
	{
		name: `"string" key with "foo" value and "string" key with "bar" value`,
		line: line(),
		log: logastic.Log{
			Trunc: 120,
			KV:    map[string]json.Marshaler{"string": logastic.String("foo")},
			Keys:  [4]string{"message"},
		},
		bytes: []byte("Hello, World!"),
		kv:    map[string]json.Marshaler{"string": logastic.String("bar")},
		expected: `{
			"message":"Hello, World!",
		  "string": "bar"
		}`,
	},
	{
		name:  `key-values is nil`,
		line:  line(),
		log:   dummy,
		bytes: []byte("Hello, World!"),
		kv:    nil,
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: `bytes appends to the "message" key with "string value"`,
		line: line(),
		log: logastic.Log{
			KV:      map[string]json.Marshaler{"message": logastic.String("string value")},
			Trunc:   120,
			Keys:    [4]string{"message", "excerpt"},
			Replace: [][2][]byte{[2][]byte{[]byte("\n"), []byte(" ")}},
		},
		bytes: []byte("\nHello, World!"),
		expected: `{
			"message":"\"string value\"\nHello, World!",
			"excerpt":"\"string value\" Hello, World!"
		}`,
	},
	{
		name:  `bytes appends to the "message" key with "string value"`,
		line:  line(),
		log:   dummy,
		bytes: []byte("\nHello, World!"),
		kv:    map[string]json.Marshaler{"message": logastic.String("string value")},
		expected: `{
			"message":"\"string value\"\nHello, World!",
			"excerpt":"\"string value\" Hello, World!"
		}`,
	},
	{
		name: `bytes is nil and "message" key with "string value"`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]json.Marshaler{"message": logastic.String("string value")},
			Trunc: 120,
			Keys:  [4]string{"message"},
		},
		bytes: nil,
		expected: `{
			"message":"string value"
		}`,
	},
	{
		name:  `bytes is nil and "message" key with "string value"`,
		line:  line(),
		log:   dummy,
		bytes: nil,
		kv:    map[string]json.Marshaler{"message": logastic.String("string value")},
		expected: `{
			"message":"string value"
		}`,
	},
	{
		name:  `bytes appends to the integer key "message"`,
		line:  line(),
		log:   dummy,
		bytes: []byte("\nHello, World!"),
		kv:    map[string]json.Marshaler{"message": logastic.Int(1)},
		expected: `{
			"message":"1\nHello, World!",
			"excerpt":"1 Hello, World!"
		}`,
	},
	{
		name:  `bytes appends to the float 32 bit key "message"`,
		line:  line(),
		log:   dummy,
		bytes: []byte("\nHello, World!"),
		kv:    map[string]json.Marshaler{"message": logastic.Float32(4.2)},
		expected: `{
			"message":"4.2\nHello, World!",
			"excerpt":"4.2 Hello, World!"
		}`,
	},
	{
		name:  `bytes appends to the float 64 bit key "message"`,
		line:  line(),
		log:   dummy,
		bytes: []byte("\nHello, World!"),
		kv:    map[string]json.Marshaler{"message": logastic.Float64(4.2)},
		expected: `{
			"message":"4.2\nHello, World!",
			"excerpt":"4.2 Hello, World!"
		}`,
	},
	{
		name:  `bytes appends to the boolean key "message"`,
		line:  line(),
		log:   dummy,
		bytes: []byte("\nHello, World!"),
		kv:    map[string]json.Marshaler{"message": logastic.Bool(true)},
		expected: `{
			"message":"true\nHello, World!",
			"excerpt":"true Hello, World!"
		}`,
	},
	{
		name:  `bytes do not appends to the nil key "message"`,
		line:  line(),
		log:   dummy,
		bytes: []byte("Hello, World!"),
		kv:    map[string]json.Marshaler{"message": nil},
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: `bytes is nil and string "message" key and default key is excerpt`,
		line: line(),
		log: logastic.Log{
			Trunc: 120,
			Keys:  [4]string{"message", "excerpt"},
			Key:   logastic.Excerpt,
		},
		bytes: nil,
		kv:    map[string]json.Marshaler{"message": logastic.String("Hello, World!")},
		expected: `{
			"excerpt":"\"Hello, World!\""
		}`,
	},
	{
		name:  `bytes is nil and bytes "message" key with json`,
		line:  line(),
		log:   dummy,
		bytes: nil,
		kv:    map[string]json.Marshaler{"message": logastic.Bytes([]byte(`{"foo":"bar"}`))},
		expected: `{
			"message":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		name:  `bytes is nil and raw "message" key with json`,
		line:  line(),
		log:   dummy,
		bytes: nil,
		kv:    map[string]json.Marshaler{"message": logastic.Raw([]byte(`{"foo":"bar"}`))},
		expected: `{
			"message":{"foo":"bar"}
		}`,
	},
}

func TestLog(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range LogTestCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			var buf bytes.Buffer

			tc.log.Output = &buf

			_, err := tc.log.Log(tc.kv, tc.bytes...)
			if err != nil {
				t.Fatalf("write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}
