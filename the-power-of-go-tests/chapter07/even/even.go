package even

// We are using this function to showcase automated mutation testing, so we
// aren't using `return n%2 == 0` so the mutation tester has enough code to
// work with.
func IsEven(n int) bool {
	if n%2 == 0 {
		return true
	}
	return false
}
