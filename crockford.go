package es

const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// Crockford is an unpadded Crockford Base32 encoding. It emits uppercase
// letters. Decoding is case-insensitive, ignores hyphens, treats I and L as 1,
// and treats O as 0.
var Crockford Encoding = newBase32Encoding("Crockford Base32", crockfordAlphabet, normalizeCrockford, true)

func normalizeCrockford(c byte) (byte, bool) {
	if c >= 'a' && c <= 'z' {
		c -= 'a' - 'A'
	}
	switch c {
	case 'I', 'L':
		return '1', true
	case 'O':
		return '0', true
	}
	for i := range crockfordAlphabet {
		if crockfordAlphabet[i] == c {
			return c, true
		}
	}
	return 0, false
}
