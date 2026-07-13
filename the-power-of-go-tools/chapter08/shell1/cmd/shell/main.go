package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"simpleshell1"
)

func main() {
	// Create a PS1
	user := os.Getenv("USER")
	hostname := os.Getenv("HOSTNAME")
	pwd := filepath.Base(os.Getenv("PWD"))
	var prompt string
	if user == "root" {
		prompt = fmt.Sprintf("[%s@%s %s] # ", user, hostname, pwd)
	} else {
		prompt = fmt.Sprintf("[%s@%s %s] $ ", user, hostname, pwd)
	}

	// Read buffered lines from stdin & print PS1
	input := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for input.Scan() {
		line := input.Text()
		if line == "exit" {
			break
		}
		cmd, err := simpleshell1.CmdFromString(line)
		if err != nil {
			continue
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Printf("%s", output)
		fmt.Print(prompt)
	}
	fmt.Println("\nGoodbye")
}
