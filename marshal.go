package logastic

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/danil/logastic/encode"
)

// Bool returns JSON marshaler for the bool type.
func Bool(v bool) json.Marshaler { return boolV{V: v} }

type boolV struct{ V bool }

func (v boolV) MarshalJSON() ([]byte, error) {
	if v.V {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}

// Bool returns JSON marshaler for the pointer to the bool type.
func Boolp(p *bool) json.Marshaler { return boolP{P: p} }

type boolP struct{ P *bool }

func (p boolP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return boolV{V: *p.P}.MarshalJSON()
}

// Byte returns JSON marshaler for the byte slice type.
func Byte(v ...byte) json.Marshaler { return byteV{V: v} }

type byteV struct{ V []byte }

func (v byteV) MarshalJSON() ([]byte, error) {
	if len(v.V) == 0 {
		return []byte("null"), nil
	}
	return append([]byte(`"`), append(encode.Bytes(v.V), []byte(`"`)...)...), nil
}

// Bytes returns JSON marshaler for the byte slice type.
func Bytes(v []byte) json.Marshaler { return bytesV{V: v} }

type bytesV struct{ V []byte }

func (v bytesV) MarshalJSON() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}
	return append([]byte(`"`), append(encode.Bytes(v.V), []byte(`"`)...)...), nil
}

// Bytesp returns JSON marshaler for the pointer to the byte slice type.
func Bytesp(p *[]byte) json.Marshaler { return bytesP{P: p} }

type bytesP struct{ P *[]byte }

func (p bytesP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return bytesV{V: *p.P}.MarshalJSON()
}

// Complex128 returns JSON marshaler for the complex128 type.
func Complex128(v complex128) json.Marshaler { return complex128V{V: v} }

type complex128V struct{ V complex128 }

func (v complex128V) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%g", v.V)
	return append([]byte(`"`), append([]byte(s[1:len(s)-1]), []byte(`"`)...)...), nil
}

// Complex64 returns JSON marshaler for the complex64 type.
func Complex64(v complex64) json.Marshaler { return complex64V{V: v} }

type complex64V struct{ V complex64 }

func (v complex64V) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%g", v.V)
	return append([]byte(`"`), append([]byte(s[1:len(s)-1]), []byte(`"`)...)...), nil
}

// Error returns JSON marshaler for the error type.
func Error(v error) json.Marshaler { return errorV{V: v} }

type errorV struct{ V error }

func (v errorV) MarshalJSON() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}
	return append([]byte(`"`), append(encode.String(v.V.Error()), []byte(`"`)...)...), nil
}

// Float32 returns JSON marshaler for the float32 type.
func Float32(v float32) json.Marshaler { return float32V{V: v} }

type float32V struct{ V float32 }

func (v float32V) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(v.V)), nil
}

// Float32p returns JSON marshaler for the pointer to the float32 type.
func Float32p(p *float32) json.Marshaler { return float32P{P: p} }

type float32P struct{ P *float32 }

func (p float32P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return float32V{V: *p.P}.MarshalJSON()
}

// Float64 returns JSON marshaler for the float64 type.
func Float64(v float64) json.Marshaler { return float64V{V: v} }

type float64V struct{ V float64 }

func (v float64V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(v.V), 'f', -1, 64)), nil
}

// Float64p returns JSON marshaler for the pointer to the float64 type.
func Float64p(p *float64) json.Marshaler { return float64P{P: p} }

type float64P struct{ P *float64 }

func (p float64P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return float64V{V: *p.P}.MarshalJSON()
}

// Int returns JSON marshaler for the int type.
func Int(v int) json.Marshaler { return intV{V: v} }

type intV struct{ V int }

func (v intV) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Intp returns JSON marshaler for the pointer to the int type.
func Intp(p *int) json.Marshaler { return intP{P: p} }

type intP struct{ P *int }

func (p intP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return intV{V: *p.P}.MarshalJSON()
}

// Int16 returns JSON marshaler for the int16 type.
func Int16(v int16) json.Marshaler { return int16V{V: v} }

type int16V struct{ V int16 }

