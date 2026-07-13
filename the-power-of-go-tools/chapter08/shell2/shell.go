package simpleshell2

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// This behaviour was taken from count.go - see it for comments.
type Session struct {
	Input              io.Reader
	Output, Error      io.Writer
	DryRun, Transcript bool
}

// This behaviour was taken from count.go - see it for comments.
func NewSession(i io.Reader, o, e io.Writer) *Session {
	s := &Session{
		Input:      i,
		Output:     o,
		Error:      e,
		DryRun:     false,
		Transcript: false,
	}
	return s
}

func (s *Session) Run() {
	// Read buffered lines from stdin & print PS1
	prompt := GetPrompt()
	fmt.Fprint(s.Output, prompt)

	input := bufio.NewScanner(s.Input)
	for input.Scan() {
		line := input.Text()
		if line == "exit" {
			break
		}

		cmd, err := CmdFromString(line)
		if err != nil {
			fmt.Fprint(s.Output, prompt)
			continue
		}

		if s.DryRun {
			fmt.Fprintf(s.Output, "%s\n", line)
		}

		if s.Transcript {
			data := []byte(line)
			err := WriteToFile("test/data/transcript.txt", data)
			if err != nil {
				fmt.Fprintln(s.Error, "error saving file: ", err)
			}
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(s.Error, "error getting command output: ", err)
		}
		fmt.Fprintf(s.Output, "%s", output)
		fmt.Fprint(s.Output, prompt)
	}
	fmt.Fprintln(s.Output, "\nGoodbye")
}

func GetPrompt() string {
	// Create a PS1
	user := os.Getenv("USER")
	hostname := os.Getenv("HOSTNAME")
	pwd := filepath.Base(os.Getenv("PWD"))
	var prompt string
	if user == "root" {
		prompt = fmt.Sprintf("[%s@%s %s] # ", user, hostname, pwd)
	} else {
		prompt = fmt.Sprintf("[%s@%s %s] $ ", user, hostname, pwd)
	}
	return prompt
}

func CmdFromString(s string) (*exec.Cmd, error) {
	// Split the string using whitespace
	args := strings.Fields(s)
	// If the string is empty we need to error out
	if len(args) < 1 {
		return nil, errors.New("empty input")
	}
	// Assume args[0] is always the command and args[1:] is the options
	return exec.Command(args[0], args[1:]...), nil
}

func Main() {
	session := NewSession(os.Stdin, os.Stdout, os.Stderr)
	session.Run()
}

// Taken from the-power-of-go-tools/chapter05/writer/writer.go - see for comments
func WriteToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0o600)
	if err != nil {
		return err
	}
	return os.Chmod(path, 0o600)
}
