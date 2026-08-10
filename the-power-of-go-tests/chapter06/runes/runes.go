package runes

import "unicode/utf8"

func FirstRune(s string) rune {
	// rune(s[0]) will return the first byte of s, but UTF-8 can have multiple
	// bytes so this will only work for the first 255 characters in UTF-8.
	// The range operator iterates over a string by runes, not bytes.
	for _, r := range s {
		return r
	}
	return utf8.RuneError
}
