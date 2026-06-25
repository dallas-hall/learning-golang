package main

import (
	"fmt"
	"time"
)

func goroutineB() {
	for i := range 10 {
		fmt.Println("Hello from goroutine B.", i)
		/*
			Fix the problem described in main1.go - the messages become interleaved because their execution overlaps.

			The output is non-deterministic, meaning it will change with each run so we cannot predict what order is displayed.
		*/
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	/*

	 */
	go goroutineB()
	for i := range 10 {
		fmt.Println("Hello from goroutine A.", i)
		time.Sleep(10 * time.Millisecond)
	}
}
