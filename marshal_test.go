package logastic_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"testing"

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
		line:  line(),
		input: map[string]json.Marshaler{"byte": logastic.Byte([]byte("Hello, World!")...)},
		expected: `{
			"byte":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"empty byte": logastic.Byte()},
		expected: `{
			"empty byte":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"byte slice with zero byte": logastic.Byte(byte(0))},
		expected: `{
			"byte slice with zero byte":"\u0000"
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
		input: map[string]json.Marshaler{"empty bytes pointer": logastic.Bytesp(nil)},
		expected: `{
			"empty bytes pointer":null
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
		input: map[string]json.Marshaler{"complex64": logastic.Complex64(complex(3, 21))},
		expected: `{
			"complex64":"3+21i"
		}`,
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
		input: map[string]json.Marshaler{"nil error": logastic.Error(nil)},
		expected: `{
			"nil error":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"float32": logastic.Float32(1.2)},
		expected: `{
			"float32":1.2
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
		line:  line(),
		input: map[string]json.Marshaler{"float64": logastic.Float64(1.2)},
		expected: `{
			"float64":1.2
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
		line:  line(),
		input: map[string]json.Marshaler{"int": logastic.Int(42)},
		expected: `{
			"int":42
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
		line:  line(),
		input: map[string]json.Marshaler{"int16": logastic.Int16(42)},
		expected: `{
			"int16":42
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
		line:  line(),
		input: map[string]json.Marshaler{"int32": logastic.Int32(42)},
		expected: `{
			"int32":42
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
		line:  line(),
		input: map[string]json.Marshaler{"int64": logastic.Int64(42)},
		expected: `{
			"int64":42
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
		line:  line(),
		input: map[string]json.Marshaler{"int8": logastic.Int8(42)},
		expected: `{
			"int8":42
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
		line:  line(),
		input: map[string]json.Marshaler{"rune": logastic.Rune([]rune("Hello, World!")...)},
		expected: `{
			"rune":"Hello, World!"
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"empty rune": logastic.Rune()},
		expected: `{
			"empty rune":null
		}`,
	},
	{
		line:  line(),
		input: map[string]json.Marshaler{"zero rune": logastic.Rune(rune(0))},
		expected: `{
			"zero rune":"\u0000"
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
		input: map[string]json.Marshaler{"uint": logastic.Uint(42)},
		expected: `{
			"uint":42
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
		input: map[string]json.Marshaler{"uint16": logastic.Uint16(42)},
		expected: `{
			"uint16":42
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
		input: map[string]json.Marshaler{"uint32": logastic.Uint32(42)},
		expected: `{
			"uint32":42
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
		input: map[string]json.Marshaler{"uint64": logastic.Uint64(42)},
		expected: `{
			"uint64":42
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
		input: map[string]json.Marshaler{"uint8": logastic.Uint8(42)},
		expected: `{
			"uint8":42
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
		input: map[string]json.Marshaler{"uintptr": logastic.Uintptr(42)},
		expected: `{
			"uintptr":42
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
}

func TestMarshal(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, tc := range MarshalTestCases {
		tc := tc
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, tc.line)

			b, err := json.Marshal(tc.input)

			if !equalastic.ErrorEqual(err, tc.error) {
				t.Fatalf("marshal error expected: %s, recieved: %s %s", tc.error, err, linkToExample)
			}

			if err == nil {
				ja := jsonassert.New(testprinter{t: t, link: linkToExample})
				ja.Assertf(string(b), tc.expected)
			}
		})
	}
}
