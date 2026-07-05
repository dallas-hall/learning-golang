package writer

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func WriteToFile(path string, data []byte) error {
	// Create a new file if it doesn't exist, or truncate it if it does.
	// We also need to pass in chmod permissions. These can be written with standard chmod 0644 or octal notation 0o644
	err := os.WriteFile(path, data, 0o600)
	if err != nil {
		return err
	}

	// Ensure correct file permissions because os.WriteFile doesn't update permissions for existing files.
	/*
		This code can be replaced with a one liner

		err = os.Chmod(path, 0o600)
		if err != nil {
			return err
		}
		return nil
	*/

	// Return nil if successful or return the error.
	return os.Chmod(path, 0o600)
}

func Main() {
	size := flag.IntP("size", "s", 0, "Write zereos into a file.")

	// Update the -h|--help output.
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-s|--size] files...\n", os.Args[0])
		fmt.Println("Write zereos into a file(s).")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	// At the very least we have:
	// 0 = program name
	// 1 = -s or --size
	// 2 = int for size
	// 3 = file path
	if len(os.Args) < 4 {
		flag.Usage()
		return
	}

	// Parse non-flag arguments.
	flag.Parse()

	// The default value for byte (uint8) is 0, so we don't need to do anything.
	data := make([]byte, *size)

	// Create a zeroed out file for every non-flag argument
	for _, path := range flag.Args() {
		err := WriteToFile(path, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
