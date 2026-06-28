package main

import (
	"helloworld"
	"os"
)

func main() {
	// Static message using fmt.Println
	helloworld.Print()
// Static message usng Fprintln which accepts any io.Writer, in this case the terminal.
	helloworld.PrintTo(os.Stdout)
}
