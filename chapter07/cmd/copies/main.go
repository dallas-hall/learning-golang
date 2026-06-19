package main

import (
	"books"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("# Copy Command")
	if len(os.Args) != 3 {
		fmt.Println("Usage: copies <BOOK ID> <NUMBER OF COPIES>")
		return
	}

	catalogue, err := books.OpenCatalogue("test-data/catalogue.json")
	if err != nil {
		fmt.Printf("opening catalogue: %v\n", err)
		return
	}

	ID := os.Args[1]
	// Don't call catalogue.GetBook(ID) here because it can lead to a race condition.
	// If 2 people open the same book object at the same time and make changes, one of those changes will be lost.
	book, ok := catalogue.GetBook(ID)
	if !ok {
		fmt.Println("Sorry I couldn't find that book in our catalogue.")
		return
	}

	copies, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("reading user copies: %v\n", err)
	}

	err = book.SetCopies(copies)
	if err != nil {
		fmt.Printf("updating copies: %v\n", err)
	}

	// TODO: save changes to disk
	// err := catalogue.Sync("test-data/catalogue.json")
	// if err != nil {
	// 	fmt.Printf("writing catalogue: %v\n", err)
	// }
	fmt.Printf("Updated book ID %q with %d copies.\n", ID, copies)
}
