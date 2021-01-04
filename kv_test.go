package logastic_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/danil/equalastic"
	"github.com/danil/logastic"
	"github.com/danil/logastic/marshalastic"
	"github.com/kinbiko/jsonassert"
)

var KVTestCases = []struct {
	line         int
	input        logastic.KV
	expected     string
	expectedText string
	expectedJSON string
	error        error
	benchmark    bool
}{
	{
		line:         line(),
		input:        logastic.StringBool("bool true", true),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool true":true
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBool("bool false", false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool false":false
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any bool false", false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"any bool false":false
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("reflect bool false", false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"reflect bool false":false
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			return logastic.StringBoolp("bool pointer to true", &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := false
			return logastic.StringBoolp("bool pointer to false", &b)
		}(),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool pointer to false":false
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBoolp("bool nil pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"bool nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			return logastic.StringAny("any bool pointer to true", &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			b2 := &b
			return logastic.StringAny("any twice pointer to bool true", &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any twice pointer to bool true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			return logastic.StringReflect("reflect bool pointer to true", &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			b2 := &b
			return logastic.StringReflect("reflect bool twice pointer to true", &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool twice pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var b *bool
			return logastic.StringReflect("reflect bool pointer to nil", b)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect bool pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("bytes", []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("bytes with quote", []byte(`Hello, "World"!`)),
		expected:     `Hello, \"World\"!`,
		expectedText: `Hello, \"World\"!`,
		expectedJSON: `{
			"bytes with quote":"Hello, \"World\"!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("bytes quote", []byte(`"Hello, World!"`)),
		expected:     `\"Hello, World!\"`,
		expectedText: `\"Hello, World!\"`,
		expectedJSON: `{
			"bytes quote":"\"Hello, World!\""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("bytes nested quote", []byte(`"Hello, "World"!"`)),
		expected:     `\"Hello, \"World\"!\"`,
		expectedText: `\"Hello, \"World\"!\"`,
		expectedJSON: `{
			"bytes nested quote":"\"Hello, \"World\"!\""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("bytes json", []byte(`{"foo":"bar"}`)),
		expected:     `{\"foo\":\"bar\"}`,
		expectedText: `{\"foo\":\"bar\"}`,
		expectedJSON: `{
			"bytes json":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("bytes json quote", []byte(`"{"foo":"bar"}"`)),
		expected:     `\"{\"foo\":\"bar\"}\"`,
		expectedText: `\"{\"foo\":\"bar\"}\"`,
		expectedJSON: `{
			"bytes json quote":"\"{\"foo\":\"bar\"}\""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("empty bytes", []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytes("nil bytes", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any bytes", []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any empty bytes", []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect bytes", []byte("Hello, Wörld!")),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect empty bytes", []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty bytes":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte("Hello, Wörld!")
			return logastic.StringBytesp("bytes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte{}
			return logastic.StringBytesp("empty bytes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringBytesp("nil bytes pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte("Hello, Wörld!")
			return logastic.StringAny("any bytes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte{}
			return logastic.StringAny("any empty bytes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte("Hello, Wörld!")
			return logastic.StringReflect("reflect bytes pointer", &p)
		}(),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes pointer":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte{}
			return logastic.StringReflect("reflect empty bytes pointer", &p)
		}(),
		expected: "",
		expectedJSON: `{
			"reflect empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringComplex128("complex128", complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any complex128", complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect complex128", complex(1, 23)),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex128 = complex(1, 23)
			return logastic.StringComplex128p("complex128 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringComplex128p("nil complex128 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex128 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex128 = complex(1, 23)
			return logastic.StringAny("any complex128 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex128 = complex(1, 23)
			return logastic.StringReflect("reflect complex128 pointer", &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        logastic.StringComplex64("complex64", complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any complex64", complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"any complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect complex64", complex(3, 21)),
		expected:     "(3+21i)",
		expectedText: "(3+21i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        logastic.StringError("error", errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringError("nil error", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil error":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any error", errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"any error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect error", errors.New("something went wrong")),
		expected:     "{something went wrong}",
		expectedText: "{something went wrong}",
		expectedJSON: `{
			"reflect error":{}
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex64 = complex(1, 23)
			return logastic.StringComplex64p("complex64 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex64 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringComplex64p("nil complex64 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex64 = complex(1, 23)
			return logastic.StringAny("any complex64 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex64 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex64 = complex(1, 23)
			return logastic.StringReflect("reflect complex64 pointer", &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex64"),
	},
	{
		line:         line(),
		input:        logastic.StringFloat32("float32", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat32("high precision float32", 0.123456789),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat32("zero float32", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any float32", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any zero float32", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect float32", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect zero float32", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float32":0
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 4.2
			return logastic.StringFloat32p("float32 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 0.123456789
			return logastic.StringFloat32p("high precision float32 pointer", &f)
		}(),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat32p("float32 nil pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float32 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 4.2
			return logastic.StringAny("any float32 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 4.2
			return logastic.StringReflect("reflect float32 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f *float32
			return logastic.StringReflect("reflect float32 pointer to nil", f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float32 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat64("float64", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat64("high precision float64", 0.123456789),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat64("zero float64", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any float64", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any zero float64", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect float64", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect zero float64", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float64":0
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 4.2
			return logastic.StringFloat64p("float64 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 0.123456789
			return logastic.StringFloat64p("high precision float64 pointer", &f)
		}(),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringFloat64p("float64 nil pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float64 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 4.2
			return logastic.StringAny("any float64 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 4.2
			return logastic.StringReflect("reflect float64 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f *float64
			return logastic.StringReflect("reflect float64 pointer to nil", f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float64 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringInt("int", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any int", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect int", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int = 42
			return logastic.StringIntp("int pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int = 42
			return logastic.StringAny("any int pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int = 42
			return logastic.StringReflect("reflect int pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringInt16("int16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any int16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect int16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int16 = 42
			return logastic.StringInt16p("int16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int16 = 42
			return logastic.StringAny("any int16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int16 = 42
			return logastic.StringReflect("reflect int16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringInt32("int32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any int32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect int32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int32 = 42
			return logastic.StringInt32p("int32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int32 = 42
			return logastic.StringAny("any int32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int32 = 42
			return logastic.StringReflect("reflect int32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringInt64("int64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any int64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect int64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int64 = 42
			return logastic.StringInt64p("int64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int64 = 42
			return logastic.StringAny("any int64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int64 = 42
			return logastic.StringReflect("reflect int64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringInt8("int8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any int8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect int8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int8 = 42
			return logastic.StringInt8p("int8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int8 = 42
			return logastic.StringAny("any int8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int8 = 42
			return logastic.StringReflect("reflect int8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRunes("runes", []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRunes("empty runes", []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRunes("nil runes", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRunes("rune slice with zero rune", []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any runes", []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any empty runes", []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any rune slice with zero rune", []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect runes", []rune("Hello, Wörld!")),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect empty runes", []rune{}),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes":[]
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect rune slice with zero rune", []rune{rune(0)}),
		expected:     "[0]",
		expectedText: "[0]",
		expectedJSON: `{
			"reflect rune slice with zero rune":[0]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune("Hello, Wörld!")
			return logastic.StringRunesp("runes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune{}
			return logastic.StringRunesp("empty runes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRunesp("nil runes pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune("Hello, Wörld!")
			return logastic.StringAny("any runes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune{}
			return logastic.StringAny("any empty runes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune("Hello, Wörld!")
			return logastic.StringReflect("reflect runes pointer", &p)
		}(),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes pointer":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune{}
			return logastic.StringReflect("reflect empty runes pointer", &p)
		}(),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes pointer":[]
		}`,
	},
	{
		line:         line(),
		input:        logastic.String("string", "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.String("empty string", ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.String("string with zero byte", string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any string", "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any empty string", ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any string with zero byte", string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect string", "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect empty string", ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect string with zero byte", string(byte(0))),
		expected:     "\u0000",
		expectedText: "\u0000",
		expectedJSON: `{
			"reflect string with zero byte":"\u0000"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := "Hello, Wörld!"
			return logastic.StringStringp("string pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := ""
			return logastic.StringStringp("empty string pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringStringp("nil string pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil string pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := "Hello, Wörld!"
			return logastic.StringAny("any string pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := ""
			return logastic.StringAny("any empty string pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string pointer":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := "Hello, Wörld!"
			return logastic.StringReflect("reflect string pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := ""
			return logastic.StringReflect("reflect empty string pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint("uint", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any uint", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect uint", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint = 42
			return logastic.StringUintp("uint pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUintp("nil uint pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint = 42
			return logastic.StringAny("any uint pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint = 42
			return logastic.StringReflect("reflect uint pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint16("uint16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any uint16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect uint16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint16 = 42
			return logastic.StringUint16p("uint16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint16p("uint16 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"uint16 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint16 = 42
			return logastic.StringAny("any uint16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint16 = 42
			return logastic.StringReflect("reflect uint16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i *uint16
			return logastic.StringReflect("reflect uint16 pointer to nil", i)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect uint16 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint32("uint32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any uint32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect uint32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint32 = 42
			return logastic.StringUint32p("uint32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint32p("nil uint32 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint32 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint32 = 42
			return logastic.StringAny("any uint32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint32 = 42
			return logastic.StringReflect("reflect uint32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint64("uint64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64":42
		}`,
	},

	{
		line:         line(),
		input:        logastic.StringAny("any uint64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect uint64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint64 = 42
			return logastic.StringUint64p("uint64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint64p("nil uint64 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint64 = 42
			return logastic.StringAny("any uint64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint64 = 42
			return logastic.StringReflect("reflect uint64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint8("uint8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any uint8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect uint8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint8 = 42
			return logastic.StringUint8p("uint8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUint8p("nil uint8 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint8 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint8 = 42
			return logastic.StringAny("any uint8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint8 = 42
			return logastic.StringReflect("reflect uint8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUintptr("uintptr", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any uintptr", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect uintptr", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uintptr = 42
			return logastic.StringUintptrp("uintptr pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringUintptrp("nil uintptr pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uintptr pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uintptr = 42
			return logastic.StringAny("any uintptr pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uintptr = 42
			return logastic.StringReflect("reflect uintptr pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringTime("time", time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any time", time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect time", time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return logastic.StringTimep("time pointer", &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var t *time.Time
			return logastic.StringTimep("nil time pointer", t)
		}(),
		expected:     "0001-01-01 00:00:00 +0000 UTC",
		expectedText: "0001-01-01T00:00:00Z",
		expectedJSON: `{
			"nil time pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return logastic.StringAny("any time pointer", &t)
		}(),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return logastic.StringReflect("reflect time pointer", &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringDuration("duration", 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any duration", 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect duration", 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			d := 42 * time.Nanosecond
			return logastic.StringDurationp("duration pointer", &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration pointer":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringDurationp("nil duration pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil duration pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			d := 42 * time.Nanosecond
			return logastic.StringAny("any duration pointer", &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration pointer":"42ns"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			d := 42 * time.Nanosecond
			return logastic.StringReflect("reflect duration pointer", &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any struct", Struct{Name: "John Doe", Age: 42}),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"any struct": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return logastic.StringAny("any struct pointer", &s)
		}(),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"any struct pointer": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("struct reflect", Struct{Name: "John Doe", Age: 42}),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"struct reflect": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return logastic.StringReflect("struct reflect pointer", &s)
		}(),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"struct reflect pointer": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRaw("raw json", []byte(`{"foo":"bar"}`)),
		expected:     `{"foo":"bar"}`,
		expectedText: `{"foo":"bar"}`,
		expectedJSON: `{
			"raw json":{"foo":"bar"}
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringRaw("raw malformed json object", []byte(`xyz{"foo":"bar"}`)),
		expected:     `xyz{"foo":"bar"}`,
		expectedText: `xyz{"foo":"bar"}`,
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character 'x' looking for beginning of value"),
	},
	{
		line:         line(),
		input:        logastic.StringRaw("raw malformed json key/value", []byte(`{"foo":"bar""}`)),
		expected:     `{"foo":"bar""}`,
		expectedText: `{"foo":"bar""}`,
		error:        errors.New(`json: error calling MarshalJSON for type json.Marshaler: invalid character '"' after object key:value pair`),
	},
	{
		line:         line(),
		input:        logastic.StringRaw("raw json with unescaped null byte", append([]byte(`{"foo":"`), append([]byte{0}, []byte(`xyz"}`)...)...)),
		expected:     "{\"foo\":\"\u0000xyz\"}",
		expectedText: "{\"foo\":\"\u0000xyz\"}",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character '\\x00' in string literal"),
	},
	{
		line:         line(),
		input:        logastic.StringRaw("raw nil", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"raw nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any byte array", [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			a := [3]byte{'f', 'o', 'o'}
			return logastic.StringAny("any byte array pointer", &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var a *[3]byte
			return logastic.StringAny("any byte array pointer to nil", a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect byte array", [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			a := [3]byte{'f', 'o', 'o'}
			return logastic.StringReflect("reflect byte array pointer", &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var a *[3]byte
			return logastic.StringReflect("reflect byte array pointer to nil", a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringAny("any untyped nil", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.StringReflect("reflect untyped nil", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBool(marshalastic.String("bool true"), true),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool true":true
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBool(marshalastic.String("bool false"), false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool false":false
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any bool false"), false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"any bool false":false
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("reflect bool false"), false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"reflect bool false":false
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			return logastic.TextBoolp(marshalastic.String("bool pointer to true"), &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := false
			return logastic.TextBoolp(marshalastic.String("bool pointer to false"), &b)
		}(),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool pointer to false":false
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBoolp(marshalastic.String("bool nil pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"bool nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			return logastic.TextAny(marshalastic.String("any bool pointer to true"), &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			b2 := &b
			return logastic.TextAny(marshalastic.String("any twice pointer to bool true"), &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any twice pointer to bool true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			return logastic.TextReflect(marshalastic.String("reflect bool pointer to true"), &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			b := true
			b2 := &b
			return logastic.TextReflect(marshalastic.String("reflect bool twice pointer to true"), &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool twice pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var b *bool
			return logastic.TextReflect(marshalastic.String("reflect bool pointer to nil"), b)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect bool pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("bytes"), []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("bytes with quote"), []byte(`Hello, "World"!`)),
		expected:     `Hello, \"World\"!`,
		expectedText: `Hello, \"World\"!`,
		expectedJSON: `{
			"bytes with quote":"Hello, \"World\"!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("bytes quote"), []byte(`"Hello, World!"`)),
		expected:     `\"Hello, World!\"`,
		expectedText: `\"Hello, World!\"`,
		expectedJSON: `{
			"bytes quote":"\"Hello, World!\""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("bytes nested quote"), []byte(`"Hello, "World"!"`)),
		expected:     `\"Hello, \"World\"!\"`,
		expectedText: `\"Hello, \"World\"!\"`,
		expectedJSON: `{
			"bytes nested quote":"\"Hello, \"World\"!\""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("bytes json"), []byte(`{"foo":"bar"}`)),
		expected:     `{\"foo\":\"bar\"}`,
		expectedText: `{\"foo\":\"bar\"}`,
		expectedJSON: `{
			"bytes json":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("bytes json quote"), []byte(`"{"foo":"bar"}"`)),
		expected:     `\"{\"foo\":\"bar\"}\"`,
		expectedText: `\"{\"foo\":\"bar\"}\"`,
		expectedJSON: `{
			"bytes json quote":"\"{\"foo\":\"bar\"}\""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("empty bytes"), []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytes(marshalastic.String("nil bytes"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any bytes"), []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any empty bytes"), []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect bytes"), []byte("Hello, Wörld!")),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect empty bytes"), []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty bytes":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte("Hello, Wörld!")
			return logastic.TextBytesp(marshalastic.String("bytes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte{}
			return logastic.TextBytesp(marshalastic.String("empty bytes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextBytesp(marshalastic.String("nil bytes pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte("Hello, Wörld!")
			return logastic.TextAny(marshalastic.String("any bytes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte{}
			return logastic.TextAny(marshalastic.String("any empty bytes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte("Hello, Wörld!")
			return logastic.TextReflect(marshalastic.String("reflect bytes pointer"), &p)
		}(),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes pointer":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []byte{}
			return logastic.TextReflect(marshalastic.String("reflect empty bytes pointer"), &p)
		}(),
		expected: "",
		expectedJSON: `{
			"reflect empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextComplex128(marshalastic.String("complex128"), complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any complex128"), complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect complex128"), complex(1, 23)),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex128 = complex(1, 23)
			return logastic.TextComplex128p(marshalastic.String("complex128 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextComplex128p(marshalastic.String("nil complex128 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex128 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex128 = complex(1, 23)
			return logastic.TextAny(marshalastic.String("any complex128 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex128 = complex(1, 23)
			return logastic.TextReflect(marshalastic.String("reflect complex128 pointer"), &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        logastic.TextComplex64(marshalastic.String("complex64"), complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any complex64"), complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"any complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect complex64"), complex(3, 21)),
		expected:     "(3+21i)",
		expectedText: "(3+21i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        logastic.TextError(marshalastic.String("error"), errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextError(marshalastic.String("nil error"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil error":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any error"), errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"any error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect error"), errors.New("something went wrong")),
		expected:     "{something went wrong}",
		expectedText: "{something went wrong}",
		expectedJSON: `{
			"reflect error":{}
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex64 = complex(1, 23)
			return logastic.TextComplex64p(marshalastic.String("complex64 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex64 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextComplex64p(marshalastic.String("nil complex64 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex64 = complex(1, 23)
			return logastic.TextAny(marshalastic.String("any complex64 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex64 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var c complex64 = complex(1, 23)
			return logastic.TextReflect(marshalastic.String("reflect complex64 pointer"), &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex64"),
	},
	{
		line:         line(),
		input:        logastic.TextFloat32(marshalastic.String("float32"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat32(marshalastic.String("high precision float32"), 0.123456789),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat32(marshalastic.String("zero float32"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any float32"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any zero float32"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect float32"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect zero float32"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float32":0
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 4.2
			return logastic.TextFloat32p(marshalastic.String("float32 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 0.123456789
			return logastic.TextFloat32p(marshalastic.String("high precision float32 pointer"), &f)
		}(),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat32p(marshalastic.String("float32 nil pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float32 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 4.2
			return logastic.TextAny(marshalastic.String("any float32 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float32 = 4.2
			return logastic.TextReflect(marshalastic.String("reflect float32 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f *float32
			return logastic.TextReflect(marshalastic.String("reflect float32 pointer to nil"), f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float32 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat64(marshalastic.String("float64"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat64(marshalastic.String("high precision float64"), 0.123456789),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat64(marshalastic.String("zero float64"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any float64"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any zero float64"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect float64"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect zero float64"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float64":0
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 4.2
			return logastic.TextFloat64p(marshalastic.String("float64 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 0.123456789
			return logastic.TextFloat64p(marshalastic.String("high precision float64 pointer"), &f)
		}(),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextFloat64p(marshalastic.String("float64 nil pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float64 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 4.2
			return logastic.TextAny(marshalastic.String("any float64 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f float64 = 4.2
			return logastic.TextReflect(marshalastic.String("reflect float64 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var f *float64
			return logastic.TextReflect(marshalastic.String("reflect float64 pointer to nil"), f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float64 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextInt(marshalastic.String("int"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any int"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect int"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int = 42
			return logastic.TextIntp(marshalastic.String("int pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int = 42
			return logastic.TextAny(marshalastic.String("any int pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int = 42
			return logastic.TextReflect(marshalastic.String("reflect int pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextInt16(marshalastic.String("int16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any int16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect int16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int16 = 42
			return logastic.TextInt16p(marshalastic.String("int16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int16 = 42
			return logastic.TextAny(marshalastic.String("any int16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int16 = 42
			return logastic.TextReflect(marshalastic.String("reflect int16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextInt32(marshalastic.String("int32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any int32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect int32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int32 = 42
			return logastic.TextInt32p(marshalastic.String("int32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int32 = 42
			return logastic.TextAny(marshalastic.String("any int32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int32 = 42
			return logastic.TextReflect(marshalastic.String("reflect int32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextInt64(marshalastic.String("int64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any int64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect int64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int64 = 42
			return logastic.TextInt64p(marshalastic.String("int64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int64 = 42
			return logastic.TextAny(marshalastic.String("any int64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int64 = 42
			return logastic.TextReflect(marshalastic.String("reflect int64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextInt8(marshalastic.String("int8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any int8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect int8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int8 = 42
			return logastic.TextInt8p(marshalastic.String("int8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int8 = 42
			return logastic.TextAny(marshalastic.String("any int8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i int8 = 42
			return logastic.TextReflect(marshalastic.String("reflect int8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRunes(marshalastic.String("runes"), []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRunes(marshalastic.String("empty runes"), []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRunes(marshalastic.String("nil runes"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRunes(marshalastic.String("rune slice with zero rune"), []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any runes"), []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any empty runes"), []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any rune slice with zero rune"), []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect runes"), []rune("Hello, Wörld!")),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect empty runes"), []rune{}),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes":[]
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect rune slice with zero rune"), []rune{rune(0)}),
		expected:     "[0]",
		expectedText: "[0]",
		expectedJSON: `{
			"reflect rune slice with zero rune":[0]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune("Hello, Wörld!")
			return logastic.TextRunesp(marshalastic.String("runes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune{}
			return logastic.TextRunesp(marshalastic.String("empty runes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRunesp(marshalastic.String("nil runes pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune("Hello, Wörld!")
			return logastic.TextAny(marshalastic.String("any runes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune{}
			return logastic.TextAny(marshalastic.String("any empty runes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune("Hello, Wörld!")
			return logastic.TextReflect(marshalastic.String("reflect runes pointer"), &p)
		}(),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes pointer":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := []rune{}
			return logastic.TextReflect(marshalastic.String("reflect empty runes pointer"), &p)
		}(),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes pointer":[]
		}`,
	},
	{
		line:         line(),
		input:        logastic.Text(marshalastic.String("string"), marshalastic.String("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.Text(marshalastic.String("empty string"), marshalastic.String("")),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.Text(marshalastic.String("string with zero byte"), marshalastic.String((string(byte(0))))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextString(marshalastic.String("string"), "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextString(marshalastic.String("empty string"), ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextString(marshalastic.String("string with zero byte"), string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any string"), "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any empty string"), ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any string with zero byte"), string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect string"), "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect empty string"), ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect string with zero byte"), string(byte(0))),
		expected:     "\u0000",
		expectedText: "\u0000",
		expectedJSON: `{
			"reflect string with zero byte":"\u0000"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := "Hello, Wörld!"
			return logastic.TextStringp(marshalastic.String("string pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := ""
			return logastic.TextStringp(marshalastic.String("empty string pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextStringp(marshalastic.String("nil string pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil string pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := "Hello, Wörld!"
			return logastic.TextAny(marshalastic.String("any string pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := ""
			return logastic.TextAny(marshalastic.String("any empty string pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string pointer":""
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := "Hello, Wörld!"
			return logastic.TextReflect(marshalastic.String("reflect string pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			p := ""
			return logastic.TextReflect(marshalastic.String("reflect empty string pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint(marshalastic.String("uint"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any uint"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect uint"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint = 42
			return logastic.TextUintp(marshalastic.String("uint pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUintp(marshalastic.String("nil uint pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint = 42
			return logastic.TextAny(marshalastic.String("any uint pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint = 42
			return logastic.TextReflect(marshalastic.String("reflect uint pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint16(marshalastic.String("uint16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any uint16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect uint16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint16 = 42
			return logastic.TextUint16p(marshalastic.String("uint16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint16p(marshalastic.String("uint16 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"uint16 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint16 = 42
			return logastic.TextAny(marshalastic.String("any uint16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint16 = 42
			return logastic.TextReflect(marshalastic.String("reflect uint16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i *uint16
			return logastic.TextReflect(marshalastic.String("reflect uint16 pointer to nil"), i)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect uint16 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint32(marshalastic.String("uint32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any uint32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect uint32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint32 = 42
			return logastic.TextUint32p(marshalastic.String("uint32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint32p(marshalastic.String("nil uint32 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint32 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint32 = 42
			return logastic.TextAny(marshalastic.String("any uint32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint32 = 42
			return logastic.TextReflect(marshalastic.String("reflect uint32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint64(marshalastic.String("uint64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64":42
		}`,
	},

	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any uint64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect uint64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint64 = 42
			return logastic.TextUint64p(marshalastic.String("uint64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint64p(marshalastic.String("nil uint64 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint64 = 42
			return logastic.TextAny(marshalastic.String("any uint64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint64 = 42
			return logastic.TextReflect(marshalastic.String("reflect uint64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint8(marshalastic.String("uint8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any uint8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect uint8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint8 = 42
			return logastic.TextUint8p(marshalastic.String("uint8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUint8p(marshalastic.String("nil uint8 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint8 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint8 = 42
			return logastic.TextAny(marshalastic.String("any uint8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uint8 = 42
			return logastic.TextReflect(marshalastic.String("reflect uint8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUintptr(marshalastic.String("uintptr"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any uintptr"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect uintptr"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uintptr = 42
			return logastic.TextUintptrp(marshalastic.String("uintptr pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextUintptrp(marshalastic.String("nil uintptr pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uintptr pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uintptr = 42
			return logastic.TextAny(marshalastic.String("any uintptr pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr pointer":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var i uintptr = 42
			return logastic.TextReflect(marshalastic.String("reflect uintptr pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextTime(marshalastic.String("time"), time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any time"), time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect time"), time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return logastic.TextTimep(marshalastic.String("time pointer"), &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var t *time.Time
			return logastic.TextTimep(marshalastic.String("nil time pointer"), t)
		}(),
		expected:     "0001-01-01 00:00:00 +0000 UTC",
		expectedText: "0001-01-01T00:00:00Z",
		expectedJSON: `{
			"nil time pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return logastic.TextAny(marshalastic.String("any time pointer"), &t)
		}(),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return logastic.TextReflect(marshalastic.String("reflect time pointer"), &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextDuration(marshalastic.String("duration"), 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any duration"), 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect duration"), 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration":42
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			d := 42 * time.Nanosecond
			return logastic.TextDurationp(marshalastic.String("duration pointer"), &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration pointer":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextDurationp(marshalastic.String("nil duration pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil duration pointer":null
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			d := 42 * time.Nanosecond
			return logastic.TextAny(marshalastic.String("any duration pointer"), &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration pointer":"42ns"
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			d := 42 * time.Nanosecond
			return logastic.TextReflect(marshalastic.String("reflect duration pointer"), &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration pointer":42
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any struct"), Struct{Name: "John Doe", Age: 42}),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"any struct": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return logastic.TextAny(marshalastic.String("any struct pointer"), &s)
		}(),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"any struct pointer": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("struct reflect"), Struct{Name: "John Doe", Age: 42}),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"struct reflect": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return logastic.TextReflect(marshalastic.String("struct reflect pointer"), &s)
		}(),
		expected:     "{John Doe 42}",
		expectedText: "{John Doe 42}",
		expectedJSON: `{
			"struct reflect pointer": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRaw(marshalastic.String("raw json"), []byte(`{"foo":"bar"}`)),
		expected:     `{"foo":"bar"}`,
		expectedText: `{"foo":"bar"}`,
		expectedJSON: `{
			"raw json":{"foo":"bar"}
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextRaw(marshalastic.String("raw malformed json object"), []byte(`xyz{"foo":"bar"}`)),
		expected:     `xyz{"foo":"bar"}`,
		expectedText: `xyz{"foo":"bar"}`,
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character 'x' looking for beginning of value"),
	},
	{
		line:         line(),
		input:        logastic.TextRaw(marshalastic.String("raw malformed json key/value"), []byte(`{"foo":"bar""}`)),
		expected:     `{"foo":"bar""}`,
		expectedText: `{"foo":"bar""}`,
		error:        errors.New(`json: error calling MarshalJSON for type json.Marshaler: invalid character '"' after object key:value pair`),
	},
	{
		line:         line(),
		input:        logastic.TextRaw(marshalastic.String("raw json with unescaped null byte"), append([]byte(`{"foo":"`), append([]byte{0}, []byte(`xyz"}`)...)...)),
		expected:     "{\"foo\":\"\u0000xyz\"}",
		expectedText: "{\"foo\":\"\u0000xyz\"}",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character '\\x00' in string literal"),
	},
	{
		line:         line(),
		input:        logastic.TextRaw(marshalastic.String("raw nil"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"raw nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any byte array"), [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			a := [3]byte{'f', 'o', 'o'}
			return logastic.TextAny(marshalastic.String("any byte array pointer"), &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var a *[3]byte
			return logastic.TextAny(marshalastic.String("any byte array pointer to nil"), a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect byte array"), [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			a := [3]byte{'f', 'o', 'o'}
			return logastic.TextReflect(marshalastic.String("reflect byte array pointer"), &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() logastic.KV {
			var a *[3]byte
			return logastic.TextReflect(marshalastic.String("reflect byte array pointer to nil"), a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextAny(marshalastic.String("any untyped nil"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        logastic.TextReflect(marshalastic.String("reflect untyped nil"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect untyped nil":null
		}`,
	},
}

func TestKV(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range KVTestCases {
		tc := tc
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			p, err := tc.input.MarshalText()
			if err != nil {
				t.Fatalf("encoding marshal text error: %s", err)
			}

			m := map[string]json.Marshaler{string(p): tc.input}

			p, err = json.Marshal(m)

			if !equalastic.ErrorEqual(err, tc.error) {
				t.Fatalf("marshal error expected: %s, recieved: %s %s", tc.error, err, linkToExample)
			}

			if err == nil {
				ja := jsonassert.New(testprinter{t: t, link: linkToExample})
				ja.Assertf(string(p), tc.expectedJSON)
			}
		})
	}
}
