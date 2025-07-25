package main

import "fmt"

/*
Go’s structs are typed collections of fields. They’re useful for grouping data together to form records.

Claude AI:

"What Are Structs?
Go structs are like Python classes but focused purely on data grouping - no inheritance, no complex methods, just fields grouped together.
Think of them as: Custom data containers that group related information"

This person struct type has name and age fields.
*/
type person struct {
	name string
	age  int
}

/*
newPerson constructs a new person struct with the given name.

Go is a garbage collected language; you can safely return a pointer to a local variable - it will only be cleaned up by the garbage collector when there are no active references to it.
*/
func newPerson(name string) *person {
	p := person{name: name}
	p.age = 42
	return &p
}

func main() {
	// This syntax creates a new struct.
	fmt.Println(person{"bob", 20}) // Must match field order

	// You can name the fields when initializing a struct.
	fmt.Println(person{name: "alice", age: 30})

	// Omitted fields will be zero-valued.
	fmt.Println(person{name: "fred"}) // age gets zero value (0)
	fmt.Println(person{age: 18})      // name gets zero value ("")
	fmt.Println(person{})             // Zero value: {name: "", age: 0}

	// An & prefix yields a pointer to the struct.
	fmt.Println(&person{name: "ann", age: 40})

	// It’s idiomatic to encapsulate new struct creation in constructor functions
	s := person{name: "sean", age: 50}
	// Access struct fields with a dot.
	fmt.Println(s.name)
	fmt.Println(s.age)

	// You can also use dots with struct pointers - the pointers are automatically dereferenced. i.e. don't need (*pointer)
	sp := &s             // Pointer to struct
	fmt.Println(sp)      // Equivalent to (*sp)
	fmt.Println(sp.name) // Equivalent to (*sp).name
	fmt.Println(sp.age)  // Equivalent to (*sp).age

	// Structs are mutable.
	sp.age = 51
	fmt.Println(sp.age)

	/*
		If a struct type is only used for a single value, we don’t have to give it a name. The value can have an anonymous struct type. This technique is commonly used for table-driven tests.

		https://gobyexample.com/testing-and-benchmarking
	*/
	dog := struct {
		name   string
		isGood bool
	}{
		"rex",
		true,
	}
	fmt.Println(dog)
}
