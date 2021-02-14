package marshal64_test

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/danil/eq32"
	"github.com/danil/log64/marshal64"
	"github.com/kinbiko/jsonassert"
)

var MarshalTestCases = []struct {
	line         int
	input        map[string]json.Marshaler
	expected     string
	expectedText string
	expectedJSON string
	error        error
	benchmark    bool
}{
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bool true": marshal64.Bool(true)},
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool true":true
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bool false": marshal64.Bool(false)},
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool false":false
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any bool false": marshal64.Any(false)},
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"any bool false":false
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect bool false": marshal64.Any(false)},
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"reflect bool false":false
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			return map[string]json.Marshaler{"bool pointer to true": marshal64.Boolp(&b)}
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := false
			return map[string]json.Marshaler{"bool pointer to false": marshal64.Boolp(&b)}
		}(),
		expected:     "false",
		expectedText: "false",
		expectedJSON: `{
			"bool pointer to false":false
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bool nil pointer": marshal64.Boolp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"bool nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			return map[string]json.Marshaler{"any bool pointer to true": marshal64.Any(&b)}
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			b2 := &b
			return map[string]json.Marshaler{"any twice pointer to bool true": marshal64.Any(&b2)}
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"any twice pointer to bool true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			return map[string]json.Marshaler{"reflect bool pointer to true": marshal64.Reflect(&b)}
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			b2 := &b
			return map[string]json.Marshaler{"reflect bool twice pointer to true": marshal64.Reflect(&b2)}
		}(),
		expected:     "true",
		expectedText: "true",
		expectedJSON: `{
			"reflect bool twice pointer to true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var b *bool
			return map[string]json.Marshaler{"reflect bool pointer to nil": marshal64.Reflect(b)}
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect bool pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bytes": marshal64.Bytes([]byte("Hello, Wörld!"))},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bytes with quote": marshal64.Bytes([]byte(`Hello, "World"!`))},
		expected:     `Hello, \"World\"!`,
		expectedText: `Hello, \"World\"!`,
		expectedJSON: `{
			"bytes with quote":"Hello, \"World\"!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bytes quote": marshal64.Bytes([]byte(`"Hello, World!"`))},
		expected:     `\"Hello, World!\"`,
		expectedText: `\"Hello, World!\"`,
		expectedJSON: `{
			"bytes quote":"\"Hello, World!\""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bytes nested quote": marshal64.Bytes([]byte(`"Hello, "World"!"`))},
		expected:     `\"Hello, \"World\"!\"`,
		expectedText: `\"Hello, \"World\"!\"`,
		expectedJSON: `{
			"bytes nested quote":"\"Hello, \"World\"!\""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bytes json": marshal64.Bytes([]byte(`{"foo":"bar"}`))},
		expected:     `{\"foo\":\"bar\"}`,
		expectedText: `{\"foo\":\"bar\"}`,
		expectedJSON: `{
			"bytes json":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"bytes json quote": marshal64.Bytes([]byte(`"{"foo":"bar"}"`))},
		expected:     `\"{\"foo\":\"bar\"}\"`,
		expectedText: `\"{\"foo\":\"bar\"}\"`,
		expectedJSON: `{
			"bytes json quote":"\"{\"foo\":\"bar\"}\""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"empty bytes": marshal64.Bytes([]byte{})},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil bytes": marshal64.Bytes(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any bytes": marshal64.Any([]byte("Hello, Wörld!"))},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any empty bytes": marshal64.Any([]byte{})},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect bytes": marshal64.Reflect([]byte("Hello, Wörld!"))},
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect empty bytes": marshal64.Reflect([]byte{})},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty bytes":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte("Hello, Wörld!")
			return map[string]json.Marshaler{"bytes pointer": marshal64.Bytesp(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte{}
			return map[string]json.Marshaler{"empty bytes pointer": marshal64.Bytesp(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil bytes pointer": marshal64.Bytesp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil bytes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte("Hello, Wörld!")
			return map[string]json.Marshaler{"any bytes pointer": marshal64.Any(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any bytes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte{}
			return map[string]json.Marshaler{"any empty bytes pointer": marshal64.Any(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty bytes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte("Hello, Wörld!")
			return map[string]json.Marshaler{"reflect bytes pointer": marshal64.Reflect(&p)}
		}(),
		expected:     "SGVsbG8sIFfDtnJsZCE=",
		expectedText: "SGVsbG8sIFfDtnJsZCE=",
		expectedJSON: `{
			"reflect bytes pointer":"SGVsbG8sIFfDtnJsZCE="
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte{}
			return map[string]json.Marshaler{"reflect empty bytes pointer": marshal64.Reflect(&p)}
		}(),
		expected: "",
		expectedJSON: `{
			"reflect empty bytes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"complex128": marshal64.Complex128(complex(1, 23))},
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any complex128": marshal64.Any(complex(1, 23))},
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect complex128": marshal64.Reflect(complex(1, 23))},
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex128 = complex(1, 23)
			return map[string]json.Marshaler{"complex128 pointer": marshal64.Complex128p(&c)}
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex128 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil complex128 pointer": marshal64.Complex128p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex128 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex128 = complex(1, 23)
			return map[string]json.Marshaler{"any complex128 pointer": marshal64.Any(&c)}
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex128 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex128 = complex(1, 23)
			return map[string]json.Marshaler{"reflect complex128 pointer": marshal64.Reflect(&c)}
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"complex64": marshal64.Complex64(complex(3, 21))},
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any complex64": marshal64.Any(complex(3, 21))},
		expected:     "3+21i",
		expectedText: "3+21i",
		expectedJSON: `{
			"any complex64":"3+21i"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect complex64": marshal64.Reflect(complex(3, 21))},
		expected:     "(3+21i)",
		expectedText: "(3+21i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"error": marshal64.Error(errors.New("something went wrong"))},
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil error": marshal64.Error(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil error":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any error": marshal64.Any(errors.New("something went wrong"))},
		expected:     "something went wrong",
		expectedText: "something went wrong",
		expectedJSON: `{
			"any error":"something went wrong"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect error": marshal64.Reflect(errors.New("something went wrong"))},
		expected:     "{something went wrong}",
		expectedText: "{something went wrong}",
		expectedJSON: `{
			"reflect error":{}
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex64 = complex(1, 23)
			return map[string]json.Marshaler{"complex64 pointer": marshal64.Complex64p(&c)}
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"complex64 pointer":"1+23i"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil complex64 pointer": marshal64.Complex64p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil complex64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex64 = complex(1, 23)
			return map[string]json.Marshaler{"any complex64 pointer": marshal64.Any(&c)}
		}(),
		expected:     "1+23i",
		expectedText: "1+23i",
		expectedJSON: `{
			"any complex64 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex64 = complex(1, 23)
			return map[string]json.Marshaler{"reflect complex64 pointer": marshal64.Reflect(&c)}
		}(),
		expected:     "(1+23i)",
		expectedText: "(1+23i)",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex64"),
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"float32": marshal64.Float32(4.2)},
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"high precision float32": marshal64.Float32(0.123456789)},
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"zero float32": marshal64.Float32(0)},
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any float32": marshal64.Any(4.2)},
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any zero float32": marshal64.Any(0)},
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float32":0
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect float32": marshal64.Reflect(4.2)},
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32":4.2
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect zero float32": marshal64.Reflect(0)},
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float32":0
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 4.2
			return map[string]json.Marshaler{"float32 pointer": marshal64.Float32p(&f)}
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 0.123456789
			return map[string]json.Marshaler{"high precision float32 pointer": marshal64.Float32p(&f)}
		}(),
		expected:     "0.12345679",
		expectedText: "0.12345679",
		expectedJSON: `{
			"high precision float32 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"float32 nil pointer": marshal64.Float32p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float32 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 4.2
			return map[string]json.Marshaler{"any float32 pointer": marshal64.Any(&f)}
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 4.2
			return map[string]json.Marshaler{"reflect float32 pointer": marshal64.Reflect(&f)}
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f *float32
			return map[string]json.Marshaler{"reflect float32 pointer to nil": marshal64.Reflect(f)}
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float32 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"float64": marshal64.Float64(4.2)},
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"high precision float64": marshal64.Float64(0.123456789)},
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"zero float64": marshal64.Float64(0)},
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any float64": marshal64.Any(4.2)},
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any zero float64": marshal64.Any(0)},
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"any zero float64":0
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect float64": marshal64.Reflect(4.2)},
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64":4.2
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect zero float64": marshal64.Reflect(0)},
		expected:     "0",
		expectedText: "0",
		expectedJSON: `{
			"reflect zero float64":0
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 4.2
			return map[string]json.Marshaler{"float64 pointer": marshal64.Float64p(&f)}
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 0.123456789
			return map[string]json.Marshaler{"high precision float64 pointer": marshal64.Float64p(&f)}
		}(),
		expected:     "0.123456789",
		expectedText: "0.123456789",
		expectedJSON: `{
			"high precision float64 pointer":0.123456789
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"float64 nil pointer": marshal64.Float64p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"float64 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 4.2
			return map[string]json.Marshaler{"any float64 pointer": marshal64.Any(&f)}
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"any float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 4.2
			return map[string]json.Marshaler{"reflect float64 pointer": marshal64.Reflect(&f)}
		}(),
		expected:     "4.2",
		expectedText: "4.2",
		expectedJSON: `{
			"reflect float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f *float64
			return map[string]json.Marshaler{"reflect float64 pointer to nil": marshal64.Reflect(f)}
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect float64 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"int": marshal64.Int(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any int": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect int": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int = 42
			return map[string]json.Marshaler{"int pointer": marshal64.Intp(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int = 42
			return map[string]json.Marshaler{"any int pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int = 42
			return map[string]json.Marshaler{"reflect int pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"int16": marshal64.Int16(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any int16": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect int16": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int16 = 42
			return map[string]json.Marshaler{"int16 pointer": marshal64.Int16p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int16 = 42
			return map[string]json.Marshaler{"any int16 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int16 = 42
			return map[string]json.Marshaler{"reflect int16 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"int32": marshal64.Int32(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any int32": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect int32": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int32 = 42
			return map[string]json.Marshaler{"int32 pointer": marshal64.Int32p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int32 = 42
			return map[string]json.Marshaler{"any int32 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int32 = 42
			return map[string]json.Marshaler{"reflect int32 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"int64": marshal64.Int64(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any int64": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect int64": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int64 = 42
			return map[string]json.Marshaler{"int64 pointer": marshal64.Int64p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int64 = 42
			return map[string]json.Marshaler{"any int64 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int64 = 42
			return map[string]json.Marshaler{"reflect int64 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"int8": marshal64.Int8(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any int8": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect int8": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int8 = 42
			return map[string]json.Marshaler{"int8 pointer": marshal64.Int8p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int8 = 42
			return map[string]json.Marshaler{"any int8 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int8 = 42
			return map[string]json.Marshaler{"reflect int8 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect int8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"runes": marshal64.Runes([]rune("Hello, Wörld!"))},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"empty runes": marshal64.Runes([]rune{})},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil runes": marshal64.Runes(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"rune slice with zero rune": marshal64.Runes([]rune{rune(0)})},
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any runes": marshal64.Any([]rune("Hello, Wörld!"))},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any empty runes": marshal64.Any([]rune{})},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any rune slice with zero rune": marshal64.Any([]rune{rune(0)})},
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect runes": marshal64.Reflect([]rune("Hello, Wörld!"))},
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect empty runes": marshal64.Reflect([]rune{})},
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes":[]
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect rune slice with zero rune": marshal64.Reflect([]rune{rune(0)})},
		expected:     "[0]",
		expectedText: "[0]",
		expectedJSON: `{
			"reflect rune slice with zero rune":[0]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune("Hello, Wörld!")
			return map[string]json.Marshaler{"runes pointer": marshal64.Runesp(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune{}
			return map[string]json.Marshaler{"empty runes pointer": marshal64.Runesp(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty runes pointer":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil runes pointer": marshal64.Runesp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil runes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune("Hello, Wörld!")
			return map[string]json.Marshaler{"any runes pointer": marshal64.Any(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any runes pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune{}
			return map[string]json.Marshaler{"any empty runes pointer": marshal64.Any(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty runes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune("Hello, Wörld!")
			return map[string]json.Marshaler{"reflect runes pointer": marshal64.Reflect(&p)}
		}(),
		expected:     "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedText: "[72 101 108 108 111 44 32 87 246 114 108 100 33]",
		expectedJSON: `{
			"reflect runes pointer":[72,101,108,108,111,44,32,87,246,114,108,100,33]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune{}
			return map[string]json.Marshaler{"reflect empty runes pointer": marshal64.Reflect(&p)}
		}(),
		expected:     "[]",
		expectedText: "[]",
		expectedJSON: `{
			"reflect empty runes pointer":[]
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"string": marshal64.String("Hello, Wörld!")},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"empty string": marshal64.String("")},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"string with zero byte": marshal64.String(string(byte(0)))},
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any string": marshal64.Any("Hello, Wörld!")},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any empty string": marshal64.Any("")},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any string with zero byte": marshal64.Any(string(byte(0)))},
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"any string with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect string": marshal64.Reflect("Hello, Wörld!")},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect empty string": marshal64.Reflect("")},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect string with zero byte": marshal64.Reflect(string(byte(0)))},
		expected:     "\u0000",
		expectedText: "\u0000",
		expectedJSON: `{
			"reflect string with zero byte":"\u0000"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := "Hello, Wörld!"
			return map[string]json.Marshaler{"string pointer": marshal64.Stringp(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := ""
			return map[string]json.Marshaler{"empty string pointer": marshal64.Stringp(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil string pointer": marshal64.Stringp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil string pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := "Hello, Wörld!"
			return map[string]json.Marshaler{"any string pointer": marshal64.Any(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"any string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := ""
			return map[string]json.Marshaler{"any empty string pointer": marshal64.Any(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"any empty string pointer":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := "Hello, Wörld!"
			return map[string]json.Marshaler{"reflect string pointer": marshal64.Reflect(&p)}
		}(),
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"reflect string pointer":"Hello, Wörld!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := ""
			return map[string]json.Marshaler{"reflect empty string pointer": marshal64.Reflect(&p)}
		}(),
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"reflect empty string pointer":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"text": marshal64.Text(marshal64.String("Hello, Wörld!"))},
		expected:     "Hello, Wörld!",
		expectedText: "Hello, Wörld!",
		expectedJSON: `{
			"text":"Hello, Wörld!"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"empty text": marshal64.Text(marshal64.String(""))},
		expected:     "",
		expectedText: "",
		expectedJSON: `{
			"empty text":""
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"text with zero byte": marshal64.Text(marshal64.String(string(byte(0))))},
		expected:     "\\u0000",
		expectedText: "\\u0000",
		expectedJSON: `{
			"text with zero byte":"\u0000"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uint": marshal64.Uint(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any uint": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect uint": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint = 42
			return map[string]json.Marshaler{"uint pointer": marshal64.Uintp(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil uint pointer": marshal64.Uintp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint = 42
			return map[string]json.Marshaler{"any uint pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint = 42
			return map[string]json.Marshaler{"reflect uint pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uint16": marshal64.Uint16(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any uint16": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect uint16": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint16 = 42
			return map[string]json.Marshaler{"uint16 pointer": marshal64.Uint16p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint16 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uint16 pointer": marshal64.Uint16p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"uint16 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint16 = 42
			return map[string]json.Marshaler{"any uint16 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint16 = 42
			return map[string]json.Marshaler{"reflect uint16 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i *uint16
			return map[string]json.Marshaler{"reflect uint16 pointer to nil": marshal64.Reflect(i)}
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect uint16 pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uint32": marshal64.Uint32(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any uint32": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect uint32": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint32 = 42
			return map[string]json.Marshaler{"uint32 pointer": marshal64.Uint32p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil uint32 pointer": marshal64.Uint32p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint32 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint32 = 42
			return map[string]json.Marshaler{"any uint32 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint32 = 42
			return map[string]json.Marshaler{"reflect uint32 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint32 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uint64": marshal64.Uint64(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any uint64": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect uint64": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint64 = 42
			return map[string]json.Marshaler{"uint64 pointer": marshal64.Uint64p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil uint64 pointer": marshal64.Uint64p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint64 = 42
			return map[string]json.Marshaler{"any uint64 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint64 = 42
			return map[string]json.Marshaler{"reflect uint64 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint64 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uint8": marshal64.Uint8(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any uint8": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect uint8": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint8 = 42
			return map[string]json.Marshaler{"uint8 pointer": marshal64.Uint8p(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil uint8 pointer": marshal64.Uint8p(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uint8 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint8 = 42
			return map[string]json.Marshaler{"any uint8 pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uint8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint8 = 42
			return map[string]json.Marshaler{"reflect uint8 pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uint8 pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"uintptr": marshal64.Uintptr(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any uintptr": marshal64.Any(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect uintptr": marshal64.Reflect(42)},
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uintptr = 42
			return map[string]json.Marshaler{"uintptr pointer": marshal64.Uintptrp(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil uintptr pointer": marshal64.Uintptrp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil uintptr pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uintptr = 42
			return map[string]json.Marshaler{"any uintptr pointer": marshal64.Any(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"any uintptr pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uintptr = 42
			return map[string]json.Marshaler{"reflect uintptr pointer": marshal64.Reflect(&i)}
		}(),
		expected:     "42",
		expectedText: "42",
		expectedJSON: `{
			"reflect uintptr pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"time": time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)},
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any time": marshal64.Any(time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC))},
		expected:     `1970-01-01 00:00:00.000000042 +0000 UTC`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect time": marshal64.Reflect(time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC))},
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return map[string]json.Marshaler{"time pointer": &t}
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01T00:00:00.000000042Z",
		expectedJSON: `{
			"time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var t time.Time
			return map[string]json.Marshaler{"nil time pointer": t}
		}(),
		expected:     "0001-01-01 00:00:00 +0000 UTC",
		expectedText: "0001-01-01T00:00:00Z",
		expectedJSON: `{
			"nil time pointer":"0001-01-01T00:00:00Z"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return map[string]json.Marshaler{"any time pointer": marshal64.Any(&t)}
		}(),
		expected:     `1970-01-01 00:00:00.000000042 +0000 UTC`,
		expectedText: `1970-01-01T00:00:00.000000042Z`,
		expectedJSON: `{
			"any time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return map[string]json.Marshaler{"reflect time pointer": marshal64.Reflect(&t)}
		}(),
		expected:     "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedText: "1970-01-01 00:00:00.000000042 +0000 UTC",
		expectedJSON: `{
			"reflect time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"duration": marshal64.Duration(42 * time.Nanosecond)},
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any duration": marshal64.Any(42 * time.Nanosecond)},
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect duration": marshal64.Reflect(42 * time.Nanosecond)},
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			d := 42 * time.Nanosecond
			return map[string]json.Marshaler{"duration pointer": marshal64.Durationp(&d)}
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"duration pointer":"42ns"
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"nil duration pointer": marshal64.Durationp(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"nil duration pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			d := 42 * time.Nanosecond
			return map[string]json.Marshaler{"any duration pointer": marshal64.Any(&d)}
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"any duration pointer":"42ns"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			d := 42 * time.Nanosecond
			return map[string]json.Marshaler{"reflect duration pointer": marshal64.Reflect(&d)}
		}(),
		expected:     "42ns",
		expectedText: "42ns",
		expectedJSON: `{
			"reflect duration pointer":42
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any struct": marshal64.Any(Struct{Name: "John Doe", Age: 42})},
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
		input: func() map[string]json.Marshaler {
			s := Struct{Name: "John Doe", Age: 42}
			return map[string]json.Marshaler{"any struct pointer": marshal64.Any(&s)}
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
		input:        map[string]json.Marshaler{"struct reflect": marshal64.Reflect(Struct{Name: "John Doe", Age: 42})},
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
		input: func() map[string]json.Marshaler {
			s := Struct{Name: "John Doe", Age: 42}
			return map[string]json.Marshaler{"struct reflect pointer": marshal64.Reflect(&s)}
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
		input:        map[string]json.Marshaler{"raw json": marshal64.Raw([]byte(`{"foo":"bar"}`))},
		expected:     `{"foo":"bar"}`,
		expectedText: `{"foo":"bar"}`,
		expectedJSON: `{
			"raw json":{"foo":"bar"}
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"raw malformed json object": marshal64.Raw([]byte(`xyz{"foo":"bar"}`))},
		expected:     `xyz{"foo":"bar"}`,
		expectedText: `xyz{"foo":"bar"}`,
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character 'x' looking for beginning of value"),
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"raw malformed json key/value": marshal64.Raw([]byte(`{"foo":"bar""}`))},
		expected:     `{"foo":"bar""}`,
		expectedText: `{"foo":"bar""}`,
		error:        errors.New(`json: error calling MarshalJSON for type json.Marshaler: invalid character '"' after object key:value pair`),
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"raw json with unescaped null byte": marshal64.Raw(append([]byte(`{"foo":"`), append([]byte{0}, []byte(`xyz"}`)...)...))},
		expected:     "{\"foo\":\"\u0000xyz\"}",
		expectedText: "{\"foo\":\"\u0000xyz\"}",
		error:        errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character '\\x00' in string literal"),
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"raw nil": marshal64.Raw(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"raw nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any byte array": marshal64.Any([3]byte{'f', 'o', 'o'})},
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			a := [3]byte{'f', 'o', 'o'}
			return map[string]json.Marshaler{"any byte array pointer": marshal64.Any(&a)}
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"any byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var a *[3]byte
			return map[string]json.Marshaler{"any byte array pointer to nil": marshal64.Any(a)}
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect byte array": marshal64.Reflect([3]byte{'f', 'o', 'o'})},
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			a := [3]byte{'f', 'o', 'o'}
			return map[string]json.Marshaler{"reflect byte array pointer": marshal64.Reflect(&a)}
		}(),
		expected:     "[102 111 111]",
		expectedText: "[102 111 111]",
		expectedJSON: `{
			"reflect byte array pointer":[102,111,111]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var a *[3]byte
			return map[string]json.Marshaler{"reflect byte array pointer to nil": marshal64.Reflect(a)}
		}(),
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect byte array pointer to nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"any untyped nil": marshal64.Any(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"any untyped nil":null
		}`,
	},
	{
		line:         line(),
		input:        map[string]json.Marshaler{"reflect untyped nil": marshal64.Reflect(nil)},
		expected:     "null",
		expectedText: "null",
		expectedJSON: `{
			"reflect untyped nil":null
		}`,
	},
}

func TestMarshal(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range MarshalTestCases {
		tc := tc
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			for k, v := range tc.input {
				str, ok := v.(fmt.Stringer)
				if !ok {
					t.Errorf("%q does not implement the stringer interface", k)

				} else {
					s := str.String()
					if s != tc.expected {
						t.Errorf("%q unexpected string, expected: %q, recieved: %q %s", k, tc.expected, s, linkToExample)
					}
				}

				txt, ok := v.(encoding.TextMarshaler)
				if !ok {
					t.Errorf("%q does not implement the text marshaler interface", k)

				} else {
					p, err := txt.MarshalText()
					if err != nil {
						t.Fatalf("%q encoding marshal text error: %s %s", k, err, linkToExample)
					}

					if string(p) != tc.expectedText {
						t.Errorf("%q unexpected text, expected: %q, recieved: %q %s", k, tc.expectedText, string(p), linkToExample)
					}
				}
			}

			p, err := json.Marshal(tc.input)

			if !eq32.ErrorEqual(err, tc.error) {
				t.Fatalf("marshal error expected: %s, recieved: %s %s", tc.error, err, linkToExample)
			}

			if err == nil {
				ja := jsonassert.New(testprinter{t: t, link: linkToExample})
				ja.Assertf(string(p), tc.expectedJSON)
			}
		})
	}
}

func line() int { _, _, l, _ := runtime.Caller(1); return l }

type testprinter struct {
	t    *testing.T
	link string
}

func (p testprinter) Errorf(msg string, args ...interface{}) {
	p.t.Errorf(p.link+"\n"+msg, args...)
}

type Struct struct {
	Name string
	Age  int
}
