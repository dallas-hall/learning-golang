package books_test

import (
	"books"
	"slices"
	"testing"
)

func TestBookToString_FormatsBookInfoAsString(t *testing.T) {
	// Run the tests in parallel.
	t.Parallel()
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

func TestGetAllBooks_ReturnsAllBooks(t *testing.T) {
	t.Parallel()
	want := []books.Book{
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
	got := books.GetAllBooks()
	if !slices.Equal(want, got) {
		t.Fatalf("want: %#v & got: %#v", want, got)
	}
}

func TestGetBook_FindBookInCatalogueByID(t *testing.T) {
	t.Parallel()
	want := books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	}
	got, ok := books.GetBook("abc")
	if !ok {
		t.Fatal("book not found")
	}
	if want != got {
		t.Fatalf("want: %#v & got: %#v", want, got)
	}
}

func TestGetBook_ReturnsFalseWhenBookNotFound(t *testing.T) {
	t.Parallel()
	// No want because we aren't expecting anything.
	_, ok := books.GetBook("non-existent ID")
	if ok {
		t.Fatal("want: false for non-existent ID & got: true for book found")
	}
}
