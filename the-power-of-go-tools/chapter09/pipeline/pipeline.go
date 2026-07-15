package pipeline

import (
	"io"
	"log"
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

// This method got replaced by FromString, to remove a bug and use our new constructor.
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
	p := New()
	// Read file into something
	file, err := os.Open(path)
	if err != nil {
		log.Printf("failed to open %q: %s", path, err)
		return &Pipeline{Error: err}
	}
	// Don't use `defer file.Close()` because it will close the file on exiting this function.`

	p.Input = file
	return p
}

// Copies the Input to the Output to form the "shell pipeline."
// If the default hasn't been changed, this will be Stdout.
func (p *Pipeline) Stdout() {
	if p.Error != nil {
		return
	}
	// Copies to io.Writer destination from io.Reader source.
	io.Copy(p.Output, p.Input)
}
