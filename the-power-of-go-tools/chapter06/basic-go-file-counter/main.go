package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func main() {
	// https://stackoverflow.com/a/47261747
	folder, err := filepath.Abs("../../..")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result := countGoFiles(folder, 0)
	fmt.Printf("Go files found: %d\n", result)
}

// This function walks the directory tree manually using recursion (rather than filepath.Walk or filepath.WalkDir, which would be the more idiomatic way to do this).
func countGoFiles(folder string, count int) int {
	// Lists everything (files and subdirectories) directly inside folder.
	files, err := os.ReadDir(folder)
	// Base case, give up if there is an error.
	if err != nil {
		return count
	}

	// It recurses into folder+"/"+f.Name() — treating every entry as if it might be a directory.
	// If f is actually a regular file (not a directory), os.ReadDir on that path will fail, hit the base case above, and just return the count unchanged.
	// This means the recursion is a bit wasteful — it doesn't check f.IsDir() first to decide whether to recurse, it just always tries and lets the error handle it.
	for _, f := range files {
		count = countGoFiles(folder+"/"+f.Name(), count)
		// It checks path.Ext(f.Name()) — using the path package's Ext function (which works the same as filepath.Ext, just intended for URL-style paths rather than OS paths, so filepath.Ext would be the more "correct" choice on Windows).
		// If the entry's name ends in .go, it increments the count. This check applies whether f is a file or a directory — so a directory literally named something like mystuff.go would (incorrectly) be counted too.
		if path.Ext(f.Name()) == ".go" {
			count++
		}
	}
	return count
}
