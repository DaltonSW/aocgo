package resources

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mattn/go-runewidth"
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

// Displays the puzzle's page to the user
func (p *Puzzle) Display() {
	NewPuzzleViewport(p)
}

// Answer Response types
const (
	IncorrectAnswer int = iota
	CorrectAnswer
	WarningAnswer // Not submitted, but weird behavior
	NeutralAnswer // Not submitted, but no warning
)

// SubmitAnswer takes an answer and a part to submit to.
// If no part is provided, it will be derived based on stored puzzle information.
func (p *Puzzle) SubmitAnswer(answer string, part int) (int, string) {
	if !time.Now().After(p.LockoutEnd) {
		return WarningAnswer, fmt.Sprintf("Still within lockout period of last submission. Lockout End: %s", p.LockoutEnd.Format(time.Stamp))
	}

	if p.AnswerOne != "" && answer == p.AnswerOne {
		return NeutralAnswer, "Correct answer for Part 1 (no answer submitted, already got star)."
	} else if p.AnswerTwo != "" && answer == p.AnswerTwo {
		return NeutralAnswer, "Correct answer for Part 2 (no answer submitted, already got star)."
	}

	if part == 0 {
		if p.AnswerOne == "" {
			part = 1
		} else if p.AnswerTwo == "" {
			part = 2
		} else {
			return NeutralAnswer, "You've already gotten both stars for this level."
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
		defer p.ReloadPuzzleData()

		if p.AnswerOne == "" {
			p.AnswerOne = answer
			if p.Day == 25 {
				p.AnswerTwo = "Merry Christmas!"
				return CorrectAnswer, "If you've got all 49 other stars for this year, submit again to get the 50th and complete the year!"
			} else {
				return CorrectAnswer, "First star obtained! Run `view` again to get part 2."
			}
		} else {
			p.AnswerTwo = answer
			return CorrectAnswer, "Second star obtained! That's all for today, good luck tomorrow!"
		}

	} else {
		lockoutDuration, err := utils.ParseDuration(submission.message)
		if err != nil {
			return IncorrectAnswer, submission.message + "\nUnable to parse lockout duration from message."
		}

		p.LockoutEnd = time.Now().Add(lockoutDuration)

		return IncorrectAnswer, submission.message
	}
}

// Creates a new puzzle by loading information from the server. Bypasses any cached data
func newPuzzle(year int, day int, userSession string) *Puzzle {
	URL := fmt.Sprintf(PUZZLE_URL, year, day)
	bucketID := strconv.Itoa(year) + strconv.Itoa(day)

	userInput, err := loadUserInputFromSite(URL, userSession)

	if err != nil {
		log.Fatal("Unable to load user input for the puzzle.", "error", err)
	} else if strings.Contains(string(userInput), "log in") {
		log.Fatal("Session token appears to be invalid. Login in a browser and get your new token.")
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

// Reloads puzzle information from the server
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

// GetPrettyPageData parses the puzzle's stored information and displays it in a visually pleasing way.
func (p *Puzzle) GetPrettyPageData() []string {
	var sOut []string
	sOut = append(sOut, lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Render(" - Part One - "))
	sOut = append(sOut, "\n")

	sOut = append(sOut, p.ArticleOne...)

	if p.AnswerOne != "" {
		sOut = append(sOut, "Answer: "+p.AnswerOne)
	}

	if len(p.ArticleTwo) != 0 {
		sOut = append(sOut, "\n\n"+lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Render(" - Part Two - "))
		sOut = append(sOut, "\n")
		sOut = append(sOut, p.ArticleTwo...)
		sOut = append(sOut, "\n")

		if p.AnswerTwo != "" {
			sOut = append(sOut, "Answer: "+p.AnswerTwo)
		}
	}
	return sOut
}

// Contacts the server to load the user's input for a given puzzle
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
					p.AnswerOne = styles.CodeStyle.Render(outStr)
				} else {
					log.Debug("Answer found!", "year", p.Year, "day", p.Day, "answer", outStr)
					p.AnswerTwo = styles.CodeStyle.Render(outStr) + "\n"
				}
			}
		}
	})
}

func getPrettyArticle(article *goquery.Selection) []string {
	var articleOut []string

	article.Contents().Each(func(i int, sel *goquery.Selection) {
		switch goquery.NodeName(sel) {
		case "h2":
			return
		case "p":
			paraContents := getPrettySelection(sel)

			articleOut = append(articleOut, wrapText(paraContents, ViewportWidth)+"\n\n")
		case "ul":
			sel.Find("li").Each(func(j int, s *goquery.Selection) {
				articleOut = append(articleOut, " - "+wrapText(getPrettySelection(s), ViewportWidth-2)+"\n\n")
			})
			articleOut = append(articleOut)
		case "pre":
			// Extract the <code> content
			preContent := sel.Find("code").Text()

			// Split content into lines
			lines := strings.Split(preContent, "\n")

			for _, line := range lines {
				if line != "" { // Ignore empty lines
					articleOut = append(articleOut, styles.CodeStyle.Render(line)+"\n")
				}
			}
			articleOut = append(articleOut, "\n") // Add spacing after the block
		}
	})

	return articleOut
}

func getPrettySelection(sel *goquery.Selection) string {
	selContents := ""
	sel.Contents().Each(func(j int, s *goquery.Selection) {
		switch goquery.NodeName(s) {
		case "em":
			if s.HasClass("star") {
				selContents += styles.StarStyle.Render(s.Text())
			} else {
				selContents += styles.ItalStyle.Render(s.Text())
			}
		case "code":
			selContents += styles.CodeStyle.Render(s.Text())
		case "h2":
			return
		default:
			selContents += s.Text()
		}
	})

	return selContents
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// Wraps text without accounting for encoding characters
func wrapText(line string, width int) string {
	var result string
	words := strings.Fields(line)
	lineLength := 0

	for _, word := range words {
		// Strip ANSI escape sequences for length calculation
		cleanWord := ansiRegex.ReplaceAllString(word, "")
		visibleLength := runewidth.StringWidth(cleanWord)

		if lineLength+visibleLength+1 > width {
			result += "\n"
			lineLength = 0
		}
		if lineLength > 0 {
			result += " "
			lineLength++
		}

		result += word
		lineLength += visibleLength
	}

	return result
}

// Possible submission value for a puzzle
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
)
