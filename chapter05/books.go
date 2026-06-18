package books

import (
	"fmt"
	"maps"
	"slices"
)

type Book struct {
	Title  string
	Author string
	Copies int
	ID     string
}

func BookToString(b Book) string {
	return fmt.Sprintf("%v by %v (%v copies)", b.Title, b.Author, b.Copies)
}

// This is a package-scoped variable, it is only visible within the package where it's defined. It may be incorrectly called a global variable.
/*
var catalogue = map[string]Book{
	"abc": {
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
		ID:     "abc",
	},
	"xyz": {
		Title:  "White Heat",
		Author: "Dominic Sandbrook",
		Copies: 2,
		ID:     "xyz",
	},
}
*/

// This function is replacing the package-level variable above.
func GetCatalogue() map[string]Book {
	return map[string]Book{
		"abc": {
			Title:  "In the Company of Cheerful Ladies",
			Author: "Alexander McCall Smith",
			Copies: 1,
			ID:     "abc",
		},
		"xyz": {
			Title:  "White Heat",
			Author: "Dominic Sandbrook",
			Copies: 2,
			ID:     "xyz",
		},
	}
}

func GetAllBooks(catalogue map[string]Book) []Book {
	// maps.Values returns an iterator
	// slices.Collect conumes the iterator and collects all elements into a slice.
	return slices.Collect(maps.Values(catalogue))
}

func GetBook(catalogue map[string]Book, ID string) (Book, bool) {
	book, ok := catalogue[ID]
	return book, ok
}

func AddBook(catalogue map[string]Book, book Book) {
	catalogue[book.ID] = book
}
