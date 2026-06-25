package main

import (
	"books"
	"fmt"
	"os"
)

func main() {
	fmt.Println("# Server Command")
	if len(os.Args) != 2 {
		fmt.Println("Usage: server <CATALOGUE PATH>")
		return
	}

	path := os.Args[1]
	catalogue, err := books.OpenCatalogue(path)
	if err != nil {
		fmt.Printf("opening catalogue: %v\n", err)
		return
	}

	err = books.ListenAndServe(":3000", catalogue)
	if err != nil {
		fmt.Println(err)
	}
}
