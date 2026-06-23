package books

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"sync"
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
//type Catalogue map[string]Book

// Replacing `type Catalogue map[string]Book` so we can add a mutex to stop race conditions
// Added Path so we don't have to pass it around everywhere.
// The lowercase field name means it is private, direct access is internal to the books package. Indirect access is via methods.
// The Uppercaes field name means it is public, anyone can directly access and update.
type Catalogue struct {
	// The address of the shared mutex literal
	// This mutex gives us a read lock which can be shared, and an exclusive write lock.
	mutex *sync.RWMutex
	data  map[string]Book
	// Save the path
	Path string
}

// func GetCatalogue() Catalogue {
// 	return Catalogue{
// 		"abc": {
// 			Title:  "In the Company of Cheerful Ladies",
// 			Author: "Alexander McCall Smith",
// 			Copies: 1,
// 			ID:     "abc",
// 		},
// 		"xyz": {
// 			Title:  "White Heat",
// 			Author: "Dominic Sandbrook",
// 			Copies: 2,
// 			ID:     "xyz",
// 		},
// 	}
// }

// This function is replacing the package-level  catalogue variable above. And also replacing the hardcoded GetCatalogue() above.
// Now returns a pointer to Catalogue after the mutex changes described in GetAllBooks()
func OpenCatalogue(path string) (*Catalogue, error) {
	// Try to open the path into a file handle object
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// Instead of writing file.Close() at all exit points, we can use the defer keyword to do it for us.
	// It will automatically close the file handle object whenever this function exits.
	defer file.Close()

	// Updated after the changes described in GetAllBooks()
	//catalogue := Catalogue{}
	catalogue := NewCatalogue()

	// json.NewDecoder uses the file hanndle object to create a new decoder.
	// json.NewDecoder.Decode reads the JSON from the file handle object and stores it into the pointer.
	err = json.NewDecoder(file).Decode(&catalogue.data)
	if err != nil {
		return nil, err
	}

	// Save the path if successful
	catalogue.Path = path

	// Need to return the address of the pointer now after the updates described in GetAllBooks(). The address now comes from NewCatalogue()
	return catalogue, nil
}

// Synchronises the in-memory catalogue with the catalogue on disk. Writes in-memory to disk.
// Now takes a pointer after the updates described in GetAllBooks()
func (catalogue *Catalogue) SyncCatalogues() error {
	catalogue.mutex.RLock()
	defer catalogue.mutex.RUnlock()

	// Create anew file if it doesn't exist, or truncate and write to if it exists.
	file, err := os.Create(catalogue.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// json.NewEnecoder uses the file hanndle object to create a new encoder.
	// json.NewEncoder.Encode writes the data as JSON into the file.
	err = json.NewEncoder(file).Encode(catalogue.data)
	if err != nil {
		return err
	}

	return nil
}

// Create a new catalogue that is ready to be used. Path can be set later if needed.
func NewCatalogue() *Catalogue {
	// We need to return the pointer address of the Catalogue we just created. Because of the changes described in GetAllBooks()
	return &Catalogue{
		mutex: &sync.RWMutex{},
		data:  map[string]Book{},
	}
}

// Originally this `catalogue Catalogue` but we updated it to a struct with a mutex. We must use a pointer for a mutex so every method will use the same mutex and its locks. So every other method is updated too.
func (catalogue *Catalogue) GetAllBooks() []Book {
	// Adding a read lock that allows multiple read access
	catalogue.mutex.RLock()
	defer catalogue.mutex.RUnlock()

	// maps.Values returns an iterator
	// slices.Collect conumes the iterator and collects all elements into a slice.
	return slices.Collect(maps.Values(catalogue.data))
}

func (catalogue *Catalogue) GetBook(ID string) (Book, bool) {
	catalogue.mutex.RLock()
	defer catalogue.mutex.RUnlock()
	book, ok := catalogue.data[ID]
	return book, ok
}

func (catalogue *Catalogue) AddBook(book Book) error {
	// Adding a writing lock that only allows a single write access
	catalogue.mutex.Lock()
	defer catalogue.mutex.Unlock()

	_, ok := catalogue.data[book.ID]
	if ok {
		return fmt.Errorf("ID %q already exists", book.ID)
	}
	catalogue.data[book.ID] = book
	return nil
}

func (catalogue *Catalogue) GetCopies(ID string) (int, error) {
	catalogue.mutex.RLock()
	defer catalogue.mutex.RUnlock()
	book, ok := catalogue.data[ID]
	if !ok {
		return 0, fmt.Errorf("ID %q not found", ID)
	}
	return book.Copies, nil
}

func (catalogue *Catalogue) SetCopies(ID string, copies int) error {
	catalogue.mutex.Lock()
	defer catalogue.mutex.Unlock()
	book, ok := catalogue.data[ID]
	if !ok {
		return fmt.Errorf("ID %q not found", ID)
	}

	err := book.SetCopies(copies)
	if err != nil {
		return err
	}

	catalogue.data[ID] = book
	return nil
}
