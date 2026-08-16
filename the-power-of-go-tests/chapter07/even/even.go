package even

// We are using this function to showcase automated mutation testing, so we
// aren't using `return n%2 == 0` so the mutation tester has enough code to
// work with.
// Install dependency with `go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest`
// Run with `go-mutesting .`
func IsEven(n int) bool {
	if n%2 == 0 {
		return true
	}
	return false
}
