package guess_test

import (
	"guess"
	"testing"
)

// Has to be run with `go test -fuzz .` or `go test -fuzz=FuzzGuess`
// The fuzzer t.Fuzz will randomly generate the inputs (eg input int), to the
// passed in testing function and execute the testing function with the inputs.
func FuzzGuess(f *testing.F) {
	f.Fuzz(func(t *testing.T, input int) {
		guess.Guess(input)
	})
}
