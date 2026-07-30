package double

func Double(n int) int {
	// Add a temporary change to the function to skew the test result, this will
	// check if the test is actually detecting bugs in the function. eg:
	// return n * 2 + 1
	return n*2 + 1
}
