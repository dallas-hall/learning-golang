package simplehowlong_test

import (
	"fmt"
	"os"
	"simplehowlong"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestRun_ReportsCorrectElapsedTimeForComma(t *testing.T) {
	t.Parallel()
	target := 100 * time.Millisecond
	elapsed, err := simplehowlong.Run("sleep", "0.1")
	if err != nil {
		t.Fatal(err)
	}
	epsilon := target - elapsed
	if epsilon.Abs() > 300*time.Millisecond {
		t.Fatalf("want %s & got %s (not close enough)", target, elapsed)
	}
}

func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "test/scripts",
	})
}

// TestMain in howlong_test.go
func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"howlong": howlongMain,
	}))
}

func howlongMain() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: howlong <command> [args...]")
		return 1
	}
	elapsed, err := simplehowlong.Run(os.Args[1], os.Args[2:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("time: %s\n", elapsed)
	return 0
}
