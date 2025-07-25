package main

import "fmt"

/*
Go supports pointers, allowing you to pass references to values and records within your program.

https://en.wikipedia.org/wiki/Pointer_(computer_programming)

We’ll show how pointers work in contrast to values with 2 functions: zeroval and zeroptr. zeroval has an int parameter, so arguments will be passed to it by value. zeroval will get a copy of ival distinct from the one in the calling function.
*/
func zeroval(ival int) {
	ival = 0
}

// zeroptr in contrast has an *int parameter, meaning that it takes an int pointer. The *iptr code in the function body then dereferences the pointer from its memory address to the current value at that address. Assigning a value to a dereferenced pointer changes the value at the referenced address.
func zeroptr(iptr *int) { // *int means "pointer to int"
	*iptr = 0 // *iptr means "value at this address"
}

func main() {
	i := 1
	fmt.Println("Initial value:", i)

	// Notice how the printed value didn't change to ival
	zeroval(i) // Pass value (copy)
	fmt.Println("Value after zeroval call:", i)

	/*
		The &i syntax gives the memory address of i, i.e. a pointer to i.

		Notice how the printed value changed to *iptr

		zeroval doesn’t change the i in main, but zeroptr does because it has a reference to the memory address for that variable.
	*/
	zeroptr(&i) // Pass pointer with &
	fmt.Println("Value after zeroptr call:", i)

	// Pointers can be printed too.
	fmt.Println("Pointer memory address:", &i)

	/*
		The symbols to remember:

		& = "give me the address of this"
		* = "give me the value at this address"
		*int = "this is a pointer type"

	*/
}
