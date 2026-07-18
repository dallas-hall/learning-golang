package keyvaluestore_test

import (
	"keyvaluestore"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tests if the constructor for creating a key/value store is working. It also
// tests GetPath(), SetPath(), and GetDataLength().
func TestKvs_OpenStoreReturnsKeyValueStoreCorrectly(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/lol.bin"

	kvs, err := keyvaluestore.OpenStore(path)
	if err != nil {
		t.Fatalf("error opening store: %s", err)
	}

	want := path
	got := kvs.GetPath()
	if !cmp.Equal(want, got) {
		t.Fatalf("want %q path & got %q", want, got)
	}

	length := kvs.GetLength()
	if length != 0 {
		t.Fatalf("want 0 lenght & got %d", length)
	}

	want = "test/data/wheelie.bin"
	kvs.SetPath(want)
	got = kvs.GetPath()
	if !cmp.Equal(want, got) {
		t.Fatalf("want %q path & got %q", want, got)
	}
}

// Test when a searched kv pair exists and returns value.
func TestKvs_SetValueAndGetValueAndOkWhenExists(t *testing.T) {
	t.Parallel()

	path := "test/data/read-test.bin"

	kvs, err := keyvaluestore.OpenStore(path)
	if err != nil {
		t.Fatalf("error opening store: %s", err)
	}

	// "a" already exists
	want := "B"
	kvs.Set("b", "B")
	v, ok := kvs.Get("b")
	if !ok {
		t.Fatal("not ok")
	}
	if v != want {
		t.Fatalf("want %q & got %q", want, v)
	}
}

// Test overriding an existing kv pair and returning the updated value.
func TestKvs_SetAndGetUpdatedValueAndOkCorrectly(t *testing.T) {
	t.Parallel()

	path := "test/data/read-test.bin"

	kvs, err := keyvaluestore.OpenStore(path)
	if err != nil {
		t.Fatalf("error opening store: %s", err)
	}

	// "a" already exists
	want := "B"
	kvs.Set("a", "B")
	v, ok := kvs.Get("a")
	if !ok {
		t.Fatal("not ok")
	}
	if v != want {
		t.Fatalf("want %q & got %q", want, v)
	}
}

// Test when a searched kv pair doesn't exist.
func TestKvs_GetNothingAndNotOkWhenDoesntExist(t *testing.T) {
	t.Parallel()

	path := "test/data/read-test.bin"

	kvs, err := keyvaluestore.OpenStore(path)
	if err != nil {
		t.Fatalf("error opening store: %s", err)
	}

	// "a" already exists
	v, ok := kvs.Get("b")
	if ok {
		t.Fatal("unexpected ok")
	}
	if v != "" {
		t.Fatalf("want empty string & got %q", v)
	}
}

func TestKvs_SaveAndOpenWorksCorrectly(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "write-test.bin"

	// Create a new kvs, add entries, and save to disk.
	kvs1, err := keyvaluestore.OpenStore(path)
	if err != nil {
		t.Fatalf("error opening kvs1: %s", err)
	}

	kvs1.Set("a", "A")
	kvs1.Set("b", "B")
	kvs1.Set("c", "C")
	err = kvs1.Save()
	if err != nil {
		t.Fatalf("error saving kvs1: %s", err)
	}

	// Create a new kvs using the previous kvs file on disk.
	kvs2, err := keyvaluestore.OpenStore(path)
	if err != nil {
		t.Fatalf("error opening kvs2: %s", err)
	}

	// Compare the 2 kvs's, they should match.
	if diff := cmp.Diff(kvs1.GetData(), kvs2.GetData()); diff != "" {
		t.Errorf("mismatch (-kvs1 +kvs2):\n%s", diff)
	}
}
