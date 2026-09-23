# Word-safe Base32

`es.WordSafe` implements the
[word-safe Base32 alphabet](https://en.wikipedia.org/wiki/Base32#Word-safe_alphabet)
as an unpadded, byte-oriented encoding.

```go
var WordSafe Encoding
```

## Alphabet

The encoding alphabet, in value order from 0 through 31, is:

```text
23456789CFGHJMPQRVWXcfghjmpqrvwx
```

The alphabet extends the Open Location Code Base20 alphabet and uses 8 digits,
12 uppercase letters, and 12 lowercase letters.

## Case sensitivity

Word-safe Base32 is case-sensitive. Uppercase and lowercase forms of the same
letter have different values. For example, `J` has value 12 and `j` has value
24. Decoding must preserve the original letter case.

## Encoding

```go
encoded := es.WordSafe.EncodeToString([]byte("hello"))
// encoded == "M3Wgjq5Q"
```

Output contains no padding.

Integer helpers use fixed-width big-endian representations:

```go
encoded32 := es.WordSafe.EncodeInt32(42) // "22222GR"
encoded64 := es.WordSafe.EncodeInt64(42) // "222222222224c"
```

## Decoding

The decoder accepts only exact alphabet symbols. It does not normalize case or
accept separators and padding. It also rejects invalid unpadded lengths and
non-zero trailing padding bits.

```go
decoded, err := es.WordSafe.DecodeString("M3Wgjq5Q")
// decoded == []byte("hello")
```

## API

`es.WordSafe` implements every method of [`es.Encoding`](encoding.md),
including buffer, string, and fixed-width integer methods.
