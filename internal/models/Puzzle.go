package models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"dalton.dog/aocutil/internal/api"
)

const PUZZLE_URL = "https://adventofcode.com/%v/day/%v"

type Puzzle struct {
	day   int
	year  int
	desc  string
	partA PuzzlePart
	partB PuzzlePart
	URL   string
}

func NewPuzzle(year int, day int) *Puzzle {
	return &Puzzle{
		day:  day,
		year: year,
		URL:  fmt.Sprintf(PUZZLE_URL, year, day),
	}
}

func (p *Puzzle) GetPuzzlePageData() []byte {
	// TODO: Try load from disk
	resp, err := api.NewGetReq(p.URL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	pageData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// TODO: Save to disk
	return pageData
}

func (p *Puzzle) GetUserPuzzleInput(userSession string) []byte {
	resp, err := api.NewGetReq(p.URL + "/input")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	inputData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return inputData
}

type PuzzlePart struct {
	starObtained  bool
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
