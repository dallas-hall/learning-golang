package howlong

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type HowLong struct {
	Command     string
	ExitCode    int
	RanAt       time.Time
	TimeOutputs TimeOutputs
}

type TimeOutputs struct {
	Shell  string
	Real   string
	User   string
	System string
}

// Taken from memory.go - see there for more comments.
var timeOutputRegex = regexp.MustCompile(`\s*real\s+(\S+)\s+user\s+(\S+)\s+sys\s+(\S+)\s*`)

// This method is coupled to an instance of a HowLong struct.
func (hl *HowLong) ParseTimeOutput() error {
	matches := timeOutputRegex.FindStringSubmatch(hl.TimeOutputs.Shell)
	// We except whole match + 3 groups.
	if len(matches) < 4 {
		return fmt.Errorf("failed to parse time <command> output: %q", hl.TimeOutputs.Shell)
	}

	hl.TimeOutputs.Real = matches[1]
	hl.TimeOutputs.User = matches[2]
	hl.TimeOutputs.System = matches[3]

	return nil
}

func (hl *HowLong) GetTimeOutput() error {
	hl.RanAt = time.Now()

	// Wrapping in { ...; } creates a compound command (a command group). Redirections placed after } apply to the group as a whole
	// Redirect stderr to stdout so we can capture it as time output goes to stderr.
	shellCommand := fmt.Sprintf("{ time %s; } 2>&1", hl.Command)

	// Use a Bash shell to run the `time` keyword with the passed in command command.
	data, err := exec.Command("bash", "-c", shellCommand).CombinedOutput()
	if err != nil {
		// Running `go run ./cmd/howlong/ false` shows this working.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			hl.ExitCode = exitErr.ExitCode()
		}
		return err
	}
	hl.TimeOutputs.Shell = string(data)
	hl.ExitCode = 0

	return nil
}

// Convenience wrapper
func Main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}
	// Assuming all arguments are 1 command and its arguments
	command := strings.Join(os.Args[1:], " ")

	hl := HowLong{
		Command: command,
	}

	err := hl.GetTimeOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running `{ time %s; }`: %v\n", hl.Command, err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", hl)
}
