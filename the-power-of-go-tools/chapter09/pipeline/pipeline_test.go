package pipeline_test

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"pipeline"
	"slices"
	"strings"
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

func TestPipeline_FindFilesReturnsCorrectFilesUsingZip(t *testing.T) {
	// Allow the test to run in parallel with other tests.
	t.Parallel()

	// Create what we want to see after running our code.
	want := []string{
		"apache_log/main.go",
		"pipeline/cmd/pipeline/main.go",
		"pipeline/pipeline.go",
		"pipeline/pipeline_test.go",
		"script-pkg/main.go",
		"writer/writer.go",
		"writer/writer_test.go",
	}

	// Read our test data, which is a zip file of this project's directory.
	filesystem, err := zip.OpenReader("test/data/test-data.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer filesystem.Close()

	// Grab all filenames and concatenate them with newlines.
	var sb strings.Builder
	for _, f := range filesystem.File {
		sb.WriteString(f.Name)
		sb.WriteByte('\n')
	}

	// Build a Pipeline with the final string and filter it for Go files.
	p := pipeline.FromString(sb.String()).FindFiles(".go")
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	// Grab all filenames from the Pipeline's Input
	var got []string
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}

	// Sort our slices so they can potentially match during comparison.
	slices.Sort(want)
	slices.Sort(got)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPipeline_ConcatReturnsCorrectFileContentsUsingDisk(t *testing.T) {
	t.Parallel()

	// Construct our test data by concatentating the contents of every file.
	files := []string{
		"test/data/1.go",
		"test/data/2.go",
		"test/data/3.go",
	}
	want := new(bytes.Buffer)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		want.Write(data)
	}

	// Create a Pipeline with our filenames.
	p := pipeline.FromString(strings.Join(files, "\n"))
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	// Concatenate all the files in the Pipeline.
	p = p.Concat()
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	// Compare concats
	got, err := io.ReadAll(p.Input)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want.Bytes(), got) {
		t.Errorf("want %q & got %q", want.Bytes(), got)
	}
}

func TestPipeline_RemoveBlankLinesReturnsCorrectContents(t *testing.T) {
	t.Parallel()

	// Construct our test data by concatentating the contents of every file.
	files := []string{
		"test/data/1.go",
		"test/data/2.go",
		"test/data/3.go",
	}
	want := new(bytes.Buffer)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) != "" {
				fmt.Fprintln(want, line)
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
	}

	// Create a Pipeline with our filenames.
	p := pipeline.FromString(strings.Join(files, "\n"))
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	// Concatenate all the files in the Pipeline and remove blank lines.
	p = p.Concat().RemoveBlankLines()
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	// Read the Pipeline's Input after Concat + RemoveBlankLines have run.
	got, err := io.ReadAll(p.Input)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want.Bytes(), got) {
		t.Errorf("want %q & got %q", want.Bytes(), got)
	}
}

func TestPipeline_CountReturnsCorrectValue(t *testing.T) {
	t.Parallel()

	// Construct our test data by concatentating the contents of every file.
	files := []string{
		"test/data/1.go",
		"test/data/2.go",
		"test/data/3.go",
	}
	want := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) != "" {
				want++
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal(err)
		}
	}

	// Create a Pipeline with our filenames.
	p := pipeline.FromString(strings.Join(files, "\n"))
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	// Concatenate all the files in the Pipeline, remove blank lines, and count.
	got, err := p.Concat().RemoveBlankLines().CountLines()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("want %d & got %d", want, got)
	}
}

func TestPipeline_FromFileSystemReadsDataCorrectly(t *testing.T) {
	t.Parallel()

	// Create what we want to see after running our code.
	want := []string{
		"apache_log/logs.txt",
		"apache_log/main.go",
		"pipeline/cmd/pipeline/main.go",
		"pipeline/go.mod",
		"pipeline/go.sum",
		"pipeline/pipeline.go",
		"pipeline/pipeline_test.go",
		"pipeline/test/data/test.csv",
		"pipeline/test/data/test.txt",
		"script-pkg/go.mod",
		"script-pkg/go.sum",
		"script-pkg/main.go",
		"writer/go.mod",
		"writer/go.sum",
		"writer/test/data/output.txt",
		"writer/test/scripts/test.txtar",
		"writer/writer.go",
		"writer/writer_test.go",
	}

	// Read our test data, which is a zip file of this project's directory.
	filesystem, err := zip.OpenReader("test/data/test-data.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer filesystem.Close()

	// Build a Pipeline with the final string and filter it for Go files.
	p := pipeline.FromFileSystem(filesystem)
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	// Grab all filenames from the Pipeline's Input
	var got []string
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}

	// Sort our slices so they can potentially match during comparison.
	slices.Sort(want)
	slices.Sort(got)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPipeline_ConcatReturnsCorrectFileContentsUsingFileSystem(t *testing.T) {
	t.Parallel()

	// Read our test data, which is a zip file of this project's directory.
	filesystem, err := zip.OpenReader("test/data/test-data.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer filesystem.Close()

	want := new(bytes.Buffer)
	err = fs.WalkDir(filesystem, ".", func(currentPath string, currentPathMetaData fs.DirEntry, err error) error {
		// .Info() returns an fs.FileInfo and an error.
		info, err := currentPathMetaData.Info()
		// Silently skip errors and directories
		if err != nil || info.IsDir() {
			return nil
		}
		// Read the file's contents from the zip filesystem and append them.
		data, err := fs.ReadFile(filesystem, currentPath)
		if err != nil {
			return err
		}
		want.Write(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a Pipeline with our filenames.
	p := pipeline.FromFileSystem(filesystem)
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	// Concatenate all the files in the Pipeline.
	p = p.Concat()
	if p.Error != nil {
		t.Fatal(p.Error)
	}

	// Compare concats.
	got, err := io.ReadAll(p.Input)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want.Bytes(), got) {
		t.Errorf("want %q & got %q", want.Bytes(), got)
	}
}
