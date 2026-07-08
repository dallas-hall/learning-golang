package memory_test

import (
	"memory"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMemory_DisplaysTotalRAMInBytes(t *testing.T) {
	t.Parallel()
	file := "test/data/free.txt"

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to open %q because: %s", file, err)
	}

	want := memory.Amounts{
		TotalBytes: 33_536_151_552,
	}

	got, err := memory.ParseFreeOutput(string(data))
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
