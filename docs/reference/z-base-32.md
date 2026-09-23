# z-base-32

`es.ZBase32` implements the
[z-base-32](https://philzimmermann.com/docs/human-oriented-base-32-encoding.txt)
alphabet as an unpadded, byte-oriented encoding.

```go
var ZBase32 Encoding
```

## Alphabet

The canonical encoding alphabet, in value order from 0 through 31, is:

```text
ybndrfg8ejkmcpqxot1uwisza345h769
```

Canonical output uses lowercase letters.

## Encoding

```go
encoded := es.ZBase32.EncodeToString([]byte("hello"))
// encoded == "pb1sa5dx"
```

Output contains no padding.

Integer helpers use fixed-width big-endian representations:

```go
encoded32 := es.ZBase32.EncodeInt32(42) // "yyyyyko"
encoded64 := es.ZBase32.EncodeInt64(42) // "yyyyyyyyyyynw"
```

## Decoding

Letters are case-insensitive during decoding. For example, both strings below
decode to `hello`:

```text
pb1sa5dx
PB1SA5DX
```

The decoder does not accept separators or padding. It rejects symbols outside
the alphabet, invalid unpadded lengths, and non-zero trailing padding bits.

## Bit-length scope

The original z-base-32 design can represent bitstrings whose lengths are not
multiples of 8 when their exact bit lengths are known by both parties. This
implementation is byte-oriented: encoding accepts whole bytes, and decoding
accepts only symbol lengths that map unambiguously to whole bytes.

## API

`es.ZBase32` implements every method of [`es.Encoding`](encoding.md), including
buffer, string, and fixed-width integer methods.
