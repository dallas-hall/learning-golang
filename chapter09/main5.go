package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mutex sync.Mutex
	copies := 1

	go func() {
		/*
			Fix the race condition described in main4.go by adding a mutually exclusive lock on the copies variable.

			The lock can only be obtained by one goroutine at a time, thus ensuring only one goroutine can read/write to the locked variable.

			The lock is relinquished when the function finishes.

			We can see this fixes the race condition with: go run -race main5.go
		*/
		mutex.Lock()
		defer mutex.Unlock()
		if copies > 0 {

			copies--
			fmt.Println("Customer A got the book.")
		}
	}()

	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		if copies > 0 {
			copies--
			fmt.Println("Customer B got the book.")
		}
	}()
	time.Sleep(250 * time.Millisecond)
}
