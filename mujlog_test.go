package mujlog_test

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/danil/mujlog"
	"github.com/kinbiko/jsonassert"
)

var (
	pool  = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
	dummy = mujlog.Log{
		Max:     120,
		Keys:    [4]string{"message", "preview", "file", "host"},
		Key:     mujlog.FullKey,
		Marks:   [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")},
		Replace: [][]byte{[]byte("\n"), []byte(" ")},
	}
)

func TestMujlogWriteTrailingNewLine(t *testing.T) {
	var buf bytes.Buffer

	mjl := mujlog.Log{Output: &buf}

	_, err := mjl.Write([]byte("Hello, Wrold!"))
	if err != nil {
		t.Fatalf("unexpected mujlog write error: %s", err)
	}

	if buf.Bytes()[len(buf.Bytes())-1] != '\n' {
		t.Errorf("trailing new line expected but not present: %q", buf.String())
	}
}

func line() int { _, _, l, _ := runtime.Caller(1); return l }

var WriteTestCases = []struct {
	name      string
	line      int
	log       mujlog.Log
	input     interface{}
	flag      int
	kvs       map[string]interface{}
	funcs     map[string]func() interface{}
	expected  string
	benchmark bool
}{
	{
		name: "first readme example",
		log: mujlog.Log{
			Keys:    [4]string{"message", "preview"},
			Marks:   [3][]byte{[]byte("…")},
			Max:     12,
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
		name:  `"string" field with "foo" value`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"string": "foo"},
		expected: `{
			"message":"Hello, World!",
		  "string": "foo"
		}`,
	},
	{
		name:  `"integer" field with 123 value`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"integer": 123},
		expected: `{
			"message":"Hello, World!",
		  "integer": 123
		}`,
	},
	{
		name:  `"float" field with 3.21 value`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"float": 3.21},
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
		log: mujlog.Log{
			Keys: [4]string{"message", "preview"},
			Max:  0,
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!",
			"preview":""
		}`,
	},
	{
		name: "without message key names",
		log: mujlog.Log{
			Keys: [4]string{},
			Max:  120,
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"":"Hello, World!"
		}`,
	},
	{
		name: "only full message key name",
		log: mujlog.Log{
			Keys: [4]string{"message"},
			Max:  120,
		},
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name:  `explicit byte slice as short message field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"preview": []byte("Explicit byte slice")},
		expected: `{
		  "message": "Hello, World!",
			"preview":"Explicit byte slice"
		}`,
	},
	{
		name:  `explicit string as short message field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"preview": "Explicit string"},
		expected: `{
		  "message": "Hello, World!",
			"preview":"Explicit string"
		}`,
	},
	{
		name:  `explicit integer as short message field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"preview": 42},
		expected: `{
		  "message": "Hello, World!",
			"preview":42
		}`,
	},
	{
		name:  `explicit float as short message field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"preview": 4.2},
		expected: `{
		  "message": "Hello, World!",
			"preview":4.2
		}`,
	},
	{
		name:  `explicit boolean as short message field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"preview": true},
		expected: `{
		  "message": "Hello, World!",
			"preview":true
		}`,
	},
	{
		name:  `explicit rune slice as short message field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"preview": []rune("Explicit rune slice")},
		expected: `{
		  "message": "Hello, World!",
			"preview":"Explicit rune slice"
		}`,
	},
	{
		name:  "dynamic field",
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		funcs: map[string]func() interface{}{"time": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).String() }},
		expected: `{
			"message":"Hello, World!",
			"time":"2020-10-15 18:09:00 +0000 UTC"
		}`,
	},
	{
		name:  `"standard flag" do not respects file path`,
		line:  line(),
		log:   dummy,
		input: "path/to/file1:23: Hello, World!",
		flag:  log.LstdFlags,
		expected: `{
			"message":"path/to/file1:23: Hello, World!"
		}`,
	},
	{
		name:  `"long file" flag respects file path`,
		line:  line(),
		log:   dummy,
		input: "path/to/file1:23: Hello, World!",
		flag:  log.Llongfile,
		expected: `{
			"message":"path/to/file1:23: Hello, World!",
			"preview":"Hello, World!",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name:  "file path with empty message",
		line:  line(),
		log:   dummy,
		input: "path/to/file1:23:",
		flag:  log.Llongfile,
		expected: `{
			"message":"path/to/file1:23:",
			"preview":"_EMPTY_",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name:  "file path with blank message",
		line:  line(),
		log:   dummy,
		input: "path/to/file4:56:  ",
		flag:  log.Llongfile,
		expected: `{
			"message":"path/to/file4:56:  ",
			"preview":"_BLANK_",
			"file":"path/to/file4:56"
		}`,
	},
	{
		name:  `"magic" host field`,
		line:  line(),
		log:   dummy,
		input: "Hello, World!",
		kvs:   map[string]interface{}{"host": "example.tld"},
		expected: `{
			"message":"Hello, World!",
			"preview":"example.tld Hello, World!",
			"host":"example.tld"
		}`,
	},
	{
		name:  "GELF",
		line:  line(),
		log:   mujlog.GELF(),
		input: "Hello, GELF!",
		kvs:   map[string]interface{}{"version": "1.1", "host": "example.tld"},
		funcs: map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }},
		expected: `{
			"version":"1.1",
			"short_message":"example.tld Hello, GELF!",
			"full_message":"Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340
		}`,
	},
	{
		name:  "GELF with file path",
		line:  line(),
		log:   mujlog.GELF(),
		input: "path/to/file7:89: Hello, GELF!",
		flag:  log.Llongfile,
		kvs:   map[string]interface{}{"version": "1.1", "host": "example.tld"},
		funcs: map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }},
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
			tc.log.Flag = tc.flag
			tc.log.KVs = tc.kvs
			tc.log.Funcs = tc.funcs

			_, err := fmt.Fprint(tc.log, tc.input)
			if err != nil {
				t.Fatalf("unexpected mujlog write error: %s", err)
			}

			ja := jsonassert.New(testprinter{t: t, link: linkToExample})
			ja.Assertf(buf.String(), tc.expected)
		})
	}
}

