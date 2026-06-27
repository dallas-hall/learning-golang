// What pacakge this file belongs to.
package helloworld_test

// What packages to import.
import (
	"bytes"
	"helloworld"
	"testing"
)

// The testing functions.
func TestPrintTo_PrintsHelloMessageToGivenWriter(t *testing.T) {
	// Run tests concurrently.
	t.Parallel()
	// An all purpose writer that remembers what we write.
	buffer := new(bytes.Buffer)
	// Static message usng Fprintln which accepts any io.Writer, in this case the bytes.Buffer which will remember the hardcoded string.
	helloworld.PrintTo(buffer)
	want := "Hello world.\n"
	got := buffer.String()
	if want != got {
		t.Errorf("want: %q & got: %q", want, got)
	}
}
