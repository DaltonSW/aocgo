package resources

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/utils"

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

	UserInput   []byte
	Submissions map[int][]*Submission

	LockoutEnd time.Time
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

func (p *Puzzle) Display() {
	NewPuzzleViewport(p)
}

// Answer Response types
const (
	IncorrectAnswer int = iota
	CorrectAnswer
	WarningAnswer
)

func (p *Puzzle) SubmitAnswer(answer string, part int) (int, string) {
	if !time.Now().After(p.LockoutEnd) {
		return WarningAnswer, fmt.Sprintf("Still within lockout period of last submission. Lockout End: %s", p.LockoutEnd.Format(time.Stamp))
	}

	if p.AnswerOne != "" && answer == p.AnswerOne {
		return CorrectAnswer, "You've already gotten it correct, but that IS the correct answer for part 1."
	} else if p.AnswerTwo != "" && answer == p.AnswerTwo {
		return CorrectAnswer, "You've already gotten it correct, but that IS the correct answer for part 2."
	}

	if part == 0 {
		if p.AnswerOne == "" {
			part = 1
		} else if p.AnswerTwo == "" {
			part = 2
		} else {
			return WarningAnswer, "You've already gotten both stars for this level."
		}
	}

	// TODO: Check past submissions before allowing submission
	//	- Past submissions
	//		- Too high / too low

	for _, pastSub := range p.Submissions[part] {
		if pastSub.answer == answer {
			return WarningAnswer, "You've already submitted that answer!"
		}
	}

	submissionData, err := api.SubmitAnswer(p.Year, p.Day, part, p.SessionToken, answer)
	if err != nil {
		log.Fatal(err)
	}

	submission, err := NewSubmission(submissionData, answer)
	if err != nil {
		log.Fatal(err)
	}

	if p.Submissions == nil {
		p.Submissions = make(map[int][]*Submission)
	}

	outList, ok := p.Submissions[part]
	if !ok {
		outList = make([]*Submission, 0)
	}

	outList = append(outList, submission)

	p.Submissions[part] = outList
	defer p.SaveResource()
	if submission.correct {
		p.ReloadPuzzleData()

		if p.AnswerOne == "" {
			p.AnswerOne = answer
			return CorrectAnswer, "First star obtained! Run `view` again to get part 2."
		} else {
			p.AnswerTwo = answer
			return CorrectAnswer, "Second star obtained! That's all for today, good luck tomorrow!"
		}

	} else {
		// TODO: Parse the response message for lockout period
		lockoutDuration, err := utils.ParseDuration(submission.message)
		if err != nil {
			return IncorrectAnswer, submission.message + "\nUnable to parse lockout duration from message."
		}

		p.LockoutEnd = time.Now().Add(lockoutDuration)

		return IncorrectAnswer, submission.message
	}
}

func newPuzzle(year int, day int, userSession string) *Puzzle {
	URL := fmt.Sprintf(PUZZLE_URL, year, day)
	bucketID := strconv.Itoa(year) + strconv.Itoa(day)

	userInput, err := loadUserInputFromSite(URL, userSession)

	if err != nil {
		log.Fatal("Unable to load user input for the puzzle.", "error", err)
	}

	subMap := make(map[int][]*Submission)
	subMap[1] = make([]*Submission, 0)
	subMap[2] = make([]*Submission, 0)

	newPuzzle := &Puzzle{
		Day:          day,
		Year:         year,
		BucketID:     bucketID,
		URL:          URL,
		UserInput:    userInput,
		SessionToken: userSession,
		Submissions:  subMap,
	}

	newPuzzle.loadPageData()
	newPuzzle.SaveResource()

	return newPuzzle
}

func (p *Puzzle) ReloadPuzzleData() error {
	newInput, err := loadUserInputFromSite(p.URL, p.SessionToken)
	if err != nil {
		return err
	}

	p.UserInput = newInput
	p.loadPageData()
	p.SaveResource()
	return nil
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
		sOut = append(sOut, "Answer: "+p.AnswerOne)
	}

	if len(p.ArticleTwo) != 0 {
		sOut = append(sOut, "\n\n")
		sOut = append(sOut, "\n"+lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Render(" --- Part Two --- "))
		sOut = append(sOut, "\n")
		sOut = append(sOut, p.ArticleTwo...)
		sOut = append(sOut, "\n")

		if p.AnswerTwo != "" {
			sOut = append(sOut, "Answer: "+p.AnswerTwo)
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
		if strings.Contains(s.Text(), "answer") {
			outStr := s.Find("code").Text()
			if outStr != "" {
				if p.AnswerOne == "" {
					log.Debug("Answer found!", "year", p.Year, "day", p.Day, "answer", outStr)
					p.AnswerOne = outStr
				} else {
					log.Debug("Answer found!", "year", p.Year, "day", p.Day, "answer", outStr)
					p.AnswerTwo = outStr
				}
			}
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

		articleOut = append(articleOut, wrapText(loopContents, ViewportWidth)+"\n")
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
	puzzleTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
	italStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3374"))
	starStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	linkStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
	codeStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
	wordWrap         = lipgloss.NewStyle().Width(ViewportWidth)
)
