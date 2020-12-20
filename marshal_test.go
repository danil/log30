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
	"github.com/kinbiko/jsonassert"
)

type Struct struct {
	Name string
	Age  int
}

var MarshalTestCases = []struct {
	line      int
	input     map[string]json.Marshaler
	expected  string
	error     error
	benchmark bool
}{
	{
		line:  line(),
		input: map[string]json.Marshaler{"bool true": logastic.Bool(true)},
		expected: `{
			"bool true":true
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bool false": logastic.Bool(false)},
		expected: `{
			"bool false":false
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any bool false": logastic.Any(false)},
		expected: `{
			"any bool false":false
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect bool false": logastic.Any(false)},
		expected: `{
			"reflect bool false":false
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			return map[string]json.Marshaler{"bool pointer true": logastic.Boolp(&b)}
		}(),
		expected: `{
			"bool pointer true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := false
			return map[string]json.Marshaler{"bool pointer false": logastic.Boolp(&b)}
		}(),
		expected: `{
			"bool pointer false":false
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bool pointer null": logastic.Boolp(nil)},
		expected: `{
			"bool pointer null":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			return map[string]json.Marshaler{"any bool pointer true": logastic.Any(&b)}
		}(),
		expected: `{
			"any bool pointer true":true
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			b := true
			return map[string]json.Marshaler{"reflect bool pointer true": logastic.Reflect(&b)}
		}(),
		expected: `{
			"reflect bool pointer true":true
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect bool pointer null": logastic.Reflect(nil)},
		expected: `{
			"reflect bool pointer null":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bytes": logastic.Bytes([]byte("Hello, World!"))},
		expected: `{
			"bytes":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bytes with quote": logastic.Bytes([]byte(`Hello, "World"!`))},
		expected: `{
			"bytes with quote":"Hello, \"World\"!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bytes quote": logastic.Bytes([]byte(`"Hello, World!"`))},
		expected: `{
			"bytes quote":"\"Hello, World!\""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bytes nested quote": logastic.Bytes([]byte(`"Hello, "World"!"`))},
		expected: `{
			"bytes nested quote":"\"Hello, \"World\"!\""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bytes json": logastic.Bytes([]byte(`{"foo":"bar"}`))},
		expected: `{
			"bytes json":"{\"foo\":\"bar\"}"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"bytes json quote": logastic.Bytes([]byte(`"{"foo":"bar"}"`))},
		expected: `{
			"bytes json quote":"\"{\"foo\":\"bar\"}\""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"empty bytes": logastic.Bytes([]byte{})},
		expected: `{
			"empty bytes":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil bytes": logastic.Bytes(nil)},
		expected: `{
			"nil bytes":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any bytes": logastic.Any([]byte("Hello, World!"))},
		expected: `{
			"any bytes":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any empty bytes": logastic.Any([]byte{})},
		expected: `{
			"any empty bytes":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect bytes": logastic.Reflect([]byte("Hello, World!"))},
		expected: `{
			"reflect bytes":"SGVsbG8sIFdvcmxkIQ=="
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect empty bytes": logastic.Reflect([]byte{})},
		expected: `{
			"reflect empty bytes":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte("Hello, World!")
			return map[string]json.Marshaler{"bytes pointer": logastic.Bytesp(&p)}
		}(),
		expected: `{
			"bytes pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte{}
			return map[string]json.Marshaler{"empty bytes pointer": logastic.Bytesp(&p)}
		}(),
		expected: `{
			"empty bytes pointer":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil bytes pointer": logastic.Bytesp(nil)},
		expected: `{
			"nil bytes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte("Hello, World!")
			return map[string]json.Marshaler{"any bytes pointer": logastic.Any(&p)}
		}(),
		expected: `{
			"any bytes pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte{}
			return map[string]json.Marshaler{"any empty bytes pointer": logastic.Any(&p)}
		}(),
		expected: `{
			"any empty bytes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte("Hello, World!")
			return map[string]json.Marshaler{"reflect bytes pointer": logastic.Reflect(&p)}
		}(),
		expected: `{
			"reflect bytes pointer":"SGVsbG8sIFdvcmxkIQ=="
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []byte{}
			return map[string]json.Marshaler{"reflect empty bytes pointer": logastic.Reflect(&p)}
		}(),
		expected: `{
			"reflect empty bytes pointer":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"complex128": logastic.Complex128(complex(1, 23))},
		expected: `{
			"complex128":"1+23i"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any complex128": logastic.Any(complex(1, 23))},
		expected: `{
			"any complex128":"1+23i"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect complex128": logastic.Reflect(complex(1, 23))},
		error: errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex128 = complex(1, 23)
			return map[string]json.Marshaler{"complex128 pointer": logastic.Complex128p(&c)}
		}(),
		expected: `{
			"complex128 pointer":"1+23i"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil complex128 pointer": logastic.Complex128p(nil)},
		expected: `{
			"nil complex128 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex128 = complex(1, 23)
			return map[string]json.Marshaler{"any complex128 pointer": logastic.Any(&c)}
		}(),
		expected: `{
			"any complex128 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex128 = complex(1, 23)
			return map[string]json.Marshaler{"reflect complex128 pointer": logastic.Reflect(&c)}
		}(),
		error: errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"complex64": logastic.Complex64(complex(3, 21))},
		expected: `{
			"complex64":"3+21i"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any complex64": logastic.Any(complex(3, 21))},
		expected: `{
			"any complex64":"3+21i"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect complex64": logastic.Reflect(complex(3, 21))},
		error: errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex128"),
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"error": logastic.Error(errors.New("something went wrong"))},
		expected: `{
			"error":"something went wrong"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any error": logastic.Any(errors.New("something went wrong"))},
		expected: `{
			"any error":"something went wrong"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex64 = complex(1, 23)
			return map[string]json.Marshaler{"complex64 pointer": logastic.Complex64p(&c)}
		}(),
		expected: `{
			"complex64 pointer":"1+23i"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil complex64 pointer": logastic.Complex64p(nil)},
		expected: `{
			"nil complex64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex64 = complex(1, 23)
			return map[string]json.Marshaler{"any complex64 pointer": logastic.Any(&c)}
		}(),
		expected: `{
			"any complex64 pointer":"1+23i"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var c complex64 = complex(1, 23)
			return map[string]json.Marshaler{"reflect complex64 pointer": logastic.Reflect(&c)}
		}(),
		error: errors.New("json: error calling MarshalJSON for type json.Marshaler: json: unsupported type: complex64"),
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect error": logastic.Reflect(errors.New("something went wrong"))},
		expected: `{
			"reflect error":{}
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil error": logastic.Error(nil)},
		expected: `{
			"nil error":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float32": logastic.Float32(4.2)},
		expected: `{
			"float32":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float32": logastic.Float32(0.123456789)},
		expected: `{
			"float32":0.123456789
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"zero float32": logastic.Float32(0)},
		expected: `{
			"zero float32":0
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any float32": logastic.Any(4.2)},
		expected: `{
			"any float32":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any zero float32": logastic.Any(0)},
		expected: `{
			"any zero float32":0
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect float32": logastic.Reflect(4.2)},
		expected: `{
			"reflect float32":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect zero float32": logastic.Reflect(0)},
		expected: `{
			"reflect zero float32":0
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 4.2
			return map[string]json.Marshaler{"float32 pointer": logastic.Float32p(&f)}
		}(),
		expected: `{
			"float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 0.123456789
			return map[string]json.Marshaler{"float32 pointer": logastic.Float32p(&f)}
		}(),
		expected: `{
			"float32 pointer":0.123456789
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float32 nil pointer": logastic.Float32p(nil)},
		expected: `{
			"float32 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 4.2
			return map[string]json.Marshaler{"any float32 pointer": logastic.Any(&f)}
		}(),
		expected: `{
			"any float32 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float32 = 4.2
			return map[string]json.Marshaler{"reflect float32 pointer": logastic.Reflect(&f)}
		}(),
		expected: `{
			"reflect float32 pointer":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect float32 nil pointer": logastic.Reflect(nil)},
		expected: `{
			"reflect float32 nil pointer":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float64": logastic.Float64(4.2)},
		expected: `{
			"float64":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float64": logastic.Float64(0.123456789)},
		expected: `{
			"float64":0.123456789
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"zero float64": logastic.Float64(0)},
		expected: `{
			"zero float64":0
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any float64": logastic.Any(4.2)},
		expected: `{
			"any float64":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any zero float64": logastic.Any(0)},
		expected: `{
			"any zero float64":0
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect float64": logastic.Reflect(4.2)},
		expected: `{
			"reflect float64":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect zero float64": logastic.Reflect(0)},
		expected: `{
			"reflect zero float64":0
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 4.2
			return map[string]json.Marshaler{"float64 pointer": logastic.Float64p(&f)}
		}(),
		expected: `{
			"float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 0.123456789
			return map[string]json.Marshaler{"float64 pointer": logastic.Float64p(&f)}
		}(),
		expected: `{
			"float64 pointer":0.123456789
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float64 nil pointer": logastic.Float64p(nil)},
		expected: `{
			"float64 nil pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 4.2
			return map[string]json.Marshaler{"any float64 pointer": logastic.Any(&f)}
		}(),
		expected: `{
			"any float64 pointer":4.2
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var f float64 = 4.2
			return map[string]json.Marshaler{"reflect float64 pointer": logastic.Reflect(&f)}
		}(),
		expected: `{
			"reflect float64 pointer":4.2
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect float64 nil pointer": logastic.Reflect(nil)},
		expected: `{
			"reflect float64 nil pointer":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"int": logastic.Int(42)},
		expected: `{
			"int":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any int": logastic.Any(42)},
		expected: `{
			"any int":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect int": logastic.Reflect(42)},
		expected: `{
			"reflect int":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int = 42
			return map[string]json.Marshaler{"int pointer": logastic.Intp(&i)}
		}(),
		expected: `{
			"int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int = 42
			return map[string]json.Marshaler{"any int pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any int pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int = 42
			return map[string]json.Marshaler{"reflect int pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect int pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"int16": logastic.Int16(42)},
		expected: `{
			"int16":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any int16": logastic.Any(42)},
		expected: `{
			"any int16":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect int16": logastic.Reflect(42)},
		expected: `{
			"reflect int16":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int16 = 42
			return map[string]json.Marshaler{"int16 pointer": logastic.Int16p(&i)}
		}(),
		expected: `{
			"int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int16 = 42
			return map[string]json.Marshaler{"any int16 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any int16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int16 = 42
			return map[string]json.Marshaler{"reflect int16 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect int16 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"int32": logastic.Int32(42)},
		expected: `{
			"int32":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any int32": logastic.Any(42)},
		expected: `{
			"any int32":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect int32": logastic.Reflect(42)},
		expected: `{
			"reflect int32":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int32 = 42
			return map[string]json.Marshaler{"int32 pointer": logastic.Int32p(&i)}
		}(),
		expected: `{
			"int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int32 = 42
			return map[string]json.Marshaler{"any int32 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any int32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int32 = 42
			return map[string]json.Marshaler{"reflect int32 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect int32 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"int64": logastic.Int64(42)},
		expected: `{
			"int64":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any int64": logastic.Any(42)},
		expected: `{
			"any int64":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect int64": logastic.Reflect(42)},
		expected: `{
			"reflect int64":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int64 = 42
			return map[string]json.Marshaler{"int64 pointer": logastic.Int64p(&i)}
		}(),
		expected: `{
			"int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int64 = 42
			return map[string]json.Marshaler{"any int64 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any int64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int64 = 42
			return map[string]json.Marshaler{"reflect int64 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect int64 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"int8": logastic.Int8(42)},
		expected: `{
			"int8":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any int8": logastic.Any(42)},
		expected: `{
			"any int8":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect int8": logastic.Reflect(42)},
		expected: `{
			"reflect int8":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int8 = 42
			return map[string]json.Marshaler{"int8 pointer": logastic.Int8p(&i)}
		}(),
		expected: `{
			"int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int8 = 42
			return map[string]json.Marshaler{"any int8 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any int8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i int8 = 42
			return map[string]json.Marshaler{"reflect int8 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect int8 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"runes": logastic.Runes([]rune("Hello, World!"))},
		expected: `{
			"runes":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"empty runes": logastic.Runes([]rune{})},
		expected: `{
			"empty runes":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil runes": logastic.Runes(nil)},
		expected: `{
			"nil runes":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"rune slice with zero rune": logastic.Runes([]rune{rune(0)})},
		expected: `{
			"rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any runes": logastic.Any([]rune("Hello, World!"))},
		expected: `{
			"any runes":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any empty runes": logastic.Any([]rune{})},
		expected: `{
			"any empty runes":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any rune slice with zero rune": logastic.Any([]rune{rune(0)})},
		expected: `{
			"any rune slice with zero rune":"\u0000"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect runes": logastic.Reflect([]rune("Hello, World!"))},
		expected: `{
			"reflect runes":[72,101,108,108,111,44,32,87,111,114,108,100,33]
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect empty runes": logastic.Reflect([]rune{})},
		expected: `{
			"reflect empty runes":[]
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect rune slice with zero rune": logastic.Reflect([]rune{rune(0)})},
		expected: `{
			"reflect rune slice with zero rune":[0]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune("Hello, World!")
			return map[string]json.Marshaler{"runes pointer": logastic.Runesp(&p)}
		}(),
		expected: `{
			"runes pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune{}
			return map[string]json.Marshaler{"empty runes pointer": logastic.Runesp(&p)}
		}(),
		expected: `{
			"empty runes pointer":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil runes pointer": logastic.Runesp(nil)},
		expected: `{
			"nil runes pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune("Hello, World!")
			return map[string]json.Marshaler{"any runes pointer": logastic.Any(&p)}
		}(),
		expected: `{
			"any runes pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune{}
			return map[string]json.Marshaler{"any empty runes pointer": logastic.Any(&p)}
		}(),
		expected: `{
			"any empty runes pointer":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune("Hello, World!")
			return map[string]json.Marshaler{"reflect runes pointer": logastic.Reflect(&p)}
		}(),
		expected: `{
			"reflect runes pointer":[72,101,108,108,111,44,32,87,111,114,108,100,33]
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := []rune{}
			return map[string]json.Marshaler{"reflect empty runes pointer": logastic.Reflect(&p)}
		}(),
		expected: `{
			"reflect empty runes pointer":[]
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"string": logastic.String("Hello, World!")},
		expected: `{
			"string":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"empty string": logastic.String("")},
		expected: `{
			"empty string":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"string with zero byte": logastic.String(string(byte(0)))},
		expected: `{
			"string with zero byte":"\u0000"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any string": logastic.Any("Hello, World!")},
		expected: `{
			"any string":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any empty string": logastic.Any("")},
		expected: `{
			"any empty string":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any string with zero byte": logastic.Any(string(byte(0)))},
		expected: `{
			"any string with zero byte":"\u0000"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect string": logastic.Reflect("Hello, World!")},
		expected: `{
			"reflect string":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect empty string": logastic.Reflect("")},
		expected: `{
			"reflect empty string":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect string with zero byte": logastic.Reflect(string(byte(0)))},
		expected: `{
			"reflect string with zero byte":"\u0000"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := "Hello, World!"
			return map[string]json.Marshaler{"string pointer": logastic.Stringp(&p)}
		}(),
		expected: `{
			"string pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := ""
			return map[string]json.Marshaler{"empty string pointer": logastic.Stringp(&p)}
		}(),
		expected: `{
			"empty string pointer":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil string pointer": logastic.Stringp(nil)},
		expected: `{
			"nil string pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := "Hello, World!"
			return map[string]json.Marshaler{"any string pointer": logastic.Any(&p)}
		}(),
		expected: `{
			"any string pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := ""
			return map[string]json.Marshaler{"any empty string pointer": logastic.Any(&p)}
		}(),
		expected: `{
			"any empty string pointer":""
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := "Hello, World!"
			return map[string]json.Marshaler{"reflect string pointer": logastic.Reflect(&p)}
		}(),
		expected: `{
			"reflect string pointer":"Hello, World!"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			p := ""
			return map[string]json.Marshaler{"reflect empty string pointer": logastic.Reflect(&p)}
		}(),
		expected: `{
			"reflect empty string pointer":""
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uint": logastic.Uint(42)},
		expected: `{
			"uint":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any uint": logastic.Any(42)},
		expected: `{
			"any uint":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uint": logastic.Reflect(42)},
		expected: `{
			"reflect uint":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint = 42
			return map[string]json.Marshaler{"uint pointer": logastic.Uintp(&i)}
		}(),
		expected: `{
			"uint pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil uint pointer": logastic.Uintp(nil)},
		expected: `{
			"nil uint pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint = 42
			return map[string]json.Marshaler{"any uint pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any uint pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint = 42
			return map[string]json.Marshaler{"reflect uint pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect uint pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uint16": logastic.Uint16(42)},
		expected: `{
			"uint16":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any uint16": logastic.Any(42)},
		expected: `{
			"any uint16":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uint16": logastic.Reflect(42)},
		expected: `{
			"reflect uint16":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint16 = 42
			return map[string]json.Marshaler{"uint16 pointer": logastic.Uint16p(&i)}
		}(),
		expected: `{
			"uint16 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uint16 pointer": logastic.Uint16p(nil)},
		expected: `{
			"uint16 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint16 = 42
			return map[string]json.Marshaler{"any uint16 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any uint16 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint16 = 42
			return map[string]json.Marshaler{"reflect uint16 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect uint16 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uint16 pointer": logastic.Reflect(nil)},
		expected: `{
			"reflect uint16 pointer":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uint32": logastic.Uint32(42)},
		expected: `{
			"uint32":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any uint32": logastic.Any(42)},
		expected: `{
			"any uint32":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uint32": logastic.Reflect(42)},
		expected: `{
			"reflect uint32":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint32 = 42
			return map[string]json.Marshaler{"uint32 pointer": logastic.Uint32p(&i)}
		}(),
		expected: `{
			"uint32 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil uint32 pointer": logastic.Uint32p(nil)},
		expected: `{
			"nil uint32 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint32 = 42
			return map[string]json.Marshaler{"any uint32 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any uint32 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint32 = 42
			return map[string]json.Marshaler{"reflect uint32 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect uint32 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uint64": logastic.Uint64(42)},
		expected: `{
			"uint64":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any uint64": logastic.Any(42)},
		expected: `{
			"any uint64":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uint64": logastic.Reflect(42)},
		expected: `{
			"reflect uint64":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint64 = 42
			return map[string]json.Marshaler{"uint64 pointer": logastic.Uint64p(&i)}
		}(),
		expected: `{
			"uint64 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil uint64 pointer": logastic.Uint64p(nil)},
		expected: `{
			"nil uint64 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint64 = 42
			return map[string]json.Marshaler{"any uint64 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any uint64 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint64 = 42
			return map[string]json.Marshaler{"reflect uint64 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect uint64 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uint8": logastic.Uint8(42)},
		expected: `{
			"uint8":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any uint8": logastic.Any(42)},
		expected: `{
			"any uint8":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uint8": logastic.Reflect(42)},
		expected: `{
			"reflect uint8":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint8 = 42
			return map[string]json.Marshaler{"uint8 pointer": logastic.Uint8p(&i)}
		}(),
		expected: `{
			"uint8 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil uint8 pointer": logastic.Uint8p(nil)},
		expected: `{
			"nil uint8 pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint8 = 42
			return map[string]json.Marshaler{"any uint8 pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any uint8 pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uint8 = 42
			return map[string]json.Marshaler{"reflect uint8 pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect uint8 pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"uintptr": logastic.Uintptr(42)},
		expected: `{
			"uintptr":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any uintptr": logastic.Any(42)},
		expected: `{
			"any uintptr":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect uintptr": logastic.Reflect(42)},
		expected: `{
			"reflect uintptr":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uintptr = 42
			return map[string]json.Marshaler{"uintptr pointer": logastic.Uintptrp(&i)}
		}(),
		expected: `{
			"uintptr pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil uintptr pointer": logastic.Uintptrp(nil)},
		expected: `{
			"nil uintptr pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uintptr = 42
			return map[string]json.Marshaler{"any uintptr pointer": logastic.Any(&i)}
		}(),
		expected: `{
			"any uintptr pointer":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			var i uintptr = 42
			return map[string]json.Marshaler{"reflect uintptr pointer": logastic.Reflect(&i)}
		}(),
		expected: `{
			"reflect uintptr pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"time": logastic.Time(time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC))},
		expected: `{
			"time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any time": logastic.Any(time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC))},
		expected: `{
			"any time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect time": logastic.Reflect(time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC))},
		expected: `{
			"reflect time":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return map[string]json.Marshaler{"time pointer": logastic.Timep(&t)}
		}(),
		expected: `{
			"time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil time pointer": logastic.Timep(nil)},
		expected: `{
			"nil time pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return map[string]json.Marshaler{"any time pointer": logastic.Any(&t)}
		}(),
		expected: `{
			"any time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			t := time.Date(1970, time.January, 1, 0, 0, 0, 42, time.UTC)
			return map[string]json.Marshaler{"reflect time pointer": logastic.Reflect(&t)}
		}(),
		expected: `{
			"reflect time pointer":"1970-01-01T00:00:00.000000042Z"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"duration": logastic.Duration(42)},
		expected: `{
			"duration":"42ns"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any duration": logastic.Any(42 * time.Nanosecond)},
		expected: `{
			"any duration":"42ns"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect duration": logastic.Reflect(42 * time.Nanosecond)},
		expected: `{
			"reflect duration":42
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			d := 42 * time.Nanosecond
			return map[string]json.Marshaler{"duration pointer": logastic.Durationp(&d)}
		}(),
		expected: `{
			"duration pointer":"42ns"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"nil duration pointer": logastic.Durationp(nil)},
		expected: `{
			"nil duration pointer":null
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			d := 42 * time.Nanosecond
			return map[string]json.Marshaler{"any duration pointer": logastic.Any(&d)}
		}(),
		expected: `{
			"any duration pointer":"42ns"
		}`,
	},
	{
		line: line(),
		input: func() map[string]json.Marshaler {
			d := 42 * time.Nanosecond
			return map[string]json.Marshaler{"reflect duration pointer": logastic.Reflect(&d)}
		}(),
		expected: `{
			"reflect duration pointer":42
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any struct": logastic.Any(Struct{Name: "John Doe", Age: 42})},
		expected: `{
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
			return map[string]json.Marshaler{"any struct pointer": logastic.Any(&s)}
		}(),
		expected: `{
			"any struct pointer": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"struct reflect": logastic.Reflect(Struct{Name: "John Doe", Age: 42})},
		expected: `{
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
			return map[string]json.Marshaler{"struct reflect pointer": logastic.Reflect(&s)}
		}(),
		expected: `{
			"struct reflect pointer": {
				"Name":"John Doe",
				"Age":42
			}
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"any nil": logastic.Any(nil)},
		expected: `{
			"any nil":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"reflect nil": logastic.Reflect(nil)},
		expected: `{
			"reflect nil":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"raw json": logastic.Raw([]byte(`{"foo":"bar"}`))},
		expected: `{
			"raw json":{"foo":"bar"}
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"raw malformed json object": logastic.Raw([]byte(`xyz{"foo":"bar"}`))},
		error: errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character 'x' looking for beginning of value"),
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"raw malformed json key/value": logastic.Raw([]byte(`{"foo":"bar""}`))},
		error: errors.New(`json: error calling MarshalJSON for type json.Marshaler: invalid character '"' after object key:value pair`),
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"raw json with unescaped null byte": logastic.Raw(append([]byte(`{"foo":"`), append([]byte{0}, []byte(`xyz"}`)...)...))},
		error: errors.New("json: error calling MarshalJSON for type json.Marshaler: invalid character '\\x00' in string literal"),
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"raw nil": logastic.Raw(nil)},
		expected: `{
			"raw nil":null
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

			p, err := json.Marshal(tc.input)

			if !equalastic.ErrorEqual(err, tc.error) {
				t.Fatalf("marshal error expected: %s, recieved: %s %s", tc.error, err, linkToExample)
			}

			if err == nil {
				ja := jsonassert.New(testprinter{t: t, link: linkToExample})
				ja.Assertf(string(p), tc.expected)
			}
		})
	}
}
