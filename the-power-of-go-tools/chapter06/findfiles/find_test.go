package find_test

import (
	"archive/zip"
	"find"
	"os"
	"slices"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Run with `go test`
func TestFiles_ReturnsFilesFromMapFSOlderThanGivenDuration(t *testing.T) {
	t.Parallel()

	now := time.Now()

	want := []string{
		"dir1/file1.go",
		"dir2/file2_test.go",
	}

	filesystem := fstest.MapFS{
		"main.go":            {ModTime: now},
		"dir1/file1.go":      {ModTime: now.Add(-time.Hour * 24)},
		"dir2/file2.go":      {ModTime: now},
		"dir2/file2_test.go": {ModTime: now.Add(-time.Hour * 24)},
	}

	got := find.Files(filesystem, time.Hour*24)

	slices.Sort(want)
	slices.Sort(got)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// Run with `go test -bench .`
func BenchmarkFilesInMemory(b *testing.B) {
	now := time.Now()

	filesystem := fstest.MapFS{
		"main.go":            {ModTime: now},
		"dir1/file1.go":      {ModTime: now.Add(-time.Hour * 24)},
		"dir2/file2.go":      {ModTime: now},
		"dir2/file2_test.go": {ModTime: now.Add(-time.Hour * 24)},
	}

	for b.Loop() {
		_ = find.Files(filesystem, time.Hour*24)
	}
}

func TestFiles_ReturnsFilesFromDiskOlderThanGivenDuration(t *testing.T) {
	t.Parallel()

	want := []string{
		"main.go",
		"dir1/file1.go",
		"dir2/file2.go",
		"dir2/file2_test.go",
	}

	got := find.Files(os.DirFS("test/data"), time.Hour*12)

	slices.Sort(want)
	slices.Sort(got)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func BenchmarkFilesOnDisk(b *testing.B) {
	filesystem := os.DirFS("test/data")
	for b.Loop() {
		_ = find.Files(filesystem, time.Hour*12)
	}
}

func TestFiles_ReturnsFilesFromZipOlderThanGivenDuration(t *testing.T) {
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
	got := find.Files(filesystem, time.Hour*12)

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

	for b.Loop() {
		_ = find.Files(filesystem, time.Hour*12)
	}
}
