// This is the external test package convention — the file lives in the same directory as greet.go but compiles as a separate package that imports greeting like any other consumer would.
package greeting_test

import (
	"bytes"
	"errors"
	"greeting"
	"testing"
	"testing/iotest" // utilities for constructing readers that misbehave in specific ways, useful for testing error paths.
)

func TestGreetUser_PrintsHelloMessageAndUsersNameToGivenWriter(t *testing.T) {
	t.Parallel()

	// Creates a *bytes.Buffer pre-loaded with the string "Bob". A *bytes.Buffer satisfies io.Reader (you can read bytes out of it) — this is the fake stdin standing in for a real terminal.
	input := bytes.NewBufferString("Bob")

	// Creates an empty *bytes.Buffer to act as the fake stdout. new(T) allocates a zero-valued T and returns a pointer to it — equivalent to &bytes.Buffer{}. A *bytes.Buffer also satisfies io.Writer, so writes land in memory instead of the terminal.
	output := new(bytes.Buffer)

	// Calls the function under test with the fakes. This only works because GreetUser was written against the io.Reader/io.Writer interfaces rather than concrete OS types.
	greeting.GreetUser(input, output)

	want := "What is your name?\nHello Bob!\n"
	// Reads everything that was written to the buffer back out as a string.
	got := output.String()

	if want != got {
		// Fatalf stops this test function right there, unlike Errorf which would let it continue.
		t.Fatalf("want: %q & got: %q", want, got)
	}
}

func TestGreetUser_PrintsHelloYouOnReadError(t *testing.T) {
	t.Parallel()

	// iotest.ErrReader builds an io.Reader whose Read method does nothing but immediately return the given error. So this isn't simulating "no input" (EOF) — it's simulating an actual I/O failure on every read attempt.
	input := iotest.ErrReader(errors.New("bad reader"))
	output := new(bytes.Buffer)
	greeting.GreetUser(input, output)

	// Since input.Scan() inside GreetUser will return false when the underlying read errors out, the if block in GreetUser never assigns a new name, and the default "you" survives — so the test expects the fallback greeting.
	want := "What is your name?\nHello you!\n"
	got := output.String()

	if want != got {
		t.Fatalf("want: %q & got: %q", want, got)
	}
}