func BenchmarkMujlog(b *testing.B) {
	for _, tc := range WriteTestCases {
		if !tc.benchmark {
			continue
		}
		b.Run(strconv.Itoa(tc.line), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf := pool.Get().(*bytes.Buffer)
				buf.Reset()
				defer pool.Put(buf)

				if tc.log.Max == 0 {
					tc.log.Max = 120
				}
				if tc.log.Keys == ([4]string{}) {
					tc.log.Keys = [4]string{"message", "preview", "file", "host"}
					tc.log.Key = mujlog.FullKey
				}
				if bytes.Equal(tc.log.Marks[0], []byte{}) && bytes.Equal(tc.log.Marks[1], []byte{}) && bytes.Equal(tc.log.Marks[2], []byte{}) {
					tc.log.Marks = [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")}
				}
				if tc.log.Replace == nil {
					tc.log.Replace = [][]byte{[]byte("\n"), []byte(" ")}
				}

				tc.log.Output = buf
				tc.log.Flag = tc.flag
				tc.log.KVs = tc.kvs
				tc.log.Funcs = tc.funcs

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
	log      mujlog.Log
	input    []byte
	kvs      map[string]interface{}
	kvs2     map[string]interface{}
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
		name:  `"string" field with "foo" value and "string" key with "bar" value`,
		line:  line(),
		log:   dummy,
		input: []byte("Hello, World!"),
		kvs:   map[string]interface{}{"string": "foo"},
		kvs2:  map[string]interface{}{"string": "bar"},
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
		kvs2:  nil,
		expected: `{
			"message":"Hello, World!"
		}`,
	},
	{
		name:  `input appends to the message field value "string"`,
		line:  line(),
		log:   dummy,
		input: []byte("\nHello, World!"),
		kvs:   map[string]interface{}{"message": "field string value"},
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
		kvs2:  map[string]interface{}{"message": "field string value"},
		expected: `{
			"message":"field string value\nHello, World!",
			"preview":"field string value Hello, World!"
		}`,
	},
	{
		name:  `input is nil and message field value is "string"`,
		line:  line(),
		log:   dummy,
		input: nil,
		kvs:   map[string]interface{}{"message": "string"},
		expected: `{
			"message":"string"
		}`,
	},
	{
		name:  `input is nil and message key-value is "string"`,
		line:  line(),
		log:   dummy,
		input: nil,
		kvs2:  map[string]interface{}{"message": "string"},
		expected: `{
			"message":"string"
		}`,
	},
	{
		name:  `input appends to the integer key-value "message"`,
		line:  line(),
		log:   dummy,
		input: []byte("\nHello, World!"),
		kvs2:  map[string]interface{}{"message": 1},
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
		kvs2:  map[string]interface{}{"message": 2.1},
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
		kvs2:  map[string]interface{}{"message": true},
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
		kvs2:  map[string]interface{}{"message": nil},
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
			tc.log.KVs = tc.kvs

			_, err := tc.log.Log(tc.input, tc.kvs2)
			if err != nil {
				t.Fatalf("unexpected mujlog write error: %s", err)
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
