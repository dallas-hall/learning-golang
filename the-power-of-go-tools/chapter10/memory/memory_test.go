package memory_test

import (
	"memory"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMemory_DisplaysTotalRAMInBytes(t *testing.T) {
	t.Parallel()

	// Read our test data
	file := "test/data/free.txt"
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to open %q because: %s", file, err)
	}

	// Parse test data into struct
	got, err := memory.ParseFreeOutput(string(data))
	if err != nil {
		t.Fatalf("Failed to parse `free` output because: %s", err)
	}

	// Create testing struct
	want := memory.Amounts{
		TotalBytes: 33_536_151_552,
		UsedBytes:  8_604_680_192,
		FreeBytes:  374_026_240,
	}

	// Compare
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMemory_NewReturnsStructCorrectly(t *testing.T) {
	t.Parallel()

	// Create an empty Memory struct manually
	want := &memory.Memory{
		Physical: memory.Amounts{},
	}

	// Create an empty Memory struct with our constructor
	got := memory.New()

	// Compare
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestMomory_ToJSONWorksCorrectly(t *testing.T) {
	t.Parallel()

	// Read our test data
	path := "test/data/free.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to open %q because: %s", path, err)
	}

	// Parse test data into struct
	ramTotals, err := memory.ParseFreeOutput(string(data))
	if err != nil {
		t.Fatalf("Failed to parse `free` output because: %s", err)
	}
	myMemory := memory.New()
	myMemory.Physical = ramTotals

	// Convert struct to JSON
	got := myMemory.ToJSON()

	// Read in test JSON file
	path = "test/data/free-parsed.json"
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to open %q because: %s", path, err)
	}
	want := string(data)

	// Compare
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMemory_ToYAMLWorksCorrectly(t *testing.T) {
	t.Parallel()
}