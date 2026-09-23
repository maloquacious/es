# Encoding interface

`Encoding` is the common API implemented by every binary-to-text encoding in
package `es`.

```go
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
```

The interface follows the buffer and string API shape of
[`encoding/base64.Encoding`](https://pkg.go.dev/encoding/base64#Encoding).

## Implementations

| Value | Encoding | Canonical output |
| --- | --- | --- |
| `es.Crockford` | [Crockford Base32](crockford.md) | Uppercase, unpadded |
| `es.ZBase32` | [z-base-32](z-base-32.md) | Lowercase, unpadded |
| `es.WordSafe` | [Word-safe Base32](word-safe.md) | Mixed case, unpadded |

## Byte methods

### `Encode`

```go
Encode(dst, src []byte)
```

Encodes `src` into `dst`. `dst` must have at least `EncodedLen(len(src))`
bytes. The method does not return the number of bytes written because the
encoded size is deterministic.

```go
dst := make([]byte, es.ZBase32.EncodedLen(len(src)))
es.ZBase32.Encode(dst, src)
```

### `EncodeToString`

```go
EncodeToString(src []byte) string
```

Returns the canonical, unpadded encoding of `src`.

### `EncodedLen`

```go
EncodedLen(n int) int
```

Returns the number of encoded bytes produced from `n` source bytes. For the
Base32 implementations, the result is `ceil(8n / 5)`.

### `Decode`

```go
Decode(dst, src []byte) (n int, err error)
```

Decodes `src` into `dst` and returns the number of bytes written. Allocate
`dst` with `DecodedLen(len(src))` bytes. That size is an upper bound and may be
larger than `n`, particularly when Crockford input contains ignored hyphens.

```go
dst := make([]byte, es.ZBase32.DecodedLen(len(src)))
n, err := es.ZBase32.Decode(dst, src)
decoded := dst[:n]
```

On validation errors, `n` is zero. Validation errors include invalid symbols,
invalid unpadded lengths, and non-zero trailing padding bits.

### `DecodeString`

```go
DecodeString(src string) ([]byte, error)
```

Decodes `src` into a newly allocated byte slice. It returns `nil` with an error
when validation fails.

### `DecodedLen`

```go
DecodedLen(n int) int
```

Returns the maximum number of decoded bytes for `n` encoded bytes. Use the
returned value to size a destination passed to `Decode`.

## Integer methods

The integer methods encode fixed-width, big-endian, two's-complement values.
They operate on the same byte representation as the byte methods:

```go
es.ZBase32.EncodeInt32(v)
```

is equivalent to encoding the four-byte big-endian representation of `v`.
Consequently, an encoded `int32` always contains 7 Base32 symbols, and an
encoded `int64` always contains 13.

### `EncodeInt32` and `EncodeInt64`

```go
EncodeInt32(src int32) string
EncodeInt64(src int64) string
```

Return the canonical encoding of the integer's 4-byte or 8-byte representation.
Negative values preserve their two's-complement bit patterns.

### `DecodeInt32` and `DecodeInt64`

```go
DecodeInt32(src string) (int32, error)
DecodeInt64(src string) (int64, error)
```

Decode an integer encoded by the corresponding encode method. In addition to
the selected alphabet's normal validation, these methods require the decoded
value to contain exactly 4 or 8 bytes. They reject shorter and longer byte
representations rather than padding or truncating them.

## Unpadded input rules

All implementations are byte-oriented and omit `=` padding. After any
encoding-specific normalization, a valid input length modulo 8 is one of:

```text
0, 2, 4, 5, 7
```

The unused low bits of the final symbol must be zero. This requirement prevents
multiple textual encodings from decoding to the same byte sequence.
