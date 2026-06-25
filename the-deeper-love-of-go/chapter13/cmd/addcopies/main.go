package main

import (
	"books"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("# Add Copies Command")
	if len(os.Args) != 3 {
		fmt.Println("Usage: getcopies <BOOK ID> <AMOUNT>")
		return
	}

	ID := os.Args[1]
	amount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := books.NewClient("localhost:3000")

	stock, err := client.AddCopies(ID, amount)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d copies in stock for %q\n", stock, ID)
}
