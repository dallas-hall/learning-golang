package codec_test

import (
	"codec"
	"math/rand/v2"
	"testing"
)

func Test_EncodeFollowedByDecodeGivesStartingValue(t *testing.T) {
	t.Parallel()
	input := rand.IntN(10) + 1
	want := input
	encoded := codec.Encode(input)
	t.Logf("encoded value: %#v", encoded)
	got := codec.Decode(encoded)
	if want != got {
		t.Errorf("want %d & got %d", want, got)
	}
}
