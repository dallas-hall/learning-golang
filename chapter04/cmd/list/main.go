package main

import (
	"books"
	"fmt"
)

func main() {
	fmt.Println("# Chapter 3")

	book1 := books.Book{
		Title:  "Sea Room",
		Author: "Adam Nicolson",
		Copies: 2,
	}

	fmt.Println(books.BookToString(book1))

	fmt.Println("# Chapter 4")
	fmt.Println(books.GetAllBooks())
	for _, book := range books.GetAllBooks() {
		fmt.Println(books.BookToString(book))
	}

}
