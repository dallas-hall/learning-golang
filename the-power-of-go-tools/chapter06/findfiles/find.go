package find

import (
	"io/fs"
	"time"
)

// Callers can pass in a real directory (os.DirFS(...)), a zip archive, an embedded filesystem, or a fake one in tests.
// age time.Duration lets the caller say "find files older than this much time" — e.g. 24 * time.Hour
func Files(filesystem fs.FS, age time.Duration) (matches []string) {
	// Converts "how old" into "how recent," which is what you actually need to compare against.
	// If age is 24 hours, threshold becomes "the moment 24 hours ago." Any file modified before that moment counts as old.
	// Note the -age — negating the duration to subtract it from the current time.
	threshold := time.Now().Add(-age)

	fs.WalkDir(filesystem, ".", func(currentPath string, currentPathMetaData fs.DirEntry, err error) error {
		// .Info() returns an fs.FileInfo and an error.
		info, err := currentPathMetaData.Info()
		// Silently skip errors and directories
		if err != nil || info.IsDir() {
			return nil
		}

		// Check if the file's last-modified time is before the cutoff.
		if info.ModTime().Before(threshold) {
			matches = append(matches, currentPath)
		}
		return nil
	})
	return matches
}
