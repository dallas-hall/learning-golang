package keyvaluestore

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	// https://pkg.go.dev/github.com/spf13/pflag
	flag "github.com/spf13/pflag"
)

// data holds the key/value data, which are [string]string
// path holds the location to where the key/value store can be saved.
type keyValueStore struct {
	data map[string]string
	path string
}

// Create a new key/value store.
// Taken from the-deeper-love-of-go/chapter13/books.go OpenCatalogue() - see for comments
func OpenStore(s string) (*keyValueStore, error) {
	kvs := &keyValueStore{
		data: make(map[string]string), // same as map[string]string{}
		path: s,
	}

	// Try to open existing kvs file, if present load the data into a kvs.
	// If not present then create a new kvs file and return new kvs.
	file, err := os.Open(kvs.path)
	if errors.Is(err, os.ErrNotExist) {
		return kvs, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = gob.NewDecoder(file).Decode(&kvs.data)
	if err != nil {
		return nil, err
	}
	return kvs, nil
}

// Return the key/value store's data.
func (kvs *keyValueStore) GetData() map[string]string {
	return kvs.data
}

// Return the key/value store's path.
func (kvs *keyValueStore) GetPath() string {
	return kvs.path
}

// Set the key/value store's path.
func (kvs *keyValueStore) SetPath(p string) {
	kvs.path = p
}

// Return the length of the key/value store, i.e. how many keys are stored.
func (kvs *keyValueStore) GetLength() int {
	return len(kvs.data)
}

// Save the key/value store to disk using path.
// Taken from the-deeper-love-of-go/chapter13/books.go SyncCatalogues() - see for comments
func (kvs *keyValueStore) Save() error {
	file, err := os.Create(kvs.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(kvs.data)
}

// Add a key/value pair into the map.
func (kvs *keyValueStore) Set(key string, value string) {
	kvs.data[key] = value
}

// Try to return the value associated with the given key. If found, return the
// value and true, if not found return empty string and false.
func (kvs *keyValueStore) Get(key string) (value string, ok bool) {
	value, ok = kvs.data[key]
	return value, ok
}

// Delete a key value pair. No operation performed when the key doesn't exist.
func (kvs *keyValueStore) Delete(key string) {
	delete(kvs.data, key)
}

// Remove all key/value data.
func (kvs *keyValueStore) Clear() {
	clear(kvs.data)
}

func (kvs *keyValueStore) All() map[string]string {
	return kvs.data
}

func Main() {
	// Taken from the-power-of-go-tools/chapter04/count-pflag/count.go see there for comments.
	setMode := flag.StringSliceP("set", "s", []string{}, "Add/update a key/value pair with key=value. Flag can be repeated with -a k1=v1,k2=v2 or a -a k1=v1 -a k2=v2.")
	deleteMode := flag.StringP("delete", "d", "", "Delete a key/value pair.")
	clearMode := flag.BoolP("clear", "c", false, "Delete all key/value pairs.")
	getMode := flag.StringP("get", "g", "", "Get value using its key.")
	allMode := flag.BoolP("all", "a", false, "Get all key/value pairs.")

	// Update the -h output.
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options...] [files...]\n", os.Args[0])
		fmt.Println("Count words (or lines or bytes) from stdin (or files).")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	// Check the command line for arguments and assign them to our matching flags.
	// This stops parsing as soon as it see a non-flag arg.
	flag.Parse()

	path := os.Getenv("KVS_PATH")
	if path == "" {
		path = "test/data/data.bin"
	}
	kvs, err := OpenStore(path)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case len(*setMode) > 0:
		for _, kv := range *setMode {
			k, v, found := strings.Cut(kv, "=")
			if !found {
				log.Fatal("add expects key=value")
			}
			kvs.Set(k, v)
		}
	case *deleteMode != "":
		kvs.Delete(*deleteMode)
	case *clearMode:
		kvs.Clear()
	case *getMode != "":
		v, ok := kvs.Get(*getMode)
		if !ok {
			fmt.Printf("key %q not found", *getMode)
		}
		fmt.Println(v)
	case *allMode:
		fmt.Println(kvs.All())
	default:
		flag.Usage()
	}
	err = kvs.Save()
	if err != nil {
		log.Fatalf("could not save kvs because %s", err)
	}
}
