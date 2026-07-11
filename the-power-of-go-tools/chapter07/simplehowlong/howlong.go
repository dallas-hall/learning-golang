package simplehowlong

import (
	"os"
	"os/exec"
	"time"
)

func Run(program string, args ...string) (time.Duration, error) {
	command := exec.Command(program, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	start := time.Now()
	if err := command.Run(); err != nil {
		return 0, err
	}
	return time.Since(start), nil
}
