package main

import (
	"books"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: find <BOOK ID>")
		return
	}
	catalogue := books.GetCatalogue()
	ID := os.Args[1]
	book, ok := books.GetBook(catalogue, ID)
	if !ok {
		fmt.Println("Sorry I couldn't find that book in our catalogue.")
		return
	}

	fmt.Println(books.BookToString(book))
}
