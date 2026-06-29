package count_test

import (
	"bytes"
	"count"
	"testing"

	// go get github.com/rogpeppe/go-internal/testscript
	"github.com/rogpeppe/go-internal/testscript"
)

func TestLineCounter_SuccessfullyCountsLinesFromStdin(t *testing.T) {
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

func TestLineCounter_SuccessfullyCountsLinesFromFile(t *testing.T) {
	t.Parallel()

	args := []string{"test/data/3-lines.txt"}

	c, err := count.NewCounter(
		count.WithInputFromArgs(args),
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

func TestLineCounter_IgnoresEmptyArgs(t *testing.T) {
	t.Parallel()

	inputBuffer := bytes.NewBufferString("1\n2\n3")
	c, err := count.NewCounter(
		count.WithInput(inputBuffer),
		// Simulate passing in no arguments
		count.WithInputFromArgs([]string{}),
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

// Replace TestLineCounter_SuccessfullyCountsLinesFromFile with a testscript
func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "test/scripts",
	})
}

// This function is always executed first before any tests are run.
// It sets up things that are needed for tests to run.
// We are associating the `count` command with our count.Main function, so it can be called in txtar testscripts.
// The associated command gets executed in a subprocess as an independent binary, just like if it had been compiled.
func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"count": count.Main,
	})
}
