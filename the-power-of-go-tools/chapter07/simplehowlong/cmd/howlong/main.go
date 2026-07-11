package main

import (
	"fmt"
	"os"
	"simplehowlong"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	elapsed, err := simplehowlong.Run(os.Args[1], os.Args[2:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("time: %s\n", elapsed.Round(time.Millisecond))
}
