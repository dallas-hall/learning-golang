package game_test

import (
	"game"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// This only works for items >= 2
func Test_MyListItemsGivesCorrectResultForInput(t *testing.T) {
	t.Parallel()
	input := []string{
		"a battery",
		"a key",
		"a tourist map",
	}
	want := "You can see here a battery, a key, and a tourist map."
	got := game.MyListItems(input)
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

// Instead of writing multiple tests for different inputs, we use one test that
// loops over all the different inputs.
func Test_ListItemsGivesCorrectResultForInput(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input []string
		want  string
	}
	cases := []testCase{
		{
			input: []string{},
			want:  "",
		},
		{
			input: []string{
				"a battery",
			},
			want: "You can see a battery here.",
		},
		{
			input: []string{
				"a battery",
				"a key",
			},
			want: "You can see here a battery and a key.",
		},
		{
			input: []string{
				"a battery",
				"a key",
				"a tourist map",
			},
			want: "You can see here a battery, a key, and a tourist map.",
		},
	}
	for _, tc := range cases {
		got := game.ListItems(tc.input)
		if tc.want != got {
			t.Error(cmp.Diff(tc.want, got))
		}
	}
}
