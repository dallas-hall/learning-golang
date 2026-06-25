package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Exercise 1
type Team struct {
	Name    string
	Players []string // Slice of player names
}
type League struct {
	Name  string
	Teams map[string]Team // key = team name and value = Team struct
	Wins  map[string]int  // key = team name and value = wins
}

// Exercise 2
func (t Team) ToString() string {
	return fmt.Sprintf("%s %v", t.Name, t.Players)
}

// League pointer is needed so the function updates the passed in League
func (l *League) MatchResult(team1 string, score1 int, team2 string, score2 int) {
	// Does the team exist.
	_, ok := l.Teams[team1]
	if !ok {
		return
	}

	// The classic "comma-ok" idiom check. Same as above
	if _, ok := l.Teams[team2]; !ok {
		return
	}
	// Handle draw
	if score1 == score2 {
		return
	}
	// Handle wins
	if score1 > score2 {
		l.Wins[team1]++
	} else {
		l.Wins[team2]++
	}
}

// No need for a League pointer as we as just using the values, not mutating them.
func (l League) Ranking() []string {
	// Create an empty string slice with 0 values and pre-allocated capacity equal to the number of teams.
	teamNames := make([]string, 0, len(l.Teams))
	// Loop through Teams and grab all the team names.
	for team := range l.Teams {
		teamNames = append(teamNames, team)
	}

	// Sort our teamNames slice in-place.
	// sort.Slice handles all the swapping internally. The Use an anonymous
	// function tells sort.Slice how to compare the 2 passed in items.
	// The anonymous function needs to be `func(i, j int) bool`
	// This anonymous function is also closure because it can use teamNames
	// and l even though it wasn't passed into the anonymous function, it gets
	// and uses these variables from the Ranking function.
	//
	// Think of the closure as a backpack. When the anonymous function is
	// created, it packs up any outside variables it references (teamNames, l)
	// and carries them with it wherever it goes.
	sort.Slice(teamNames, func(i, j int) bool {
		return l.Wins[teamNames[i]] > l.Wins[teamNames[j]]
	})

	return teamNames
}

// Exercise 3
// The Ranking() method matches this interface.
type Ranker interface {
	Ranking() []string
}

// io.Writer — anything you can write bytes to (a file, terminal, network connection, etc.)
func RankPrinter(r Ranker, w io.Writer) {
	// Pretty print the teams sorted by descending wins.
	rankings := r.Ranking()
	for i, team := range rankings {
		// Use i + 1 to give a ranking, does not handle ties.
		io.WriteString(w, fmt.Sprintf("%d: %s", i+1, team))
		w.Write([]byte("\n"))
	}
}

func main() {
	fmt.Println("# Chapter 7")
	fmt.Println("## Exercise 1")
	fmt.Println("Types defined.")

	fmt.Println("## Exercise 2")
	l := League{
		Name: "The Big League",
		Teams: map[string]Team{
			"Australia": {
				Name:    "Australia",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
			"New Zealand": {
				Name:    "New Zealand",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
			"Indonesia": {
				Name:    "Indonesia",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
			"Papua New Guinea": {
				Name:    "Papua New Guinea",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
		},
		Wins: map[string]int{},
	}
	fmt.Printf("Teams & Players in %s\n", l.Name)
	for team := range l.Teams {
		fmt.Printf("* %s\n", l.Teams[team].ToString())
	}

	l.MatchResult("Australia", 2, "New Zealand", 1)
	l.MatchResult("Indonesia", 1, "Papua New Guinea", 0)
	l.MatchResult("Australia", 3, "Indonesia", 1)
	l.MatchResult("New Zealand", 3, "Papua New Guinea", 1)
	l.MatchResult("Australia", 5, "Papua New Guinea", 0)
	l.MatchResult("New Zealand", 3, "Indonesia", 2)
	results := l.Ranking()
	fmt.Printf("Current rankings: %s\n", results)

	fmt.Println("## Exercise 3")
	fmt.Println("Pretty print rankings.")
	RankPrinter(l, os.Stdout)
}
