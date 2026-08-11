package runes_test

import (
	"runes"
	"testing"
	"unicode/utf8"
)

// Run with `go test -fuzz . -fuzztime=1m`
func FuzzFirstRune(f *testing.F) {
	// Seed the fuzzer's training data corpus with known good inputs. The fuzzer
	// will mutate these to generate new inputs.
	f.Add("Hello")
	f.Add("world")
	f.Add("안녕하세요")
	f.Add("🫠")
	// The inner function is called once per seed, and then repeatedly with the
	// mutated s values. The type in f.Add must match the second type within the
	// inner function.
	f.Fuzz(func(t *testing.T, s string) {
		got := runes.FirstRune(s)
		want, _ := utf8.DecodeRuneInString(s)
		// Apparently size == 0 is a better check, as utf8.RuneError is U+FFFD which
		// is a valid run.
		if want == utf8.RuneError {
			// Ignore invalid UTF-8 runes, there will be plenty from the fuzzer.
			t.Skip()
		}
		if want != got {
			// %q uses s. The cursor moves +1 to want.
			// %[1]x points to arg 1, so uses s. The cursor moves +1 to want.
			// %c uses want. The cursor is pointing to nothing, no args left.
			// %[2]x points to arg 2, so uses want.
			t.Errorf("given %q (0x%[1]x): want '%c' (0x%[2]x)", s, want)
			t.Errorf("got '%c' (0x%[1]x)", got)
			// The above is equivalent to:
			// t.Errorf("given %q (bytes: % x): want %q (0x%x), got %q (0x%x)", s, s, want, want, got, got)
		}
	})
}
