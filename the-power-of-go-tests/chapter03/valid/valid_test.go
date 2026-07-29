package valid_test

import (
	"testing"
	"valid"

	"github.com/google/go-cmp/cmp"
)

// The point of exercise was to ensure testing the expected behvaiour of the
// function. And also to use ACE with names: Action, Condition, Expectation.

func Test_ValidReturnsTrueForValidInput(t *testing.T) {
	t.Parallel()

	want := true

	input := "valid input"
	got := valid.Valid(input)

	if !cmp.Equal(want, got) {
		t.Errorf("want: %v & got: %v", want, got)
	}
}

func Test_ValidReturnsFalseForInvalidInput(t *testing.T) {
	t.Parallel()

	want := false

	input := "invalid input"
	got := valid.Valid(input)

	if !cmp.Equal(want, got) {
		t.Errorf("want: %v & got: %v", want, got)
	}
}
