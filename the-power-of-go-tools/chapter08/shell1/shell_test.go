package simpleshell1_test

import (
	"simpleshell1"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestShell_CmdFromStringReturnsExpectedCommand(t *testing.T) {
	t.Parallel()
	cmd, err := simpleshell1.CmdFromString("/usr/bin/ls -Adhl /\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/usr/bin/ls", "-Adhl", "/"}
	got := cmd.Args
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestShell_CmdFromStringErrorsOnEmptyInput(t *testing.T) {
	t.Parallel()
	_, err := simpleshell1.CmdFromString("")
	if err == nil {
		t.Fatal("want error on empty input & got nil")
	}
}
