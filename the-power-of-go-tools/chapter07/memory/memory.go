package memory

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type memory struct {
	Physical Amounts
}

type Amounts struct {
	TotalBytes int
	UsedBytes  int
	FreeBytes  int
}

// MustCompile is an expensive operation, so instead of doing it inside the function on every call, we do it once here.
// Using a `raw string` so we don't have to double escape.
var freeMemRegex = regexp.MustCompile(`Mem:\s+(\d+)\s+(\d+)\s+(\d+) .*`)

func ParseFreeOutput(text string) (Amounts, error) {
	// FindString finds first string that matches this
	// FindStringSubmatch returns the entire matched string and the contents of the capture group
	matches := freeMemRegex.FindStringSubmatch(text)
	// We expect at least 2 elements in the returned slice
	if len(matches) < 2 {
		return Amounts{}, fmt.Errorf("failed to parse free -b output: %q", text)
	}

	// Convert total bytes to a string
	total, err := strconv.Atoi(matches[1])
	if err != nil {
		return Amounts{}, fmt.Errorf("failed to parse total memory: %q", matches[1])
	}

	return Amounts{
		TotalBytes: total,
	}, nil

}

func GetFreeOutput() (string, error) {
	data, err := exec.Command("/usr/bin/free", "-b").CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(data), nil
}