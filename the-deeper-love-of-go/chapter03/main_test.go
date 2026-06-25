package books_test

import (
	"testing"
	"books"
)

func TestBookToString_FormatsBookInfoAsString(t *testing.T) {
	input := books.Book{
		Title:  "Sea Room",
		Author: "Adam Nicolson",
		Copies: 2,
	}

	want := "Sea Room by Adam Nicolson (2 copies)"
	got := books.BookToString(input)
	if want != got {
		//panic("BookToString: wrong result") // Print and exit the program.
		//t.Fatal("BookToString: wrong result") // Run with go test
		t.Fatalf("want: %q & got: %q", want, got)
	}
}
