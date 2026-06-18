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

// This function requires a book to be passed in.
/*
func BookToString(b Book) string {
	return fmt.Sprintf("%v by %v (%v copies)", b.Title, b.Author, b.Copies)
}
*/

// This method is coupled to an instance of a book.
// fmt.Println will automatically call a method named String if it exists.
// e.g. fmt.Println(book) instead of fmt.Println(book.String())
func (b Book) String() string {
	return fmt.Sprintf("%v by %v (%v copies)", b.Title, b.Author, b.Copies)
}

func (b *Book) SetCopies(copies int) error {
	if copies < 0 {
		return fmt.Errorf("negative number of copies: %d", copies)
	}
	// We don't need to manually dereference with:
	// (*b).Copies = copies
	// As Go will automatically deference it with:
	b.Copies = copies
	return nil
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

// Creating our own type so we can add methods to it.
// Conceptually similar to an object in OOP.
type Catalogue map[string]Book

// This function is replacing the package-level variable above.
func GetCatalogue() Catalogue {
	return Catalogue{
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

func (catalogue Catalogue) GetAllBooks() []Book {
	// maps.Values returns an iterator
	// slices.Collect conumes the iterator and collects all elements into a slice.
	return slices.Collect(maps.Values(catalogue))
}

func (catalogue Catalogue) GetBook(ID string) (Book, bool) {
	book, ok := catalogue[ID]
	return book, ok
}

func (catalogue Catalogue) AddBook(book Book) {
	catalogue[book.ID] = book
}
