package logastic_test

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/danil/logastic"
	"github.com/kinbiko/jsonassert"
)

var (
	pool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

	dummy = logastic.Log{
		Trunc:   120,
		Keys:    [4]string{"message", "preview", "file", "host"},
		Key:     logastic.FullKey,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][]byte{[]byte("\n"), []byte(" ")},
	}

	gelf = func() logastic.Log {
		l := logastic.GELF()
		l.Funcs = map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }}
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

func line() int { _, _, l, _ := runtime.Caller(1); return l }

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
			Keys:    [4]string{"message", "preview"},
			Marks:   [3][]byte{[]byte("…")},
			Replace: [][]byte{[]byte("\n"), []byte(" ")},
		},
		line:  line(),
		input: "Hello,\nWorld!",
		expected: `{
			"preview":"Hello, World…",
			"message":"Hello,\nWorld!"
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
		name:  "integer type appears in the short messages as a string",
		line:  line(),
		log:   dummy,
		input: 123,
		expected: `{
			"message":"123"
		}`,
	},
	{
		name:  "float type appears in the short messages as a string",
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
			"preview":"_EMPTY_"
		}`,
	},
	{
		name:  "blank message",
		line:  line(),
		log:   dummy,
		input: " ",
		expected: `{
	    "message":" ",
			"preview":"_BLANK_"
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
			"preview":"Hello, World!"
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
		name: `"string" field with "foo" value`,
		line: line(),
		log: logastic.Log{
			KV:   map[string]interface{}{"string": "foo"},
			Keys: [4]string{"message"},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
		  "string": "foo"
		}`,
	},
	{
		name: `"integer" field with 123 value`,
		line: line(),
		log: logastic.Log{
			KV:   map[string]interface{}{"integer": 123},
			Keys: [4]string{"message"},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
		  "integer": 123
		}`,
	},
	{
		name: `"float" field with 3.21 value`,
		line: line(),
		log: logastic.Log{
			KV:   map[string]interface{}{"float": 3.21},
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
		input: "Hello,\nWorld!",
		expected: `{
			"message":"Hello,\nWorld!",
			"preview":"Hello, World!"
		}`,
	},
	{
		name:  "long string",
		line:  line(),
		log:   dummy,
		input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"message":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"preview":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliq…"
		}`,
	},
	{
		name:  "multiline long string with leading spaces",
		line:  line(),
		log:   dummy,
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"preview":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliq…"
		}`,
	},
	{
		name:  "multiline long string with leading spaces and multibyte character",
		line:  line(),
		log:   dummy,
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
		expected: `{
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
			"preview":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna Ää…"
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
		name: "only full message key name",
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
		name: `explicit byte slice as short message field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"preview": []byte("Explicit byte slice")},
			Trunc: 120,
			Keys:  [4]string{"message", "preview"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"preview":"Explicit byte slice"
		}`,
	},
	{
		name: `explicit string as short message field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"preview": "Explicit string"},
			Trunc: 120,
			Keys:  [4]string{"message", "preview"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"preview":"Explicit string"
		}`,
	},
	{
		name: `explicit integer as short message field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"preview": 42},
			Trunc: 120,
			Keys:  [4]string{"message", "preview"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"preview":42
		}`,
	},
	{
		name: `explicit float as short message field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"preview": 4.2},
			Trunc: 120,
			Keys:  [4]string{"message", "preview"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"preview":4.2
		}`,
	},
	{
		name: `explicit boolean as short message field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"preview": true},
			Trunc: 120,
			Keys:  [4]string{"message", "preview"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"preview":true
		}`,
	},
	{
		name: `explicit rune slice as short message field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"preview": []rune("Explicit rune slice")},
			Trunc: 120,
			Keys:  [4]string{"message", "preview"},
		},
		input: "Hello, World!",
		expected: `{
		  "message": "Hello, World!",
			"preview":"Explicit rune slice"
		}`,
	},
	{
		name: "dynamic field",
		line: line(),
		log: logastic.Log{
			Funcs: map[string]func() interface{}{"time": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).String() }},
			Keys:  [4]string{"message"},
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
			Keys:  [4]string{"message", "preview", "file"},
		},
		input: "path/to/file1:23: Hello, World!",
		expected: `{
			"message":"path/to/file1:23: Hello, World!",
			"preview":"Hello, World!",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name: "file path with empty message",
		line: line(),
		log: logastic.Log{
			Flag:  log.Llongfile,
			Trunc: 120,
			Keys:  [4]string{"message", "preview", "file"},
			Marks: [3][]byte{[]byte("…"), []byte("_EMPTY_")},
		},
		input: "path/to/file1:23:",
		expected: `{
			"message":"path/to/file1:23:",
			"preview":"_EMPTY_",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name: "file path with blank message",
		line: line(),
		log: logastic.Log{
			Flag:  log.Llongfile,
			Trunc: 120,
			Keys:  [4]string{"message", "preview", "file"},
			Marks: [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		},
		input: "path/to/file4:56:  ",
		expected: `{
			"message":"path/to/file4:56:  ",
			"preview":"_BLANK_",
			"file":"path/to/file4:56"
		}`,
	},
	{
		name: `"magic" host field`,
		line: line(),
		log: logastic.Log{
			KV:    map[string]interface{}{"host": "example.tld"},
			Trunc: 120,
			Keys:  [4]string{"message", "preview", "file", "host"},
		},
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
			"preview":"example.tld Hello, World!",
			"host":"example.tld"
		}`,
	},
	{
		name: "GELF",
		line: line(),
		log: func() logastic.Log {
			l := logastic.GELF()
			l.Funcs = map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }}
			l.KV = map[string]interface{}{"version": "1.1", "host": "example.tld"}
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
			l.Funcs = map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }}
			l.KV = map[string]interface{}{"version": "1.1", "host": "example.tld"}
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

			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			defer pool.Put(buf)

			tc.log.Output = buf

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
				buf := pool.Get().(*bytes.Buffer)
				buf.Reset()
				defer pool.Put(buf)

				tc.log.Output = buf

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
	input    []byte
	kv       map[string]interface{}
	expected string
}{
	{
		name:  "nil",
		line:  line(),
		input: nil,
		log:   dummy,
		expected: `{
	    "message":"",
			"preview":"_EMPTY_"
		}`,
	},
	{
		name: `"string" field with "foo" value and "string" key with "bar" value`,
		line: line(),
		log: logastic.Log{
			Trunc: 120,
			KV:    map[string]interface{}{"string": "foo"},
			Keys:  [4]string{"message"},
		},
		input: []byte("Hello, World!"),
		kv:    map[string]interface{}{"string": "bar"},
		expected: `{
			"message":"Hello, World!",
		  "string": "bar"
		}`,
	},
	{
		name:  `key-values is nil`,
		line:  line(),
		log:   dummy,
		input: []byte("Hello, World!"),
		kv:    nil,
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name: `input appends to the message field value "string"`,
		line: line(),
		log: logastic.Log{
			KV:      map[string]interface{}{"message": "field string value"},
			Trunc:   120,
			Keys:    [4]string{"message", "preview"},
			Replace: [][]byte{[]byte("\n"), []byte(" ")},
		},
		input: []byte("\nHello, World!"),
		expected: `{
			"message":"field string value\nHello, World!",
			"preview":"field string value Hello, World!"
		}`,
	},
	{
		name:  `input appends to the message key-field value "string"`,
		line:  line(),
		log:   dummy,
		input: []byte("\nHello, World!"),
		kv:    map[string]interface{}{"message": "field string value"},
		expected: `{
			"message":"field string value\nHello, World!",
			"preview":"field string value Hello, World!"
		}`,
	},
	{
		name: `input is nil and message field value is "string"`,
		line: line(),
		log: logastic.Log{
			KV:      map[string]interface{}{"message": "string"},
			Trunc:   120,
			Keys:    [4]string{"message", "preview"},
			Replace: [][]byte{[]byte("\n"), []byte(" ")},
		},
		input: nil,
		expected: `{
			"message":"string"
		}`,
	},
	{
		name:  `input is nil and message key-value is "string"`,
		line:  line(),
		log:   dummy,
		input: nil,
		kv:    map[string]interface{}{"message": "string"},
		expected: `{
			"message":"string"
		}`,
	},
	{
		name:  `input appends to the integer key-value "message"`,
		line:  line(),
		log:   dummy,
		input: []byte("\nHello, World!"),
		kv:    map[string]interface{}{"message": 1},
		expected: `{
			"message":"1\nHello, World!",
			"preview":"1 Hello, World!"
		}`,
	},
	{
		name:  `input appends to the float key-value "message"`,
		line:  line(),
		log:   dummy,
		input: []byte("\nHello, World!"),
		kv:    map[string]interface{}{"message": 2.1},
		expected: `{
			"message":"2.1\nHello, World!",
			"preview":"2.1 Hello, World!"
		}`,
	},
	{
		name:  `input appends to the boolean key-value "message"`,
		line:  line(),
		log:   dummy,
		input: []byte("\nHello, World!"),
		kv:    map[string]interface{}{"message": true},
		expected: `{
			"message":"true\nHello, World!",
			"preview":"true Hello, World!"
		}`,
	},
	{
		name:  `input do not appends to the nil key-value "message"`,
		line:  line(),
		log:   dummy,
		input: []byte("Hello, World!"),
		kv:    map[string]interface{}{"message": nil},
		expected: `{
			"message":"Hello, World!"
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

			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			defer pool.Put(buf)

			tc.log.Output = buf

			_, err := tc.log.Log(tc.input, tc.kv)
			if err != nil {
				t.Fatalf("write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}

type testprinter struct {
	t    *testing.T
	link string
}

func (p testprinter) Errorf(msg string, args ...interface{}) {
	p.t.Errorf(p.link+"\n"+msg, args...)
}
