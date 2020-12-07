package encode

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

			out := Bytes([]byte(string(in)))

			if !bytes.Equal(out, expected) {
				t.Fatalf("expected: %s, recieved: %s", expected, out)
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

			out := Runes([]rune{in})

			if !bytes.Equal(out, expected) {
				t.Fatalf("expected: %s, recieved: %s", expected, out)
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

			out := String(string(in))

			if !bytes.Equal(out, expected) {
				t.Fatalf("expected: %s, recieved: %s", expected, out)
			}
		})
	}
}
