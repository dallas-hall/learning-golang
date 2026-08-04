package reader_test

import (
	"errors"
	"io"
	"reader"
	"testing"
	"testing/iotest"
)

// Create our own trivial reader spewing errors for testing.
type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func Test_ReadAllReturnsErrorWithOurCustomErrReader(t *testing.T) {
	input := errReader{}
	_, err := reader.ReadAll(input)
	if err == nil {
		t.Errorf("want error from custom errReader, got nil")
	}
}

func Test_ReadAllReturnsErrorWithIOTestErrReader(t *testing.T) {
	t.Parallel()
	// https://pkg.go.dev/testing/iotest#ErrReader
	input := iotest.ErrReader(errors.New("my custom reader error"))
	// Throw away the result because it will be nil or something useless.
	_, err := reader.ReadAll(input)
	if err == nil {
		t.Error("want error from iotest.ErrReader, got nil")
	}
}
