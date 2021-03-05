// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log0

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/danil/log0/encode0"
)

// Bool returns stringer/JSON marshaler interface implementation for the bool type.
func Bool(v bool) boolV { return boolV{V: v} }

type boolV struct{ V bool }

func (v boolV) String() string {
	p, _ := v.MarshalText()
	return string(p)
}

func (v boolV) MarshalText() ([]byte, error) {
	if v.V {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

func (v boolV) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Bool returns stringer/JSON marshaler interface implementation for the pointer to the bool type.
func Boolp(p *bool) boolP { return boolP{P: p} }

type boolP struct{ P *bool }

func (p boolP) String() string {
	t, _ := p.MarshalText()
	return string(t)
}

func (p boolP) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return boolV{V: *p.P}.MarshalText()
}

func (p boolP) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

var bufPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

// Bytes returns stringer/JSON marshaler interface implementation for the byte slice type.
func Bytes(v []byte) bytesV { return bytesV{V: v} }

type bytesV struct{ V []byte }

func (v bytesV) String() string {
	p, _ := v.MarshalText()
	return string(p)
}

func (v bytesV) MarshalText() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer

	err := encode0.Bytes(&buf, v.V)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v bytesV) MarshalJSON() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}

	p, err := v.MarshalText()
	if err != nil {
		return nil, err
	}

	return append([]byte(`"`), append(p, []byte(`"`)...)...), nil
}

// Bytesp returns stringer/JSON marshaler interface implementation for the pointer to the byte slice type.
func Bytesp(p *[]byte) bytesP { return bytesP{P: p} }

type bytesP struct{ P *[]byte }

func (p bytesP) String() string {
	t, _ := p.MarshalText()
	return string(t)
}

func (p bytesP) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return bytesV{V: *p.P}.MarshalText()
}

func (p bytesP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return bytesV{V: *p.P}.MarshalJSON()
}

// Complex128 returns stringer/JSON marshaler interface implementation for the complex128 type.
func Complex128(v complex128) complex128V { return complex128V{V: v} }

type complex128V struct{ V complex128 }

func (v complex128V) String() string {
	s := fmt.Sprintf("%g", v.V)
	return s[1 : len(s)-1]
}

func (v complex128V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v complex128V) MarshalJSON() ([]byte, error) {
	return append([]byte(`"`), append([]byte(v.String()), []byte(`"`)...)...), nil
}

// Complex128p returns stringer/JSON marshaler interface implementation for the pointer to the complex128 type.
func Complex128p(p *complex128) complex128P { return complex128P{P: p} }

type complex128P struct{ P *complex128 }

func (p complex128P) String() string {
	if p.P == nil {
		return "null"
	}
	return complex128V{V: *p.P}.String()
}

func (p complex128P) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return complex128V{V: *p.P}.MarshalText()
}

func (p complex128P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return complex128V{V: *p.P}.MarshalJSON()
}

// Complex64 returns stringer/JSON marshaler interface implementation for the complex64 type.
func Complex64(v complex64) complex64V { return complex64V{V: v} }

type complex64V struct{ V complex64 }

func (v complex64V) String() string {
	s := fmt.Sprintf("%g", v.V)
	return s[1 : len(s)-1]
}

func (v complex64V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v complex64V) MarshalJSON() ([]byte, error) {
	return append([]byte(`"`), append([]byte(v.String()), []byte(`"`)...)...), nil
}

// Complex64p returns stringer/JSON marshaler interface implementation for the pointer to the complex64 type.
func Complex64p(p *complex64) complex64P { return complex64P{P: p} }

type complex64P struct{ P *complex64 }

func (p complex64P) String() string {
	if p.P == nil {
		return "null"
	}
	return complex64V{V: *p.P}.String()
}

func (p complex64P) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return complex64V{V: *p.P}.MarshalText()
}

func (p complex64P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return complex64V{V: *p.P}.MarshalJSON()
}

// Error returns stringer/JSON marshaler interface implementation for the error type.
func Error(v error) errorV { return errorV{V: v} }

