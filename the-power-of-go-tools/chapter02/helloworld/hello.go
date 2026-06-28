package helloworld

import (
	"fmt"
	"io"
	"os"
)

// Creating a struct ensures this is concurrency safe.
// Each instance of Printer has its own Output.
type Printer struct {
	Output io.Writer
}

// Using a constructor the set the default values for Printer. Because Go has no way of doing that natively.
func NewPrinter() *Printer {
	return &Printer{
		Output: os.Stdout,
	}
}

func (p *Printer) Print() {
	fmt.Fprintln(p.Output, "Hello world.")
}

// This is a convenience wrapper so user's don't need to make multiple function calls. They just make one, this one.
// net/http package uses this approach. e.g:
//
// client := http.Client{
// Timeout: 10 * time.Second,
// }
// resp, err := client.Get("https://example.com")
func Main() {
	NewPrinter().Print()
}
