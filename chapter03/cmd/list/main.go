package main

import (
	"fmt"
	"books"
)

func main() {
	fmt.Println("# Chapter 3")

	book1 := books.Book{
		Title:  "Sea Room",
		Author: "Adam Nicolson",
		Copies: 2,
	}

	fmt.Println(books.BookToString(book1))
}
