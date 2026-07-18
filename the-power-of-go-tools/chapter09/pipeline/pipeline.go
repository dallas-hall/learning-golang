package pipeline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// Allows for reading and writing in a pipeline, much like a shell pipeline.
// We can accept any io.Reader input. We can pass on any io.Writer output.
// Error is used to stop a pipeline.
// Filesystem is used to optionally store the inputs FS object. If nil, local
// disk is used.
type Pipeline struct {
	Input      io.Reader
	Output     io.Writer
	Error      error
	Filesystem fs.FS
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

// Create and return a Pipeline object with the Input set from all files in
// the filesystem. The Output is Stdout. Error is nil.
func FromFileSystem(filesystem fs.FS) *Pipeline {
	files := new(bytes.Buffer)
	fs.WalkDir(filesystem, ".", func(currentPath string, currentPathMetaData fs.DirEntry, err error) error {
		// .Info() returns an fs.FileInfo and an error.
		info, err := currentPathMetaData.Info()
		// Silently skip errors and directories
		if err != nil || info.IsDir() {
			return nil
		}
		// Store the exact path.
		fmt.Fprintln(files, currentPath)
		return nil
	})

	p := New()
	p.Input = files
	p.Filesystem = filesystem
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
		p.Error = fmt.Errorf("bad column, must be positive: %d", column)
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
	if err := scanner.Err(); err != nil {
		p.Error = err
		p.Input = strings.NewReader("")
		return p
	}

	p.Input = result
	return p
}

// Pass in a string to do a simple strings.Contains search.
// We always check for a Pipeline Error first and clear the reader if one exists.
func (p *Pipeline) FindFiles(want string) *Pipeline {
	if p.Error != nil {
		p.Input = strings.NewReader("")
		return p
	}

	result := new(bytes.Buffer)
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		file := scanner.Text()
		if strings.Contains(file, want) {
			fmt.Fprintln(result, file)
		}
	}
	if err := scanner.Err(); err != nil {
		p.Error = err
		p.Input = strings.NewReader("")
		return p
	}

	p.Input = result
	return p
}

// Reads paths from the Input, one by one, and then joins their contents as-is.
// We always check for a Pipeline Error first and clear the reader if one exists.
func (p *Pipeline) Concat() *Pipeline {
	if p.Error != nil {
		p.Input = strings.NewReader("")
		return p
	}

	result := new(bytes.Buffer)
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		file := scanner.Text()
		var data []byte
		var err error
		if p.Filesystem != nil {
			data, err = fs.ReadFile(p.Filesystem, file)
		} else {
			data, err = os.ReadFile(file)
		}
		if err != nil {
			p.Input = strings.NewReader("")
			p.Error = err
			return p
		}
		result.Write(data)
	}
	if err := scanner.Err(); err != nil {
		p.Error = err
		p.Input = strings.NewReader("")
		return p
	}

	p.Input = result
	return p
}

// Removes all empty lines and lines with whitspace only.
// We always check for a Pipeline Error first and clear the reader if one exists.
func (p *Pipeline) RemoveBlankLines() *Pipeline {
	if p.Error != nil {
		p.Input = strings.NewReader("")
		return p
	}

	result := new(bytes.Buffer)
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			fmt.Fprintln(result, line)
		}
	}
	if err := scanner.Err(); err != nil {
		p.Error = err
		p.Input = strings.NewReader("")
		return p
	}

	p.Input = result
	return p
}

// Counts the current lines in the Pipeline Input and return the result as an int.
// We always check for a Pipeline Error first and clear the reader if one exists.
func (p *Pipeline) CountLines() (int, error) {
	if p.Error != nil {
		p.Input = strings.NewReader("")
		return 0, p.Error
	}

	lines := 0
	scanner := bufio.NewScanner(p.Input)
	for scanner.Scan() {
		lines++
	}
	if err := scanner.Err(); err != nil {
		p.Error = err
		p.Input = strings.NewReader("")
		return 0, p.Error
	}

	return lines, nil
}
