package main

import (
	"fmt"
)

func hello(s string) {
	fmt.Println(s)
}

func main() {
	fmt.Println("# Chapter 1")
	fmt.Println("Hello world, from main()")
	var msg = "Hello world, from hello()"
	hello(msg)
	msg = "Another hello world, from hello()"
	hello(msg)
}
