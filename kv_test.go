package log30_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/danil/equal4"
	"github.com/danil/log30"
	"github.com/danil/log30/marshal30"
	"github.com/kinbiko/jsonassert"
)

var KVTestCases = []struct {
	line         int
	input        log30.KV
	expected     string
	expectedText string
	expectedJSON string
	error        error
	benchmark    bool
}{
	{
		line:         line(),
		input:        log30.StringBool("bool true", true),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool true":true
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBool("bool false", false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool false":false
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any bool false", false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"any bool false":false
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("reflect bool false", false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"reflect bool false":false
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			return log30.StringBoolp("bool pointer to true", &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := false
			return log30.StringBoolp("bool pointer to false", &b)
		}(),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool pointer to false":false
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBoolp("bool nil pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"bool nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			return log30.StringAny("any bool pointer to true", &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			b2 := &b
			return log30.StringAny("any twice pointer to bool true", &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any twice pointer to bool true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			return log30.StringReflect("reflect bool pointer to true", &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			b2 := &b
			return log30.StringReflect("reflect bool twice pointer to true", &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool twice pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var b *bool
			return log30.StringReflect("reflect bool pointer to nil", b)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect bool pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("bytes", []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("bytes with quote", []byte(`Hello, "World"!`)),
		expected:     `Hello, \"World\"!`,
		expectedText: `Hello, \"World\"!`,
		expectedJSON: `{
			"bytes with quote":"Hello, \"World\"!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("bytes quote", []byte(`"Hello, World!"`)),
		expected:     `\"Hello, World!\"`,
		expectedText: `\"Hello, World!\"`,
		expectedJSON: `{
			"bytes quote":"\"Hello, World!\""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("bytes nested quote", []byte(`"Hello, "World"!"`)),
		expected:     `\"Hello, \"World\"!\"`,
		expectedText: `\"Hello, \"World\"!\"`,
		expectedJSON: `{
			"bytes nested quote":"\"Hello, \"World\"!\""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("bytes json", []byte(`{"foo":"bar"}`)),
		expected:     `{\"foo\":\"bar\"}`,
		expectedText: `{\"foo\":\"bar\"}`,
		expectedJSON: `{
			"bytes json":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("bytes json quote", []byte(`"{"foo":"bar"}"`)),
		expected:     `\"{\"foo\":\"bar\"}\"`,
		expectedText: `\"{\"foo\":\"bar\"}\"`,
		expectedJSON: `{
			"bytes json quote":"\"{\"foo\":\"bar\"}\""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("empty bytes", []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytes("nil bytes", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any bytes", []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any empty bytes", []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect bytes", []byte("Hello, Wörld!")),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect empty bytes", []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty bytes":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte("Hello, Wörld!")
			return log30.StringBytesp("bytes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte{}
			return log30.StringBytesp("empty bytes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringBytesp("nil bytes pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte("Hello, Wörld!")
			return log30.StringAny("any bytes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte{}
			return log30.StringAny("any empty bytes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte("Hello, Wörld!")
			return log30.StringReflect("reflect bytes pointer", &p)
		}(),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes pointer":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte{}
			return log30.StringReflect("reflect empty bytes pointer", &p)
		}(),
		expected: "",
		expectedJSON: `{
			"reflect empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringComplex128("complex128", complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any complex128", complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect complex128", complex(1, 23)),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex128 = complex(1, 23)
			return log30.StringComplex128p("complex128 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringComplex128p("nil complex128 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex128 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex128 = complex(1, 23)
			return log30.StringAny("any complex128 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex128 = complex(1, 23)
			return log30.StringReflect("reflect complex128 pointer", &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        log30.StringComplex64("complex64", complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any complex64", complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"any complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect complex64", complex(3, 21)),
		expected:     "(3+21i)",
		expectedText: "(3+21i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        log30.StringError("error", errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringError("nil error", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil error":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any error", errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"any error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect error", errors.New("something went wrong")),
		expected:     "{something went wrong}",
		expectedText: "{something went wrong}",
		expectedJSON: `{
			"reflect error":{}
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex64 = complex(1, 23)
			return log30.StringComplex64p("complex64 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex64 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringComplex64p("nil complex64 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex64 = complex(1, 23)
			return log30.StringAny("any complex64 pointer", &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex64 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex64 = complex(1, 23)
			return log30.StringReflect("reflect complex64 pointer", &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex64"),
	},
	{
		line:         line(),
		input:        log30.StringFloat32("float32", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat32("high precision float32", 0.123456789),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat32("zero float32", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any float32", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any zero float32", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect float32", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect zero float32", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float32":0
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 4.2
			return log30.StringFloat32p("float32 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 0.123456789
			return log30.StringFloat32p("high precision float32 pointer", &f)
		}(),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat32p("float32 nil pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float32 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 4.2
			return log30.StringAny("any float32 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 4.2
			return log30.StringReflect("reflect float32 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f *float32
			return log30.StringReflect("reflect float32 pointer to nil", f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float32 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat64("float64", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat64("high precision float64", 0.123456789),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat64("zero float64", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any float64", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any zero float64", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect float64", 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect zero float64", 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float64":0
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 4.2
			return log30.StringFloat64p("float64 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 0.123456789
			return log30.StringFloat64p("high precision float64 pointer", &f)
		}(),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.StringFloat64p("float64 nil pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float64 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 4.2
			return log30.StringAny("any float64 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 4.2
			return log30.StringReflect("reflect float64 pointer", &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f *float64
			return log30.StringReflect("reflect float64 pointer to nil", f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float64 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringInt("int", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any int", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect int", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int = 42
			return log30.StringIntp("int pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int = 42
			return log30.StringAny("any int pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int = 42
			return log30.StringReflect("reflect int pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringInt16("int16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any int16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect int16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int16 = 42
			return log30.StringInt16p("int16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int16 = 42
			return log30.StringAny("any int16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int16 = 42
			return log30.StringReflect("reflect int16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringInt32("int32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any int32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect int32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int32 = 42
			return log30.StringInt32p("int32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int32 = 42
			return log30.StringAny("any int32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int32 = 42
			return log30.StringReflect("reflect int32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringInt64("int64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any int64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect int64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int64 = 42
			return log30.StringInt64p("int64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int64 = 42
			return log30.StringAny("any int64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int64 = 42
			return log30.StringReflect("reflect int64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringInt8("int8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any int8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect int8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int8 = 42
			return log30.StringInt8p("int8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int8 = 42
			return log30.StringAny("any int8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int8 = 42
			return log30.StringReflect("reflect int8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringRunes("runes", []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringRunes("empty runes", []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringRunes("nil runes", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringRunes("rune slice with zero rune", []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any runes", []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any empty runes", []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any rune slice with zero rune", []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect runes", []rune("Hello, Wörld!")),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect empty runes", []rune{}),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes":[]
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect rune slice with zero rune", []rune{rune(0)}),
		expected:     "[0]",
		expectedText: "[0]",
		expectedJSON: `{
			"reflect rune slice with zero rune":[0]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune("Hello, Wörld!")
			return log30.StringRunesp("runes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune{}
			return log30.StringRunesp("empty runes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringRunesp("nil runes pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune("Hello, Wörld!")
			return log30.StringAny("any runes pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune{}
			return log30.StringAny("any empty runes pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune("Hello, Wörld!")
			return log30.StringReflect("reflect runes pointer", &p)
		}(),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes pointer":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune{}
			return log30.StringReflect("reflect empty runes pointer", &p)
		}(),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes pointer":[]
		}`,
	},
	{
		line:         line(),
		input:        log30.String("string", "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.String("empty string", ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.String("string with zero byte", string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any string", "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any empty string", ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any string with zero byte", string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect string", "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect empty string", ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect string with zero byte", string(byte(0))),
		expected:     "\u0000",
		expectedText: "\u0000",
		expectedJSON: `{
			"reflect string with zero byte":"\u0000"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := "Hello, Wörld!"
			return log30.StringStringp("string pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := ""
			return log30.StringStringp("empty string pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringStringp("nil string pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil string pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := "Hello, Wörld!"
			return log30.StringAny("any string pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := ""
			return log30.StringAny("any empty string pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string pointer":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := "Hello, Wörld!"
			return log30.StringReflect("reflect string pointer", &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := ""
			return log30.StringReflect("reflect empty string pointer", &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint("uint", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any uint", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect uint", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint = 42
			return log30.StringUintp("uint pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUintp("nil uint pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint = 42
			return log30.StringAny("any uint pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint = 42
			return log30.StringReflect("reflect uint pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint16("uint16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any uint16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect uint16", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint16 = 42
			return log30.StringUint16p("uint16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint16p("uint16 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"uint16 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint16 = 42
			return log30.StringAny("any uint16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint16 = 42
			return log30.StringReflect("reflect uint16 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i *uint16
			return log30.StringReflect("reflect uint16 pointer to nil", i)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect uint16 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint32("uint32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any uint32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect uint32", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint32 = 42
			return log30.StringUint32p("uint32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint32p("nil uint32 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint32 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint32 = 42
			return log30.StringAny("any uint32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint32 = 42
			return log30.StringReflect("reflect uint32 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint64("uint64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64":42
		}`,
	},

	{
		line:         line(),
		input:        log30.StringAny("any uint64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect uint64", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint64 = 42
			return log30.StringUint64p("uint64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint64p("nil uint64 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint64 = 42
			return log30.StringAny("any uint64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint64 = 42
			return log30.StringReflect("reflect uint64 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint8("uint8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any uint8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect uint8", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint8 = 42
			return log30.StringUint8p("uint8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUint8p("nil uint8 pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint8 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint8 = 42
			return log30.StringAny("any uint8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint8 = 42
			return log30.StringReflect("reflect uint8 pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUintptr("uintptr", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any uintptr", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect uintptr", 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uintptr = 42
			return log30.StringUintptrp("uintptr pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringUintptrp("nil uintptr pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uintptr pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uintptr = 42
			return log30.StringAny("any uintptr pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uintptr = 42
			return log30.StringReflect("reflect uintptr pointer", &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringTime("time", time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any time", time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect time", time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return log30.StringTimep("time pointer", &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var t *time.Time
			return log30.StringTimep("nil time pointer", t)
		}(),
		expected:     "0001-01-01 00:00:00 +0000 UTC",
		expectedText: "0001-01-01T00:00:00Z",
		expectedJSON: `{
			"nil time pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			return log30.StringFunc("function", func() json.Marshaler {
				t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
				return marshal30.Time(t)
			})
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"function":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return log30.StringAny("any time pointer", &t)
		}(),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return log30.StringReflect("reflect time pointer", &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringDuration("duration", 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any duration", 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect duration", 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			d := 42 * time.Nanosecond
			return log30.StringDurationp("duration pointer", &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration pointer":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        log30.StringDurationp("nil duration pointer", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil duration pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			d := 42 * time.Nanosecond
			return log30.StringAny("any duration pointer", &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration pointer":"42ns"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			d := 42 * time.Nanosecond
			return log30.StringReflect("reflect duration pointer", &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any struct", Struct{Name: "John Doe", Age: 42}),
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
		input: func() log30.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return log30.StringAny("any struct pointer", &s)
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
		input:        log30.StringReflect("struct reflect", Struct{Name: "John Doe", Age: 42}),
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
		input: func() log30.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return log30.StringReflect("struct reflect pointer", &s)
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
		input:        log30.StringRaw("raw json", []byte(`{"foo":"bar"}`)),
		expected:     `{"foo":"bar"}`,
		expectedText: `{"foo":"bar"}`,
		expectedJSON: `{
			"raw json":{"foo":"bar"}
		}`,
	},
	{
		line:         line(),
		input:        log30.StringRaw("raw malformed json object", []byte(`xyz{"foo":"bar"}`)),
		expected:     `xyz{"foo":"bar"}`,
		expectedText: `xyz{"foo":"bar"}`,
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character 'x' looking for beginning of value"),
	},
	{
		line:         line(),
		input:        log30.StringRaw("raw malformed json key/value", []byte(`{"foo":"bar""}`)),
		expected:     `{"foo":"bar""}`,
		expectedText: `{"foo":"bar""}`,
		error:        errors.New(`json: error calling MarshalJSON for type json.Marshaler: invalid character '"' after object key:value pair`),
	},
	{
		line:         line(),
		input:        log30.StringRaw("raw json with unescaped null byte", append([]byte(`{"foo":"`), append([]byte{0}, []byte(`xyz"}`)...)...)),
		expected:     "{\"foo\":\"\u0000xyz\"}",
		expectedText: "{\"foo\":\"\u0000xyz\"}",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character '\\x00' in string literal"),
	},
	{
		line:         line(),
		input:        log30.StringRaw("raw nil", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"raw nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any byte array", [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			a := [3]byte{'f', 'o', 'o'}
			return log30.StringAny("any byte array pointer", &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var a *[3]byte
			return log30.StringAny("any byte array pointer to nil", a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect byte array", [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			a := [3]byte{'f', 'o', 'o'}
			return log30.StringReflect("reflect byte array pointer", &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var a *[3]byte
			return log30.StringReflect("reflect byte array pointer to nil", a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringAny("any untyped nil", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.StringReflect("reflect untyped nil", nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBool(marshal30.String("bool true"), true),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool true":true
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBool(marshal30.String("bool false"), false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool false":false
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any bool false"), false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"any bool false":false
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("reflect bool false"), false),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"reflect bool false":false
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			return log30.TextBoolp(marshal30.String("bool pointer to true"), &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := false
			return log30.TextBoolp(marshal30.String("bool pointer to false"), &b)
		}(),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool pointer to false":false
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBoolp(marshal30.String("bool nil pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"bool nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			return log30.TextAny(marshal30.String("any bool pointer to true"), &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			b2 := &b
			return log30.TextAny(marshal30.String("any twice pointer to bool true"), &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any twice pointer to bool true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			return log30.TextReflect(marshal30.String("reflect bool pointer to true"), &b)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			b := true
			b2 := &b
			return log30.TextReflect(marshal30.String("reflect bool twice pointer to true"), &b2)
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool twice pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var b *bool
			return log30.TextReflect(marshal30.String("reflect bool pointer to nil"), b)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect bool pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("bytes"), []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("bytes with quote"), []byte(`Hello, "World"!`)),
		expected:     `Hello, \"World\"!`,
		expectedText: `Hello, \"World\"!`,
		expectedJSON: `{
			"bytes with quote":"Hello, \"World\"!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("bytes quote"), []byte(`"Hello, World!"`)),
		expected:     `\"Hello, World!\"`,
		expectedText: `\"Hello, World!\"`,
		expectedJSON: `{
			"bytes quote":"\"Hello, World!\""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("bytes nested quote"), []byte(`"Hello, "World"!"`)),
		expected:     `\"Hello, \"World\"!\"`,
		expectedText: `\"Hello, \"World\"!\"`,
		expectedJSON: `{
			"bytes nested quote":"\"Hello, \"World\"!\""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("bytes json"), []byte(`{"foo":"bar"}`)),
		expected:     `{\"foo\":\"bar\"}`,
		expectedText: `{\"foo\":\"bar\"}`,
		expectedJSON: `{
			"bytes json":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("bytes json quote"), []byte(`"{"foo":"bar"}"`)),
		expected:     `\"{\"foo\":\"bar\"}\"`,
		expectedText: `\"{\"foo\":\"bar\"}\"`,
		expectedJSON: `{
			"bytes json quote":"\"{\"foo\":\"bar\"}\""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("empty bytes"), []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytes(marshal30.String("nil bytes"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any bytes"), []byte("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any empty bytes"), []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect bytes"), []byte("Hello, Wörld!")),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect empty bytes"), []byte{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty bytes":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte("Hello, Wörld!")
			return log30.TextBytesp(marshal30.String("bytes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte{}
			return log30.TextBytesp(marshal30.String("empty bytes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextBytesp(marshal30.String("nil bytes pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte("Hello, Wörld!")
			return log30.TextAny(marshal30.String("any bytes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte{}
			return log30.TextAny(marshal30.String("any empty bytes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte("Hello, Wörld!")
			return log30.TextReflect(marshal30.String("reflect bytes pointer"), &p)
		}(),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes pointer":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []byte{}
			return log30.TextReflect(marshal30.String("reflect empty bytes pointer"), &p)
		}(),
		expected: "",
		expectedJSON: `{
			"reflect empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextComplex128(marshal30.String("complex128"), complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any complex128"), complex(1, 23)),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect complex128"), complex(1, 23)),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex128 = complex(1, 23)
			return log30.TextComplex128p(marshal30.String("complex128 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextComplex128p(marshal30.String("nil complex128 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex128 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex128 = complex(1, 23)
			return log30.TextAny(marshal30.String("any complex128 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex128 = complex(1, 23)
			return log30.TextReflect(marshal30.String("reflect complex128 pointer"), &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        log30.TextComplex64(marshal30.String("complex64"), complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any complex64"), complex(3, 21)),
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"any complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect complex64"), complex(3, 21)),
		expected:     "(3+21i)",
		expectedText: "(3+21i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        log30.TextError(marshal30.String("error"), errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextError(marshal30.String("nil error"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil error":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any error"), errors.New("something went wrong")),
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"any error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect error"), errors.New("something went wrong")),
		expected:     "{something went wrong}",
		expectedText: "{something went wrong}",
		expectedJSON: `{
			"reflect error":{}
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex64 = complex(1, 23)
			return log30.TextComplex64p(marshal30.String("complex64 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex64 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextComplex64p(marshal30.String("nil complex64 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex64 = complex(1, 23)
			return log30.TextAny(marshal30.String("any complex64 pointer"), &c)
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex64 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var c complex64 = complex(1, 23)
			return log30.TextReflect(marshal30.String("reflect complex64 pointer"), &c)
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex64"),
	},
	{
		line:         line(),
		input:        log30.TextFloat32(marshal30.String("float32"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat32(marshal30.String("high precision float32"), 0.123456789),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat32(marshal30.String("zero float32"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any float32"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any zero float32"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect float32"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect zero float32"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float32":0
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 4.2
			return log30.TextFloat32p(marshal30.String("float32 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 0.123456789
			return log30.TextFloat32p(marshal30.String("high precision float32 pointer"), &f)
		}(),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat32p(marshal30.String("float32 nil pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float32 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 4.2
			return log30.TextAny(marshal30.String("any float32 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float32 = 4.2
			return log30.TextReflect(marshal30.String("reflect float32 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f *float32
			return log30.TextReflect(marshal30.String("reflect float32 pointer to nil"), f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float32 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat64(marshal30.String("float64"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat64(marshal30.String("high precision float64"), 0.123456789),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat64(marshal30.String("zero float64"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any float64"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any zero float64"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect float64"), 4.2),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect zero float64"), 0),
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float64":0
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 4.2
			return log30.TextFloat64p(marshal30.String("float64 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 0.123456789
			return log30.TextFloat64p(marshal30.String("high precision float64 pointer"), &f)
		}(),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        log30.TextFloat64p(marshal30.String("float64 nil pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float64 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 4.2
			return log30.TextAny(marshal30.String("any float64 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f float64 = 4.2
			return log30.TextReflect(marshal30.String("reflect float64 pointer"), &f)
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var f *float64
			return log30.TextReflect(marshal30.String("reflect float64 pointer to nil"), f)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float64 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextInt(marshal30.String("int"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any int"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect int"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int = 42
			return log30.TextIntp(marshal30.String("int pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int = 42
			return log30.TextAny(marshal30.String("any int pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int = 42
			return log30.TextReflect(marshal30.String("reflect int pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextInt16(marshal30.String("int16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any int16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect int16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int16 = 42
			return log30.TextInt16p(marshal30.String("int16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int16 = 42
			return log30.TextAny(marshal30.String("any int16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int16 = 42
			return log30.TextReflect(marshal30.String("reflect int16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextInt32(marshal30.String("int32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any int32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect int32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int32 = 42
			return log30.TextInt32p(marshal30.String("int32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int32 = 42
			return log30.TextAny(marshal30.String("any int32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int32 = 42
			return log30.TextReflect(marshal30.String("reflect int32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextInt64(marshal30.String("int64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any int64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect int64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int64 = 42
			return log30.TextInt64p(marshal30.String("int64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int64 = 42
			return log30.TextAny(marshal30.String("any int64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int64 = 42
			return log30.TextReflect(marshal30.String("reflect int64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextInt8(marshal30.String("int8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any int8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect int8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int8 = 42
			return log30.TextInt8p(marshal30.String("int8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int8 = 42
			return log30.TextAny(marshal30.String("any int8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i int8 = 42
			return log30.TextReflect(marshal30.String("reflect int8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextRunes(marshal30.String("runes"), []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextRunes(marshal30.String("empty runes"), []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextRunes(marshal30.String("nil runes"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextRunes(marshal30.String("rune slice with zero rune"), []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any runes"), []rune("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any empty runes"), []rune{}),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any rune slice with zero rune"), []rune{rune(0)}),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect runes"), []rune("Hello, Wörld!")),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect empty runes"), []rune{}),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes":[]
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect rune slice with zero rune"), []rune{rune(0)}),
		expected:     "[0]",
		expectedText: "[0]",
		expectedJSON: `{
			"reflect rune slice with zero rune":[0]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune("Hello, Wörld!")
			return log30.TextRunesp(marshal30.String("runes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune{}
			return log30.TextRunesp(marshal30.String("empty runes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextRunesp(marshal30.String("nil runes pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune("Hello, Wörld!")
			return log30.TextAny(marshal30.String("any runes pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune{}
			return log30.TextAny(marshal30.String("any empty runes pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune("Hello, Wörld!")
			return log30.TextReflect(marshal30.String("reflect runes pointer"), &p)
		}(),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes pointer":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := []rune{}
			return log30.TextReflect(marshal30.String("reflect empty runes pointer"), &p)
		}(),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes pointer":[]
		}`,
	},
	{
		line:         line(),
		input:        log30.Text(marshal30.String("string"), marshal30.String("Hello, Wörld!")),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.Text(marshal30.String("empty string"), marshal30.String("")),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.Text(marshal30.String("string with zero byte"), marshal30.String((string(byte(0))))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextString(marshal30.String("string"), "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextString(marshal30.String("empty string"), ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextString(marshal30.String("string with zero byte"), string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any string"), "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any empty string"), ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any string with zero byte"), string(byte(0))),
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect string"), "Hello, Wörld!"),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect empty string"), ""),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect string with zero byte"), string(byte(0))),
		expected:     "\u0000",
		expectedText: "\u0000",
		expectedJSON: `{
			"reflect string with zero byte":"\u0000"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := "Hello, Wörld!"
			return log30.TextStringp(marshal30.String("string pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := ""
			return log30.TextStringp(marshal30.String("empty string pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextStringp(marshal30.String("nil string pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil string pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := "Hello, Wörld!"
			return log30.TextAny(marshal30.String("any string pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := ""
			return log30.TextAny(marshal30.String("any empty string pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string pointer":""
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := "Hello, Wörld!"
			return log30.TextReflect(marshal30.String("reflect string pointer"), &p)
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			p := ""
			return log30.TextReflect(marshal30.String("reflect empty string pointer"), &p)
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint(marshal30.String("uint"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any uint"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect uint"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint = 42
			return log30.TextUintp(marshal30.String("uint pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUintp(marshal30.String("nil uint pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint = 42
			return log30.TextAny(marshal30.String("any uint pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint = 42
			return log30.TextReflect(marshal30.String("reflect uint pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint16(marshal30.String("uint16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any uint16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect uint16"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint16 = 42
			return log30.TextUint16p(marshal30.String("uint16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint16p(marshal30.String("uint16 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"uint16 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint16 = 42
			return log30.TextAny(marshal30.String("any uint16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint16 = 42
			return log30.TextReflect(marshal30.String("reflect uint16 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i *uint16
			return log30.TextReflect(marshal30.String("reflect uint16 pointer to nil"), i)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect uint16 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint32(marshal30.String("uint32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any uint32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect uint32"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint32 = 42
			return log30.TextUint32p(marshal30.String("uint32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint32p(marshal30.String("nil uint32 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint32 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint32 = 42
			return log30.TextAny(marshal30.String("any uint32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint32 = 42
			return log30.TextReflect(marshal30.String("reflect uint32 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint64(marshal30.String("uint64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64":42
		}`,
	},

	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any uint64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect uint64"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint64 = 42
			return log30.TextUint64p(marshal30.String("uint64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint64p(marshal30.String("nil uint64 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint64 = 42
			return log30.TextAny(marshal30.String("any uint64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint64 = 42
			return log30.TextReflect(marshal30.String("reflect uint64 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint8(marshal30.String("uint8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any uint8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect uint8"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint8 = 42
			return log30.TextUint8p(marshal30.String("uint8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUint8p(marshal30.String("nil uint8 pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint8 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint8 = 42
			return log30.TextAny(marshal30.String("any uint8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uint8 = 42
			return log30.TextReflect(marshal30.String("reflect uint8 pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUintptr(marshal30.String("uintptr"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any uintptr"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect uintptr"), 42),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uintptr = 42
			return log30.TextUintptrp(marshal30.String("uintptr pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextUintptrp(marshal30.String("nil uintptr pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uintptr pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uintptr = 42
			return log30.TextAny(marshal30.String("any uintptr pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr pointer":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var i uintptr = 42
			return log30.TextReflect(marshal30.String("reflect uintptr pointer"), &i)
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextTime(marshal30.String("time"), time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any time"), time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect time"), time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return log30.TextTimep(marshal30.String("time pointer"), &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var t *time.Time
			return log30.TextTimep(marshal30.String("nil time pointer"), t)
		}(),
		expected:     "0001-01-01 00:00:00 +0000 UTC",
		expectedText: "0001-01-01T00:00:00Z",
		expectedJSON: `{
			"nil time pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			return log30.TextFunc(marshal30.String("function"), func() json.Marshaler {
				t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
				return marshal30.Time(t)
			})
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"function":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return log30.TextAny(marshal30.String("any time pointer"), &t)
		}(),
		expected:     `"1970-01-01T00:00:00.000000042Z"`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return log30.TextReflect(marshal30.String("reflect time pointer"), &t)
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextDuration(marshal30.String("duration"), 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any duration"), 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect duration"), 42*time.Nanosecond),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration":42
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			d := 42 * time.Nanosecond
			return log30.TextDurationp(marshal30.String("duration pointer"), &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration pointer":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        log30.TextDurationp(marshal30.String("nil duration pointer"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil duration pointer":null
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			d := 42 * time.Nanosecond
			return log30.TextAny(marshal30.String("any duration pointer"), &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration pointer":"42ns"
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			d := 42 * time.Nanosecond
			return log30.TextReflect(marshal30.String("reflect duration pointer"), &d)
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration pointer":42
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any struct"), Struct{Name: "John Doe", Age: 42}),
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
		input: func() log30.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return log30.TextAny(marshal30.String("any struct pointer"), &s)
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
		input:        log30.TextReflect(marshal30.String("struct reflect"), Struct{Name: "John Doe", Age: 42}),
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
		input: func() log30.KV {
			s := Struct{Name: "John Doe", Age: 42}
			return log30.TextReflect(marshal30.String("struct reflect pointer"), &s)
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
		input:        log30.TextRaw(marshal30.String("raw json"), []byte(`{"foo":"bar"}`)),
		expected:     `{"foo":"bar"}`,
		expectedText: `{"foo":"bar"}`,
		expectedJSON: `{
			"raw json":{"foo":"bar"}
		}`,
	},
	{
		line:         line(),
		input:        log30.TextRaw(marshal30.String("raw malformed json object"), []byte(`xyz{"foo":"bar"}`)),
		expected:     `xyz{"foo":"bar"}`,
		expectedText: `xyz{"foo":"bar"}`,
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character 'x' looking for beginning of value"),
	},
	{
		line:         line(),
		input:        log30.TextRaw(marshal30.String("raw malformed json key/value"), []byte(`{"foo":"bar""}`)),
		expected:     `{"foo":"bar""}`,
		expectedText: `{"foo":"bar""}`,
		error:        errors.New(`json: error calling MarshalJSON for type json.Marshaler: invalid character '"' after object key:value pair`),
	},
	{
		line:         line(),
		input:        log30.TextRaw(marshal30.String("raw json with unescaped null byte"), append([]byte(`{"foo":"`), append([]byte{0}, []byte(`xyz"}`)...)...)),
		expected:     "{\"foo\":\"\u0000xyz\"}",
		expectedText: "{\"foo\":\"\u0000xyz\"}",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character '\\x00' in string literal"),
	},
	{
		line:         line(),
		input:        log30.TextRaw(marshal30.String("raw nil"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"raw nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any byte array"), [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			a := [3]byte{'f', 'o', 'o'}
			return log30.TextAny(marshal30.String("any byte array pointer"), &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var a *[3]byte
			return log30.TextAny(marshal30.String("any byte array pointer to nil"), a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect byte array"), [3]byte{'f', 'o', 'o'}),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			a := [3]byte{'f', 'o', 'o'}
			return log30.TextReflect(marshal30.String("reflect byte array pointer"), &a)
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() log30.KV {
			var a *[3]byte
			return log30.TextReflect(marshal30.String("reflect byte array pointer to nil"), a)
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextAny(marshal30.String("any untyped nil"), nil),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        log30.TextReflect(marshal30.String("reflect untyped nil"), nil),
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

			if !equal4.ErrorEqual(err, tc.error) {
				t.Fatalf("marshal error expected: %s, recieved: %s %s", tc.error, err, linkToExample)
			}

			if err == nil {
				ja := jsonassert.New(testprinter{t: t, link: linkToExample})
				ja.Assertf(string(p), tc.expectedJSON)
			}
		})
	}
}
