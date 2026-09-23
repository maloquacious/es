# Crockford Base32

`es.Crockford` implements an unpadded, byte-oriented form of
[Crockford's Base32](https://www.crockford.com/base32.html).

```go
var Crockford Encoding
```

## Alphabet

The canonical encoding alphabet, in value order from 0 through 31, is:

```text
0123456789ABCDEFGHJKMNPQRSTVWXYZ
```

Canonical output uses uppercase letters. The alphabet omits `I`, `L`, `O`, and
`U`.

## Encoding

```go
encoded := es.Crockford.EncodeToString([]byte("foobar"))
// encoded == "CSQPYRK1E8"
```

Output contains no separators or padding.

Integer helpers use fixed-width big-endian representations:

```go
encoded32 := es.Crockford.EncodeInt32(42) // "00000AG"
encoded64 := es.Crockford.EncodeInt64(42) // "000000000002M"
```

## Decoding

Decoding applies the following normalization rules before validating the
Base32 data:

| Input | Decoded as |
| --- | --- |
| `a`–`z` | Corresponding uppercase letter |
| `I`, `i`, `L`, `l` | `1` |
| `O`, `o` | `0` |
| `-` | Ignored |

For example, each input below decodes to `foobar`:

```text
CSQPYRK1E8
csqpyrk1e8
CSQPY-RK1E8
csqpyrkle8
csqpyrkie8
```

Other characters are rejected. The decoder also rejects invalid unpadded
lengths and non-zero trailing padding bits.

## Check symbols

Crockford's optional modulo-37 check symbol is not implemented. The additional
check-symbol characters `*`, `~`, `$`, `=`, and `U` are invalid input.

## API

`es.Crockford` implements every method of [`es.Encoding`](encoding.md),
including buffer, string, and fixed-width integer methods.
