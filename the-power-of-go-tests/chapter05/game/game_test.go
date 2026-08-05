package game_test

import (
	"game"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

// Using a map with t.Run and the subtest function will print the name of any
// subtest that is failing. func(t *testing.T) { ... } is the subtest. Spaces
// are replaced with underscore.
// We excluded a word after For because a failure will look like:
// FAIL: Test_ListItemsGivesCorrectResultFor/two_items
func Test_ListItemsGivesCorrectResultFor(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input []string
		want  string
	}{
		// The map's key is often the name of the test used by the subtest function.
		"no items": {
			input: []string{},
			want:  "",
		},
		"one item": {
			input: []string{"a battery"},
			want:  "You can see a battery here.",
		},
		"two items": {
			input: []string{
				"a battery",
				"a key",
			},
			want: "You can see here a battery and a key.",
		},
		"three items": {
			input: []string{
				"a battery",
				"a key",
				"a tourist map",
			},
			want: "You can see here a battery, a key, and a tourist map.",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := game.ListItems(tc.input)
			if tc.want != got {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}
