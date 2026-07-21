package memory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"go.yaml.in/yaml/v3"
)

type Memory struct {
	Physical Amounts `yaml:"Physical" json:"Physical"`
}

type Amounts struct {
	TotalBytes int `yaml:"TotalBytes" json:"TotalBytes"`
	UsedBytes  int `yaml:"UsedBytes" json:"UsedBytes"`
	FreeBytes  int `yaml:"FreeBytes" json:"FreeBytes"`
}

func New() *Memory {
	return &Memory{
		Physical: Amounts{},
	}
}

// MustCompile is an expensive operation, so instead of doing it inside the function on every call, we do it once here.
// Using a `raw string` so we don't have to double escape.
var freeMemRegex = regexp.MustCompile(`Mem:\s+(\d+)\s+(\d+)\s+(\d+) .*`)

// Accepts a string and uses a hardcoded regex to parse the output of `free -b“
func ParseFreeOutput(text string) (Amounts, error) {
	// FindString finds first string that matches this
	// FindStringSubmatch returns the entire matched string and the contents of the capture group
	matches := freeMemRegex.FindStringSubmatch(text)
	// We expect at least 2 elements in the returned slice
	if len(matches) < 2 {
		return Amounts{}, fmt.Errorf("failed to parse free -b output: %q", text)
	}

	// Convert total, used, and free bytes to a string
	total, err := strconv.Atoi(matches[1])
	if err != nil {
		return Amounts{}, fmt.Errorf("failed to parse total memory: %q", matches[1])
	}

	used, err := strconv.Atoi(matches[2])
	if err != nil {
		return Amounts{}, fmt.Errorf("failed to parse used memory: %q", matches[2])
	}

	free, err := strconv.Atoi(matches[3])
	if err != nil {
		return Amounts{}, fmt.Errorf("failed to parse free memory: %q", matches[3])
	}

	return Amounts{
		TotalBytes: total,
		UsedBytes:  used,
		FreeBytes:  free,
	}, nil

}

func GetFreeOutput() (string, error) {
	data, err := exec.Command("/usr/bin/free", "-b").CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Returns a string of the Memory struct in JSON.
// Panics if invalid input.
func (m *Memory) ToJSON() string {
	// No prefix and indent 2 spaces.
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		// Panic on JSON marshalling errors so the user fixes the input
		panic(err)
	}
	return string(data)

}

// Returns a string of the Memory struct in YAML.
// Panics if invalid input.
func (m *Memory) ToYAML() string {
	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	err := encoder.Encode(m)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
