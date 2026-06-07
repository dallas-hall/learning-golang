package main

import "fmt"

func main() {
	fmt.Println("# Chapter 3")
	fmt.Println("## Exercise 1")
	greetings := []string{
		"Hello",
		"Hola",
		"नमस्कार",
		"こんにちは",
		"Привіт",
	}

	fmt.Println("### Entire Slice")
	for _, greeting := range greetings {
		fmt.Println(greeting)
	}

	fmt.Println("### First 2 Values")
	greetings2 := greetings[:2]
	for _, greeting := range greetings2 {
		fmt.Println(greeting)
	}

	fmt.Println("### 2nd, 3rd, and 4th Values")
	greetings3 := greetings[1:4]
	for _, greeting := range greetings3 {
		fmt.Println(greeting)
	}

	fmt.Println("### 4th and 5th Values")
	greetings4 := greetings[3:]
	for _, greeting := range greetings4 {
		fmt.Println(greeting)
	}

	fmt.Println("## Exercise 2")
	var message string = "Hi 👦 and 👧"
	fmt.Println(message)

	runes := []rune(message)
	for i := range runes {
		// Print all the rune's bytes, 1-4 depending on UTF-8 code point.
		fmt.Println(runes[i])
		// Convert bytes to UTF-8 string.
		fmt.Println(string(runes[i]))
	}

	fmt.Println("## Exercise 3")
	type Employee struct {
		firstName string
		lastName  string
		id        int
	}

	// struct literal without names
	p1 := Employee{"Alice", "Cooper", 1}

	// struct literal with names
	p2 := Employee{
		firstName: "Billy",
		lastName:  "Joel",
		id:        2,
	}

	// var declaration
	var p3 Employee

	p3.firstName = "Frank"
	p3.lastName = "Sinatra"
	p3.id = 3

	fmt.Printf("p1: %v\n", p1)
	fmt.Printf("p2: %v\n", p2)
	fmt.Printf("p2: %v\n", p3)

	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p2: %+v\n", p2)
	fmt.Printf("p2: %+v\n", p3)

}
