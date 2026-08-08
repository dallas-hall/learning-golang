package square_test

import (
	"math/rand"
	"square"
	"strconv"
	"testing"
)

// This is a property based test. We don't know what the square of a random
// number will be, but we know it cannot be negative. So we test for that.
func Test_SquareGivesNonNegativeResult(t *testing.T) {
	t.Parallel()
	// Get randomly ordered slice of ints from 0 to 99
	inputs := rand.Perm(100)
	for _, n := range inputs {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			got := square.Square(n)
			if got < 0 {
				t.Errorf("Square(%d) is negative: %d", n, got)
			}
		})
	}
}
