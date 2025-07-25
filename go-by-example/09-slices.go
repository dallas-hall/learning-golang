package main

import (
	"fmt"
	"slices"
)

func main() {
	/*
		Slices are an important data type in Go, giving a more powerful interface to sequences than arrays.

		Unlike arrays, slices are typed only by the elements they contain (not the number of elements). An uninitialized slice equals to nil and has length 0.

		Claude AI: In Go, arrays are fixed sized collections, slices are dynamic sized collections.
	*/
	var s []string
	//fmt.Println("Uninitialised slice:", s, s == nil, len(s) == 0)
	fmt.Println("Uninitialised slice:", s)
	fmt.Println("s == nil? ", s == nil)
	fmt.Println("len(s) == 0? ", len(s) == 0)

	/*
	  To create a slice with non-zero length, use the builtin make. Here we make a slice of strings of length 3 (initially zero-valued).

	  By default a new slice’s capacity is equal to its length; if we know the slice is going to grow ahead of time, it’s possible to pass a capacity explicitly as an additional parameter to make.
	*/
	s = make([]string, 3)
	//fmt.Println("Slice:", s, "Length:", len(s), "Capacity:", cap(s))
	fmt.Println("Slice:", s)
	fmt.Println("Length:", len(s))
	fmt.Println("Capacity:", cap(s))

	// We can set and get just like with arrays.
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("Set:", s)
	fmt.Println("Get:", s[2])

	// len returns the length of the slice as expected.
	fmt.Println("Length:", len(s))

	/*
		In addition to these basic operations, slices support several more that make them richer than arrays.

		One is the builtin append, which returns a slice containing one or more new values.

		Note that we need to accept a return value from append as we may get a new slice value.
	*/
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("Appended slice:", s)

	// Slices can also be copy’d. Here we create an empty slice c of the same length as s and copy into c from s.
	s_copy := make([]string, len(s))
	copy(s_copy, s)
	fmt.Println("Slice copy:", s_copy)

	/*
		Slices support a “slice” operator with the syntax slice[low:high].

		e.g. this gets a slice of the elements s[2], s[3], and s[4].
	*/
	s_slice := s[2:5]
	fmt.Println("A slice of an existing slice:", s_slice)

	// This slices up to (but excluding) s[5].
	s_slice = s[:5]
	fmt.Println("A slice of an existing slice:", s_slice)

	// And this slices up from (and including) s[2].
	s_slice = s[2:]
	fmt.Println("A slice of an existing slice:", s_slice)

	// We can declare and initialize a variable for slice in a single line as well.
	s2 := []string{"x", "y", "z"}
	fmt.Println("Declared slice:", s2)

	// The slices package contains a number of useful utility functions for slices.
	s3 := []string{"x", "y", "z"}
	if slices.Equal(s2, s3) {
		fmt.Println("s2 == s3? ", slices.Equal(s2, s3))
	}

	// Slices can be composed into multi-dimensional data structures. The length of the inner slices can vary, unlike with multi-dimensional arrays.

	slice_2d := make([][]int, 3)
	for i := range 3 {
		length := i + 1
		slice_2d[i] = make([]int, length)
		for j := range length {
			slice_2d[i][j] = i + j
		}
	}
	fmt.Println("2D slice:", slice_2d)

	/*
		Note that while slices are different types than arrays, they are rendered similarly by fmt.Println.

		Check out the great blog post by the Go team for more details on the design and implementation of slices in Go. https://go.dev/blog/slices-intro
	*/
}
