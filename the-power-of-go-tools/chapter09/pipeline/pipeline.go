package pipeline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Allows for reading and writing in a pipeline, much like a shell pipeline.
// We can accept any io.Reader input. We can pass on any io.Writer output.
// Error is used to stop a pipline.
type Pipeline struct {
	Input  io.Reader
	Output io.Writer
	Error  error
}

// Creates and returns a new Pipeline object with Stdout as Output.
func New() *Pipeline {
	return &Pipeline{
		Output: os.Stdout,
	}
}

// This method got replaced by FromString. It removes a bug of no Output and
// uses our new constructor.
func FromStringOld(s string) *Pipeline {
	return &Pipeline{
		Input: strings.NewReader(s),
	}
}

// Create and return a Pipeline object with the Input set from a String.
// The Output is Stdout. Error is nil.
func FromString(s string) *Pipeline {
	p := New()
	p.Input = strings.NewReader(s)
	return p
}

// Create and return a Pipeline object with the Input set from a file.
// The Output is Stdout. Error is nil.
func FromFile(path string) *Pipeline {
	file, err := os.Open(path)
	if err != nil {
		return &Pipeline{Error: err}
	}
	// Don't use this because it will close the file on exiting this function.`
	// defer file.Close()

	p := New()
	p.Input = file
	return p
}

// Copies the Input to the Output to form the "shell pipeline."
// If the default hasn't been changed, this will be Stdout.
// We always check for a Pipeline Error first and do nothing is one exists.
func (p *Pipeline) Stdout() {
	if p.Error != nil {
		return
	}
	// Copies to io.Writer destination from io.Reader source.
	io.Copy(p.Output, p.Input)
}

// Read everything from the Pipeline's input and return it as a string.
// We always check for a Pipeline Error first and do nothing is one exists.
func (p *Pipeline) String() (string, error) {
	if p.Error != nil {
		return "", p.Error
	}

	data, err := io.ReadAll(p.Input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Read everything from the Pipeline's input, split the text using the delimiter
// and return the supplied column from each line.
// We always check for a Pipeline Error first and clear the reader if one exists.
func (p *Pipeline) Column(column int, delimiter string) *Pipeline {
	if p.Error != nil {
		p.Input = strings.NewReader("")
		return p
	}
	if column < 1 {
		p.Error = fmt.Errorf("bad column, must be positive: %d\n", column)
		return p
	}

	result := new(bytes.Buffer)
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), delimiter)
		if len(fields) < column {
			continue
		}
		fmt.Fprintln(result, fields[column-1])
	}

	return &Pipeline{
		Input: result,
	}
}
