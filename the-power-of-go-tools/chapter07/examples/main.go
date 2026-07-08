package main

import (
	"os"
	"os/exec"
)

func lsNoArgs() {
	// cmd will hold a *exec.Cmd
	cmd := exec.Command("/usr/bin/ls")
	// We need this to see output
	cmd.Stdout = os.Stdout
	// Execute the command
	cmd.Run()
}

func lsArgs() {

	cmd := exec.Command("/usr/bin/ls", "-l")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	lsNoArgs()
	lsArgs()

}
