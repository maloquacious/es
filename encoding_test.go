package es_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/maloquacious/es"
)

func TestEncodingBufferAPI(t *testing.T) {
	tests := []struct {
		name    string
		codec   es.Encoding
		plain   []byte
		encoded string
	}{
		{"Crockford", es.Crockford, []byte("foobar"), "CSQPYRK1E8"},
		{"z-base-32", es.ZBase32, []byte("hello"), "pb1sa5dx"},
		{"word-safe", es.WordSafe, []byte("hello"), "M3Wgjq5Q"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := make([]byte, test.codec.EncodedLen(len(test.plain)))
			test.codec.Encode(encoded, test.plain)
			if string(encoded) != test.encoded {
				t.Fatalf("Encode(%q) = %q, want %q", test.plain, encoded, test.encoded)
			}

			decoded := make([]byte, test.codec.DecodedLen(len(encoded)))
			n, err := test.codec.Decode(decoded, encoded)
			if err != nil {
				t.Fatalf("Decode(%q): %v", encoded, err)
			}
			if !bytes.Equal(decoded[:n], test.plain) {
				t.Fatalf("Decode(%q) = %q, want %q", encoded, decoded[:n], test.plain)
			}
		})
	}
}

func TestEncodingIntegerVectors(t *testing.T) {
	tests := []struct {
		name      string
		codec     es.Encoding
		encoded32 string
		encoded64 string
	}{
		{"Crockford", es.Crockford, "00000AG", "000000000002M"},
		{"z-base-32", es.ZBase32, "yyyyyko", "yyyyyyyyyyynw"},
		{"word-safe", es.WordSafe, "22222GR", "222222222224c"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.codec.EncodeInt32(42); got != test.encoded32 {
				t.Errorf("EncodeInt32(42) = %q, want %q", got, test.encoded32)
			}
			if got := test.codec.EncodeInt64(42); got != test.encoded64 {
				t.Errorf("EncodeInt64(42) = %q, want %q", got, test.encoded64)
			}
		})
	}
}

func TestEncodingIntegerRoundTrip(t *testing.T) {
	codecs := []struct {
		name  string
		codec es.Encoding
	}{
		{"Crockford", es.Crockford},
		{"z-base-32", es.ZBase32},
		{"word-safe", es.WordSafe},
	}

	for _, test := range codecs {
		t.Run(test.name, func(t *testing.T) {
			for _, want := range []int32{math.MinInt32, -1, 0, 1, 42, math.MaxInt32} {
				got, err := test.codec.DecodeInt32(test.codec.EncodeInt32(want))
				if err != nil {
					t.Fatalf("DecodeInt32 for %d: %v", want, err)
				}
				if got != want {
					t.Errorf("DecodeInt32(EncodeInt32(%d)) = %d", want, got)
				}
			}

			for _, want := range []int64{math.MinInt64, -1, 0, 1, 42, math.MaxInt64} {
				got, err := test.codec.DecodeInt64(test.codec.EncodeInt64(want))
				if err != nil {
					t.Fatalf("DecodeInt64 for %d: %v", want, err)
				}
				if got != want {
					t.Errorf("DecodeInt64(EncodeInt64(%d)) = %d", want, got)
				}
			}
		})
	}
}

func TestDecodeIntegerRejectsWrongWidth(t *testing.T) {
	for _, test := range []struct {
		name  string
		codec es.Encoding
	}{
		{"Crockford", es.Crockford},
		{"z-base-32", es.ZBase32},
		{"word-safe", es.WordSafe},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got, err := test.codec.DecodeInt32(test.codec.EncodeToString([]byte{1})); err == nil {
				t.Errorf("DecodeInt32 returned %d, want width error", got)
			}
			if got, err := test.codec.DecodeInt64(test.codec.EncodeToString([]byte{1, 2, 3, 4})); err == nil {
				t.Errorf("DecodeInt64 returned %d, want width error", got)
			}
		})
	}
}
