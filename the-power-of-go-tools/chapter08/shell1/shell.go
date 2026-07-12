package simpleshell1

import (
	"errors"
	"os/exec"
	"strings"
)

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
