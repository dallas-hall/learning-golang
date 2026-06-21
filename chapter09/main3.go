package main

import (
	"fmt"
	"time"
)

// Introduce a package-level variable read by all goroutines.
var msg = "Hello"

func goroutineB() {
	for i := range 10 {
		/*
			Update the package-level variable to introduce a future race condition. A race condition occurs when:
			* Multiple goroutines access the same variable
			* At least one of them writes to that variable.
		*/
		msg = "Goodbye"
		fmt.Println(msg, "from goroutine B.", i)

		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	go goroutineB()
	for i := range 10 {
		fmt.Println(msg, "from goroutine A.", i)
		time.Sleep(10 * time.Millisecond)
	}
}
