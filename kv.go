package log64

import (
	"encoding"
	"encoding/json"
	"time"

	"github.com/danil/log64/marshal8"
)

// kvp is a key-value pair.
type kvp struct {
	K encoding.TextMarshaler
	V json.Marshaler
}

func (kv kvp) MarshalText() (text []byte, err error) { return kv.K.MarshalText() }
func (kv kvp) MarshalJSON() ([]byte, error)          { return kv.V.MarshalJSON() }

func StringBool(k string, v bool) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Bool(v)}
}

func StringBoolp(k string, v *bool) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Boolp(v)}
}

func StringBytes(k string, v []byte) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Bytes(v)}
}

func StringBytesp(k string, v *[]byte) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Bytesp(v)}
}

func StringComplex128(k string, v complex128) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Complex128(v)}
}

func StringComplex128p(k string, v *complex128) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Complex128p(v)}
}

func StringComplex64(k string, v complex64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Complex64(v)}
}

func StringComplex64p(k string, v *complex64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Complex64p(v)}
}

func StringError(k string, v error) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Error(v)}
}

func StringFloat32(k string, v float32) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Float32(v)}
}

func StringFloat32p(k string, v *float32) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Float32p(v)}
}

func StringFloat64(k string, v float64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Float64(v)}
}

func StringFloat64p(k string, v *float64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Float64p(v)}
}

func StringInt(k string, v int) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int(v)}
}

func StringIntp(k string, v *int) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Intp(v)}
}

func StringInt16(k string, v int16) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int16(v)}
}

func StringInt16p(k string, v *int16) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int16p(v)}
}

func StringInt32(k string, v int32) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int32(v)}
}

func StringInt32p(k string, v *int32) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int32p(v)}
}

func StringInt64(k string, v int64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int64(v)}
}

func StringInt64p(k string, v *int64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int64p(v)}
}

func StringInt8(k string, v int8) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int8(v)}
}

func StringInt8p(k string, v *int8) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Int8p(v)}
}

func StringRunes(k string, v []rune) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Runes(v)}
}

func StringRunesp(k string, v *[]rune) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Runesp(v)}
}

func String(a ...string) kvp {
	if len(a) == 0 {
		kv := marshal8.String("")
		return kvp{K: kv, V: kv}
	}
	if len(a) == 1 {
		return kvp{K: marshal8.String(a[0]), V: marshal8.String("")}
	}
	return kvp{K: marshal8.String(a[0]), V: marshal8.String(a[1])}
}

func StringStringp(k string, v *string) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Stringp(v)}
}

func StringUint(k string, v uint) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint(v)}
}

func StringUintp(k string, v *uint) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uintp(v)}
}

func StringUint16(k string, v uint16) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint16(v)}
}

func StringUint16p(k string, v *uint16) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint16p(v)}
}

func StringUint32(k string, v uint32) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint32(v)}
}

func StringUint32p(k string, v *uint32) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint32p(v)}
}

func StringUint64(k string, v uint64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint64(v)}
}

func StringUint64p(k string, v *uint64) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint64p(v)}
}

func StringUint8(k string, v uint8) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint8(v)}
}

func StringUint8p(k string, v *uint8) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uint8p(v)}
}

func StringUintptr(k string, v uintptr) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uintptr(v)}
}

func StringUintptrp(k string, v *uintptr) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Uintptrp(v)}
}

func StringDuration(k string, v time.Duration) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Duration(v)}
}

func StringDurationp(k string, v *time.Duration) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Durationp(v)}
}

func StringTime(k string, v time.Time) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Time(v)}
}

func StringTimep(k string, v *time.Time) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Timep(v)}
}

func StringRaw(k string, v []byte) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Raw(v)}
}

func StringAny(k string, v interface{}) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Any(v)}
}

func StringReflect(k string, v interface{}) kvp {
	return kvp{K: marshal8.String(k), V: marshal8.Reflect(v)}
}

func TextBool(k encoding.TextMarshaler, v bool) kvp {
	return kvp{K: k, V: marshal8.Bool(v)}
}

func TextBoolp(k encoding.TextMarshaler, v *bool) kvp {
	return kvp{K: k, V: marshal8.Boolp(v)}
}

func TextBytes(k encoding.TextMarshaler, v []byte) kvp {
	return kvp{K: k, V: marshal8.Bytes(v)}
}

func TextBytesp(k encoding.TextMarshaler, v *[]byte) kvp {
	return kvp{K: k, V: marshal8.Bytesp(v)}
}

func TextComplex128(k encoding.TextMarshaler, v complex128) kvp {
	return kvp{K: k, V: marshal8.Complex128(v)}
}

func TextComplex128p(k encoding.TextMarshaler, v *complex128) kvp {
	return kvp{K: k, V: marshal8.Complex128p(v)}
}

func TextComplex64(k encoding.TextMarshaler, v complex64) kvp {
	return kvp{K: k, V: marshal8.Complex64(v)}
}

func TextComplex64p(k encoding.TextMarshaler, v *complex64) kvp {
	return kvp{K: k, V: marshal8.Complex64p(v)}
}

func TextError(k encoding.TextMarshaler, v error) kvp {
	return kvp{K: k, V: marshal8.Error(v)}
}

func TextFloat32(k encoding.TextMarshaler, v float32) kvp {
	return kvp{K: k, V: marshal8.Float32(v)}
}

func TextFloat32p(k encoding.TextMarshaler, v *float32) kvp {
	return kvp{K: k, V: marshal8.Float32p(v)}
}

