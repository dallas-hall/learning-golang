package count

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Use a struct to provide concurrent friendly variables that can be used for default values.
type counter struct {
	input  io.Reader
	output io.Writer
}

// Creating an option type that is just a function that sets something inside the counter struct.
type option func(*counter) error

// Users can only set our counter struct fields by calling validating functions in the constructor arguments.
func WithInput(input io.Reader) option {
	return func(c *counter) error {
		if input == nil {
			return errors.New("nil input reaader")
		}
		c.input = input
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
	// Accepts user input until EOF from ^D
	input := bufio.NewScanner(c.input)
	for input.Scan() {
		lines++
	}
	return lines
}

// Convenience wrapper
func Main() {
	c, err := NewCounter()
	if err != nil {
		// You typically want to return the error to the caller. But since we are using Main() as an extention of main(), we are deciding to panic instead.
		panic(err)
	}
	fmt.Println(c.Lines())
}
