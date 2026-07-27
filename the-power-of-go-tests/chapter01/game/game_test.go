package game_test

import (
	"game"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ListItemsGivesCorrectResultForInput(t *testing.T) {
	t.Parallel()
	input := []string{
		"a battery",
		"a key",
		"a tourist map",
	}
	want := "You can see here a battery, a key, and a tourist map."
	got := game.ListItems(input)
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_MyListItemsGivesCorrectResultForInput(t *testing.T) {
	t.Parallel()
	input := []string{
		"a battery",
		"a key",
		"a tourist map",
	}
	want := "You can see here a battery, a key, and a tourist map."
	got := game.MyListItems(input)
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}