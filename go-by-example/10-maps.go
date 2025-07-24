package main

import (
	"fmt"
	"maps"
)

func main() {
	/*
		Maps are Go’s built-in associative data type (sometimes called hashes or dicts in other languages).

		https://en.wikipedia.org/wiki/Associative_array

		Note that maps appear in the form map[k:v k:v] when printed with fmt.Println.
	*/

	// To create an empty map, use the builtin make: make(map[key-type]val-type).
	m := make(map[string]int)

	// Set key/value pairs using typical name[key] = val syntax.
	m["key1"] = 10
	m["key2"] = 20

	// Printing a map with e.g. fmt.Println will show all of its key/value pairs.
	fmt.Println("New map:", m)

	// Get a value for a key with name[key].
	v1 := m["key1"]
	fmt.Println("Value 1:", v1)

	/*
		If the key doesn’t exist, the zero value of the value type is returned. No value is added

		https://go.dev/ref/spec#The_zero_value
	*/
	v3 := m["k3"]
	fmt.Println("Value 3:", v3)

	// The builtin len returns the number of key/value pairs when called on a map.
	fmt.Println("Length: ", len(m))

	// The builtin delete removes key/value pairs from a map.
	delete(m, "key2")
	fmt.Println("Map after deleting a key:", m)

	// To remove all key/value pairs from a map, use the clear builtin.
	clear(m)
	fmt.Println("Map after clear:", m)

	// 	The optional second return value when getting a value from a map indicates if the key was present in the map. This can be used to disambiguate between missing keys and keys with zero values like 0 or "". Here we didn’t need the value itself, so we ignored it with the blank identifier _.
	_, present_key := m["key2"]
	fmt.Println("Is key2 present?", present_key)

	// You can also declare and initialize a new map in the same line with this syntax.
	m2 := map[string]int{"foo": 100, "bar": 200}
	fmt.Println("New Map", m2)

	// The maps package contains a number of useful utility functions for maps.
	m3 := map[string]int{"foo": 100, "bar": 200}
	if maps.Equal(m2, m3) {
		fmt.Println("m2 == m3?", maps.Equal(m2, m3))
	}
}
