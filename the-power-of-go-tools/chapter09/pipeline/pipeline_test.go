package pipeline_test

import (
	"bytes"
	"errors"
	"io"
	"pipeline"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPipeline_StdoutPrintsStringToOutput(t *testing.T) {
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

func TestPipeline_StdoutPrintsNothingOnError(t *testing.T) {
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

func TestPipeline_FromFileReadsDataCorrectly(t *testing.T) {
	t.Parallel()
	want := []byte("Helloworld 123\n")

	p := pipeline.FromFile("test/data/test.txt")
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

func TestPipeline_FromFileSetsErrorGivenNonexistentFile(t *testing.T) {
	t.Parallel()
	p := pipeline.FromFile("doesnt-exist.txt")
	if p.Error == nil {
		t.Fatal("want error opening non-existent file, got nil")
	}
}

func TestPipeline_StringReturnsPipeContentsAsString(t *testing.T) {
	t.Parallel()
	want := "Helloworld 123\n"
	p := pipeline.FromString(want)
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("want %q & got %q", want, got)
	}
}

func TestPipeline_StringReturnsErrorWhenPipelineErrorSet(t *testing.T) {
	t.Parallel()
	p := pipeline.FromString("Helloworld 123\n")
	p.Error = errors.New("uh oh spaghettio, Pipeline error!")
	_, err := p.String()
	if err == nil {
		t.Errorf("want Pipeline error from String but not nil")
	}
}

func TestPipeline_ColumnGetsDataCorrectly(t *testing.T) {
	t.Parallel()
	want := "2\n2\n2\n"
	p := pipeline.FromFile("test/data/test.csv")
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	got, err := p.Column(2, ",").String()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("want %q & got %q", want, got)
	}
}

func TestPipeline_ColumnProducesNothingWhenPipelineErrorSet(t *testing.T) {
	t.Parallel()
	p := pipeline.FromFile("test/data/test.csv")
	p.Error = errors.New("uh oh spaghettio, Pipeline error!")
	data, err := io.ReadAll(p.Column(2, ",").Input)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > 0 {
		t.Errorf("want Pipeline error from Column but got %q", data)
	}
}

func TestPipeline_ColumnProducesNothingWhenGivenInvalidColumn(t *testing.T) {
	t.Parallel()
	p := pipeline.FromFile("test/data/test.csv")
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	p.Column(-1, ",")
	if p.Error == nil {
		t.Error("want Pipeline error from Column but got nil")
	}

	data, err := io.ReadAll(p.Column(1, ",").Input)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > 0 {
		t.Errorf("want no output from Column but got %q", data)
	}
}
