package es_test

import (
	"bytes"
	"testing"

	"github.com/maloquacious/es"
)

func TestWordSafeVectors(t *testing.T) {
	tests := []struct {
		plain   string
		encoded string
	}{
		{"", ""},
		{"f", "Jj"},
		{"fo", "JmhR"},
		{"foo", "Jmhgw"},
		{"foob", "JmhgwjR"},
		{"fooba", "JmhgwjX3"},
		{"foobar", "JmhgwjX3PC"},
	}

	for _, test := range tests {
		t.Run(test.plain, func(t *testing.T) {
			if got := es.WordSafe.EncodeToString([]byte(test.plain)); got != test.encoded {
				t.Fatalf("WordSafe.EncodeToString(%q) = %q, want %q", test.plain, got, test.encoded)
			}
			got, err := es.WordSafe.DecodeString(test.encoded)
			if err != nil {
				t.Fatalf("WordSafe.DecodeString(%q): %v", test.encoded, err)
			}
			if string(got) != test.plain {
				t.Fatalf("WordSafe.DecodeString(%q) = %q, want %q", test.encoded, got, test.plain)
			}
		})
	}
}

func TestWordSafeIsCaseSensitive(t *testing.T) {
	upper, err := es.WordSafe.DecodeString("JmhgwjX3PC")
	if err != nil {
		t.Fatalf("Decode uppercase symbol: %v", err)
	}
	lower, err := es.WordSafe.DecodeString("jmhgwjX3PC")
	if err != nil {
		t.Fatalf("Decode lowercase symbol: %v", err)
	}
	if bytes.Equal(upper, lower) {
		t.Fatal("changing J to j did not change the decoded value")
	}
}

func TestWordSafeRoundTrip(t *testing.T) {
	for length := 0; length <= 64; length++ {
		src := make([]byte, length)
		for i := range src {
			src[i] = byte(i*37 + length)
		}
		got, err := es.WordSafe.DecodeString(es.WordSafe.EncodeToString(src))
		if err != nil {
			t.Fatalf("length %d: %v", length, err)
		}
		if !bytes.Equal(got, src) {
			t.Errorf("length %d: got %x, want %x", length, got, src)
		}
	}
}

func TestWordSafeRejectsMalformedInput(t *testing.T) {
	for _, encoded := range []string{
		"A2",      // Invalid character.
		"C",       // Invalid encoded length.
		"CCC",     // Invalid encoded length.
		"23",      // Non-zero trailing bits in a two-character tail.
		"2223",    // Non-zero trailing bits in a four-character tail.
		"22223",   // Non-zero trailing bits in a five-character tail.
		"2222223", // Non-zero trailing bits in a seven-character tail.
	} {
		t.Run(encoded, func(t *testing.T) {
			got, err := es.WordSafe.DecodeString(encoded)
			if err == nil {
				t.Fatalf("WordSafe.DecodeString(%q) = %x, want error", encoded, got)
			}
			if got != nil {
				t.Errorf("WordSafe.DecodeString(%q) returned partial output %x", encoded, got)
			}
		})
	}
}
