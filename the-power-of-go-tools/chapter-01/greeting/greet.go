// The separate package makes this testable.
package greeting

import (
	"bufio" // Buffered I/O helpers
	"fmt"   // Formatters
	"io"    // I/O interfaces
)

// This function accepts the interfaces io.Reader and io.Writer as parameters. That means a test can pass in a strings.Reader and a bytes.Buffer instead of the real terminal, and check what got written without any actual I/O happening.
func GreetUser(stdin io.Reader, stdout io.Writer) {
	name := "you"

	// Writes the prompt followed by a newline to whatever stdout is — could be the terminal, could be an in-memory buffer in a test
	fmt.Fprintln(stdout, "What is your name?")

	// Wraps the stdin reader in a Scanner, which by default splits input into lines.
	input := bufio.NewScanner(stdin)

	// Advances to the next token (a line, by default) and returns true if it found one, false if it hit EOF or an error. Using it inside an if rather than the more familiar for input.Scan() loop is deliberate: this function only ever wants one line of input, not the whole stream. If Scan() returns false — e.g., the user piped in an empty file, or just hit Ctrl-D with no input — the body is skipped and name keeps its "you" default.
	if input.Scan() {
		// Returns the line Scan() just read, as a string, with the trailing newline stripped.
		name = input.Text()
	}

	fmt.Fprintf(stdout, "Hello %s!\n", name)
}
