package models

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/tui"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// Base URL for a single day's puzzle
const PUZZLE_URL = "https://adventofcode.com/%v/day/%v"

// Puzzle represents a single day's puzzle.
// Consists of user info as well as page display info.
type Puzzle struct {
	SessionToken string

	Day      int
	Year     int
	BucketID string
	URL      string

	Title      string
	ArticleOne []string
	AnswerOne  string
	ArticleTwo []string
	AnswerTwo  string

	UserInput []byte
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
		return puzzle
	}

	return newPuzzle(year, day, userSession)
}

func newPuzzle(year int, day int, userSession string) *Puzzle {
	URL := fmt.Sprintf(PUZZLE_URL, year, day)
	bucketID := strconv.Itoa(year) + strconv.Itoa(day)

	userInput, err := loadUserInputFromSite(URL, userSession)

	if err != nil {
		log.Fatal("Unable to load user input for the puzzle.", "error", err)
	}

	newPuzzle := &Puzzle{
		Day:          day,
		Year:         year,
		BucketID:     bucketID,
		URL:          URL,
		UserInput:    userInput,
		SessionToken: userSession,
	}

	newPuzzle.loadPageData()
	newPuzzle.SaveResource()

	return newPuzzle
}

// GetUserInput returns the input for the associated puzzle.
func (p *Puzzle) GetUserInput() ([]byte, error) {
	if p.UserInput != nil {
		return p.UserInput, nil
	}

	input, err := loadUserInputFromSite(p.URL, p.SessionToken)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (p *Puzzle) GetPrettyPageData() []string {
	sOut := p.ArticleOne

	if p.AnswerOne != "" {
		sOut = append(sOut, "\n"+p.AnswerOne)
	}

	if len(p.ArticleTwo) != 0 {
		sOut = append(sOut, "\n\n")
		sOut = append(sOut, "\n"+titleStyle.Render("--- Part Two ---"))
		sOut = append(sOut, "\n")
		sOut = append(sOut, p.ArticleTwo...)
		sOut = append(sOut, "\n")

		if p.AnswerTwo != "" {
			sOut = append(sOut, p.AnswerTwo)
		}
	}
	return sOut
}

func loadUserInputFromSite(URL, userSession string) ([]byte, error) {
	resp, err := api.NewGetReq(URL+"/input", userSession)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	inputData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return inputData, nil
}

// loadPageData will make the HTTP request and pass it off to be parsed.
func (p *Puzzle) loadPageData() {
	resp, err := api.NewGetReq(p.URL, p.SessionToken)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error constructing new PageData.", "error", err)
	}

	mainContents := doc.Find("main")

	p.processPageContents(mainContents)
}

// processPageContents will go through the <main> tag
// of the page and extract the relevant information.
func (p *Puzzle) processPageContents(mainContents *goquery.Selection) {
	// Clearing out existing parsed info to ensure data is up to date
	p.AnswerOne = ""
	p.AnswerTwo = ""
	p.ArticleOne = make([]string, 0)
	p.ArticleTwo = make([]string, 0)
	p.Title = ""

	p.Title = mainContents.Find("h2").First().Text()

	mainContents.Find("article").Each(func(i int, s *goquery.Selection) {
		if len(p.ArticleOne) == 0 {
			p.ArticleOne = getPrettyArticle(s)

		} else {
			p.ArticleTwo = getPrettyArticle(s)
		}
	})

	// This should only grab "Your puzzle answer was: " tags
	mainContents.Find("article + p").Each(func(i int, s *goquery.Selection) {
		outStr := s.Text()
		if p.AnswerOne == "" {
			p.AnswerOne = outStr
		} else {
			p.AnswerTwo = outStr
		}
	})
}

func getPrettyArticle(article *goquery.Selection) []string {
	var articleOut []string

	article.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "h2" {
			return
		}

		loopContents := ""
		sel.Contents().Each(func(j int, s *goquery.Selection) {
			// TODO: Try to fix links. Maybe try "termlink" module

			// if goquery.NodeName(s) == "a" {
			// 	href, exists := s.Attr("href")
			// 	if exists {
			// 		// Links get made blue with an underline
			// 		articleOut += createLink(href, linkStyle.Render(s.Text()))
			// 		// articleOut += linkStyle.Render(s.Text() + "(" + href + ")")
			// 	}
			// } else

			if goquery.NodeName(s) == "em" {
				parent := s.Parent()
				if goquery.NodeName(parent) == "code" {
					// Emphatic code should get rendered as code and emphasis
					loopContents += italStyle.Render(codeStyle.Render(s.Text()))
				} else if s.HasClass("star") {
					loopContents += starStyle.Render(s.Text())
				} else {
					loopContents += italStyle.Render(s.Text())
				}
			} else if goquery.NodeName(s) == "code" {
				loopContents += codeStyle.Render(s.Text())
			} else if goquery.NodeName(s) != "h2" {
				loopContents += s.Text()
			}
		})

		articleOut = append(articleOut, wrapText(loopContents, tui.ViewportWidth)+"\n")
	})

	return articleOut
}

func wrapText(line string, width int) string {
	var result string
	words := strings.Fields(line)
	lineLength := 0

	for _, word := range words {
		if lineLength+len(word)+1 > width {
			result += "\n"
			lineLength = 0
		}
		if lineLength > 0 {
			result += " "
			lineLength++
		}

		result += word
		lineLength += len(word)
	}

	return result
}

// type PuzzlePart struct {
// 	starObtained  bool
// 	example       string
// 	isPartB       bool
// 	submissions   []*Submission
// 	CorrectAnswer Value
// }
//
// type Submission struct {
// 	submissionVal Value
// 	correct       bool
// 	timeSubmitted time.Time
// 	feedback      string
// }

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

// // TODO: answerStyle
var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
	italStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3374"))
	starStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
	codeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
	wordWrap   = lipgloss.NewStyle().Width(tui.ViewportWidth)
)
