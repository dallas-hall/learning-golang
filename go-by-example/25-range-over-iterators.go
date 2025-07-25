package main

import (
	"fmt"
	"iter"
	"slices"
)

/*
Starting with version 1.23, Go has added support for iterators, which lets us range over pretty much anything!
https://go.dev/blog/range-functions

An iterator in Go is a function with this signature:
func(yield func(T) bool)

It takes a yield function as a parameter and calls yield for each item it wants to produce. The yield function returns bool - if it returns false, the iterator should stop.

Let’s look at the List type from the previous example again. In that example we had an AllElements method that returned a slice of all elements in the list. With Go iterators, we can do it better - as shown below.
*/
type List[T any] struct {
	head, tail *element[T] // Pointers to first AND last elements
}

type element[T any] struct {
	next  *element[T] // Pointer to next element
	value T           // The actual data
}

func (list *List[T]) Push(v T) {
	if list.tail == nil { // Empty list
		list.head = &element[T]{value: v} // Head is the new element.
		list.tail = list.head             // Tail is now head, the new element.
	} else {
		list.tail.next = &element[T]{value: v} // Create a new element and attach it to the current tail's next pointer.
		list.tail = list.tail.next             // Move tail to the new element.
	}
}

/*
All returns an iterator, which in Go is a function with a special signature.

https://pkg.go.dev/iter#Seq

The iterator function takes another function as a parameter, called yield by convention (but the name can be arbitrary). It will call yield for every element we want to iterate over, and note yield’s return value for a potential early termination.
*/
func (list *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := list.head; e != nil; e = e.next {
			if !yield(e.value) { // Call yield with each value
				return // Stop if yield returns false
			}
		}
	}
}

// Iteration doesn’t require an underlying data structure, and doesn’t even have to be finite! Here’s a function returning an iterator over Fibonacci numbers: it keeps running as long as yield keeps returning true.
func genFib() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 1, 1

		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func main() {
	list := List[int]{}
	list.Push(10)
	list.Push(20)
	list.Push(30)

	// Since List.All returns an iterator, we can use it in a regular range loop.
	for e := range list.All() {
		fmt.Println(e)
	}

	/*
		Packages like slices have a number of useful functions to work with iterators. For example, Collect takes any iterator and collects all its values into a slice.

		https://pkg.go.dev/slices
	*/
	all := slices.Collect(list.All())
	fmt.Println("All:", all)

	for n := range genFib() {
		if n >= 10 {
			// Once the loop hits break or an early return, the yield function passed to the iterator will return false.
			break
		}
		fmt.Println(n)
	}
}
