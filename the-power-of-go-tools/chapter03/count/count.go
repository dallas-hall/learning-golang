package count

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Use a struct to provide concurrent friendly variables that can be used for default values.
//
// The `files` field holds onto each individual opened file as a separate io.Reader, and `input` eventually holds all of them combined into one reader.
// Files need to be remembered separately so they can be closed later. You can't call Close() on the result of io.MultiReader, only on the original *os.File values.
//
// counter is lowercase — unexported. Code outside this package can't refer to the type counter directly; it only ever sees *counter values returned from NewCounter, and only ever touches them through exported methods like Lines(). The struct itself is private; the capability is public.
type counter struct {
	input  io.Reader
	output io.Writer
	files  []io.Reader
}

// Creating an option type that is just a function that sets something inside the counter struct.
type option func(*counter) error

// Users can only set our counter struct fields by calling validating functions in the constructor arguments.
func WithInput(input io.Reader) option {
	return func(c *counter) error {
		if input == nil {
			return errors.New("nil input reader")
		}
		c.input = input
		return nil
	}
}

func WithInputFromArgs(args []string) option {
	return func(c *counter) error {
		// If no user args are passed in, we just use the default c.input which is os.Stdin
		if len(args) < 1 {
			return nil
		}

		// Create a new slice for all the passed in paths.
		c.files = make([]io.Reader, len(args))

		for i, path := range args {
			// os.Open returns an *os.File, which satisfies io.Reader and io.Closer in Main()
			f, err := os.Open(path)
			// Any error opening a file will break everything.
			if err != nil {
				return err
			}
			c.files[i] = f
		}
		// Once all files are opened, io.MultiReader glues them all together into a single io.Reader that reads through them in entirely and in sequence.
		c.input = io.MultiReader(c.files...)
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(c *counter) error {
		if output == nil {
			return errors.New("nil output writer")
		}
		c.output = output
		return nil
	}
}

/*
Add concurrent friendly default values in our constructor. This function is variadic because if accepts `...type`, which means it can accept 0 to many number of arguments.

We created the `option` type which allows us to validate the user provided options.

We provide the user with validating functions (eg WithInput) that they can use to pass in options to the constructor. eg:

c, err := count.NewCounter(

	count.WithInput(inputBuffer),

)

The pattern is called functional options, since we are passing in functions to the constructor.
*/
func NewCounter(opts ...option) (*counter, error) {
	c := &counter{
		input:  os.Stdin,
		output: os.Stdout,
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

// The line counting logic.
func (c *counter) Lines() int {
	lines := 0
	// Accepts user input until EOF from ^D. It wraps the reader in a Scanner for line-based reading.
	input := bufio.NewScanner(c.input)
	// Keep reading lines until Scan() returns false (EOF or error), incrementing the counter each time.
	for input.Scan() {
		lines++
	}

	for _, f := range c.files {
		// After scanning is done, loop over every individually-opened file and close it.
		// We are using type assertion f.(io.Closer) because io.Reader doesn't have a Close()
		f.(io.Closer).Close()
	}
	return lines
}

// Convenience wrapper
func Main() {
	c, err := NewCounter(
		// Read all args and pass them along, no user args is handled inside.
		// The usual Go idiom for "just the arguments the user actually typed"
		WithInputFromArgs(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(c.Lines())
}
