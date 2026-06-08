package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("# Chapter 4")
	fmt.Println("## Exercise 1")
	//myInts := []int{}
	myInts := make([]int, 0, 100)
	for i := 0; i < 16; i++ {
		prn := rand.Intn(100)
		myInts = append(myInts, prn)
	}
	fmt.Println(myInts)

	fmt.Println("## Exercise 2")
	fmt.Println("### For Loop + If")
	for i := 0; i < len(myInts); i++ {
		if myInts[i]%2 == 0 {
			fmt.Println("Two!")
		} else if myInts[i]%3 == 0 {
			fmt.Println("Three!")
		} else if myInts[i]%6 == 0 {
			fmt.Println("Six!")
		}
		fmt.Println("Nevermind.")
	}

	fmt.Println("### For Range Loop + Blank Switch")
	for _, v := range myInts {
		switch {
		case v%2 == 0:
			fmt.Println("Two!")
		case v%3 == 0:
			fmt.Println("Three!")
		case v%6 == 0:
			fmt.Println("Six!")
		default:
			fmt.Println("Nevermind.")
		}
	}

	fmt.Println("## Exercise 3")
	var total int // Same as `var total int = 0`
	fmt.Println("### Total Pre-loop")
	fmt.Println(total)
	fmt.Println("### Shadow Total During Loop")
	for i := 1; i < 10; i++ {
		//Shadow total values will be los
		total := total + i
		fmt.Println(total)
	}
	fmt.Println("### Total Post-loop")
	fmt.Println(total)
}
