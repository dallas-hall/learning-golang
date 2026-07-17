package main

import (
	"fmt"
	"log"
	"os"
	"pipeline"
)

func main() {
	// Create a Pipeline from a string and print it.
	fmt.Println("# String Based Pipeline")
	pipeline.FromString("hello from pipeline.FromString\n").Stdout()

	// Create a Pipeline from a file and print it.
	fmt.Println("# File Based Pipeline")
	pipeline.FromFile("cmd/pipeline/main.go").Stdout()

	// Create a PipeLine from a ZIP file and count all non-whitespace only or
	// empty lines inside of all Go files.
	fmt.Println("# Filesystem Based Pipeline")
	filesystem := os.DirFS(".")
	lines, err := pipeline.FromFileSystem(filesystem).FindFiles(".go").Concat().RemoveBlankLines().CountLines()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You have typed %d lines of code.\n", lines)
}
