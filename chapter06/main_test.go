package books_test

import (
	"books"
	"cmp"
	"slices"
	"testing"
)

func getTestCatalogue() books.Catalogue {
	return books.Catalogue{
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
	got := input.String()
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
	got := catalogue.GetAllBooks()
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
	got, ok := catalogue.GetBook("abc")
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
	_, ok := catalogue.GetBook("non-existent ID")
	if ok {
		t.Fatal("want: false for non-existent ID & got: true for book found")
	}
}

func TestAddBook_AddBookToCatalogue(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	// PRE-CONDITION: The book shouldn't exist until we add it.
	_, ok := catalogue.GetBook("123")
	if ok {
		t.Fatal("book already exists")
	}

	// POST-CONDITION: The book should exist after we add it.
	catalogue.AddBook(books.Book{
		ID:     "123",
		Title:  "The Prize of all the Oceans",
		Author: "Glyn Williams",
		Copies: 2,
	})

	_, ok = catalogue.GetBook("123")
	if !ok {
		t.Fatal("added book not found")
	}
}

func TestSetCopies_SetsNumberOfCopiesToGivenValue(t *testing.T) {
	t.Parallel()
	book := books.Book{
		// Ignore other values as they are useless in this test
		Copies: 5,
	}
	err := book.SetCopies(12)
	if err != nil {
		t.Fatal(err)
	}
	if book.Copies != 12 {
		t.Errorf("want: 12 copies & got: %d copies", book.Copies)
	}
}

func TestSetCopies_ReturnsErrorIfCopiesNegative(t *testing.T) {
	t.Parallel()
	// Ingore initial values as we don't need them.
	book := books.Book {}
	err := book.SetCopies(-1)
	if err == nil {
		t.Error("want: error for negative copies & got: nil")
	}
}