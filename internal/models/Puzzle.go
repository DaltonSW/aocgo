package models

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/cache"
)

const PUZZLE_URL = "https://adventofcode.com/%v/day/%v"

// Puzzle represents a single day's puzzle.
// It consists of the URL's page data and two correct answers, if they exist.
type Puzzle struct {
	Day      int
	Year     int
	BucketID string
	PageData *PageData
	partA    PuzzlePart
	partB    PuzzlePart
	URL      string
}

func (p *Puzzle) GetID() string                { return p.BucketID }
func (p *Puzzle) GetBucketName() string        { return cache.PUZZLES }
func (p *Puzzle) MarshalData() ([]byte, error) { return json.Marshal(p) }
func (p *Puzzle) SaveResource()                { cache.SaveResource(p) }

// LoadOrCreatePuzzle attempts to load the requested puzzle from
// storage. If it's unable to be loaded, it will attempt to be
// created, loading the information from the website.
func LoadOrCreatePuzzle(year int, day int, userSession string) *Puzzle {
	bucketID := strconv.Itoa(year) + strconv.Itoa(day)
	puzzleData := cache.LoadResource(cache.PUZZLES, bucketID)
	if puzzleData != nil {
		var puzzle *Puzzle
		json.Unmarshal(puzzleData, &puzzle)
		pageData := LoadOrCreatePageData(year, day, userSession, puzzle.URL)
		puzzle.PageData = pageData
		return puzzle
	}

	return NewPuzzle(year, day, userSession)
}

func NewPuzzle(year int, day int, userSession string) *Puzzle {
	URL := fmt.Sprintf(PUZZLE_URL, year, day)
	puzzlePageData := LoadOrCreatePageData(year, day, userSession, URL)

	newPuzzle := &Puzzle{
		Day:      day,
		Year:     year,
		BucketID: strconv.Itoa(year) + strconv.Itoa(day),
		PageData: puzzlePageData,
		URL:      URL,
	}

	newPuzzle.SaveResource()

	return newPuzzle
}

func (p *Puzzle) GetPageDataContent() []string {
	return p.PageData.GetPageDataPrettyString()
}

func (p *Puzzle) GetUserPuzzleInput(userSession string) ([]byte, error) {
	data := cache.LoadResource(cache.USER_INPUTS, p.BucketID)

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

	cache.SaveGenericResource(cache.USER_INPUTS, p.BucketID, inputData)

	return inputData, nil
}

type PuzzlePart struct {
	starObtained  bool
	example       string
	isPartB       bool
	submissions   []*Submission
	CorrectAnswer Value
}

type Submission struct {
	submissionVal Value
	correct       bool
	timeSubmitted time.Time
	feedback      string
}

type Value struct {
	number int
	string string
}

func (v Value) GetValue() string {
	if v.string != "" {
		return v.string
	}
	return strconv.Itoa(v.number)
}
