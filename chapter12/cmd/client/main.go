package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Make a HTTP GET request to our http-server.go
	// PROTOCOL://HOST:PORT/URI
	response, err := http.Get("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	// Handle errors (eg not 200 OK)
	if response.StatusCode != http.StatusOK {
		panic(response.StatusCode)
	}

	// Read the response from the server and print it on the client
	// We need a reader because response.Body is in bytes.
	msg, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", msg)
}
