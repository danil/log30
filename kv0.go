// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log0

import (
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// kvjt is a key-value pair implements json/text marshaler.
type kvjt struct {
	K encoding.TextMarshaler
	V json.Marshaler
}

func (kv kvjt) MarshalText() (text []byte, err error) { return kv.K.MarshalText() }
func (kv kvjt) MarshalJSON() ([]byte, error)          { return kv.V.MarshalJSON() }

func StringBool(k string, v bool) kvjt {
	return kvjt{K: String(k), V: Bool(v)}
}

func StringBoolp(k string, v *bool) kvjt {
	return kvjt{K: String(k), V: Boolp(v)}
}

func StringBytes(k string, v []byte) kvjt {
	return kvjt{K: String(k), V: Bytes(v)}
}

func StringBytesp(k string, v *[]byte) kvjt {
	return kvjt{K: String(k), V: Bytesp(v)}
}

func StringComplex128(k string, v complex128) kvjt {
	return kvjt{K: String(k), V: Complex128(v)}
}

func StringComplex128p(k string, v *complex128) kvjt {
	return kvjt{K: String(k), V: Complex128p(v)}
}

func StringComplex64(k string, v complex64) kvjt {
	return kvjt{K: String(k), V: Complex64(v)}
}

func StringComplex64p(k string, v *complex64) kvjt {
	return kvjt{K: String(k), V: Complex64p(v)}
}

func StringError(k string, v error) kvjt {
	return kvjt{K: String(k), V: Error(v)}
}

func StringFloat32(k string, v float32) kvjt {
	return kvjt{K: String(k), V: Float32(v)}
}

func StringFloat32p(k string, v *float32) kvjt {
	return kvjt{K: String(k), V: Float32p(v)}
}

func StringFloat64(k string, v float64) kvjt {
	return kvjt{K: String(k), V: Float64(v)}
}

func StringFloat64p(k string, v *float64) kvjt {
	return kvjt{K: String(k), V: Float64p(v)}
}

func StringInt(k string, v int) kvjt {
	return kvjt{K: String(k), V: Int(v)}
}

func StringIntp(k string, v *int) kvjt {
	return kvjt{K: String(k), V: Intp(v)}
}

func StringInt16(k string, v int16) kvjt {
	return kvjt{K: String(k), V: Int16(v)}
}

func StringInt16p(k string, v *int16) kvjt {
	return kvjt{K: String(k), V: Int16p(v)}
}

func StringInt32(k string, v int32) kvjt {
	return kvjt{K: String(k), V: Int32(v)}
}

func StringInt32p(k string, v *int32) kvjt {
	return kvjt{K: String(k), V: Int32p(v)}
}

func StringInt64(k string, v int64) kvjt {
	return kvjt{K: String(k), V: Int64(v)}
}

func StringInt64p(k string, v *int64) kvjt {
	return kvjt{K: String(k), V: Int64p(v)}
}

func StringInt8(k string, v int8) kvjt {
	return kvjt{K: String(k), V: Int8(v)}
}

func StringInt8p(k string, v *int8) kvjt {
	return kvjt{K: String(k), V: Int8p(v)}
}

func StringRunes(k string, v []rune) kvjt {
	return kvjt{K: String(k), V: Runes(v)}
}

func StringRunesp(k string, v *[]rune) kvjt {
	return kvjt{K: String(k), V: Runesp(v)}
}

func Strings(a ...string) kvjt {
	if len(a) == 0 {
		kv := String("")
		return kvjt{K: kv, V: kv}
	}
	if len(a) == 1 {
		return kvjt{K: String(a[0]), V: String("")}
	}
	return kvjt{K: String(a[0]), V: String(a[1])}
}

func StringStringp(k string, v *string) kvjt {
	return kvjt{K: String(k), V: Stringp(v)}
}

func StringUint(k string, v uint) kvjt {
	return kvjt{K: String(k), V: Uint(v)}
}

func StringUintp(k string, v *uint) kvjt {
	return kvjt{K: String(k), V: Uintp(v)}
}

func StringUint16(k string, v uint16) kvjt {
	return kvjt{K: String(k), V: Uint16(v)}
}

func StringUint16p(k string, v *uint16) kvjt {
	return kvjt{K: String(k), V: Uint16p(v)}
}

func StringUint32(k string, v uint32) kvjt {
	return kvjt{K: String(k), V: Uint32(v)}
}

func StringUint32p(k string, v *uint32) kvjt {
	return kvjt{K: String(k), V: Uint32p(v)}
}

func StringUint64(k string, v uint64) kvjt {
	return kvjt{K: String(k), V: Uint64(v)}
}

func StringUint64p(k string, v *uint64) kvjt {
	return kvjt{K: String(k), V: Uint64p(v)}
}

func StringUint8(k string, v uint8) kvjt {
	return kvjt{K: String(k), V: Uint8(v)}
}

func StringUint8p(k string, v *uint8) kvjt {
	return kvjt{K: String(k), V: Uint8p(v)}
}

func StringUintptr(k string, v uintptr) kvjt {
	return kvjt{K: String(k), V: Uintptr(v)}
}

func StringUintptrp(k string, v *uintptr) kvjt {
	return kvjt{K: String(k), V: Uintptrp(v)}
}

func StringDuration(k string, v time.Duration) kvjt {
	return kvjt{K: String(k), V: Duration(v)}
}

func StringDurationp(k string, v *time.Duration) kvjt {
	return kvjt{K: String(k), V: Durationp(v)}
}

func StringTime(k string, v time.Time) kvjt {
	return kvjt{K: String(k), V: Time(v)}
}

func StringTimep(k string, v *time.Time) kvjt {
	return kvjt{K: String(k), V: Timep(v)}
}

func StringFunc(k string, v func() KV) kvjt {
	return kvjt{K: String(k), V: Func(v)}
}

func StringRaw(k string, v []byte) kvjt {
	return kvjt{K: String(k), V: Raw(v)}
}

func StringAny(k string, v interface{}) kvjt {
	return kvjt{K: String(k), V: Any(v)}
}

func StringReflect(k string, v interface{}) kvjt {
	return kvjt{K: String(k), V: Reflect(v)}
}

func TextBool(k encoding.TextMarshaler, v bool) kvjt {
	return kvjt{K: k, V: Bool(v)}
}

func TextBoolp(k encoding.TextMarshaler, v *bool) kvjt {
	return kvjt{K: k, V: Boolp(v)}
}

func TextBytes(k encoding.TextMarshaler, v []byte) kvjt {
	return kvjt{K: k, V: Bytes(v)}
}

func TextBytesp(k encoding.TextMarshaler, v *[]byte) kvjt {
	return kvjt{K: k, V: Bytesp(v)}
}

func TextComplex128(k encoding.TextMarshaler, v complex128) kvjt {
	return kvjt{K: k, V: Complex128(v)}
}

func TextComplex128p(k encoding.TextMarshaler, v *complex128) kvjt {
	return kvjt{K: k, V: Complex128p(v)}
}

func TextComplex64(k encoding.TextMarshaler, v complex64) kvjt {
	return kvjt{K: k, V: Complex64(v)}
}

func TextComplex64p(k encoding.TextMarshaler, v *complex64) kvjt {
	return kvjt{K: k, V: Complex64p(v)}
}

func TextError(k encoding.TextMarshaler, v error) kvjt {
	return kvjt{K: k, V: Error(v)}
}

func TextFloat32(k encoding.TextMarshaler, v float32) kvjt {
	return kvjt{K: k, V: Float32(v)}
}

func TextFloat32p(k encoding.TextMarshaler, v *float32) kvjt {
	return kvjt{K: k, V: Float32p(v)}
}

func TextFloat64(k encoding.TextMarshaler, v float64) kvjt {
	return kvjt{K: k, V: Float64(v)}
}

func TextFloat64p(k encoding.TextMarshaler, v *float64) kvjt {
	return kvjt{K: k, V: Float64p(v)}
}

func TextInt(k encoding.TextMarshaler, v int) kvjt {
	return kvjt{K: k, V: Int(v)}
}

func TextIntp(k encoding.TextMarshaler, v *int) kvjt {
	return kvjt{K: k, V: Intp(v)}
}

func TextInt16(k encoding.TextMarshaler, v int16) kvjt {
	return kvjt{K: k, V: Int16(v)}
}

func TextInt16p(k encoding.TextMarshaler, v *int16) kvjt {
	return kvjt{K: k, V: Int16p(v)}
}

func TextInt32(k encoding.TextMarshaler, v int32) kvjt {
	return kvjt{K: k, V: Int32(v)}
}

func TextInt32p(k encoding.TextMarshaler, v *int32) kvjt {
	return kvjt{K: k, V: Int32p(v)}
}

func TextInt64(k encoding.TextMarshaler, v int64) kvjt {
	return kvjt{K: k, V: Int64(v)}
}

func TextInt64p(k encoding.TextMarshaler, v *int64) kvjt {
	return kvjt{K: k, V: Int64p(v)}
}

func TextInt8(k encoding.TextMarshaler, v int8) kvjt {
	return kvjt{K: k, V: Int8(v)}
}

func TextInt8p(k encoding.TextMarshaler, v *int8) kvjt {
	return kvjt{K: k, V: Int8p(v)}
}

func TextRunes(k encoding.TextMarshaler, v []rune) kvjt {
	return kvjt{K: k, V: Runes(v)}
}

func TextRunesp(k encoding.TextMarshaler, v *[]rune) kvjt {
	return kvjt{K: k, V: Runesp(v)}
}

func Texts(a ...encoding.TextMarshaler) kvjt {
	if len(a) == 0 {
		kv := String("")
		return kvjt{K: kv, V: kv}
	}
	if len(a) == 1 {
		return kvjt{K: a[0], V: String("")}
	}
	return kvjt{K: a[0], V: Text(a[1])}
}

func TextString(k encoding.TextMarshaler, v string) kvjt {
	return kvjt{K: k, V: String(v)}
}

func TextStringp(k encoding.TextMarshaler, v *string) kvjt {
	return kvjt{K: k, V: Stringp(v)}
}

func TextUint(k encoding.TextMarshaler, v uint) kvjt {
	return kvjt{K: k, V: Uint(v)}
}

func TextUintp(k encoding.TextMarshaler, v *uint) kvjt {
	return kvjt{K: k, V: Uintp(v)}
}

func TextUint16(k encoding.TextMarshaler, v uint16) kvjt {
	return kvjt{K: k, V: Uint16(v)}
}

func TextUint16p(k encoding.TextMarshaler, v *uint16) kvjt {
	return kvjt{K: k, V: Uint16p(v)}
}

func TextUint32(k encoding.TextMarshaler, v uint32) kvjt {
	return kvjt{K: k, V: Uint32(v)}
}

func TextUint32p(k encoding.TextMarshaler, v *uint32) kvjt {
	return kvjt{K: k, V: Uint32p(v)}
}

func TextUint64(k encoding.TextMarshaler, v uint64) kvjt {
	return kvjt{K: k, V: Uint64(v)}
}

func TextUint64p(k encoding.TextMarshaler, v *uint64) kvjt {
	return kvjt{K: k, V: Uint64p(v)}
}

func TextUint8(k encoding.TextMarshaler, v uint8) kvjt {
	return kvjt{K: k, V: Uint8(v)}
}

func TextUint8p(k encoding.TextMarshaler, v *uint8) kvjt {
	return kvjt{K: k, V: Uint8p(v)}
}

func TextUintptr(k encoding.TextMarshaler, v uintptr) kvjt {
	return kvjt{K: k, V: Uintptr(v)}
}

func TextUintptrp(k encoding.TextMarshaler, v *uintptr) kvjt {
	return kvjt{K: k, V: Uintptrp(v)}
}

func TextDuration(k encoding.TextMarshaler, v time.Duration) kvjt {
	return kvjt{K: k, V: Duration(v)}
}

func TextDurationp(k encoding.TextMarshaler, v *time.Duration) kvjt {
	return kvjt{K: k, V: Durationp(v)}
}

func TextTime(k encoding.TextMarshaler, v time.Time) kvjt {
	return kvjt{K: k, V: Time(v)}
}

func TextTimep(k encoding.TextMarshaler, v *time.Time) kvjt {
	return kvjt{K: k, V: Timep(v)}
}

func TextFunc(k encoding.TextMarshaler, v func() json.Marshaler) kvjt {
	return kvjt{K: k, V: v()}
}

func TextRaw(k encoding.TextMarshaler, v []byte) kvjt {
	return kvjt{K: k, V: Raw(v)}
}

func TextAny(k encoding.TextMarshaler, v interface{}) kvjt {
	return kvjt{K: k, V: Any(v)}
}

func TextReflect(k encoding.TextMarshaler, v interface{}) kvjt {
	return kvjt{K: k, V: Reflect(v)}
}

// kvjts is a key-value pair implements text/json marshaler and stringer.
// String method intends to indicate severity level.
// For example syslog levels: "0" emergency;
//               						  "1" alert;
//               						  "2" critical;
//               						  "3" error;
//               						  "4" warning;
//               						  "5" notice;
//               						  "6" informational;
//               						  "7" debug;
// (https://en.wikipedia.org/wiki/Syslog#Severity_level).
type kvjts struct {
	K encoding.TextMarshaler
	V json.Marshaler
	S fmt.Stringer
}

func (kv kvjts) MarshalText() (text []byte, err error) { return kv.K.MarshalText() }
func (kv kvjts) MarshalJSON() ([]byte, error)          { return kv.V.MarshalJSON() }
func (kv kvjts) String() string                        { return kv.S.String() }

func StringSeverity(k string, v string) kvjts {
	return kvjts{K: String(k), V: String(v), S: String(v)}
}
