package main

import (
	"greeting"
	"os"
)

func main() {
	// Only the outermost Main function touches concrete OS types; everything else operates on interfaces, which is exactly the seam you need for unit testing without spawning real processes or faking stdin at the OS level.
	greeting.GreetUser(os.Stdin, os.Stdout)
}
