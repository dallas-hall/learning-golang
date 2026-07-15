package main

import (
	"io"
	"log"
	"os"
)

func main() {
	writeBad(os.Stdout)
	write(os.Stdout)

	path := "test/data/output.txt"
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("error opening %q: %s", path, err)
	}
	defer file.Close()

	err = write(file)
	if err != nil {
		log.Fatalf("error writing to %q: %s", path, err)
	}
}

// Started with this monster simulating a pipline, e.g. someObject.writeBad().writeBad().writeBad()...
func writeBad(w io.Writer) error {
	metadata := []byte("hello from writeBad\n")
	_, err := w.Write(metadata)
	if err != nil {
		return err
	}

	_, err = w.Write(metadata)
	if err != nil {
		return err
	}

	_, err = w.Write(metadata)
	if err != nil {
		return err
	}

	return nil
}

// Made a true pipeline of method calls
type safeWriter struct {
	Writer io.Writer
	Reader io.Reader
	Error  error
}

func (sw *safeWriter) Write(data []byte) {
	if sw.Error != nil {
		return
	}
	_, err := sw.Writer.Write(data)
	if err != nil {
		sw.Error = err
	}

}

func write(w io.Writer) error {
	metadata := []byte("hello from write\n")
	sw := safeWriter{Writer: w}
	sw.Write(metadata)
	sw.Write(metadata)
	sw.Write(metadata)
	return sw.Error
}
