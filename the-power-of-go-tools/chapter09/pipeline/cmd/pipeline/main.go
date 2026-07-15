package main

import "pipeline"

func main() {
	pipeline.Main()
	pipeline.FromString("hello from pipeline.FromString\n").Stdout()
}
