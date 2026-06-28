package helloworld_test

import (
	"bytes"
	"helloworld"
	"testing"
)

func TestPrintTo_PrintsHelloMessageToGivenWriter(t *testing.T) {
	t.Parallel()
	buffer := new(bytes.Buffer)
	// Call our new constructor
	printer := helloworld.NewPrinter()
	// Update this test's printer instance to use our buffer instead of os.Stdout default.
	printer.Output = buffer
	// Run the usual tests.
	printer.Print()
	want := "Hello world.\n"
	got := buffer.String()
	if want != got {
		t.Errorf("want: %q & got: %q", want, got)
	}
}
