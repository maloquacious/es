package es

const zBase32Alphabet = "ybndrfg8ejkmcpqxot1uwisza345h769"

// ZBase32 is an unpadded z-base-32 encoding. It emits lowercase letters and
// accepts letters in either case when decoding.
var ZBase32 Encoding = newBase32Encoding("z-base-32", zBase32Alphabet, normalizeZBase32, false)

func normalizeZBase32(c byte) (byte, bool) {
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	for i := range zBase32Alphabet {
		if zBase32Alphabet[i] == c {
			return c, true
		}
	}
	return 0, false
}
