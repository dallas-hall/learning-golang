package pipeline

import (
	"io"
	"strings"
)

type Pipeline struct {
	Input  io.Reader
	Output io.Writer
	Error  error
}

func FromString(s string) *Pipeline {
	return &Pipeline{
		Input: strings.NewReader(s),
	}
}

func (p *Pipeline) Stdout() {
	if p.Error != nil {
		return
	}
	// Copies to io.Writer destination from io.Reader source.
	io.Copy(p.Output, p.Input)
}
