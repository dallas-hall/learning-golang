package findgo

import (
	"io/fs"
	"path/filepath"
)

// We are defining our return variable in the signature.
// We replaced the string with a filesystem. Before we needed to use os.DirFS() to convert the string to a filesystem..
func Files(filesystem fs.FS) (goFiles []string) {

	// Walk the filesystem starting at the top, and call our anonymous function to check for .go files.
	fs.WalkDir(filesystem, ".", func(currentPath string, currentPathMetaData fs.DirEntry, err error) error {
		// There should be a currentPathMetaData.IsDir() check here too, as directories containing ".go" are currently being counted.
		if filepath.Ext(currentPath) == ".go" {
			goFiles = append(goFiles, currentPath)
		}
		// Returning an error would break recursion.
		return nil
	})
	return goFiles
}
