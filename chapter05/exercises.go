package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Exercise 1
func add(i int, j int) (int, error) { return i + j, nil }
func sub(i int, j int) (int, error) { return i - j, nil }
func mul(i int, j int) (int, error) { return i * j, nil }
func div(i int, j int) (int, error) {
	if j == 0 {
		return 0, errors.New("division by zero")
	}
	return i / j, nil
}

var opMap = map[string]func(int, int) (int, error){
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
}

// Exercise 2
func fileLen(fileName string) (int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	data := make([]byte, 2048)
	total := 0
	for {
		count, err := f.Read(data)
		total += count
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			break
		}
	}
	return total, nil
}

// Exercise 3
func prefixer(prefix string) func(string) string {
	return func(body string) string {
		return prefix + " " + body
	}
}

func main() {
	fmt.Println("# Chapter 5")
	fmt.Println("## Exercise 1")
	expressions := [][]string{
		[]string{"2", "+", "3"},
		[]string{"2", "-", "3"},
		[]string{"2", "*", "3"},
		[]string{"2", "/", "3"},
		[]string{"2", "%", "3"},
		[]string{"two", "+", "three"},
		[]string{"5"},
		[]string{"10", "/", "0"},
	}
	for _, expression := range expressions {
		if len(expression) != 3 {
			fmt.Println("invalid expression:", expression)
			continue
		}
		p1, err := strconv.Atoi(expression[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		op := expression[1]
		opFunc, ok := opMap[op]
		if !ok {
			fmt.Println("unsupported operator:", op)
			continue
		}
		p2, err := strconv.Atoi(expression[2])
		if err != nil {
			fmt.Println(err)
			continue
		}
		result, err := opFunc(p1, p2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(result)
	}

	fmt.Println("## Exercise 2")
	count, err := fileLen("exercises")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)

	fmt.Println("## Exercise 3")
	helloPrefix := prefixer("Hello")
	fmt.Println(helloPrefix("Bob"))   // should print Hello Bob
	fmt.Println(helloPrefix("Maria")) // should print Hello Maria
}
