// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encode0

import (
	"io"
	"sync"
	"unicode/utf8"
)

var pool = sync.Pool{New: func() interface{} { p := make([]byte, utf8.UTFMax); return &p }}

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

func Bytes(dst io.Writer, src []byte) error {
	idx := 0
	var oldr rune

	for {
		r, n := utf8.DecodeRune(src[idx:])
		if n == 0 {
			break
		}

		p, ok := codec[r]

		if ok {
			_, err := dst.Write(p)
			if err != nil {
				return err
			}

		} else if r == '"' && oldr != '\\' {
			_, err := dst.Write(append([]byte("\\"), src[idx:idx+n]...))
			if err != nil {
				return err
			}

		} else {
			_, err := dst.Write(src[idx : idx+n])
			if err != nil {
				return err
			}
		}

		oldr = r
		idx += n
	}

	return nil
}

func Runes(dst io.Writer, src []rune) error {
	idx := 0

	var (
		oldr rune
		err  error
	)

	for _, r := range src {
		idx, err = enc(dst, idx, r, oldr)
		if err != nil {
			return err
		}

		oldr = r
	}

	return nil
}

func String(dst io.Writer, src string) error {
	var (
		idx  int
		oldr rune
		err  error
	)

	for _, r := range src {
		idx, err = enc(dst, idx, r, oldr)
		if err != nil {
			return err
		}

		oldr = r
	}

	return nil
}

func enc(dst io.Writer, idx int, r, oldr rune) (int, error) {
	p, ok := codec[r]

	if ok {
		_, err := dst.Write(p)
		if err != nil {
			return 0, err
		}

		idx += len(p)
	} else if r == '"' && oldr != '\\' {
		_, err := dst.Write(append([]byte("\\"), '"'))
		if err != nil {
			return 0, err
		}

		idx += 2
	} else {
		p := *pool.Get().(*[]byte)
		defer pool.Put(&p)

		n := utf8.EncodeRune(p, r)

		_, err := dst.Write(p[:n])
		if err != nil {
			return 0, err
		}

		idx += n
	}

	return idx, nil
}
