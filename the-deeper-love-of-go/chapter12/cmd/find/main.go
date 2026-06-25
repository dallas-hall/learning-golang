package main

import (
	"books"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("# Find Command")
	if len(os.Args) != 2 {
		fmt.Println("Usage: find <BOOK ID>")
		return
	}

	response, err := http.Get("http://localhost:3000/v1/find/" + os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status %d\n", response.StatusCode)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	book := books.Book{}
	err = json.Unmarshal(data, &book)
	if err != nil {
		fmt.Printf("%v in %q", err, data)
		return
	}

	fmt.Println(book)
}
