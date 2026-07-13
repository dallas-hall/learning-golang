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
	Input         io.Reader
	Output, Error io.Writer
	DryRun        bool
}

// This behaviour was taken from count.go - see it for comments.
func NewSession(i io.Reader, o, e io.Writer) *Session {
	s := &Session{
		Input:  i,
		Output: o,
		Error:  e,
		DryRun: false,
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

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(s.Error, "error: ", err)
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
