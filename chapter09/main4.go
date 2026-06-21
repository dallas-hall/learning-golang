package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a variable to be shared across multiple go routines.
	copies := 1

	/*
	 Create multiple anonymous function literals (AKA closures) and run them as a goroutines. Since both goroutines are trying to read and write to the same variable, this creates a race condition.

	 We can see that race condition with: `go run -race main4.go`

	 It shows the attempted reading and writing by multiple goroutines.
	*/
	go func() {
		if copies > 0 {
			copies--
			fmt.Println("Customer A got the book.")
		}
	}()

	go func() {
		if copies > 0 {
			copies--
			fmt.Println("Customer B got the book.")
		}
	}()
	time.Sleep(250 * time.Millisecond)
}
