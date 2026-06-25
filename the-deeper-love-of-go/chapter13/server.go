// Include this in the books package, even though it is a separate source file to books.go
package books

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
			ID := r.PathValue("id")
			book, ok := catalogue.GetBook(ID)
			if !ok {
				http.Error(w, fmt.Sprintf("%q not found", ID), http.StatusNotFound)
				return
			}
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				panic(err)
			}
		})

	multiplexer.HandleFunc(version+"/getcopies/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			ID := r.PathValue("id")
			copies, err := catalogue.GetCopies(ID)
			if err != nil {
				http.Error(w, fmt.Sprintf("%q not found", ID), http.StatusNotFound)
				return
			}

			err = json.NewEncoder(w).Encode(copies)
			if err != nil {
				panic(err)
			}
		})

	multiplexer.HandleFunc(version+"/addcopies/{id}/{amount}",
		func(w http.ResponseWriter, r *http.Request) {
			ID := r.PathValue("id")
			amount, err := strconv.Atoi(r.PathValue("amount"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			stock, err := catalogue.AddCopies(ID, amount)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			err = catalogue.SyncCatalogues()
			if err != nil {
				// Use panic so the error goes to the server logs and not the user.
				panic(err)
			}

			err = json.NewEncoder(w).Encode(stock)
			if err != nil {
				// Use panic here because http.ResponseWriter is probably having connectivity issues.
				panic(err)
			}
		})

	multiplexer.HandleFunc(version+"/subcopies/{id}/{amount}",
		func(w http.ResponseWriter, r *http.Request) {
			ID := r.PathValue("id")
			amount, err := strconv.Atoi(r.PathValue("amount"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			stock, err := catalogue.SubCopies(ID, amount)
			if err != nil {
				if errors.Is(err, ErrorNotEnoughStock) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				} else if errors.Is(err, ErrorNotEnoughStock) {
					http.Error(w, err.Error(), http.StatusNotFound)
				} else {
					http.Error(w, err.Error(), http.StatusNotImplemented)
				}
				return
			}

			err = catalogue.SyncCatalogues()
			if err != nil {
				panic(err)
			}

			err = json.NewEncoder(w).Encode(stock)
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
