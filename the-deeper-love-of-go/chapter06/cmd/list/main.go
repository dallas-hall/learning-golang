package main

import (
	"books"
	"fmt"
)

func main() {
	catalogue := books.GetCatalogue()

	for _, book := range catalogue.GetAllBooks() {
		// implicitly call book.String()
		fmt.Println(book)
	}
}
