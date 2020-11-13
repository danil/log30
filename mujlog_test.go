package mujlog_test

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/danil/mujlog"
	"github.com/kinbiko/jsonassert"
)

func TestMujlogWriteTrailingNewLine(t *testing.T) {
	var buf bytes.Buffer

	mjl := mujlog.Mujlog{Output: &buf}

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
	input     string
	flag      int
	fields    map[string]string
	functions map[string]func() interface{}
	expected  string
}{
	{
		name:  "string",
		line:  line(),
		input: "Hello, World!",
		expected: `{
			"short_message":"Hello, World!"
		}`,
	},
	{
		name:  "empty message",
		line:  line(),
		input: "",
		expected: `{
			"short_message":"_EMPTY_",
	    "full_message":""
		}`,
	},
	{
		name:  "blank message",
		line:  line(),
		input: " ",
		expected: `{
			"short_message":"_BLANK_",
	    "full_message":" "
		}`,
	},
	{
		name:  "double quotes",
		line:  line(),
		input: `foo "bar"`,
		expected: `{
			"short_message":"foo \"bar\""
		}`,
	},
	{
		name:  `leading/trailing "spaces"`,
		line:  line(),
		input: " \nHello, World! \n",
		expected: `{
			"short_message":"Hello, World!",
			"full_message":" \nHello, World! \n"
		}`,
	},
	{
		name:  "JSON string",
		line:  line(),
		input: `{"foo":"bar"}`,
		expected: `{
			"short_message":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		name:  "multiline string",
		line:  line(),
		input: "Hello, World!\npath/to/file1:23\npath/to/file4:56",
		expected: `{
			"short_message":"Hello, World!",
			"full_message":"Hello, World!\npath/to/file1:23\npath/to/file4:56"
		}`,
	},
	{
		name:  `"standard flag" do not respects file path`,
		line:  line(),
		input: "path/to/file1:23: Hello, World!",
		flag:  log.LstdFlags,
		expected: `{
			"short_message":"path/to/file1:23: Hello, World!"
		}`,
	},
	{
		name:  `"long file" flag respects file path`,
		line:  line(),
		input: "path/to/file1:23: Hello, World!",
		flag:  log.Llongfile,
		expected: `{
			"short_message":"Hello, World!",
			"full_message":"path/to/file1:23: Hello, World!",
			"_file":"path/to/file1:23"
		}`,
	},
	{
		name:  "file path with empty message",
		line:  line(),
		input: "path/to/file1:23:",
		flag:  log.Llongfile,
		expected: `{
			"short_message":"_BLANK_",
			"full_message":"path/to/file1:23:",
			"_file":"path/to/file1:23"
		}`,
	},
	{
		name:  "file path with blank message",
		line:  line(),
		input: "path/to/file4:56:  ",
		flag:  log.Llongfile,
		expected: `{
			"short_message":"_BLANK_",
			"full_message":"path/to/file4:56:  ",
			"_file":"path/to/file4:56"
		}`,
	},
	{
		name:   `"environment" field with "production" value`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]string{"environment": "production"},
		expected: `{
			"short_message":"Hello, World!",
		  "environment": "production"
		}`,
	},
	{
		name:   `"magic" host field`,
		line:   line(),
		input:  "Hello, World!",
		fields: map[string]string{"host": "example.tld"},
		expected: `{
			"short_message":"example.tld Hello, World!",
			"host":"example.tld"
		}`,
	},
	{
		name:      "dynamic field",
		line:      line(),
		input:     "Hello, World!",
		functions: map[string]func() interface{}{"time": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).String() }},
		expected: `{
			"short_message":"Hello, World!",
			"time":"2020-10-15 18:09:00 +0000 UTC"
		}`,
	},
	{
		name:      "JSON like GELF",
		line:      line(),
		input:     "Hello, GELF!",
		fields:    map[string]string{"version": "1.1", "host": "example.tld"},
		functions: map[string]func() interface{}{"timestamp": func() interface{} { return time.Date(2020, time.October, 15, 18, 9, 0, 0, time.UTC).Unix() }},
		expected: `{
			"version":"1.1",
			"short_message":"example.tld Hello, GELF!",
			"host":"example.tld",
			"timestamp":1602785340
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

			var buf bytes.Buffer

			mjl := mujlog.Mujlog{
				Output:    &buf,
				Flag:      tc.flag,
				Fields:    tc.fields,
				Functions: tc.functions,
			}

			_, err := mjl.Write([]byte(tc.input))
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
