package simpleshell2_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"simpleshell2"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestShell_CmdFromStringReturnsExpectedCommand(t *testing.T) {
	t.Parallel()
	cmd, err := simpleshell2.CmdFromString("/usr/bin/ls -Adhl /\n")
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
	_, err := simpleshell2.CmdFromString("")
	if err == nil {
		t.Fatal("want error on empty input & got nil")
	}
}

func TestShell_SuccessfullyCreateNewShellSession(t *testing.T) {
	t.Parallel()

	want := simpleshell2.Session{
		Input:  os.Stdin,
		Output: os.Stdout,
		Error:  os.Stderr,
	}

	got := *simpleshell2.NewSession(os.Stdin, os.Stdout, os.Stderr)

	if want != got {
		t.Errorf("want %v & got %v", want, got)
	}
}

func TestShell_SuccessfullyRunNewShellWithInputPassedToCommandAndCheckOutput(t *testing.T) {
	t.Parallel()

	s := "echo 1 2 3"
	// Create an io.Reader with our command to pass into the shell
	inputBuffer := strings.NewReader(s)
	// Create an io.Writer to capture the output from the shell
	outputBuffer := new(bytes.Buffer)

	// Create a new shell with our custom io.Reader and io.Writers.
	session := simpleshell2.NewSession(inputBuffer, outputBuffer, io.Discard)
	session.Run()

	// We see the commands output as it gets executed. But this is brittle, what is `echo` doesn't exist inside of a scratch container running the test?
	// We don't see the inputted command though, e.g. no `echo 1 2 3` - even though it got executed?
	prompt := simpleshell2.GetPrompt()
	want := fmt.Sprintf("%s1 2 3\n%s\nGoodbye\n", prompt, prompt)
	got := outputBuffer.String()

	if want != got {
		// Fatalf stops this test function right there, unlike Errorf which would let it continue.
		t.Fatalf("want: %q & got: %q", want, got)
	}
}

func TestShell_SuccessfullyRunNewShellWithInputThatsPassedToOutput(t *testing.T) {
	t.Parallel()

	s := "echo 1 2 3"
	inputBuffer := strings.NewReader(s)
	outputBuffer := new(bytes.Buffer)

	session := simpleshell2.NewSession(inputBuffer, outputBuffer, io.Discard)
	session.DryRun = true
	session.Run()

	// We see the entire output, even though it wasn't executed.
	prompt := simpleshell2.GetPrompt()
	want := fmt.Sprintf("%s%s\n1 2 3\n%s\nGoodbye\n", prompt, s, prompt)
	got := outputBuffer.String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestShell_SuccessfullyRunNewShellWithInputThatsSavedToFile(t *testing.T) {
	t.Parallel()
	data := "echo hello world 123"
	inputBuffer := strings.NewReader(data)
	outputBuffer := new(bytes.Buffer)

	session := simpleshell2.NewSession(inputBuffer, outputBuffer, io.Discard)
	session.Transcript = true
	session.Run()

	path := "test/data/transcript.txt"
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to open %q because: %s", path, err)
	}

	want := []byte(data)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}

	// Check permissions are correct. We hardcoded 0600 inside of WriteToFile
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	permission := stat.Mode().Perm()
	if permission != 0o600 {
		t.Errorf("want: file permission 0o600 & got: file permission 0o%o", permission)
	}
}

// Taken from the-power-of-go-tools/chapter04/count-pflag/count_test.go - see for comments
func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "test/scripts",
	})
}

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"shell": simpleshell2.Main,
	})
}
