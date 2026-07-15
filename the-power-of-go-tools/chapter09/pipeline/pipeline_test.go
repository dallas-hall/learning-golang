package pipeline_test

import (
	"bytes"
	"errors"
	"io"
	"pipeline"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPipeline_StdoutPrintsToOutputFromString(t *testing.T) {
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

func TestPipeline_StdoutPrintsToOutputFromFile(t *testing.T) {
	t.Parallel()
	want := []byte("Helloworld 123\n")

	p := pipeline.FromFile("test/data/text.txt")
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	got, err := io.ReadAll(p.Input)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("want %q & got %q", want, got)
	}
}

func TestPipeline_StdoutFailsGracefullyOnPipelineError(t *testing.T) {
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
