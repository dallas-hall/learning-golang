package match

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// This struct holds the input and output mechanisms.
type matcher struct {
	input  io.Reader
	output io.Writer
	text   string
}

// Use the functional options pattern which allows users to pass in options that will be validated by functions we provide.
type option func(*matcher) error

func WithInput(input io.Reader) option {
	return func(m *matcher) error {
		if input == nil {
			return errors.New("nil input reaader")
		}
		m.input = input
		return nil
	}
}

func WithSearchStringFromArgs(args []string) option {
	return func(m *matcher) error {
		// Use default search text if nothing is supplied.
		if len(args) < 1 {
			return nil
		}

		// Use a string builder to efficiently concatenate multiple arguments to a string.
		var sb strings.Builder
		for _, arg := range args {
			sb.WriteString(arg + "\n")
		}
		m.text = sb.String()
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(m *matcher) error {
		if output == nil {
			return errors.New("nil output writer")
		}
		m.output = output
		return nil
	}
}

// This variadic constructor allows for 0 to many options being passed in. Passed in options will be validated. Omitted options get given defaults.
func NewMatcher(opts ...option) (*matcher, error) {
	c := &matcher{
		input:  os.Stdin,
		output: os.Stdout,
		text:   "The quick brown fox jumps over the lazy dog\nQuick nymph bugs vex fjord waltz.\nWaltz, bad nymph, for quick jigs vex.\nSphinx of black quartz, judge my vow.\nHow vexingly quick daft zebras jump!\nThe five boxing wizards jump quickly.\nPack my box with five dozen liquor jugs.",
	}

	// Handle the 0 to many options passed in
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (m *matcher) PrintMatchingLines() {
	// Accepts user input until EOF from ^D
	match := "quick"
	fmt.Fprintf(m.output, "Match: ")

	input := bufio.NewScanner(m.input)
	// Accept one line of input only with newline stripped. We have a backup if user passes nothing.
	if input.Scan() {
		match = input.Text()
	}

	// https://stackoverflow.com/a/33162487
	// Split a string on the newline, remove the newline, and iterate over all lines.
	for _, line := range strings.Split(strings.TrimSuffix(m.text, "\n"), "\n") {
		if strings.Contains(line, match) {
			fmt.Fprintf(m.output, "Matched: %q\n", line)
		}
	}
}

// Convenience wrapper
func Main() {
	m, err := NewMatcher(
		WithSearchStringFromArgs(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	m.PrintMatchingLines()
}
