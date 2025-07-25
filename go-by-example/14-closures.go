package main

import "fmt"

/*
Go supports anonymous functions, which can form closures. Anonymous functions are useful when you want to define a function inline without having to name it.

From Claude AI
"What's a Closure?
A closure is a function that "captures" variables from its surrounding scope, even after that scope has finished executing.
Think of it as: A function that "remembers" variables from where it was created."

https://en.wikipedia.org/wiki/Anonymous_function
https://en.wikipedia.org/wiki/Closure_(computer_science)

This function intSeq returns another function, which we define anonymously in the body of intSeq. The returned function closes over the variable i to form a closure.
*/
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	// We call intSeq, assigning the result (a function) to nextInt. This function value captures its own i value, which will be updated each time we call nextInt.
	nextInt := intSeq()

	// See the effect of the closure by calling nextInt a few times.
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	// To confirm that the state is unique to that particular function, create and test a new one.
	nextInt2 := intSeq()
	fmt.Println(nextInt2())
}
