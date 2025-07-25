package main

import "fmt"

/*
Starting with version 1.18, Go has added support for generics, also known as type parameters.

Generics let you create containers or functions that can work with any type you specify, while still being type-safe.

As an example of a generic function, SlicesIndex takes a slice of any comparable type and an element of that type and returns the index of the first occurrence of v in s, or -1 if not present. The comparable constraint means that we can compare values of this type with the == and != operators. For a more thorough explanation of this type signature, see this blog post. Note that this function exists in the standard library as slices.Index.

https://go.dev/blog/deconstructing-type-parameters

https://pkg.go.dev/slices#Index

Type Parameters: [S ~[]E, E comparable]
* S ~[]E: S is a type that has the underlying type of slice of E (~ means "underlying type")
* E comparable: E must be a comparable type (can use == and !=)

Function signature: Takes a slice of type S and a value of type E, returns the index or -1.
*/
func SlicesIndex[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

/*
As an example of a generic type, List is a singly-linked list with values of any type.

Type Parameter: [T any]
* T can be any type (any is the least restrictive constraint)
* Both List and element are parameterized by type T
*/
type List[T any] struct {
	head, tail *element[T] // Pointers to first AND last elements
}

type element[T any] struct {
	next  *element[T] // Pointer to next element
	value T           // The actual data
}

// We can define methods on generic types just like we do on regular types, but we have to keep the type parameters in place. The type is List[T], not List.
func (list *List[T]) Push(v T) {
	if list.tail == nil { // Empty list
		list.head = &element[T]{value: v} // Head is the new element.
		list.tail = list.head             // Tail is now head, the new element.
	} else {
		list.tail.next = &element[T]{value: v} // Create a new element and attach it to the current tail's next pointer.
		list.tail = list.tail.next             // Move tail to the new element.
	}
}

// AllElements returns all the List elements as a slice. In the next example we’ll see a more idiomatic way of iterating over all elements of custom types.
func (list *List[T]) AllElements() []T {
	var elements []T
	for e := list.head; e != nil; e = e.next {
		elements = append(elements, e.value)
	}
	return elements
}

func main() {
	var s = []string{"foo", "bar", "zoo"}

	// // Type inference - compiler figures out the types
	fmt.Println("Index of zoo:", SlicesIndex(s, "zoo"))

	/*
		When invoking generic functions, we can often rely on type inference. Note that we don’t have to specify the types for S and E when calling SlicesIndex - the compiler infers them automatically.

		… though we could also specify them explicitly.
	*/
	_ = SlicesIndex[[]string, string](s, "zoo") // Explicit type specification (optional)
	list := List[int]{}                         // Specify T = int
	list.Push(10)                               // Push accepts int values
	list.Push(20)
	list.Push(30)
	fmt.Println("List:", list.AllElements())
}
