package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Start a HTTP server listening on port 3000
	// It calls http.HandlerFunc when a request comes in
	// This will run until a fatal error or it is closed.
	http.ListenAndServe(":3000", http.HandlerFunc(hello))
}

// http.HandlerFunc always takes a http.ResponseWriter and *http.Request. We aren't using the request though.
// http.ResponseWriter replies to the client making the HTTP request.

func hello(w http.ResponseWriter, r *http.Request) {
	// Fprintln prints to any destination that we provide.
	fmt.Fprintln(w, "Hello world.")
}
