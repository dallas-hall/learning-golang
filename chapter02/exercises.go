package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("# Exercise 1")
	var i = 20
	var f = float64(i)

	fmt.Printf("i: %d\n", i)
	fmt.Printf("f: %f\n", f)

	fmt.Println("# Exercise 2")
	const n = 10
	i2 := n
	var f2 float64 = n

	fmt.Printf("i2: %d\n", i2)
	fmt.Printf("f2: %f\n", f2)

	fmt.Println("# Exercise 3")

	fmt.Println("## Max Values")
	var b byte = math.MaxUint8
	var smallInt int32 = math.MaxInt32
	var bigInt uint64 = math.MaxInt64
	fmt.Printf("b: %d\n", b)
	fmt.Printf("smallInt: %d\n", smallInt)
	fmt.Printf("bigInt: %d\n", bigInt)

	fmt.Println("## Max Values + 1")
	b += 1
	smallInt += 1
	bigInt += 1
	fmt.Printf("b: %d\n", b)
	fmt.Printf("smallInt: %d\n", smallInt)
	fmt.Printf("bigInt: %d\n", bigInt)
}
