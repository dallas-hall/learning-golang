package pipeline_test

import (
	"bytes"
	"errors"
	"pipeline"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriter_StdoutPrintsToOutput(t *testing.T) {
	t.Parallel()
	want := "Helloworld 123\n"
	p := pipeline.FromString(want)
	buffer := new(bytes.Buffer)
	p.Output = buffer
	p.Stdout()
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	got := buffer.String()
	if !cmp.Equal(want, got) {
		t.Errorf("want %q & got %q", want, got)
	}
}

func TestWriter_StdoutFailsGracefullyOnPipelineError(t *testing.T) {
	t.Parallel()
	want := "Helloworld 123\n"
	p := pipeline.FromString(want)
	buffer := new(bytes.Buffer)
	p.Output = buffer
	p.Error = errors.New("oops")
	p.Stdout()
	got := buffer.String()
	if got != "" {
		t.Errorf("want no output because of error & got %q", got)
	}
}
