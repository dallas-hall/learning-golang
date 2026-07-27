package game

import "strings"

func ListItems(items []string) string {
	result := "You can see here "
	result += strings.Join(items[:len(items)-1], ", ")
	result += ", and "
	result += items[len(items)-1]
	result += "."
	return result
}

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
