package books_test

import (
	"books"
	"cmp"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"slices"
	"testing"
)

// Needs a pointer after the changes describe in GetAllBooks()
func getTestCatalogue() *books.Catalogue {
	catalogue := books.NewCatalogue()
	err := catalogue.AddBook(books.Book{
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
		ID:     "abc",
	})
	if err != nil {
		panic(err)
	}

	err = catalogue.AddBook(books.Book{
		Title:  "White Heat",
		Author: "Dominic Sandbrook",
		Copies: 2,
		ID:     "xyz",
	})
	if err != nil {
		panic(err)
	}

	return catalogue
}

func assertTestBooks(t *testing.T, got []books.Book) {
	// This makes Go print the error line from the test, not this helper function.
	t.Helper()
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

func randomLocalAddress(t *testing.T) string {
	t.Helper()
	// Omitting the address will inject "localhost"
	// We use port 0 so the O/S will choose a random free port for us.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	return listener.Addr().String()
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
	books := catalogue.GetAllBooks()
	assertTestBooks(t, books)
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
	err := catalogue.AddBook(books.Book{
		ID:     "123",
		Title:  "The Prize of all the Oceans",
		Author: "Glyn Williams",
		Copies: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, ok = catalogue.GetBook("123")
	if !ok {
		t.Fatal("added book not found")
	}
}

func TestAddBook_ReturnsErrorIfIDExists(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	_, ok := catalogue.GetBook("abc")
	if !ok {
		t.Fatal("book not found")
	}

	err := catalogue.AddBook(books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	})

	if err == nil {
		t.Fatal("want: error for duplicate ID & got: nil")
	}
}

// Original modified a book directly, now uses catalogue instead.
func TestSetCopies_SetsNumberOfCopiesToGivenValue(t *testing.T) {
	t.Parallel()

	catalogue := getTestCatalogue()
	book, ok := catalogue.GetBook("abc")
	if !ok {
		t.Fatal("book not found")
	}

	if book.Copies != 1 {
		t.Fatalf("want: 1 copy before change & got: %d copies", book.Copies)
	}

	err := catalogue.SetCopies("abc", 2)
	if err != nil {
		t.Fatal(err)
	}

	book, ok = catalogue.GetBook("abc")
	if !ok {
		t.Fatal("book not found")
	}

	if book.Copies != 2 {
		t.Fatalf("want: 2 copies after change & got: %d copies", book.Copies)
	}

}

func TestSetCopies_ReturnsErrorIfCopiesNegative(t *testing.T) {
	t.Parallel()
	// Ingore initial values as we don't need them.
	book := books.Book{}
	err := book.SetCopies(-1)
	if err == nil {
		t.Error("want: error for negative copies & got: nil")
	}

	catalogue := getTestCatalogue()
	err = catalogue.SetCopies("abc", -1)
	if err == nil {
		t.Error("want: error for negative copies & got: nil")
	}
}

// Combined TestOpenCatalogue_LoadsCatalogueDataFromFile with TestSyncCatalogue_OverwriteFileWithDataFromMemory.
// Because OpenCatalogue reads the file from disk and SyncCatalogue writes to disk.
// Each function acts a test for each other.
func TestOpenCatalogueAndSyncCatalogues_ReadDataWrittenBySync(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()

	// Create a temporary directory, unique to this test instance, and delete it once the test is complete.
	catalogue.Path = t.TempDir() + "/catalogue-test-data.json"
	err := catalogue.SyncCatalogues()
	if err != nil {
		t.Fatal(err)
	}

	// Test reading catalogue from disk
	newCatalogue, err := books.OpenCatalogue(catalogue.Path)
	if err != nil {
		t.Fatal(err)
	}

	books := newCatalogue.GetAllBooks()
	assertTestBooks(t, books)
}

func TestSetCopies_IsRaceFree(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	// Run an anoymous function as a goroutine, so we can test for race conditions by concurrently reading and writing to the same variable.
	// Run with `go test -race -run SetCopies_IsRaceFree`
	go func() {
		for range 100 {
			err := catalogue.SetCopies("abc", 0)
			if err != nil {
				// Can't use t.Fatal(err) here because it cannot stop other goroutines in the same test.
				// They are running independently via the Go Scheduler.
				panic(err)
			}
		}
	}()

	// Run this for loop in the initial goroutine, creating 2 goroutines with the closure above.
	for range 100 {
		_, err := catalogue.GetCopies("abc")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestNewCatalogue_CreateEmptyCatalogue(t *testing.T) {
	t.Parallel()
	catalogue := books.NewCatalogue()
	books := catalogue.GetAllBooks()
	if len(books) > 0 {
		t.Errorf("want: empty catalogue & got: %#v", books)
	}
}

func TestServerListAllBooks(t *testing.T) {
	t.Parallel()
	catalogue := getTestCatalogue()
	catalogue.Path = t.TempDir() + "/catalogue.json"

	// Get a random address:port combination from our helper
	// This means we can run multiple HTTP servers with different ports for different test instances without conflicts.
	address := randomLocalAddress(t)

	// Create a concurrent HTTP server
	go func() {
		err := books.ListenAndServe(address, catalogue)
		if err != nil {
			panic(err)
		}
	}()

	// Send a HTTP client request
	response, err := http.Get("http://" + address + "/v1/list")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	// Check the server's response
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected HTTP status %d", response.StatusCode)
	}

	// Decode the JSON into a slice of books
	books := []books.Book{}

	// Reading consumes the contents of the reader and you can't get it back. So we save it first in case there is an error and we need to review it.
	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Try to read the JSON into our slice
	err = json.Unmarshal(data, &books)
	if err != nil {
		t.Fatalf("%v in %q", err, data)
	}

	assertTestBooks(t, books)
}

func TestServerGetBook_FindBookInCatalogueByID(t *testing.T) {
	t.Parallel()

	catalogue := getTestCatalogue()
	catalogue.Path = t.TempDir() + "/catalogue.json"

	address := randomLocalAddress(t)
	go func() {
		err := books.ListenAndServe(address, catalogue)
		if err != nil {
			panic(err)
		}
	}()

	response, err := http.Get("http://" + address + "/v1/find/abc")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected HTTP status %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	got := books.Book{}
	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("%v in %q", err, data)
	}

	want := books.Book{
		ID:     "abc",
		Title:  "In the Company of Cheerful Ladies",
		Author: "Alexander McCall Smith",
		Copies: 1,
	}

	if want != got {
		t.Fatalf("want: %#v & got: %#v", want, got)
	}
}

func TestFindReturnsNotFoundWhenBookNotFound(t *testing.T) {
	t.Parallel()

	catalogue := getTestCatalogue()
	catalogue.Path = t.TempDir() + "/catalogue.json"

	address := randomLocalAddress(t)
	go func() {
		err := books.ListenAndServe(address, catalogue)
		if err != nil {
			panic(err)
		}
	}()

	response, err := http.Get("http://" + address + "/v1/find/bogus")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected HTTP status %d", response.StatusCode)
	}

}
