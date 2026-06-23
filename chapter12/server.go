// Include this in the books package, even though it is a separate source file to books.go
package books

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ListenAndServe(address string, catalogue *Catalogue) error {
	// A Multiplexer can distinguish between different kinds of API requests and route each to a specific handler.
	multiplexer := http.NewServeMux()

	version := "/v1"

	// Create a handler for each endpoint
	multiplexer.HandleFunc(version+"/health",
		func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintln(writer, "200 OK")
		})

	// The anonymous closure function inside of http.HandlerFunc can see and use the catalogue pointer being passed in!
	multiplexer.HandleFunc(version+"/list",
		func(w http.ResponseWriter, r *http.Request) {
			books := catalogue.GetAllBooks()
			err := json.NewEncoder(w).Encode(books)
			if err != nil {
				panic(err)
			}
		})

	multiplexer.HandleFunc(version+"/find/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			book, ok := catalogue.GetBook(r.PathValue("id"))
			if !ok {
				http.Error(w, fmt.Sprintf("%q not found", r.PathValue("id")), http.StatusNotFound)
				return
			}
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				panic(err)
			}
		})

	return http.ListenAndServe(address, multiplexer)

	// The anonymous closure function inside of http.HandlerFunc can see and use the catalogue pointer being passed in!
	// return http.ListenAndServe(address, http.HandlerFunc(
	//
	//	func(writer http.ResponseWriter, r *http.Request) {
	//		books := catalogue.GetAllBooks()
	//		err := json.NewEncoder(writer).Encode(books)
	//		if err != nil {
	//			panic(err)
	//		}
	//	}))
}
