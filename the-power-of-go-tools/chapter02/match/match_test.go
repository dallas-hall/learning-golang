package match_test

import (
	"bytes"
	"match"
	"testing"
)

func TestPrintMatchingLines_PrintsMatchingLinesFromDefaultSearchText(t *testing.T) {
	t.Parallel()

	output := new(bytes.Buffer)
	input := bytes.NewBufferString("fox")

	m, err := match.NewMatcher(
		match.WithInput(input),
		match.WithOutput(output),
	)
	if err != nil {
		t.Fatal(err)
	}

	// search input string for match string, return bool for found
	m.PrintMatchingLines()
	got := output.String()
	// We need to add "Match: " here because that gets printed when asking the user for input.
	want := "Match: Matched: \"The quick brown fox jumps over the lazy dog\"\n"
	if want != got {
		t.Errorf("want: %q & got: %q", want, got)
	}
}

func TestPrintMatchingLines_PrintsMatchingLinesFromPassedInSearchText(t *testing.T) {
	t.Parallel()

	output := new(bytes.Buffer)
	input := bytes.NewBufferString("qui")

	m, err := match.NewMatcher(
		match.WithInput(input),
		match.WithOutput(output),
		match.WithSearchStringFromArgs(
			[]string{
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
				"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
				"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
				"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			}),
	)
	if err != nil {
		t.Fatal(err)
	}

	// search input string for match string, return bool for found
	m.PrintMatchingLines()
	got := output.String()
	// We need to add "Match: " here because that gets printed when asking the user for input.
	want := "Match: Matched: \"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\"\nMatched: \"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\"\n"
	if want != got {
		t.Errorf("want: %q & got: %q", want, got)
	}
}
