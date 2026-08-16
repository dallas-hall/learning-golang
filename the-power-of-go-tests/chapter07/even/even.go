package even

import "sync"

// Store previous calls to IsEvenCached into a memory cache.
var (
	mutex = sync.Mutex{}
	cache = map[int]bool{}
)

// We are using this function to showcase automated mutation testing, so we
// aren't using `return n%2 == 0` so the mutation tester has enough code to
// work with.
// Install dependency with `go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest`
// Run with `go-mutesting .`
// The flag `--do-not-remove-tmp-folder` can be used to stop deleting the files
// inside /tmp after the run.
func IsEven(n int) bool {
	if n%2 == 0 {
		return true
	}
	return false
}

func IsEvenCache(n int) (even bool) {
	// Need a mutex lock so t.Parrallel() works correctly
	mutex.Lock()
	defer mutex.Unlock()
	even = n%2 == 0
	if _, ok := cache[n]; !ok {
		cache[n] = n%2 == 0
		// This even is redundant, it is here so the mutation testing can show us
		// it is redundant. ie there will be failing mutation tests showing the
		// first or second even assignment statements being removed.
		even = n%2 == 0
	}
	return even
}
