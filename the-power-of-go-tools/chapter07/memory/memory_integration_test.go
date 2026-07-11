//go:build integration

package memory_test

import (
	"bytes"
	"memory"
	"os/exec"
	"testing"
)

// This test has a build tag which means it won't be run unless that tag is defined.
// The "integration" tag is not a predefined tag so we have to pass the tag in for it to run.
// `go test -tags=integration -v` is the only way to run this test.
func TestGetFreeOutput_CapturesCommandOutput(t *testing.T) {
	t.Parallel()

	// exec.Command().Output returns a string of stdout
	// exec.Command().OutputCombined returns a string of stdout + stderr
	data, err := exec.Command("/usr/bin/free", "-b").CombinedOutput()
	if err != nil {
		t.Skipf("unable to run `free` command: %v", err)
	}

	if !bytes.Contains(data, []byte("Mem:")) {
		t.Skip("No RAM detected in output")
	}

	text, err := memory.GetFreeOutput()
	if err != nil {
		t.Fatal(err)
	}

	total, err := memory.ParseFreeOutput(text)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Total memory: %d", total)

}
