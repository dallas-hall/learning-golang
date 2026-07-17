package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/bitfield/script"
)

func main() {
	// https://github.com/bitfield/script#unix-equivalents - there is quite a lot!

	// ls
	dotFiles := regexp.MustCompile(`^\.`)
	script.ListFiles(".").RejectRegexp(dotFiles).Stdout()

	// xargs
	// Run the `echo` command for each file found with *
	script.ListFiles("*").ExecForEach("echo {{.}}").Stdout()

	// If there doesn’t happen to be a built-in script method that does what we
	// want, we can just write our own, using Filter:
	script.Echo("hello world").Filter(func(r io.Reader, w io.Writer) error {
		n, err := io.Copy(w, r)
		fmt.Fprintf(w, "\nfiltered %d bytes\n", n)
		return err
	}).Stdout()

	// grep --include=*.go -r -c '.' $PWD | cut -d ':' -f 2 | paste -sd+ - | bc

	// grep --include=*.go -r -c '.' $PWD
	// Count non-blank lines in go files in $PWD

	// cut -d ':' -f 2
	// Splits each line on : and get the second column

	// paste -sd+ -
	// From stdin, join all lines together with + delimiter

	// bc
	// Calculate the total

	//re := regexp.MustCompile(`.*\.go`)
	lines, err := script.FindFiles(".").
		MatchRegexp(regexp.MustCompile(".go$")).
		Concat().
		RejectRegexp(regexp.MustCompile(`^$`)).
		CountLines()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("You've written %d lines of Go in this project. Nice work!\n", lines)
}
