package main

import "fmt"

func main() {
	/*
		In Go, an array is a numbered sequence of elements of a specific length. In typical Go code, slices are much more common; arrays are useful in some special scenarios.

		Here we create an array a that will hold exactly 5 ints. The type of elements and length are both part of the array’s type. By default an array is zero-valued, which for ints means 0s.
	*/
	var a [5]int
	fmt.Println("Array 1:", a)

	// We can set a value at an index using the array[index] = value syntax, and get a value with array[index].
	a[4] = 100
	fmt.Println("Set:", a)
	fmt.Println("Get:", a[4])

	// The builtin len returns the length of an array.
	fmt.Println("Length:", len(a))

	// Use this syntax to declare and initialize an array in one line.
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Array 2:", b)

	// You can also have the compiler count the number of elements for you with ...
	b = [...]int{1, 2, 3, 4, 5}
	fmt.Println("Array 2:", b)

	// If you specify the index with : then the elements in between will be zeroed.
	b = [...]int{100, 3: 400, 500}
	fmt.Println("Array 2:", b)

	// Array types are one-dimensional, but you can compose types to build multi-dimensional data structures.
	var c [2][3]int
	for i := range 2 {
		for j := range 3 {
			c[i][j] = i + j
		}
	}
	fmt.Println("2d Array 1: ", c)

	// You can create and initialize multi-dimensional arrays at once too.
	c = [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}

	// Note that arrays appear in the form [v1 v2 v3 ...] when printed with fmt.Println.
	fmt.Println("2d Array 2: ", c)
}
