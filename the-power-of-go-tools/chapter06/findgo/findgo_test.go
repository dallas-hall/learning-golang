package findgo_test

import (
	"archive/zip"
	"findgo"
	"os"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func TestFiles_CorrectlyListFilesInTree(t *testing.T) {
	t.Parallel()

	want := []string{
		"main.go",
		"dir1/file1.go",
		"dir2/file2.go",
		"dir2/file2_test.go",
	}

	// Need to pass in a filesystem
	got := findgo.Files(os.DirFS("test/data"))

	// Sort our slices so they can potentially match
	slices.Sort(want)
	slices.Sort(got)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// Run with `go test -bench .` The dot is a regex wildcard character and will run all benchmarks in the package.
// b *testing.B controls benchmarking from the testing package.
func BenchmarkFilesOnDisk(b *testing.B) {
	filesystem := os.DirFS("test/data")
	// This executes as many times as necessary to make the overall benchmark time about 10 seconds.
	for b.Loop() {
		_ = findgo.Files(filesystem)
	}
}

func TestFiles_CorrectlyListFilesInMapFS(t *testing.T) {
	t.Parallel()

	want := []string{
		"main.go",
		"dir1/file1.go",
		"dir2/file2.go",
		"dir2/file2_test.go",
	}

	// Create a MapFS based filesystem, where the key is filesystem path and the value is the metadata representing the file.
	// The folders aren't included as they are implied from the pathnames with files.
	filesystem := fstest.MapFS{
		"main.go":            {},
		"dir1/file1.go":      {},
		"dir2/file2.go":      {},
		"dir2/file2_test.go": {},
	}

	got := findgo.Files(filesystem)

	// Sort our slices so they can potentially match
	slices.Sort(want)
	slices.Sort(got)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func BenchmarkFilesInMemory(b *testing.B) {
	filesystem := fstest.MapFS{
		"main.go":            {},
		"dir1/file1.go":      {},
		"dir2/file2.go":      {},
		"dir2/file2_test.go": {},
	}

	for b.Loop() {
		_ = findgo.Files(filesystem)
	}
}

func TestFiles_CorrectlyListFilesInZIP(t *testing.T) {
	t.Parallel()

	filesystem, err := zip.OpenReader("test/data/test-data.zip")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"main.go",
		"dir1/file1.go",
		"dir2/file2.go",
		"dir2/file2_test.go",
	}

	// Need to pass in a filesystem, and zip.OpenReader creates one.
	got := findgo.Files(filesystem)

	// Sort our slices so they can potentially match
	slices.Sort(want)
	slices.Sort(got)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func BenchmarkFilesInZIP(b *testing.B) {
	filesystem, err := zip.OpenReader("test/data/test-data.zip")
	if err != nil {
		b.Fatal(err)
	}

	// This executes as many times as necessary to make the overall benchmark time about 10 seconds.
	for b.Loop() {
		_ = findgo.Files(filesystem)
	}
}
