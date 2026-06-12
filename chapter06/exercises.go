package main

import (
	"fmt"
)

// Exercise 1
type Person struct {
	Firstname string
	LastName  string
	Age       int
}

func MakePerson(firstName, lastName string, age int) Person {
	return Person{
		Firstname: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

func MakePersonPointer(firstName, lastName string, age int) *Person {
	return &Person{
		Firstname: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// Exercise 2
func UpdateSlice(a []string, b string) {
	a[len(a)-1] = b
	fmt.Println("Inside UpdateSlice:", a)
}

func GrowSlice(a []string, b string) {
	a = append(a, b)
	fmt.Println("Inside GrowSlice:", a)
}

func main() {
	fmt.Println("# Chapter 6")
	fmt.Println("## Exercise 1")
	p1 := MakePerson("John", "Doe", 32)
	p2 := MakePersonPointer("Jane", "Doe", 24)
	fmt.Printf("p1: %v\n", p1)
	fmt.Printf("p2: %v\n", p2)

	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p2: %+v\n", p2)

	// Run with:
	// go build -gcflags="-m"
	// Results explained https://github.com/learning-go-book-2e/ch06/blob/main/exercise_solutions/ex1/README.md

	fmt.Println("## Exercise 2")
	a := []string{"a", "b", "c"}
	b := "C"
	fmt.Println("### UpdateSlice")
	fmt.Println("Before UpdateSlice: ", a)
	UpdateSlice(a, b) // I think we are changing the original slice because we are passing the pointer in?
	fmt.Println("After UpdateSlice: ", a)

	fmt.Println("### GrowSlice")
	c := []string{"x", "y", "z"}
	d := "Z"
	fmt.Println("Before GrowSlice: ", c)
	GrowSlice(c, d) // I think we aren't changing the original slice because we create a copy of it inside the function and never return it?
	fmt.Println("After GrowSlice: ", c)

	fmt.Println("## Exercise 3")
	limit := 10_000_000
	people := make([]Person, 0, limit)
	for i := 0; i < limit; i++ {
		people = append(people, MakePerson("John", "Doe", 32))
	}
}
