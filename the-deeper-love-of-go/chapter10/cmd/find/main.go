package main

import (
	"books"
	"fmt"
	"os"
)

func main() {
	fmt.Println("# Find Command")
	if len(os.Args) != 2 {
		fmt.Println("Usage: find <BOOK ID>")
		return
	}
	catalogue, err := books.OpenCatalogue("test-data/catalogue.json")
	if err != nil {
		fmt.Printf("opening catalogue: %v\n", err)
		return
	}

	ID := os.Args[1]
	book, ok := catalogue.GetBook(ID)
	if !ok {
		fmt.Println("Sorry I couldn't find that book in our catalogue.")
		return
	}

	fmt.Println(book)
}
