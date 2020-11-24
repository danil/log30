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

var pool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

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
		name:  "string",
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"shortMessage":"Hello, World!"
		}`,
	},
	{
		name:  "integer type appears in the short messages as a string",
		line:  line(),
		input: 123,
		expected: `{
			"shortMessage":"123"
		}`,
	},
	{
		name:  "float type appears in the short messages as a string",
		line:  line(),
		input: 3.21,
		expected: `{
			"shortMessage":"3.21"
		}`,
	},
	{
		name:  "empty message",
		line:  line(),
		input: "",
		expected: `{
			"shortMessage":"_EMPTY_",
	    "message":""
		}`,
	},
	{
		name:  "blank message",
		line:  line(),
		input: " ",
		expected: `{
			"shortMessage":"_BLANK_",
	    "message":" "
		}`,
	},
	{
		name:  "single quotes",
		line:  line(),
		input: "foo 'bar'",
		expected: `{
			"shortMessage":"foo 'bar'"
		}`,
	},
	{
		name:  "double quotes",
		line:  line(),
		input: `foo "bar"`,
		expected: `{
			"shortMessage":"foo \"bar\""
		}`,
	},
	{
		name:  `leading/trailing "spaces"`,
		line:  line(),
		input: " \n\tHello, World! \t\n",
		expected: `{
			"shortMessage":"Hello, World!",
			"message":" \n\tHello, World! \t\n"
		}`,
	},
	{
		name:  "JSON string",
		line:  line(),
		input: `{"foo":"bar"}`,
		expected: `{
			"shortMessage":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		name:  `"string" field with "foo" value`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"string": "foo"},
		expected: `{
			"shortMessage":"Hello, World!",
		  "string": "foo"
		}`,
	},
	{
		name:  `"integer" field with 123 value`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"integer": 123},
		expected: `{
			"shortMessage":"Hello, World!",
		  "integer": 123
		}`,
	},
	{
		name:  `"float" field with 3.21 value`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"float": 3.21},
		expected: `{
			"shortMessage":"Hello, World!",
		  "float": 3.21
		}`,
	},
	{
		name:  "fmt.Fprint prints nil as <nil>",
		line:  line(),
		input: nil,
		expected: `{
			"shortMessage":"<nil>"
		}`,
	},
	{
		name:  "multiline string",
		line:  line(),
		input: "Hello,\nWorld!",
		expected: `{
			"shortMessage":"Hello, World!",
			"message":"Hello,\nWorld!"
		}`,
	},
	{
		name:  "long string",
		line:  line(),
		input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"shortMessage":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna ali…",
			"message":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		}`,
	},
	{
		name:  "multiline long string with leading spaces",
		line:  line(),
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"shortMessage":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna ali…",
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		}`,
	},
	{
		name:  "multiline long string with leading spaces and multibyte character",
		line:  line(),
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää.",
		expected: `{
			"shortMessage":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna Ää…",
			"message":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna Ää."
		}`,
		benchmark: true,
	},
	{
		name:  `explicit byte slice as short message field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"shortMessage": []byte("Explicit byte slice")},
		expected: `{
			"shortMessage":"Explicit byte slice",
		  "message": "Hello, World!"
		}`,
	},
	{
		name:  `explicit string as short message field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"shortMessage": "Explicit string"},
		expected: `{
			"shortMessage":"Explicit string",
		  "message": "Hello, World!"
		}`,
	},
	{
		name:  `explicit integer as short message field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"shortMessage": 42},
		expected: `{
			"shortMessage":42,
		  "message": "Hello, World!"
		}`,
	},
	{
		name:  `explicit float as short message field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"shortMessage": 4.2},
		expected: `{
			"shortMessage":4.2,
		  "message": "Hello, World!"
		}`,
	},
	{
		name:  `explicit boolean as short message field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"shortMessage": true},
		expected: `{
			"shortMessage":true,
		  "message": "Hello, World!"
		}`,
	},
	{
		name:  `explicit rune slice as short message field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"shortMessage": []rune("Explicit rune slice")},
		expected: `{
			"shortMessage":"Explicit rune slice",
		  "message": "Hello, World!"
		}`,
	},
	{
		name:  "dynamic field",
		line:  line(),
		input: "Hello, World!",
		funcs: map[string]func() interface{}{"time": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).String() }},
		expected: `{
			"shortMessage":"Hello, World!",
			"time":"2020-10-15 18:09:00 +0000 UTC"
		}`,
	},
	{
		name:  `"standard flag" do not respects file path`,
		line:  line(),
		input: "path/to/file1:23: Hello, World!",
		flag:  log.LstdFlags,
		expected: `{
			"shortMessage":"path/to/file1:23: Hello, World!"
		}`,
	},
	{
		name:  `"long file" flag respects file path`,
		line:  line(),
		input: "path/to/file1:23: Hello, World!",
		flag:  log.Llongfile,
		expected: `{
			"shortMessage":"Hello, World!",
			"message":"path/to/file1:23: Hello, World!",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name:  "file path with empty message",
		line:  line(),
		input: "path/to/file1:23:",
		flag:  log.Llongfile,
		expected: `{
			"shortMessage":"_EMPTY_",
			"message":"path/to/file1:23:",
			"file":"path/to/file1:23"
		}`,
	},
	{
		name:  "file path with blank message",
		line:  line(),
		input: "path/to/file4:56:  ",
		flag:  log.Llongfile,
		expected: `{
			"shortMessage":"_BLANK_",
			"message":"path/to/file4:56:  ",
			"file":"path/to/file4:56"
		}`,
	},
	{
		name:  `"magic" host field`,
		line:  line(),
		input: "Hello, World!",
		kvs:   map[string]interface{}{"host": "example.tld"},
		expected: `{
			"shortMessage":"example.tld Hello, World!",
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

			if tc.log.Max == 0 {
				tc.log.Max = 120
			}
			if tc.log.Keys == ([4]string{}) {
				tc.log.Keys = [4]string{"message", "shortMessage", "file", "host"}
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
					tc.log.Keys = [4]string{"message", "shortMessage", "file", "host"}
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
		expected: `{
			"shortMessage":"_EMPTY_",
      "message":""
		}`,
	},
	{
		name:  `"string" field with "foo" value and "string" key with "bar" value`,
		line:  line(),
		input: []byte("Hello, World!"),
		kvs:   map[string]interface{}{"string": "foo"},
		kvs2:  map[string]interface{}{"string": "bar"},
		expected: `{
			"shortMessage":"Hello, World!",
		  "string": "bar"
		}`,
	},
	{
		name:  `key-values is nil`,
		line:  line(),
		input: []byte("Hello, World!"),
		kvs2:  nil,
		expected: `{
			"shortMessage":"Hello, World!"
		}`,
	},
	{
		name:  `input appends to the message field value "string"`,
		line:  line(),
		input: []byte("\nHello, World!"),
		kvs:   map[string]interface{}{"message": "field string value"},
		expected: `{
			"shortMessage":"field string value Hello, World!",
			"message":"field string value\nHello, World!"
		}`,
	},
	{
		name:  `input appends to the message key-field value "string"`,
		line:  line(),
		input: []byte("\nHello, World!"),
		kvs2:  map[string]interface{}{"message": "field string value"},
		expected: `{
			"shortMessage":"field string value Hello, World!",
			"message":"field string value\nHello, World!"
		}`,
	},
	{
		name:  `input is nil and message field value is "string"`,
		line:  line(),
		input: nil,
		kvs:   map[string]interface{}{"message": "string"},
		expected: `{
			"shortMessage":"string"
		}`,
	},
	{
		name:  `input is nil and message key-value is "string"`,
		line:  line(),
		input: nil,
		kvs2:  map[string]interface{}{"message": "string"},
		expected: `{
			"shortMessage":"string"
		}`,
	},
	{
		name:  `input appends to the integer key-value "message"`,
		line:  line(),
		input: []byte("\nHello, World!"),
		kvs2:  map[string]interface{}{"message": 1},
		expected: `{
			"shortMessage":"1 Hello, World!",
			"message":"1\nHello, World!"
		}`,
	},
	{
		name:  `input appends to the float key-value "message"`,
		line:  line(),
		input: []byte("\nHello, World!"),
		kvs2:  map[string]interface{}{"message": 2.1},
		expected: `{
			"shortMessage":"2.1 Hello, World!",
			"message":"2.1\nHello, World!"
		}`,
	},
	{
		name:  `input appends to the boolean key-value "message"`,
		line:  line(),
		input: []byte("\nHello, World!"),
		kvs2:  map[string]interface{}{"message": true},
		expected: `{
			"shortMessage":"true Hello, World!",
			"message":"true\nHello, World!"
		}`,
	},
	{
		name:  `input do not appends to the nil key-value "message"`,
		line:  line(),
		input: []byte("Hello, World!"),
		kvs2:  map[string]interface{}{"message": nil},
		expected: `{
			"shortMessage":"Hello, World!"
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

			tc.log.Keys = [4]string{"message", "shortMessage", "file", "host"}
			tc.log.Marks = [3][]byte{[]byte("…"), []byte("_EMPTY_"), []byte("_BLANK_")}
			tc.log.Replace = [][]byte{[]byte("\n"), []byte(" ")}

			if tc.log.Max == 0 {
				tc.log.Max = 120
			}

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
