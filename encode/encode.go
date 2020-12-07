package encode

import (
	"unicode/utf8"
)

var codec = map[rune][]byte{
	'\x00': []byte("\\u0000"),
	'\x01': []byte("\\u0001"),
	'\x02': []byte("\\u0002"),
	'\x03': []byte("\\u0003"),
	'\x04': []byte("\\u0004"),
	'\x05': []byte("\\u0005"),
	'\x06': []byte("\\u0006"),
	'\x07': []byte("\\u0007"),
	'\x08': []byte("\\u0008"),
	'\x09': []byte("\\t"),
	'\x0a': []byte("\\n"),
	'\x0b': []byte("\\u000b"),
	'\x0c': []byte("\\u000c"),
	'\x0d': []byte("\\r"),
	'\x0e': []byte("\\u000e"),
	'\x0f': []byte("\\u000f"),
	'\x10': []byte("\\u0010"),
	'\x11': []byte("\\u0011"),
	'\x12': []byte("\\u0012"),
	'\x13': []byte("\\u0013"),
	'\x14': []byte("\\u0014"),
	'\x15': []byte("\\u0015"),
	'\x16': []byte("\\u0016"),
	'\x17': []byte("\\u0017"),
	'\x18': []byte("\\u0018"),
	'\x19': []byte("\\u0019"),
	'\x1a': []byte("\\u001a"),
	'\x1b': []byte("\\u001b"),
	'\x1c': []byte("\\u001c"),
	'\x1d': []byte("\\u001d"),
	'\x1e': []byte("\\u001e"),
	'\x1f': []byte("\\u001f"),
}

func Bytes(raw []byte) []byte {
	esc := make([]byte, 0, len(raw))
	idx := 0
	var oldr rune

	for {
		r, n := utf8.DecodeRune(raw[idx:])
		if n == 0 {
			break
		}

		p, ok := codec[r]

		if ok {
			esc = append(esc, p...)
		} else if r == '"' && oldr != '\\' {
			esc = append(esc, append([]byte("\\"), raw[idx:idx+n]...)...)
		} else {
			esc = append(esc, raw[idx:idx+n]...)
		}

		oldr = r
		idx += n
	}

	return esc
}

func Runes(raw []rune) []byte {
	esc := make([]byte, len(raw)*utf8.UTFMax)
	idx := 0
	var oldr rune

	for _, r := range raw {
		idx, esc = enc(idx, esc, r, oldr)
		oldr = r
	}

	esc = esc[:idx]
	return esc
}

func String(raw string) []byte {
	esc := make([]byte, len(raw)*utf8.UTFMax)
	idx := 0
	var oldr rune

	for _, r := range raw {
		idx, esc = enc(idx, esc, r, oldr)
		oldr = r
	}

	esc = esc[:idx]
	return esc
}

func enc(idx int, esc []byte, r, oldr rune) (int, []byte) {
	p, ok := codec[r]

	if ok {
		esc = append(esc[:idx], p...)
		idx += len(p)
	} else if r == '"' && oldr != '\\' {
		esc = append(esc[:idx], append([]byte("\\"), '"')...)
		idx += 2
	} else {
		idx += utf8.EncodeRune(esc[idx:], r)
	}

	return idx, esc
}
