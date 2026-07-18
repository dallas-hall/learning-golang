package keyvaluestore

import (
	"encoding/gob"
	"errors"
	"os"
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
