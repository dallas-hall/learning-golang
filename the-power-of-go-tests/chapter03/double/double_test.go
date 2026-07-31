package double_test

import (
	"double"
	"fmt"
	"testing"
)

func Test_Double2Returns4(t *testing.T) {
	t.Parallel()

	want := 4
	got := double.Double(2)
	if want != got {
		t.Errorf("want %d & got %d", want, got)
	}
}

// This function gets included in the package documentation on pkg.go.dev and
// it will be available to run in a sandbox.
//
// The `// Output:` has to be correct or this will fail when running `go test`.
//
// You can also use `// Unordered output:` for printing things like a map, where
// the output is not deterministic.
//
// The function after Example must exist, as that is the function used by the
// documentation tool.
func ExampleDouble() {
	fmt.Println(double.Double(2))
	// Output:
	// 4
}

// You can give multiple examples of the same function, add a suffix after the
// function name.
func ExampleDouble_with2() {
	fmt.Println(double.Double(2))
	// Output:
	// 4
}

func ExampleDouble_with3() {
	fmt.Println(double.Double(3))
	// Output:
	// 6
}

// This demonstrates how to use the entire package
func Example() {

}

// Examples work with types as well
func ExampleNumber() {
	n := double.Number{Value: 2}
	fmt.Println(n)
	// Output:
	// {2}
}

// Examples work on methods on types. The syntax is:
// Example + Type Name + underscore + method name
func ExampleNumber_ValueString() {
	n := double.Number{Value: 2}
	fmt.Println(n.ValueString())
	// Output:
	// 2
}
