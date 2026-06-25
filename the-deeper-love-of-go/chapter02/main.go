package main

import (
	"fmt"
)

type Book struct {
	Title  string
	Author string
	Copies int
}

func printBook(b Book) {
	// Print all arguments separated by spaces.
	fmt.Println(b.Title, "by", b.Author, "-", b.Copies, "copies")
}

func printBook2(b Book) {
	fmt.Printf("%v by %v - %v copies\n", b.Title, b.Author, b.Copies)
}

func main() {
	fmt.Println("# Chapter 2")
	book1 := Book{
		Title:  "Sea Room",
		Author: "Adam Nicolson",
		Copies: 2,
	}

	printBook(book1)
	printBook2(book1)
}
