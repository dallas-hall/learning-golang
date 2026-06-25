package main

import (
	"fmt"
	"time"
)

/*
A goroutine is a lightweight thread of execution, much cheaper than OS threads. You can easily create thousands of them. It is Go's way of doing concurrent programming.
*/
func f(from string) {
	for i := range 10 {
		fmt.Println(from, ":", i)
	}
}

func main() {
	/*
		Suppose we have a function call f(s). Here’s how we’d call that in the usual way, running it synchronously.

		It runs and completes before moving to next line
	*/
	f("direct")

	/*
		To invoke this function in a goroutine, use go f(s). This new goroutine will execute concurrently with the calling one.

		It starts running in background, doesn't wait.
	*/
	go f("goroutine")

	// You can also start a goroutine for an anonymous function call.
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	/*
		Our two function calls are running asynchronously in separate goroutines now. Wait for them to finish (for a more robust approach, use a WaitGroup).

		https://gobyexample.com/waitgroups


		When we run this program, we see the output of the blocking call first, then the output of the two goroutines. The goroutines’ output may be interleaved, because goroutines are being run concurrently by the Go runtime.
	*/
	time.Sleep(time.Second)
	fmt.Println("done")
}
