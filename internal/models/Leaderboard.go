package models

import "fmt"

type Leaderboard struct {
	year   int
	day    int
	places []*Placing
}

type Placing struct {
	score    int
	username string
	userID   int
	href     string
	placing  int
}

func (l *Leaderboard) Display() {
	// Try load from storage. If can't, load from API

}

func (l *Leaderboard) tryLoadFromDisk() {

}

func (l *Leaderboard) tryLoadFromWeb() {
	URL := fmt.Sprintf("https://adventofcode.com/%v/leaderboard", l.year)
	if l.day != 0 {
		URL += fmt.Sprintf("/day/%v", l.day)
	}

}

func (l *Leaderboard) saveToDisk() {

}