type errorV struct{ V error }

func (v errorV) String() string {
	p, _ := v.MarshalText()
	return string(p)
}

func (v errorV) MarshalText() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer

	err := encode0.String(&buf, v.V.Error())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v errorV) MarshalJSON() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}

	p, err := v.MarshalText()
	if err != nil {
		return nil, err
	}

	return append([]byte(`"`), append(p, []byte(`"`)...)...), nil
}

// Float32 returns stringer/JSON marshaler interface implementation for the float32 type.
func Float32(v float32) float32V { return float32V{V: v} }

type float32V struct{ V float32 }

func (v float32V) String() string {
	return fmt.Sprint(v.V)
}

func (v float32V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v float32V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Float32p returns stringer/JSON marshaler interface implementation for the pointer to the float32 type.
func Float32p(p *float32) float32P { return float32P{P: p} }

type float32P struct{ P *float32 }

func (p float32P) String() string {
	if p.P == nil {
		return "null"
	}
	return float32V{V: *p.P}.String()
}

func (p float32P) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return float32V{V: *p.P}.MarshalText()
}

func (p float32P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Float64 returns stringer/JSON marshaler interface implementation for the float64 type.
func Float64(v float64) float64V { return float64V{V: v} }

type float64V struct{ V float64 }

func (v float64V) String() string {
	return strconv.FormatFloat(float64(v.V), 'f', -1, 64)
}

func (v float64V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v float64V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Float64p returns stringer/JSON marshaler interface implementation for the pointer to the float64 type.
func Float64p(p *float64) float64P { return float64P{P: p} }

type float64P struct{ P *float64 }

func (p float64P) String() string {
	if p.P == nil {
		return "null"
	}
	return float64V{V: *p.P}.String()
}

func (p float64P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p float64P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Int returns stringer/JSON marshaler interface implementation for the int type.
func Int(v int) intV { return intV{V: v} }

type intV struct{ V int }

func (v intV) String() string {
	return strconv.Itoa(v.V)
}

func (v intV) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v intV) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Intp returns stringer/JSON marshaler interface implementation for the pointer to the int type.
func Intp(p *int) intP { return intP{P: p} }

type intP struct{ P *int }

func (p intP) String() string {
	if p.P == nil {
		return "null"
	}
	return intV{V: *p.P}.String()
}

func (p intP) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p intP) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Int16 returns stringer/JSON marshaler interface implementation for the int16 type.
func Int16(v int16) int16V { return int16V{V: v} }

type int16V struct{ V int16 }

func (v int16V) String() string {
	return strconv.Itoa(int(v.V))
}

func (v int16V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v int16V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Int16p returns stringer/JSON marshaler interface implementation for the pointer to the int16 type.
func Int16p(p *int16) int16P { return int16P{P: p} }

type int16P struct{ P *int16 }

func (p int16P) String() string {
	if p.P == nil {
		return "null"
	}
	return int16V{V: *p.P}.String()
}

func (p int16P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p int16P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Int32 returns stringer/JSON marshaler interface implementation for the int32 type.
func Int32(v int32) int32V { return int32V{V: v} }

type int32V struct{ V int32 }

func (v int32V) String() string {
	return strconv.Itoa(int(v.V))
}

func (v int32V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v int32V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Int32p returns stringer/JSON marshaler interface implementation for the pointer to the int32 type.
func Int32p(p *int32) int32P { return int32P{P: p} }

type int32P struct{ P *int32 }

func (p int32P) String() string {
	if p.P == nil {
		return "null"
	}
	return int32V{V: *p.P}.String()
}

func (p int32P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p int32P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Int64 returns stringer/JSON marshaler interface implementation for the int64 type.
func Int64(v int64) int64V { return int64V{V: v} }

type int64V struct{ V int64 }

func (v int64V) String() string {
	return strconv.FormatInt(int64(v.V), 10)
}

func (v int64V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v int64V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Int64p returns stringer/JSON marshaler interface implementation for the pointer to the int64 type.
func Int64p(p *int64) int64P { return int64P{P: p} }

type int64P struct{ P *int64 }

func (p int64P) String() string {
	if p.P == nil {
		return "null"
	}
	return int64V{V: *p.P}.String()
}

func (p int64P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p int64P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Int8 returns stringer/JSON marshaler interface implementation for the int8 type.
func Int8(v int8) int8V { return int8V{V: v} }

type int8V struct{ V int8 }

func (v int8V) String() string {
	return strconv.Itoa(int(v.V))
}

func (v int8V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v int8V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Int8p returns stringer/JSON marshaler interface implementation for the pointer to the int8 type.
func Int8p(p *int8) int8P { return int8P{P: p} }

type int8P struct{ P *int8 }

func (p int8P) String() string {
	if p.P == nil {
		return "null"
	}
	return int8V{V: *p.P}.String()
}

func (p int8P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p int8P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Runes returns stringer/JSON marshaler interface implementation for the rune slice type.
func Runes(v []rune) runesV { return runesV{V: v} }

type runesV struct{ V []rune }

func (v runesV) String() string {
	if v.V == nil {
		return "null"
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	err := encode0.Runes(buf, v.V)
	if err != nil {
		return ""
	}

	return buf.String()
}

func (v runesV) MarshalText() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer

	err := encode0.Runes(&buf, v.V)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v runesV) MarshalJSON() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}

	p, err := v.MarshalText()
	if err != nil {
		return nil, err
	}

	return append([]byte(`"`), append(p, []byte(`"`)...)...), nil
}

// Runesp returns stringer/JSON marshaler interface implementation for the pointer to the rune slice type.
func Runesp(p *[]rune) runesP { return runesP{P: p} }

type runesP struct{ P *[]rune }

func (p runesP) String() string {
	if p.P == nil {
		return "null"
	}
	return runesV{V: *p.P}.String()
}

func (p runesP) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return runesV{V: *p.P}.MarshalText()
}

func (p runesP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return runesV{V: *p.P}.MarshalJSON()
}

// String returns stringer/JSON marshaler interface implementation for the string type.
func String(v string) stringV { return stringV{V: v} }

type stringV struct{ V string }

func (v stringV) String() string {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	err := encode0.String(buf, v.V)
	if err != nil {
		return ""
	}

	return buf.String()
}

func (v stringV) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	err := encode0.String(&buf, v.V)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v stringV) MarshalJSON() ([]byte, error) {
	p, err := v.MarshalText()
	if err != nil {
		return nil, err
	}
	return append([]byte(`"`), append(p, []byte(`"`)...)...), nil
}

// Stringp returns stringer/JSON marshaler interface implementation for the pointer to the string type.
func Stringp(p *string) stringP { return stringP{P: p} }

type stringP struct{ P *string }

func (p stringP) String() string {
	if p.P == nil {
		return "null"
	}
	return stringV{V: *p.P}.String()
}

func (p stringP) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return stringV{V: *p.P}.MarshalText()
}

func (p stringP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return stringV{V: *p.P}.MarshalJSON()
}

// Text returns stringer/JSON marshaler interface implementation for the encoding.TextMarshaler type.
func Text(v encoding.TextMarshaler) textV { return textV{V: v} }

type textV struct{ V encoding.TextMarshaler }

func (v textV) String() string {
	p, err := v.V.MarshalText()
	if err != nil {
		return ""
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	err = encode0.Bytes(buf, p)
	if err != nil {
		return ""
	}

	return buf.String()
}

func (v textV) MarshalText() ([]byte, error) {
	p, err := v.V.MarshalText()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = encode0.Bytes(&buf, p)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (v textV) MarshalJSON() ([]byte, error) {
	p, err := v.MarshalText()
	if err != nil {
		return nil, err
	}
	return append([]byte(`"`), append(p, []byte(`"`)...)...), nil
}

// Uint returns stringer/JSON marshaler interface implementation for the uint type.
func Uint(v uint) uintV { return uintV{V: v} }

type uintV struct{ V uint }

func (v uintV) String() string {
	return strconv.FormatUint(uint64(v.V), 10)
}

func (v uintV) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v uintV) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Uintp returns stringer/JSON marshaler interface implementation for the pointer to the uint type.
func Uintp(p *uint) uintP { return uintP{P: p} }

type uintP struct{ P *uint }

func (p uintP) String() string {
	if p.P == nil {
		return "null"
	}
	return uintV{V: *p.P}.String()
}

func (p uintP) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p uintP) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Uint16 returns stringer/JSON marshaler interface implementation for the uint16 type.
func Uint16(v uint16) uint16V { return uint16V{V: v} }

type uint16V struct{ V uint16 }

func (v uint16V) String() string {
	return strconv.FormatUint(uint64(v.V), 10)
}

func (v uint16V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v uint16V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Uint16p returns stringer/JSON marshaler interface implementation for the pointer to the uint16 type.
func Uint16p(p *uint16) uint16P { return uint16P{P: p} }

type uint16P struct{ P *uint16 }

func (p uint16P) String() string {
	if p.P == nil {
		return "null"
	}
	return uint16V{V: *p.P}.String()
}

func (p uint16P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p uint16P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Uint32 returns stringer/JSON marshaler interface implementation for the uint32 type.
func Uint32(v uint32) uint32V { return uint32V{V: v} }

type uint32V struct{ V uint32 }

func (v uint32V) String() string {
	return strconv.FormatUint(uint64(v.V), 10)
}

func (v uint32V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v uint32V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Uint32p returns stringer/JSON marshaler interface implementation for the pointer to the uint32 type.
func Uint32p(p *uint32) uint32P { return uint32P{P: p} }

type uint32P struct{ P *uint32 }

func (p uint32P) String() string {
	if p.P == nil {
		return "null"
	}
	return uint32V{V: *p.P}.String()
}

func (p uint32P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p uint32P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Uint64 returns stringer/JSON marshaler interface implementation for the uint64 type.
func Uint64(v uint64) uint64V { return uint64V{V: v} }

type uint64V struct{ V uint64 }

func (v uint64V) String() string {
	return strconv.FormatUint(v.V, 10)
}

func (v uint64V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v uint64V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Uint64p returns stringer/JSON marshaler interface implementation for the pointer to the uint64 type.
func Uint64p(p *uint64) uint64P { return uint64P{P: p} }

type uint64P struct{ P *uint64 }

func (p uint64P) String() string {
	if p.P == nil {
		return "null"
	}
	return uint64V{V: *p.P}.String()
}

func (p uint64P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p uint64P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Uint8 returns stringer/JSON marshaler interface implementation for the uint8 type.
func Uint8(v uint8) uint8V { return uint8V{V: v} }

type uint8V struct{ V uint8 }

func (v uint8V) String() string {
	return strconv.FormatUint(uint64(v.V), 10)
}

func (v uint8V) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v uint8V) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Uint8p returns stringer/JSON marshaler interface implementation for the pointer to the uint8 type.
func Uint8p(p *uint8) uint8P { return uint8P{P: p} }

type uint8P struct{ P *uint8 }

func (p uint8P) String() string {
	if p.P == nil {
		return "null"
	}
	return uint8V{V: *p.P}.String()
}

func (p uint8P) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p uint8P) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Uintptr returns stringer/JSON marshaler interface implementation for the uintptr type.
func Uintptr(v uintptr) uintptrV { return uintptrV{V: v} }

type uintptrV struct{ V uintptr }

func (v uintptrV) String() string {
	return strconv.FormatUint(uint64(v.V), 10)
}

func (v uintptrV) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v uintptrV) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

// Uintptrp returns stringer/JSON marshaler interface implementation for the pointer to the uintptr type.
func Uintptrp(p *uintptr) uintptrP { return uintptrP{P: p} }

type uintptrP struct{ P *uintptr }

func (p uintptrP) String() string {
	if p.P == nil {
		return "null"
	}
	return uintptrV{V: *p.P}.String()
}

func (p uintptrP) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p uintptrP) MarshalJSON() ([]byte, error) {
	return p.MarshalText()
}

// Time returns stringer/JSON marshaler interface implementation for the time time type.
func Time(v time.Time) timeV { return timeV{V: v} }

type timeV struct{ V time.Time }

func (v timeV) String() string {
	return v.V.String()
}

func (v timeV) MarshalText() ([]byte, error) {
	return v.V.MarshalText()
}

func (v timeV) MarshalJSON() ([]byte, error) {
	return v.V.MarshalJSON()
}

// Timep returns stringer/JSON marshaler interface implementation for the pointer to the time time type.
func Timep(p *time.Time) timeP { return timeP{P: p} }

type timeP struct{ P *time.Time }

func (p timeP) String() string {
	if p.P == nil {
		return "null"
	}
	return timeV{V: *p.P}.String()
}

func (p timeP) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return timeV{V: *p.P}.MarshalText()
}

func (p timeP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return timeV{V: *p.P}.MarshalJSON()
}

// Duration returns stringer/JSON marshaler interface implementation for the time duration type.
func Duration(v time.Duration) durationV { return durationV{V: v} }

type durationV struct{ V time.Duration }

func (v durationV) String() string {
	return v.V.String()
}

func (v durationV) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v durationV) MarshalJSON() ([]byte, error) {
	return append([]byte(`"`), append([]byte(v.String()), []byte(`"`)...)...), nil
}

// Durationp returns stringer/JSON marshaler interface implementation for the pointer to the time duration type.
func Durationp(p *time.Duration) durationP { return durationP{P: p} }

type durationP struct{ P *time.Duration }

func (p durationP) String() string {
	if p.P == nil {
		return "null"
	}
	return durationV{V: *p.P}.String()
}

func (p durationP) MarshalText() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return durationV{V: *p.P}.MarshalText()
}

func (p durationP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return durationV{V: *p.P}.MarshalJSON()
}

// Func returns stringer/JSON marshaler interface implementation for the custom func type.
func Func(v func() KV) funcV { return funcV{V: v} }

type funcV struct{ V func() KV }

func (v funcV) String() string {
	p, _ := v.V().MarshalText()
	return string(p)
}

func (v funcV) MarshalText() ([]byte, error) {
	return v.V().MarshalText()
}

func (v funcV) MarshalJSON() ([]byte, error) {
	return v.V().MarshalJSON()
}

// Raw returns stringer/JSON marshaler interface implementation for the raw byte slice.
func Raw(v []byte) rawV { return rawV{V: v} }

type rawV struct{ V []byte }

func (v rawV) String() string {
	if v.V == nil {
		return "null"
	}
	return string(v.V)
}

func (v rawV) MarshalText() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}
	return v.V, nil
}

func (v rawV) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

func Any(v interface{}) anyV { return anyV{V: v} }

type anyV struct{ V interface{} }

func (v anyV) String() string {
	switch x := v.V.(type) {
	case bool:
		return Bool(x).String()
	case *bool:
		return Boolp(x).String()
	case []byte:
		return Bytes(x).String()
	case *[]byte:
		return Bytesp(x).String()
	case complex128:
		return Complex128(x).String()
	case *complex128:
		return Complex128p(x).String()
	case complex64:
		return Complex64(x).String()
	case *complex64:
		return Complex64p(x).String()
	case error:
		return Error(x).String()
	case float32:
		return Float32(x).String()
	case *float32:
		return Float32p(x).String()
	case float64:
		return Float64(x).String()
	case *float64:
		return Float64p(x).String()
	case int:
		return Int(x).String()
	case *int:
		return Intp(x).String()
	case int16:
		return Int16(x).String()
	case *int16:
		return Int16p(x).String()
	case int32:
		return Int32(x).String()
	case *int32:
		return Int32p(x).String()
	case int64:
		return Int64(x).String()
	case *int64:
		return Int64p(x).String()
	case int8:
		return Int8(x).String()
	case *int8:
		return Int8p(x).String()
	case []rune:
		return Runes(x).String()
	case *[]rune:
		return Runesp(x).String()
	case string:
		return String(x).String()
	case *string:
		return Stringp(x).String()
	case uint:
		return Uint(x).String()
	case *uint:
		return Uintp(x).String()
	case uint16:
		return Uint16(x).String()
	case *uint16:
		return Uint16p(x).String()
	case uint32:
		return Uint32(x).String()
	case *uint32:
		return Uint32p(x).String()
	case uint64:
		return Uint64(x).String()
	case *uint64:
		return Uint64p(x).String()
	case uint8:
		return Uint8(x).String()
	case *uint8:
		return Uint8p(x).String()
	case uintptr:
		return Uintptr(x).String()
	case *uintptr:
		return Uintptrp(x).String()
	case time.Time:
		return Time(x).String()
	case *time.Time:
		return Timep(x).String()
	case time.Duration:
		return Duration(x).String()
	case *time.Duration:
		return Durationp(x).String()
	case encoding.TextMarshaler:
		return Text(x).String()
	case json.Marshaler:
		p, _ := x.MarshalJSON()
		return string(p)
	default:
		return Reflect(x).String()
	}
}

func (v anyV) MarshalText() ([]byte, error) {
	switch x := v.V.(type) {
	case bool:
		return Bool(x).MarshalText()
	case *bool:
		return Boolp(x).MarshalText()
	case []byte:
		return Bytes(x).MarshalText()
	case *[]byte:
		return Bytesp(x).MarshalText()
	case complex128:
		return Complex128(x).MarshalText()
	case *complex128:
		return Complex128p(x).MarshalText()
	case complex64:
		return Complex64(x).MarshalText()
	case *complex64:
		return Complex64p(x).MarshalText()
	case error:
		return Error(x).MarshalText()
	case float32:
		return Float32(x).MarshalText()
	case *float32:
		return Float32p(x).MarshalText()
	case float64:
		return Float64(x).MarshalText()
	case *float64:
		return Float64p(x).MarshalText()
	case int:
		return Int(x).MarshalText()
	case *int:
		return Intp(x).MarshalText()
	case int16:
		return Int16(x).MarshalText()
	case *int16:
		return Int16p(x).MarshalText()
	case int32:
		return Int32(x).MarshalText()
	case *int32:
		return Int32p(x).MarshalText()
	case int64:
		return Int64(x).MarshalText()
	case *int64:
		return Int64p(x).MarshalText()
	case int8:
		return Int8(x).MarshalText()
	case *int8:
		return Int8p(x).MarshalText()
	case []rune:
		return Runes(x).MarshalText()
	case *[]rune:
		return Runesp(x).MarshalText()
	case string:
		return String(x).MarshalText()
	case *string:
		return Stringp(x).MarshalText()
	case uint:
		return Uint(x).MarshalText()
	case *uint:
		return Uintp(x).MarshalText()
	case uint16:
		return Uint16(x).MarshalText()
	case *uint16:
		return Uint16p(x).MarshalText()
	case uint32:
		return Uint32(x).MarshalText()
	case *uint32:
		return Uint32p(x).MarshalText()
	case uint64:
		return Uint64(x).MarshalText()
	case *uint64:
		return Uint64p(x).MarshalText()
	case uint8:
		return Uint8(x).MarshalText()
	case *uint8:
		return Uint8p(x).MarshalText()
	case uintptr:
		return Uintptr(x).MarshalText()
	case *uintptr:
		return Uintptrp(x).MarshalText()
	case time.Time:
		return Time(x).MarshalText()
	case *time.Time:
		return Timep(x).MarshalText()
	case time.Duration:
		return Duration(x).MarshalText()
	case *time.Duration:
		return Durationp(x).MarshalText()
	case encoding.TextMarshaler:
		return x.MarshalText()
	default:
		return Reflect(x).MarshalText()
	}
}

func (v anyV) MarshalJSON() ([]byte, error) {
	switch x := v.V.(type) {
	case bool:
		return Bool(x).MarshalJSON()
	case *bool:
		return Boolp(x).MarshalJSON()
	case []byte:
		return Bytes(x).MarshalJSON()
	case *[]byte:
		return Bytesp(x).MarshalJSON()
	case complex128:
		return Complex128(x).MarshalJSON()
	case *complex128:
		return Complex128p(x).MarshalJSON()
	case complex64:
		return Complex64(x).MarshalJSON()
	case *complex64:
		return Complex64p(x).MarshalJSON()
	case error:
		return Error(x).MarshalJSON()
	case float32:
		return Float32(x).MarshalJSON()
	case *float32:
		return Float32p(x).MarshalJSON()
	case float64:
		return Float64(x).MarshalJSON()
	case *float64:
		return Float64p(x).MarshalJSON()
	case int:
		return Int(x).MarshalJSON()
	case *int:
		return Intp(x).MarshalJSON()
	case int16:
		return Int16(x).MarshalJSON()
	case *int16:
		return Int16p(x).MarshalJSON()
	case int32:
		return Int32(x).MarshalJSON()
	case *int32:
		return Int32p(x).MarshalJSON()
	case int64:
		return Int64(x).MarshalJSON()
	case *int64:
		return Int64p(x).MarshalJSON()
	case int8:
		return Int8(x).MarshalJSON()
	case *int8:
		return Int8p(x).MarshalJSON()
	case []rune:
		return Runes(x).MarshalJSON()
	case *[]rune:
		return Runesp(x).MarshalJSON()
	case string:
		return String(x).MarshalJSON()
	case *string:
		return Stringp(x).MarshalJSON()
	case uint:
		return Uint(x).MarshalJSON()
	case *uint:
		return Uintp(x).MarshalJSON()
	case uint16:
		return Uint16(x).MarshalJSON()
	case *uint16:
		return Uint16p(x).MarshalJSON()
	case uint32:
		return Uint32(x).MarshalJSON()
	case *uint32:
		return Uint32p(x).MarshalJSON()
	case uint64:
		return Uint64(x).MarshalJSON()
	case *uint64:
		return Uint64p(x).MarshalJSON()
	case uint8:
		return Uint8(x).MarshalJSON()
	case *uint8:
		return Uint8p(x).MarshalJSON()
	case uintptr:
		return Uintptr(x).MarshalJSON()
	case *uintptr:
		return Uintptrp(x).MarshalJSON()
	case time.Time:
		return Time(x).MarshalJSON()
	case *time.Time:
		return Timep(x).MarshalJSON()
	case time.Duration:
		return Duration(x).MarshalJSON()
	case *time.Duration:
		return Durationp(x).MarshalJSON()
	case encoding.TextMarshaler:
		return Text(x).MarshalJSON()
	case json.Marshaler:
		return x.MarshalJSON()
	default:
		return Reflect(x).MarshalJSON()
	}
}

func Reflect(v interface{}) reflectV { return reflectV{V: v} }

type reflectV struct{ V interface{} }

func (v reflectV) String() string {
	if v.V == nil {
		return "null"
	}

	val := reflect.ValueOf(v.V)

	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if val.IsNil() {
			return "null"

		} else if val.Kind() == reflect.Ptr {
			return reflectV{V: val.Elem().Interface()}.String()

		} else if val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8 { // Byte slice.
			buf := bufPool.Get().(*bytes.Buffer)
			buf.Reset()
			defer bufPool.Put(buf)

			p := val.Bytes()
			enc := base64.NewEncoder(base64.StdEncoding, buf)
			_, _ = enc.Write(p)
			enc.Close()

			return buf.String()
		}
	}

	return fmt.Sprint(v.V)
}

func (v reflectV) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v reflectV) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.V)
}
