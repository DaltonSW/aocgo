package models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/cache"
)

const PUZZLE_URL = "https://adventofcode.com/%v/day/%v"

type Puzzle struct {
	day      int
	year     int
	bucketID string
	pageData *PageData
	partA    PuzzlePart
	partB    PuzzlePart
	URL      string
}

func NewPuzzle(year int, day int) *Puzzle {
	return &Puzzle{
		day:      day,
		year:     year,
		bucketID: strconv.Itoa(day) + strconv.Itoa(year),
		URL:      fmt.Sprintf(PUZZLE_URL, year, day),
	}
}

func (p *Puzzle) GetPuzzlePageData(userSession string) PageData {
	if p.pageData != nil {
		return *p.pageData
	}
	// TODO: Try load from disk
	resp, err := api.NewGetReq(p.URL, userSession)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	rawPage, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	pageData := NewPageData(rawPage)

	p.pageData = pageData

	// TODO: Save to disk
	return *pageData
}

func (p *Puzzle) GetUserPuzzleInput(userSession string) ([]byte, error) {
	data := cache.LoadSubResource(userSession, cache.USER_INPUTS, p.bucketID)

	if data != nil {
		return data, nil
	}

	resp, err := api.NewGetReq(p.URL+"/input", userSession)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	inputData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cache.SaveSubResource(userSession, cache.USER_INPUTS, p.bucketID, inputData)

	return inputData, nil
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