func (v int16V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Int16p returns JSON marshaler for the pointer to the int16 type.
func Int16p(p *int16) json.Marshaler { return int16P{P: p} }

type int16P struct{ P *int16 }

func (p int16P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return int16V{V: *p.P}.MarshalJSON()
}

// Int32 returns JSON marshaler for the int32 type.
func Int32(v int32) json.Marshaler { return int32V{V: v} }

type int32V struct{ V int32 }

func (v int32V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Int32p returns JSON marshaler for the pointer to the int32 type.
func Int32p(p *int32) json.Marshaler { return int32P{P: p} }

type int32P struct{ P *int32 }

func (p int32P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return int32V{V: *p.P}.MarshalJSON()
}

// Int64 returns JSON marshaler for the int64 type.
func Int64(v int64) json.Marshaler { return int64V{V: v} }

type int64V struct{ V int64 }

func (v int64V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Int64p returns JSON marshaler for the pointer to the int64 type.
func Int64p(p *int64) json.Marshaler { return int64P{P: p} }

type int64P struct{ P *int64 }

func (p int64P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return int64V{V: *p.P}.MarshalJSON()
}

// Int8 returns JSON marshaler for the int8 type.
func Int8(v int8) json.Marshaler { return int8V{V: v} }

type int8V struct{ V int8 }

func (v int8V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Int8p returns JSON marshaler for the pointer to the int8 type.
func Int8p(p *int8) json.Marshaler { return int8P{P: p} }

type int8P struct{ P *int8 }

func (p int8P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return int8V{V: *p.P}.MarshalJSON()
}

// Rune returns JSON marshaler for the rune type.
func Rune(v ...rune) json.Marshaler { return runeV{V: v} }

type runeV struct{ V []rune }

func (v runeV) MarshalJSON() ([]byte, error) {
	if len(v.V) == 0 {
		return []byte("null"), nil
	}
	return append([]byte(`"`), append(encode.Runes(v.V), []byte(`"`)...)...), nil
}

// Runes returns JSON marshaler for the rune slice type.
func Runes(v []rune) json.Marshaler { return runesV{V: v} }

type runesV struct{ V []rune }

func (v runesV) MarshalJSON() ([]byte, error) {
	if v.V == nil {
		return []byte("null"), nil
	}
	return append([]byte(`"`), append(encode.Runes(v.V), []byte(`"`)...)...), nil
}

// String returns JSON marshaler for the string type.
func String(v string) json.Marshaler { return stringV{V: v} }

type stringV struct{ V string }

func (v stringV) MarshalJSON() ([]byte, error) {
	return append([]byte(`"`), append(encode.String(v.V), []byte(`"`)...)...), nil
}

// Uint returns JSON marshaler for the uint type.
func Uint(v uint) json.Marshaler { return uintV{V: v} }

type uintV struct{ V uint }

func (v uintV) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Uintp returns JSON marshaler for the pointer to the uint type.
func Uintp(p *uint) json.Marshaler { return uintP{P: p} }

type uintP struct{ P *uint }

func (p uintP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return uintV{V: *p.P}.MarshalJSON()
}

// Uint16 returns JSON marshaler for the uint16 type.
func Uint16(v uint16) json.Marshaler { return uint16V{V: v} }

type uint16V struct{ V uint16 }

func (v uint16V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Uint16p returns JSON marshaler for the pointer to the uint16 type.
func Uint16p(p *uint16) json.Marshaler { return uint16P{P: p} }

type uint16P struct{ P *uint16 }

func (p uint16P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return uint16V{V: *p.P}.MarshalJSON()
}

// Uint32 returns JSON marshaler for the uint32 type.
func Uint32(v uint32) json.Marshaler { return uint32V{V: v} }

type uint32V struct{ V uint32 }

func (v uint32V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Uint32p returns JSON marshaler for the pointer to the uint32 type.
func Uint32p(p *uint32) json.Marshaler { return uint32P{P: p} }

type uint32P struct{ P *uint32 }

func (p uint32P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return uint32V{V: *p.P}.MarshalJSON()
}

// Uint64 returns JSON marshaler for the uint64 type.
func Uint64(v uint64) json.Marshaler { return uint64V{V: v} }

type uint64V struct{ V uint64 }

func (v uint64V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Uint64p returns JSON marshaler for the pointer to the uint64 type.
func Uint64p(p *uint64) json.Marshaler { return uint64P{P: p} }

type uint64P struct{ P *uint64 }

func (p uint64P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return uint64V{V: *p.P}.MarshalJSON()
}

// Uint8 returns JSON marshaler for the uint8 type.
func Uint8(v uint8) json.Marshaler { return uint8V{V: v} }

type uint8V struct{ V uint8 }

func (v uint8V) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Uint8p returns JSON marshaler for the pointer to the uint8 type.
func Uint8p(p *uint8) json.Marshaler { return uint8P{P: p} }

type uint8P struct{ P *uint8 }

func (p uint8P) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return uint8V{V: *p.P}.MarshalJSON()
}

// Uintptr returns JSON marshaler for the uintptr type.
func Uintptr(v uintptr) json.Marshaler { return uintptrV{V: v} }

type uintptrV struct{ V uintptr }

func (v uintptrV) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(v.V))), nil
}

// Uintptrp returns JSON marshaler for the pointer to the uintptr type.
func Uintptrp(p *uintptr) json.Marshaler { return uintptrP{P: p} }

type uintptrP struct{ P *uintptr }

func (p uintptrP) MarshalJSON() ([]byte, error) {
	if p.P == nil {
		return []byte("null"), nil
	}
	return uintptrV{V: *p.P}.MarshalJSON()
}
