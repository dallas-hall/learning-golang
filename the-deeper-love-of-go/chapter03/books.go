package books

import (
	"fmt"
)

type Book struct {
	Title  string
	Author string
	Copies int
}

func BookToString(b Book) string {
	return fmt.Sprintf("%v by %v (%v copies)", b.Title, b.Author, b.Copies)
}
