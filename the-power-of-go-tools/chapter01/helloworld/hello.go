package helloworld

/*
1. Write this program first.
2. Run `go mod init helloworld` to create `go.mod`
3. Run `go mod tidy`
4. Write cmd/main.go which imports `helloworld` package and calls this fucntion.
5. Run with `go run ./cmd/`
*/

import (
	"fmt"
	"io"
)

// Untestable.
func Print() {
	fmt.Println("Hello world.")
}

// Testable.
func PrintTo(w io.Writer) {
	fmt.Fprintln(w, "Hello world.")
}