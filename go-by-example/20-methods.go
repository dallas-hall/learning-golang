package main

import "fmt"

type rect struct {
	width, height int
}

/*
Go supports methods defined on struct types, and defined outside the struct types. You explicitly choose whether the method works on a copy or the original using pointer vs value receivers.

This area method has a receiver type of *rect.

The (r *rect) syntax means "this method belongs to the rect type and receives a pointer to a rect instance."
*/
func (r *rect) area() int {
	return r.width * r.height
}

// Methods can be defined for either pointer or value receiver types. Here’s an example of a value receiver which works on a value copy.
func (r rect) perim() int {
	return (2 * r.width) + (2 * r.height)
}

// Can modify the original struct with a pointer receiever.
func (r *rect) scale(factor int) {
	r.width *= factor  // Modifies original
	r.height *= factor // Modifies original
}

func main() {
	r := rect{width: 10, height: 5}
	fmt.Println("Rectangle:", r)
	// Here we call the 2 methods defined for our struct.
	fmt.Println("Area:", r.area())       // Direct value call
	fmt.Println("Perimeter:", r.perim()) // Go converts r to &r automatically

	// Go automatically handles conversion between values and pointers for method calls. You may want to use a pointer receiver type to avoid copying on method calls or to allow the method to mutate the receiving struct.
	rp := &r
	rp.scale(2)
	fmt.Println("Rectangle after scale:", rp)
	fmt.Println("Area:", rp.area())       // Direct pointer call
	fmt.Println("Perimeter:", rp.perim()) // Go converts *rp to value automatically
}
