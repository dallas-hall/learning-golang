package main

import (
	"books"
	"fmt"
)

func main() {
	fmt.Println("# List Command")
	//catalogue := books.GetCatalogue()

	// This replaces the hardcoded catalogue data from above.
	catalogue, err := books.OpenCatalogue("test-data/catalogue.json")
	if err != nil {
		fmt.Printf("opening catalogue: %v\n", err)
		return
	}

	for _, book := range catalogue.GetAllBooks() {
		// implicitly call book.String()
		fmt.Println(book)
	}
}
