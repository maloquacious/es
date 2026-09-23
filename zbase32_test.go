package es_test

import (
	"bytes"
	"testing"

	"github.com/maloquacious/es"
)

func TestZBase32Vectors(t *testing.T) {
	tests := []struct {
		plain   []byte
		encoded string
	}{
		{nil, ""},
		{[]byte{0x00}, "yy"},
		{[]byte{0xf0, 0xbf, 0xc7}, "6n9hq"},
		{[]byte{0xd4, 0x7a, 0x04}, "4t7ye"},
	}

	for _, test := range tests {
		t.Run(test.encoded, func(t *testing.T) {
			if got := es.ZBase32.EncodeToString(test.plain); got != test.encoded {
				t.Fatalf("ZBase32.EncodeToString(%x) = %q, want %q", test.plain, got, test.encoded)
			}
			got, err := es.ZBase32.DecodeString(test.encoded)
			if err != nil {
				t.Fatalf("ZBase32.DecodeString(%q): %v", test.encoded, err)
			}
			if !bytes.Equal(got, test.plain) {
				t.Fatalf("ZBase32.DecodeString(%q) = %x, want %x", test.encoded, got, test.plain)
			}
		})
	}
}

func TestDecodeZBase32IsCaseInsensitive(t *testing.T) {
	got, err := es.ZBase32.DecodeString("6N9HQ")
	if err != nil {
		t.Fatalf("ZBase32.DecodeString: %v", err)
	}
	if want := []byte{0xf0, 0xbf, 0xc7}; !bytes.Equal(got, want) {
		t.Errorf("ZBase32.DecodeString = %x, want %x", got, want)
	}
}

func TestZBase32RoundTrip(t *testing.T) {
	for length := 0; length <= 64; length++ {
		src := make([]byte, length)
		for i := range src {
			src[i] = byte(i*37 + length)
		}
		got, err := es.ZBase32.DecodeString(es.ZBase32.EncodeToString(src))
		if err != nil {
			t.Fatalf("length %d: %v", length, err)
		}
		if !bytes.Equal(got, src) {
			t.Errorf("length %d: got %x, want %x", length, got, src)
		}
	}
}

func TestDecodeZBase32RejectsMalformedInput(t *testing.T) {
	for _, encoded := range []string{
		"0y",      // Invalid character.
		"y",       // Invalid encoded length.
		"yyy",     // Invalid encoded length.
		"yb",      // Non-zero trailing bits in a two-character tail.
		"yyyb",    // Non-zero trailing bits in a four-character tail.
		"yyyyb",   // Non-zero trailing bits in a five-character tail.
		"yyyyyyb", // Non-zero trailing bits in a seven-character tail.
	} {
		t.Run(encoded, func(t *testing.T) {
			got, err := es.ZBase32.DecodeString(encoded)
			if err == nil {
				t.Fatalf("ZBase32.DecodeString(%q) = %x, want error", encoded, got)
			}
			if got != nil {
				t.Errorf("ZBase32.DecodeString(%q) returned partial output %x", encoded, got)
			}
		})
	}
}
