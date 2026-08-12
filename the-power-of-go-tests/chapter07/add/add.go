package add

// Deliberately returns the wrong result to show the bad test still passes.
func BadAdd(x, y int) int {
	return x * y
}
