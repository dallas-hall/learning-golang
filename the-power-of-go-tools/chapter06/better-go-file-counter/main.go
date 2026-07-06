package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	var count int

	startFolder, err := filepath.Abs("../../..")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wrap the real OS directory at startFolder into an fs.FS value.
	filesystem := os.DirFS(startFolder)

	// fs.WalkDir takes three things: the filesystem to walk, a starting path ("." means "the root of filesystem"), and a callback function. It handles all the recursion internally — visiting every file and directory in the tree — and calls your function once per entry.
	fs.WalkDir(filesystem, ".", func(currentPath string, currentPathMetaData fs.DirEntry, err error) error {
		// There should be a currentPathMetaData.IsDir() check here too, as directories containing ".go" are currently being counted.
		if filepath.Ext(currentPath) == ".go" {
			count++
		}
		// Returning an actual error would stop the entire walk early and propagate that error out of fs.WalkDir. i.e. breaking recursion.
		return nil
	})

	fmt.Printf("Go files found: %d\n", count)
}
