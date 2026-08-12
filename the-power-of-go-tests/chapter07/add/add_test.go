package add_test

import (
	"add"
	"testing"
)

// Bad test as it can't fail, as there is no call to t.Error or t.Fatal. It is
// also not testing the function, just calling it. The name doesn't describe
// what the test is trying to achieve either.
// `go test` passes
// `go test -cover` shows 100% coverage, even though the result is wrong.
func TestBadAdd(t *testing.T) {
	t.Parallel()
	add.BadAdd(2, 2)
}
