package runes

import "unicode/utf8"

func FirstRune(s string) rune {
	// `return rune(s[0])`` returns the first byte of s, but UTF-8 can have
	// multiple bytes so this only works for the first 255 characters.
	// The range operator iterates over a string by runes, not bytes.
	for _, r := range s {
		return r
	}
	return utf8.RuneError
}
