package es_test

import (
	"bytes"
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
