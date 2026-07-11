//go:build integration

package howlong_test

import (
	"bytes"
	"howlong"
	"os/exec"
	"testing"
)

// Taken from memory_integration_text.go
// `go test -tags=integration -v` is the only way to run this test.
func TestGetTimeOutput_CapturesSleep1ShellOutput(t *testing.T) {
	t.Parallel()

	// This is commented inside of howlong.go GetTimeOutput()
	// We are checking if the Bash shell, the `time`` keyword, and the `sleep`` command are available.
	data, err := exec.Command("bash", "-c", "{ time sleep 1; } 2>&1").CombinedOutput()
	if err != nil {
		t.Skipf("unable to run `{ time sleep 1; }` because: %v", err)
	}

	// Make sure expected output is present
	if !bytes.Contains(data, []byte("real")) {
		t.Skip("No real time detected in output")
	}
	if !bytes.Contains(data, []byte("user")) {
		t.Skip("No user time detected in output")
	}
	if !bytes.Contains(data, []byte("sys")) {
		t.Skip("No sys time detected in output")
	}

	// Create our test object
	hl := howlong.HowLong{
		Command: "sleep 1",
	}

	err = hl.GetTimeOutput()
	if err != nil {
		t.Fatal(err)
	}

	err = hl.ParseTimeOutput()
	if err != nil {
		t.Fatal(err)
	}
}
