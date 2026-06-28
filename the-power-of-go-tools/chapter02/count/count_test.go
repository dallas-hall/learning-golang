package count_test

import (
	"bytes"
	"count"
	"testing"
)

func TestLineCounter_SuccessfullyCountsLines(t *testing.T) {
	t.Parallel()

	// Test data with 3 lines
	inputBuffer := bytes.NewBufferString("1\n2\n3")

	// Call our fancy variadic constructor that accepts 0 to many options, and those options are validaed by the functions we provided in count.go
	c, err := count.NewCounter(
		count.WithInput(inputBuffer),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := 3
	got := c.Lines()
	if want != got {
		t.Errorf("want: %d & got: %d", want, got)
	}
}
