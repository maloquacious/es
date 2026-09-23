package es

const wordSafeAlphabet = "23456789CFGHJMPQRVWXcfghjmpqrvwx"

// WordSafe is an unpadded word-safe Base32 encoding. Its mixed-case alphabet
// is case-sensitive because uppercase and lowercase letters represent
// different values.
var WordSafe Encoding = newBase32Encoding("word-safe Base32", wordSafeAlphabet, normalizeWordSafe, false)

func normalizeWordSafe(c byte) (byte, bool) {
	for i := range wordSafeAlphabet {
		if wordSafeAlphabet[i] == c {
			return c, true
		}
	}
	return 0, false
}
