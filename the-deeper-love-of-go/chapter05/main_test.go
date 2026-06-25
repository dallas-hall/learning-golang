package books_test

import (
	"books"
	"cmp"
	"slices"
	"testing"
)

func getTestCatalogue() map[string]books.Book {
	return map[string]books.Book{
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
	catalogue := getTestCatalogue()
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
	got := books.GetAllBooks(catalogue)
	// slices.SortFunc takes a slice to sort and then the sorting function.
	// We are comparing to Book structs and returning an it of the order.
	// cmp.Compare does the comparison for us.
	slices.SortFunc(got, func(a, b books.Book) int {
		return cmp.Compare(a.Author, b.Author)
	})
	if !slices.Equal(want, got) {
		t.Fatalf("want: %#v & got: %#v", want, got)
	}
}

func TestGetBook_FindBookInCatalogueByID(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	want := books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	}
	got, ok := books.GetBook(catalogue, "abc")
	if !ok {
		t.Fatal("book not found")
	}
	if want != got {
		t.Fatalf("want: %#v & got: %#v", want, got)
	}
}

func TestGetBook_ReturnsFalseWhenBookNotFound(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	// No want because we aren't expecting anything.
	_, ok := books.GetBook(catalogue, "non-existent ID")
	if ok {
		t.Fatal("want: false for non-existent ID & got: true for book found")
	}
}

func TestAddBook_AddBookToCatalogue(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	// PRE-CONDITION: The book shouldn't exist until we add it.
	_, ok := books.GetBook(catalogue, "123")
	if ok {
		t.Fatal("book already exists")
	}

	// POST-CONDITION: The book should exist after we add it.
	books.AddBook(catalogue, books.Book{
		ID:     "123",
		Title:  "The Prize of all the Oceans",
		Author: "Glyn Williams",
		Copies: 2,
	})

	_, ok = books.GetBook(catalogue, "123")
	if !ok {
		t.Fatal("added book not found")
	}
}
