package models

import (
	"io"
	"strconv"
	"time"

	"dalton.dog/aocutil/internal/api"
)

type Puzzle struct {
	day   int
	year  int
	partA PuzzlePart
	partB PuzzlePart
	URL   string
}

func (p *Puzzle) GetPuzzlePageData() []byte {
	resp, err := api.NewGetReq(p.URL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	pageData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

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
