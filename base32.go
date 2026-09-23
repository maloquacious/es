package es

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
)

// Encoding is a binary-to-text encoding scheme. Its API follows the shape of
// encoding/base64.Encoding.
//
// Encode and Decode require destination buffers large enough to hold the
// result. EncodedLen and DecodedLen calculate the required maximum sizes.
// The integer helpers encode fixed-width big-endian two's-complement values;
// their decoders reject input that does not represent exactly 4 or 8 bytes.
type Encoding interface {
	Encode(dst, src []byte)
	EncodeToString(src []byte) string
	EncodeInt32(src int32) string
	EncodeInt64(src int64) string
	EncodedLen(n int) int
	Decode(dst, src []byte) (n int, err error)
	DecodeString(src string) ([]byte, error)
	DecodeInt32(src string) (int32, error)
	DecodeInt64(src string) (int64, error)
	DecodedLen(n int) int
}

type base32Encoding struct {
	name       string
	alphabet   string
	encoding   *base32.Encoding
	normalize  func(byte) (byte, bool)
	skipHyphen bool
}

func newBase32Encoding(name, alphabet string, normalize func(byte) (byte, bool), skipHyphen bool) *base32Encoding {
	return &base32Encoding{
		name:       name,
		alphabet:   alphabet,
		encoding:   base32.NewEncoding(alphabet).WithPadding(base32.NoPadding),
		normalize:  normalize,
		skipHyphen: skipHyphen,
	}
}

func (enc *base32Encoding) Encode(dst, src []byte) {
	enc.encoding.Encode(dst, src)
}

func (enc *base32Encoding) EncodeToString(src []byte) string {
	return enc.encoding.EncodeToString(src)
}

func (enc *base32Encoding) EncodeInt32(src int32) string {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(src))
	return enc.EncodeToString(buf[:])
}

func (enc *base32Encoding) EncodeInt64(src int64) string {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(src))
	return enc.EncodeToString(buf[:])
}

func (enc *base32Encoding) EncodedLen(n int) int {
	return enc.encoding.EncodedLen(n)
}

func (enc *base32Encoding) Decode(dst, src []byte) (int, error) {
	normalized := make([]byte, 0, len(src))
	positions := make([]int, 0, len(src))
	for i := 0; i < len(src); i++ {
		if enc.skipHyphen && src[i] == '-' {
			continue
		}

		c, ok := enc.normalize(src[i])
		if !ok {
			return 0, fmt.Errorf("es: invalid %s character %q at byte %d", enc.name, src[i], i)
		}
		normalized = append(normalized, c)
		positions = append(positions, i)
	}

	if err := enc.validateTail(normalized, positions); err != nil {
		return 0, err
	}

	n, err := enc.encoding.Decode(dst, normalized)
	if err != nil {
		return n, fmt.Errorf("es: invalid %s encoding: %w", enc.name, err)
	}
	return n, nil
}

func (enc *base32Encoding) DecodeString(src string) ([]byte, error) {
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, []byte(src))
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

func (enc *base32Encoding) DecodeInt32(src string) (int32, error) {
	decoded, err := enc.DecodeString(src)
	if err != nil {
		return 0, err
	}
	if len(decoded) != 4 {
		return 0, fmt.Errorf("es: decoded %s Int32 has %d bytes, want 4", enc.name, len(decoded))
	}
	return int32(binary.BigEndian.Uint32(decoded)), nil
}

func (enc *base32Encoding) DecodeInt64(src string) (int64, error) {
	decoded, err := enc.DecodeString(src)
	if err != nil {
		return 0, err
	}
	if len(decoded) != 8 {
		return 0, fmt.Errorf("es: decoded %s Int64 has %d bytes, want 8", enc.name, len(decoded))
	}
	return int64(binary.BigEndian.Uint64(decoded)), nil
}

func (enc *base32Encoding) DecodedLen(n int) int {
	return enc.encoding.DecodedLen(n)
}

func (enc *base32Encoding) validateTail(src []byte, positions []int) error {
	var mask byte
	switch len(src) % 8 {
	case 0:
		return nil
	case 2:
		mask = 0x03
	case 4:
		mask = 0x0f
	case 5:
		mask = 0x01
	case 7:
		mask = 0x07
	default:
		return fmt.Errorf("es: invalid %s length %d", enc.name, len(src))
	}

	value := byte(0)
	for i := range enc.alphabet {
		if enc.alphabet[i] == src[len(src)-1] {
			value = byte(i)
			break
		}
	}
	if value&mask != 0 {
		return fmt.Errorf("es: non-zero %s padding bits at byte %d", enc.name, positions[len(positions)-1])
	}
	return nil
}
