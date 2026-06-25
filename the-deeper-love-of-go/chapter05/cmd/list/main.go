package main

import (
	"books"
	"fmt"
)

func main() {
	catalogue := books.GetCatalogue()

	for _, book := range books.GetAllBooks(catalogue) {
		fmt.Println(books.BookToString(book))
	}
}
