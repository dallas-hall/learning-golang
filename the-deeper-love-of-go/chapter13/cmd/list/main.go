package main

import (
	"books"
	"fmt"
)

func main() {
	fmt.Println("# List Command")
	
	client := books.NewClient("localhost:3000")

	books, err := client.GetAllBooks()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, book := range books {
		// implicitly call book.String()
		fmt.Println(book)
	}
}
