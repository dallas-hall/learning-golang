package main

import (
	"log"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

// Taken from the-power-of-go-tools/chapter04/count-pflag/count_test.go - see for comments
func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "test/scripts",
	})
}

// Taken from the-power-of-go-tools/chapter07/simplehowlong/howlong_test.go - see for comments
func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"write": writeMain,
	}))
}

// Must return an int.
func writeMain() int {
	// This file is inside of testscripts isolated container which is created for
	// every txtar file processed. e.g., /tmp/testscript12345
	path := "output.txt"
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("error opening %q: %s", path, err)
	}
	defer file.Close()

	err = write(file)
	if err != nil {
		log.Fatalf("error writing to %q: %s", path, err)
	}

	path2 := "output2.txt"
	file2, err := os.Create(path2)
	if err != nil {
		log.Fatalf("error opening %q: %s", path2, err)
	}
	defer file2.Close()

	err = writeBad(file2)
	if err != nil {
		log.Fatalf("error writing to %q: %s", path2, err)
	}

	return 0
}