func TextFloat64(k encoding.TextMarshaler, v float64) kvp {
	return kvp{K: k, V: marshal8.Float64(v)}
}

func TextFloat64p(k encoding.TextMarshaler, v *float64) kvp {
	return kvp{K: k, V: marshal8.Float64p(v)}
}

func TextInt(k encoding.TextMarshaler, v int) kvp {
	return kvp{K: k, V: marshal8.Int(v)}
}

func TextIntp(k encoding.TextMarshaler, v *int) kvp {
	return kvp{K: k, V: marshal8.Intp(v)}
}

func TextInt16(k encoding.TextMarshaler, v int16) kvp {
	return kvp{K: k, V: marshal8.Int16(v)}
}

func TextInt16p(k encoding.TextMarshaler, v *int16) kvp {
	return kvp{K: k, V: marshal8.Int16p(v)}
}

func TextInt32(k encoding.TextMarshaler, v int32) kvp {
	return kvp{K: k, V: marshal8.Int32(v)}
}

func TextInt32p(k encoding.TextMarshaler, v *int32) kvp {
	return kvp{K: k, V: marshal8.Int32p(v)}
}

func TextInt64(k encoding.TextMarshaler, v int64) kvp {
	return kvp{K: k, V: marshal8.Int64(v)}
}

func TextInt64p(k encoding.TextMarshaler, v *int64) kvp {
	return kvp{K: k, V: marshal8.Int64p(v)}
}

func TextInt8(k encoding.TextMarshaler, v int8) kvp {
	return kvp{K: k, V: marshal8.Int8(v)}
}

func TextInt8p(k encoding.TextMarshaler, v *int8) kvp {
	return kvp{K: k, V: marshal8.Int8p(v)}
}

func TextRunes(k encoding.TextMarshaler, v []rune) kvp {
	return kvp{K: k, V: marshal8.Runes(v)}
}

func TextRunesp(k encoding.TextMarshaler, v *[]rune) kvp {
	return kvp{K: k, V: marshal8.Runesp(v)}
}

func Text(a ...encoding.TextMarshaler) kvp {
	if len(a) == 0 {
		kv := marshal8.String("")
		return kvp{K: kv, V: kv}
	}
	if len(a) == 1 {
		return kvp{K: a[0], V: marshal8.String("")}
	}
	return kvp{K: a[0], V: marshal8.Text(a[1])}
}

func TextString(k encoding.TextMarshaler, v string) kvp {
	return kvp{K: k, V: marshal8.String(v)}
}

func TextStringp(k encoding.TextMarshaler, v *string) kvp {
	return kvp{K: k, V: marshal8.Stringp(v)}
}

func TextUint(k encoding.TextMarshaler, v uint) kvp {
	return kvp{K: k, V: marshal8.Uint(v)}
}

func TextUintp(k encoding.TextMarshaler, v *uint) kvp {
	return kvp{K: k, V: marshal8.Uintp(v)}
}

func TextUint16(k encoding.TextMarshaler, v uint16) kvp {
	return kvp{K: k, V: marshal8.Uint16(v)}
}

func TextUint16p(k encoding.TextMarshaler, v *uint16) kvp {
	return kvp{K: k, V: marshal8.Uint16p(v)}
}

func TextUint32(k encoding.TextMarshaler, v uint32) kvp {
	return kvp{K: k, V: marshal8.Uint32(v)}
}

func TextUint32p(k encoding.TextMarshaler, v *uint32) kvp {
	return kvp{K: k, V: marshal8.Uint32p(v)}
}

func TextUint64(k encoding.TextMarshaler, v uint64) kvp {
	return kvp{K: k, V: marshal8.Uint64(v)}
}

func TextUint64p(k encoding.TextMarshaler, v *uint64) kvp {
	return kvp{K: k, V: marshal8.Uint64p(v)}
}

func TextUint8(k encoding.TextMarshaler, v uint8) kvp {
	return kvp{K: k, V: marshal8.Uint8(v)}
}

func TextUint8p(k encoding.TextMarshaler, v *uint8) kvp {
	return kvp{K: k, V: marshal8.Uint8p(v)}
}

func TextUintptr(k encoding.TextMarshaler, v uintptr) kvp {
	return kvp{K: k, V: marshal8.Uintptr(v)}
}

func TextUintptrp(k encoding.TextMarshaler, v *uintptr) kvp {
	return kvp{K: k, V: marshal8.Uintptrp(v)}
}

func TextDuration(k encoding.TextMarshaler, v time.Duration) kvp {
	return kvp{K: k, V: marshal8.Duration(v)}
}

func TextDurationp(k encoding.TextMarshaler, v *time.Duration) kvp {
	return kvp{K: k, V: marshal8.Durationp(v)}
}

func TextTime(k encoding.TextMarshaler, v time.Time) kvp {
	return kvp{K: k, V: marshal8.Time(v)}
}

func TextTimep(k encoding.TextMarshaler, v *time.Time) kvp {
	return kvp{K: k, V: marshal8.Timep(v)}
}

func TextRaw(k encoding.TextMarshaler, v []byte) kvp {
	return kvp{K: k, V: marshal8.Raw(v)}
}

func TextAny(k encoding.TextMarshaler, v interface{}) kvp {
	return kvp{K: k, V: marshal8.Any(v)}
}

func TextReflect(k encoding.TextMarshaler, v interface{}) kvp {
	return kvp{K: k, V: marshal8.Reflect(v)}
}
