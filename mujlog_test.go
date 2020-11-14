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

var MujlogWriteTestCases = []struct {
	name      string
	line      int
	log       mujlog.Log
	input     interface{}
	flag      int
	fields    map[string]interface{}
	functions map[string]func() interface{}
	metadata  map[string]string
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
	    "fullMessage":""
		}`,
	},
	{
		name:  "blank message",
		line:  line(),
		input: " ",
		expected: `{
			"shortMessage":"_BLANK_",
	    "fullMessage":" "
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
			"fullMessage":" \n\tHello, World! \t\n"
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
		name:   `"string" field with "foo" value`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]interface{}{"string": "foo"},
		expected: `{
			"shortMessage":"Hello, World!",
		  "string": "foo"
		}`,
	},
	{
		name:   `"integer" field with 123 value`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]interface{}{"integer": 123},
		expected: `{
			"shortMessage":"Hello, World!",
		  "integer": 123
		}`,
	},
	{
		name:   `"float" field with 3.21 value`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]interface{}{"float": 3.21},
		expected: `{
			"shortMessage":"Hello, World!",
		  "float": 3.21
		}`,
	},
	{
		name:  "multiline string",
		line:  line(),
		input: "Hello,\nWorld!",
		expected: `{
			"shortMessage":"Hello, World!",
			"fullMessage":"Hello,\nWorld!"
		}`,
	},
	{
		name:  "long string",
		line:  line(),
		input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"shortMessage":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna ali…",
			"fullMessage":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		}`,
	},
	{
		name:  "multiline long string with leading spaces",
		line:  line(),
		input: " \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		expected: `{
			"shortMessage":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna ali…",
			"fullMessage":" \n \tLorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		}`,
		benchmark: true,
	},
	{
		name:   `explicit short message field`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]interface{}{"shortMessage": "Explicit short message"},
		expected: `{
			"shortMessage":"Explicit short message",
		  "fullMessage": "Hello, World!"
		}`,
	},
	{
		name:      "dynamic field",
		line:      line(),
		input:     "Hello, World!",
		functions: map[string]func() interface{}{"time": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).String() }},
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
			"fullMessage":"path/to/file1:23: Hello, World!",
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
			"fullMessage":"path/to/file1:23:",
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
			"fullMessage":"path/to/file4:56:  ",
			"file":"path/to/file4:56"
		}`,
	},
	{
		name:   `"magic" host field`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]interface{}{"host": "example.tld"},
		expected: `{
			"shortMessage":"example.tld Hello, World!",
			"host":"example.tld"
		}`,
	},
	{
		name:      "GELF",
		line:      line(),
		log:       mujlog.GELF(),
		input:     "Hello, GELF!",
		fields:    map[string]interface{}{"version": "1.1", "host": "example.tld"},
		functions: map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }},
		expected: `{
			"version":"1.1",
			"short_message":"example.tld Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340
		}`,
	},
	{
		name:      "GELF with file path",
		line:      line(),
		log:       mujlog.GELF(),
		input:     "path/to/file7:89: Hello, GELF!",
		flag:      log.Llongfile,
		fields:    map[string]interface{}{"version": "1.1", "host": "example.tld"},
		functions: map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }},
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

func TestMujlogWrite(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range MujlogWriteTestCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			defer pool.Put(buf)

			if tc.log.Short == "" || tc.log.Full == "" || tc.log.File == "" || tc.log.Truncate == 0 {
				tc.log.Short = "shortMessage"
				tc.log.Full = "fullMessage"
				tc.log.File = "file"
				tc.log.Truncate = 120
			}

			tc.log.Output = buf
			tc.log.Flag = tc.flag
			tc.log.Fields = tc.fields
			tc.log.Functions = tc.functions

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
	for _, tc := range MujlogWriteTestCases {
		if !tc.benchmark {
			continue
		}
		b.Run(strconv.Itoa(tc.line), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf := pool.Get().(*bytes.Buffer)
				buf.Reset()
				defer pool.Put(buf)

				if tc.log.Short == "" || tc.log.Full == "" || tc.log.File == "" || tc.log.Truncate == 0 {
					tc.log.Short = "shortMessage"
					tc.log.Full = "fullMessage"
					tc.log.File = "file"
					tc.log.Truncate = 120
				}

				tc.log.Output = buf
				tc.log.Flag = tc.flag
				tc.log.Fields = tc.fields
				tc.log.Functions = tc.functions

				_, err := fmt.Fprint(tc.log, tc.input)
				if err != nil {
					fmt.Println(err)
				}
			}
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
