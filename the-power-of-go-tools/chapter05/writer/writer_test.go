package writer_test

import (
	"os"
	"testing"
	"writer"

	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestWriteToFile_WritesGivenDataToFile(t *testing.T) {
	t.Parallel()

	// Delete our test file when the test function exits at any point.
	path := "test/data/write_test.txt"
	defer os.Remove(path)

	// https://stackoverflow.com/a/66405130 cautions against this approach.
	_, err := os.Stat(path)
	if err == nil {
		t.Fatalf("test artifact not cleaned up: %q", path)
	}

	want := []byte{1, 2, 3}
	err = writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestWriteToFile_WritesGivenDataToTemporaryFile(t *testing.T) {
	t.Parallel()

	// Create a temporary directory, unique to this test instance, and delete it once the test is complete.
	path := t.TempDir() + "/write_test.txt"
	want := []byte{1, 2, 3}
	err := writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}

	// Check permissions are correct. We hardcoded 0600 inside of WriteToFile
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	permission := stat.Mode().Perm()
	if permission != 0o600 {
		t.Errorf("want: file permission 0o600 & got: file permission 0o%o", permission)
	}
}

func TestWriteToFile_ReturnsErrorForUnwriteableFile(t *testing.T) {
	t.Parallel()
	path := "fake/path/to/test/data/write_test.txt"
	err := writer.WriteToFile(path, []byte{})
	if err == nil {
		t.Fatal("want: error for file not writeable & got: writeable file")
	}
}

func TestWriteToFile_OverwritesFileSuccessfully(t *testing.T) {
	t.Parallel()

	// Create a new file manually.
	path := t.TempDir() + "/clobber_test.txt"
	err := os.WriteFile(path, []byte{4, 5, 6}, 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Overwrite existing file.
	want := []byte{1, 2, 3}
	err = writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestWriteToFile_ExistingFilePermissionsAreCorrect(t *testing.T) {
	t.Parallel()

	// Create a new insecure file manually.
	path := t.TempDir() + "/permissions_test.txt"

	// Looks like os.WriteFile changes 0o777 to 0o755 which is a nice touch.
	err := os.WriteFile(path, []byte{4, 5, 6}, 0o777)
	if err != nil {
		t.Fatal(err)
	}

	// Overwrite existing file.
	want := []byte{1, 2, 3}
	err = writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}

	// Check permissions are correct. We hardcoded 0600 inside of WriteToFile
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	permission := stat.Mode().Perm()
	if permission != 0o600 {
		t.Errorf("want: file permission 0o600 & got: file permission 0o%o", permission)
	}
}

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"writefile": writer.Main,
	})
}

func Test(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "test/scripts",
	})
}
