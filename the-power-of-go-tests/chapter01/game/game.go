package game

import "strings"

func ListItems(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return "You can see " + items[0] + " here."
	case 2:
		return "You can see here " + items[0] + " and " + items[1] + "."
	default:
		return "You can see here " +
			strings.Join(items[:len(items)-1], ", ") +
			", and " + items[len(items)-1] + "."
	}
}

// This only works for items >= 2
func MyListItems(items []string) string {
	result := "You can see here "
	for i, item := range items {
		if i+1 < len(items) {
			result += item + ", "
		} else {
			result += "and " + item + "."
		}
	}
	return result
}
