// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encode0

import (
	"bytes"
	"testing"
)

func TestBytes(t *testing.T) {
	for in, expected := range codec {
		in := in
		expected := expected
		t.Run(string(in), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			err := Bytes(&buf, []byte(string(in)))
			if err != nil {
				t.Fatalf("encode bytes write error: %s", err)
			}

			if !bytes.Equal(buf.Bytes(), expected) {
				t.Errorf("expected: %s, recieved: %s", expected, buf.String())
			}
		})
	}
}

func TestRunes(t *testing.T) {
	for in, expected := range codec {
		in := in
		expected := expected
		t.Run(string(in), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			err := Runes(&buf, []rune{in})
			if err != nil {
				t.Fatalf("encode runes write error: %s", err)
			}

			if !bytes.Equal(buf.Bytes(), expected) {
				t.Errorf("expected: %s, recieved: %s", expected, buf.String())
			}
		})
	}
}

func TestString(t *testing.T) {
	for in, expected := range codec {
		in := in
		expected := expected
		t.Run(string(in), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			err := String(&buf, string(in))
			if err != nil {
				t.Fatalf("encode string write error: %s", err)
			}

			if !bytes.Equal(buf.Bytes(), expected) {
				t.Errorf("expected: %s, recieved: %s", expected, buf.String())
			}
		})
	}
}
