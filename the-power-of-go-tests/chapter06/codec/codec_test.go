package codec_test

import (
	"codec"
	"math/rand/v2"
	"testing"
)

func Test_EncodeFollowedByDecodeGivesStartingValue(t *testing.T) {
	t.Parallel()
	// Create a random number between 1 and 10 (0 to 9 + 1)
	input := rand.IntN(10) + 1
	want := input
	encoded := codec.Encode(input)
	t.Logf("encoded value: %#v", encoded)
	got := codec.Decode(encoded)
	if want != got {
		t.Errorf("want %d & got %d", want, got)
	}
}
