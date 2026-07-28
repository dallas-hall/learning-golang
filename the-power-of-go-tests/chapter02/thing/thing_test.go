package thing_test

import (
	"testing"
	"thing"

	"github.com/google/go-cmp/cmp"
)

// This test is useless as it doesn't actually test if NewThing works. ie we
// are ignoring whatever is being returned by NewThing and not checking if
// the input we passed in is being used.
func Test_UselessThingTest(t *testing.T) {
	t.Parallel()
	_, err := thing.NewThing(1, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
}

// This test actually tests that NewThing uses the inputs that we passed in.
// It also doesn't use t.Fatal on when checking got vs want, because t.Fatal
// will stop the test from completing, and we want to see which cominbation of
// input is failng as that will give us debugging clues.
func Test_NewThingReturnsThingWithGivenInput(t *testing.T) {
	t.Parallel()
	x, y, z := 1, 2, 3

	got, err := thing.NewThing(x, y, z)
	if err != nil {
		t.Fatal(err)
	}

	// This manual checking of a struct is tedious
	if got.X != x {
		t.Errorf("want X: %v & got X: %v", x, got.X)
	}

	if got.Y != y {
		t.Errorf("want Y: %v & got Y: %v", y, got.Y)
	}

	if got.Z != z {
		t.Errorf("want Z: %v & got Z: %v", z, got.Z)
	}
}

// This is the Go way to test complicated structs.
func Test_NewThingReturnsThingWithGivenInputImproved(t *testing.T) {
	t.Parallel()

	x, y, z := 1, 2, 3
	want := &thing.Thing{
		X: x,
		Y: y,
		Z: z,
	}

	got, err := thing.NewThing(x, y, z)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
