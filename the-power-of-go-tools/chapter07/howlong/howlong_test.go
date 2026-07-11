package howlong_test

import (
	"howlong"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHowLong_ReturnsTimeOutputAsString(t *testing.T) {
	t.Parallel()
	// Generated with `{ time sleep 1; } 2>&1 test/data/time-sleep-1.txt` which is explained in howlong.go GetTimeOutput()
	file := "test/data/time-sleep-1.txt"

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to open %q because: %s", file, err)
	}

	want := howlong.HowLong{
		TimeOutputs: howlong.TimeOutputs{
			Shell: `
real	0m1.001s
user	0m0.000s
sys	0m0.001s
`,
			Real:   "0m1.001s",
			User:   "0m0.000s",
			System: "0m0.001s",
		},
	}

	got := howlong.HowLong{
		TimeOutputs: howlong.TimeOutputs{
			Shell: string(data),
		},
	}

	err = got.ParseTimeOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
