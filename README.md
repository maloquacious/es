# es

Simple binary-to-text encoding schemes for Go.

The package supports unpadded byte-oriented encoders and decoders for:

- [Crockford's Base32](https://www.crockford.com/base32.html)
- [z-base-32](https://philzimmermann.com/docs/human-oriented-base-32-encoding.txt)

```go
encoded := es.Crockford.EncodeToString([]byte("foobar")) // "CSQPYRK1E8"
decoded, err := es.Crockford.DecodeString(encoded)

zEncoded := es.ZBase32.EncodeToString([]byte("hello")) // "pb1sa5dx"
zDecoded, err := es.ZBase32.DecodeString(zEncoded)
```

Both values implement `es.Encoding`, whose buffer, string, and length methods
follow the API shape of Go's `encoding/base64.Encoding`.

Crockford decoding is case-insensitive, ignores hyphens, maps `I` and `L` to
`1`, and maps `O` to `0`. z-base-32 decoding is case-insensitive. Both decoders
reject padding, invalid lengths, and non-zero trailing padding bits.
