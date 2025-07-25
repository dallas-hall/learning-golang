package main

import "fmt"

/*
Go supports embedding of structs and interfaces to express a more seamless composition of types. This is not to be confused with //go:embed which is a go directive introduced in Go version 1.16+ to embed files and folders into the application binary.

https://gobyexample.com/embed-directive

This is Go's way of achieving composition and a form of inheritance
*/

type base struct {
	num int
}

func (b base) describe() string {
	return fmt.Sprintf("Base with number: %v", b.num)
}

// A container embeds a base. An embedding looks like a field without a name.
type container struct {
	base // This is embedding - no field name, just the type
	str  string
}

func main() {
	// When creating structs with literals, we have to initialize the embedding explicitly; here the embedded type serves as the field name.
	c := container{
		base: base{
			num: 1,
		},
		str: "My base's name",
	}

	// We can access the base’s fields directly on c, e.g. c.num.
	fmt.Printf("Container: {num: %v, str: %v}\n", c.num, c.str)

	// Alternatively, we can spell out the full path using the embedded type name.
	fmt.Println("c.base.num:", c.base.num)

	/*
	 Since container embeds base, the methods of base also become methods of a container. Here we invoke a method that was embedded from base directly on container c.

	 Since container embeds base, and base has a describe() method, container automatically satisfies the describer interface.
	*/
	fmt.Println("Describe:", c.describe())

	type describer interface {
		describe() string
	}

	// Embedding structs with methods may be used to bestow interface implementations onto other structs. Here we see that a container now implements the describer interface because it embeds base.
	var d describer = c
	fmt.Println("Describer:", d.describe())

	/*
		Key Benefits

		* Composition over inheritance: You compose types by embedding rather than traditional inheritance
		* Automatic method promotion: Embedded methods become available on the outer type
		* Interface satisfaction: Embedding can make types satisfy interfaces they don't explicitly implement
		* Clean syntax: Access embedded fields as if they were direct fields

	*/
}
