package es_test

import (
	"bytes"
	"testing"

	"github.com/maloquacious/es"
)

func TestCrockfordVectors(t *testing.T) {
	tests := []struct {
		plain   string
		encoded string
	}{
		{"", ""},
		{"f", "CR"},
		{"fo", "CSQG"},
		{"foo", "CSQPY"},
		{"foob", "CSQPYRG"},
		{"fooba", "CSQPYRK1"},
		{"foobar", "CSQPYRK1E8"},
		{"test", "EHJQ6X0"},
	}

	for _, test := range tests {
		t.Run(test.plain, func(t *testing.T) {
			if got := es.Crockford.EncodeToString([]byte(test.plain)); got != test.encoded {
				t.Fatalf("Crockford.EncodeToString(%q) = %q, want %q", test.plain, got, test.encoded)
			}
			got, err := es.Crockford.DecodeString(test.encoded)
			if err != nil {
				t.Fatalf("Crockford.DecodeString(%q): %v", test.encoded, err)
			}
			if string(got) != test.plain {
				t.Fatalf("Crockford.DecodeString(%q) = %q, want %q", test.encoded, got, test.plain)
			}
		})
	}
}

func TestDecodeCrockfordHumanInput(t *testing.T) {
	for _, encoded := range []string{"csqpyrk1e8", "CSQPY-RK1E8", "csqpyrkle8", "csqpyrkie8"} {
		got, err := es.Crockford.DecodeString(encoded)
		if err != nil {
			t.Fatalf("Crockford.DecodeString(%q): %v", encoded, err)
		}
		if string(got) != "foobar" {
			t.Errorf("Crockford.DecodeString(%q) = %q, want %q", encoded, got, "foobar")
		}
	}

	got, err := es.Crockford.DecodeString("ehjq6xo")
	if err != nil {
		t.Fatalf("Crockford.DecodeString alias O: %v", err)
	}
	if string(got) != "test" {
		t.Errorf("Crockford.DecodeString alias O = %q, want %q", got, "test")
	}
}

func TestCrockfordRoundTrip(t *testing.T) {
	for length := 0; length <= 64; length++ {
		src := make([]byte, length)
		for i := range src {
			src[i] = byte(i*37 + length)
		}
		got, err := es.Crockford.DecodeString(es.Crockford.EncodeToString(src))
		if err != nil {
			t.Fatalf("length %d: %v", length, err)
		}
		if !bytes.Equal(got, src) {
			t.Errorf("length %d: got %x, want %x", length, got, src)
		}
	}
}

func TestDecodeCrockfordRejectsMalformedInput(t *testing.T) {
	for _, encoded := range []string{
		"U0",      // Invalid character.
		"C",       // Invalid encoded length.
		"CCC",     // Invalid encoded length.
		"CS",      // Non-zero trailing bits in a two-character tail.
		"0001",    // Non-zero trailing bits in a four-character tail.
		"00001",   // Non-zero trailing bits in a five-character tail.
		"0000001", // Non-zero trailing bits in a seven-character tail.
	} {
		t.Run(encoded, func(t *testing.T) {
			got, err := es.Crockford.DecodeString(encoded)
			if err == nil {
				t.Fatalf("Crockford.DecodeString(%q) = %x, want error", encoded, got)
			}
			if got != nil {
				t.Errorf("Crockford.DecodeString(%q) returned partial output %x", encoded, got)
			}
		})
	}
}
