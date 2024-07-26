package models

import (
	"strconv"
	"time"
)

type Year struct {
	numStars    int
	puzzles     []*Puzzle
	leaderboard Leaderboard
}

type Puzzle struct {
	day   int
	year  int
	partA PuzzlePart
	partB PuzzlePart
	URL   string
}

type PuzzlePart struct {
	starObtained  bool
	description   string
	example       string
	isPartB       bool
	submissions   []*Submission
	correctAnswer SubValue
}

type Submission struct {
	submissionVal SubValue
	correct       bool
	timeSubmitted time.Time
	feedback      string
}

type SubValue struct {
	number int
	string string
}

func (v SubValue) GetValue() string {
	if v.string != "" {
		return v.string
	}
	return strconv.Itoa(v.number)
}

type Leaderboard struct {
	year   int
	day    int
	places []*Placing
}

type Placing struct {
	score    int
	username string
}
