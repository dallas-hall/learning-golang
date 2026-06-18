package books

import (
	"fmt"
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
var catalogue = []Book{
	{
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
		ID:     "abc",
	},
	{
		Title:  "White Heat",
		Author: "Dominic Sandbrook",
		Copies: 2,
		ID:     "xyz",
	},
}

func GetAllBooks() []Book {
	return catalogue
}

func GetBook(ID string) (Book, bool) {
	for _, book := range catalogue {
		if book.ID == ID {
			return book, true
		}
	}
	return Book{}, false
}
