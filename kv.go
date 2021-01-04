package logastic

import (
	"encoding"
	"encoding/json"
	"time"
)

// kvp is a key-value pair.
type kvp struct {
	K encoding.TextMarshaler
	V json.Marshaler
}

func (kv kvp) MarshalText() (text []byte, err error) { return kv.K.MarshalText() }
func (kv kvp) MarshalJSON() ([]byte, error)          { return kv.V.MarshalJSON() }

func StringBool(k string, v bool) kvp                { return kvp{K: String(k), V: Bool(v)} }
func StringBoolp(k string, v *bool) kvp              { return kvp{K: String(k), V: Boolp(v)} }
func StringBytes(k string, v []byte) kvp             { return kvp{K: String(k), V: Bytes(v)} }
func StringBytesp(k string, v *[]byte) kvp           { return kvp{K: String(k), V: Bytesp(v)} }
func StringComplex128(k string, v complex128) kvp    { return kvp{K: String(k), V: Complex128(v)} }
func StringComplex128p(k string, v *complex128) kvp  { return kvp{K: String(k), V: Complex128p(v)} }
func StringComplex64(k string, v complex64) kvp      { return kvp{K: String(k), V: Complex64(v)} }
func StringComplex64p(k string, v *complex64) kvp    { return kvp{K: String(k), V: Complex64p(v)} }
func StringError(k string, v error) kvp              { return kvp{K: String(k), V: Error(v)} }
func StringFloat32(k string, v float32) kvp          { return kvp{K: String(k), V: Float32(v)} }
func StringFloat32p(k string, v *float32) kvp        { return kvp{K: String(k), V: Float32p(v)} }
func StringFloat64(k string, v float64) kvp          { return kvp{K: String(k), V: Float64(v)} }
func StringFloat64p(k string, v *float64) kvp        { return kvp{K: String(k), V: Float64p(v)} }
func StringInt(k string, v int) kvp                  { return kvp{K: String(k), V: Int(v)} }
func StringIntp(k string, v *int) kvp                { return kvp{K: String(k), V: Intp(v)} }
func StringInt16(k string, v int16) kvp              { return kvp{K: String(k), V: Int16(v)} }
func StringInt16p(k string, v *int16) kvp            { return kvp{K: String(k), V: Int16p(v)} }
func StringInt32(k string, v int32) kvp              { return kvp{K: String(k), V: Int32(v)} }
func StringInt32p(k string, v *int32) kvp            { return kvp{K: String(k), V: Int32p(v)} }
func StringInt64(k string, v int64) kvp              { return kvp{K: String(k), V: Int64(v)} }
func StringInt64p(k string, v *int64) kvp            { return kvp{K: String(k), V: Int64p(v)} }
func StringInt8(k string, v int8) kvp                { return kvp{K: String(k), V: Int8(v)} }
func StringInt8p(k string, v *int8) kvp              { return kvp{K: String(k), V: Int8p(v)} }
func StringRunes(k string, v []rune) kvp             { return kvp{K: String(k), V: Runes(v)} }
func StringRunesp(k string, v *[]rune) kvp           { return kvp{K: String(k), V: Runesp(v)} }
func StringString(k string, v string) kvp            { return kvp{K: String(k), V: String(v)} }
func StringStringp(k string, v *string) kvp          { return kvp{K: String(k), V: Stringp(v)} }
func StringUint(k string, v uint) kvp                { return kvp{K: String(k), V: Uint(v)} }
func StringUintp(k string, v *uint) kvp              { return kvp{K: String(k), V: Uintp(v)} }
func StringUint16(k string, v uint16) kvp            { return kvp{K: String(k), V: Uint16(v)} }
func StringUint16p(k string, v *uint16) kvp          { return kvp{K: String(k), V: Uint16p(v)} }
func StringUint32(k string, v uint32) kvp            { return kvp{K: String(k), V: Uint32(v)} }
func StringUint32p(k string, v *uint32) kvp          { return kvp{K: String(k), V: Uint32p(v)} }
func StringUint64(k string, v uint64) kvp            { return kvp{K: String(k), V: Uint64(v)} }
func StringUint64p(k string, v *uint64) kvp          { return kvp{K: String(k), V: Uint64p(v)} }
func StringUint8(k string, v uint8) kvp              { return kvp{K: String(k), V: Uint8(v)} }
func StringUint8p(k string, v *uint8) kvp            { return kvp{K: String(k), V: Uint8p(v)} }
func StringUintptr(k string, v uintptr) kvp          { return kvp{K: String(k), V: Uintptr(v)} }
func StringUintptrp(k string, v *uintptr) kvp        { return kvp{K: String(k), V: Uintptrp(v)} }
func StringDuration(k string, v time.Duration) kvp   { return kvp{K: String(k), V: Duration(v)} }
func StringDurationp(k string, v *time.Duration) kvp { return kvp{K: String(k), V: Durationp(v)} }
func StringTime(k string, v time.Time) kvp           { return kvp{K: String(k), V: Time(v)} }
func StringTimep(k string, v *time.Time) kvp         { return kvp{K: String(k), V: Timep(v)} }
func StringRaw(k string, v []byte) kvp               { return kvp{K: String(k), V: Raw(v)} }
func StringAny(k string, v interface{}) kvp          { return kvp{K: String(k), V: Any(v)} }
func StringReflect(k string, v interface{}) kvp      { return kvp{K: String(k), V: Reflect(v)} }
