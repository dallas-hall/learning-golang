package main

import "fmt"

func goroutineB() {
	for i := range 10 {
		fmt.Println("Hello from goroutine B.", i)
	}
}

func main() {
	/*
		This call to goroutineB never gets executed because it gets "starved" by the goroutine run by main(), which immediately exits the program when it finishes.

		The Go scheduler looks for goroutines in the Ready, Running, or Blocked state.
		* Ready Queue = Run when possible
		* Running = currently running
		* Blocked Queue = was running but got blocked by something (eg waiting for I/O). Go back to ready queue when unblocked.

		This gives us concurrency. Parallelism is acheived by running multiple goroutines simultaneously across multiple CPU cores/threads.
	*/
	go goroutineB()
	for i := range 10 {
		fmt.Println("Hello from goroutine A.", i)
	}
}
