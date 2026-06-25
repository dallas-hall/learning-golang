package main

import (
	"books"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("# List Command")

	response, err := http.Get("http://localhost:3000/v1/list")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status %d", response.StatusCode)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	books := []books.Book{}
	err = json.Unmarshal(data, &books)
	if err != nil {
		fmt.Printf("%v in %q", err, data)
		return
	}

	for _, book := range books {
		// implicitly call book.String()
		fmt.Println(book)
	}
}
