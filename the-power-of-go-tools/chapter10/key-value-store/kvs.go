package keyvaluestore

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"maps"
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

// Create a new key/value store. This type is unexported so you must use this
// constructor to get a keyValueStore with sane defaults. The fields are
// unexported too, so you must use the exported methods (e.g. SetPath, Set)
// This is trying to ake illegal states unreachable through this type system.
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

// Return a copy of the key/value store's data.
func (kvs *keyValueStore) Data() map[string]string {
	return maps.Clone(kvs.data)
}

// Return the key/value store's path.
func (kvs *keyValueStore) Path() string {
	return kvs.path
}

// Set the key/value store's path.
func (kvs *keyValueStore) SetPath(p string) {
	kvs.path = p
}

// Return the length of the key/value store, i.e. how many keys are stored.
func (kvs *keyValueStore) Length() int {
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

func Main() {
	// Taken from the-power-of-go-tools/chapter04/count-pflag/count.go see there for comments.
	setMode := flag.StringSliceP("set", "s", []string{}, "Add/update a key/value pair with key=value. Flag can be repeated with -s k1=v1,k2=v2 or -s k1=v1 -s k2=v2.")
	deleteMode := flag.StringP("delete", "d", "", "Delete a key/value pair.")
	clearMode := flag.BoolP("clear", "c", false, "Delete all key/value pairs.")
	getMode := flag.StringP("get", "g", "", "Get value using its key.")
	allMode := flag.BoolP("all", "a", false, "Get all key/value pairs.")

	// Update the -h output.
	flag.Usage = func() {
		fmt.Printf("Usage: %s option\n", os.Args[0])
		fmt.Println("Create and manipulate a key/value store.")
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

	var save bool
	switch {
	case len(*setMode) > 0:
		for _, kv := range *setMode {
			k, v, found := strings.Cut(kv, "=")
			if !found {
				log.Fatal("add expects key=value")
			}
			kvs.Set(k, v)
		}
		save = true
	case *deleteMode != "":
		kvs.Delete(*deleteMode)
		save = true
	case *clearMode:
		kvs.Clear()
		save = true
	case *getMode != "":
		v, ok := kvs.Get(*getMode)
		if !ok {
			fmt.Printf("key %q not found\n", *getMode)
			return
		}
		fmt.Println(v)
	case *allMode:
		fmt.Println(kvs.Data())
	default:
		flag.Usage()
	}
	if save {
		err = kvs.Save()
		if err != nil {
			log.Fatalf("could not save kvs because %s", err)
		}
	}

}
