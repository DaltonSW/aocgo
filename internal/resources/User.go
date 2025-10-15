package resources

import (
	"errors"
	"strings"
	"time"

	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/PuerkitoBio/goquery"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// User represents a session token and accompanying puzzles.
type User struct {
	Label       string
	DisplayName string
	NumStars    int
	Years       map[int][]*Puzzle
	SessionTok  string
}

// GetToken returns the user's session token.
func (u *User) GetToken() string {
	return u.SessionTok
}

// GetLabel returns the user's configured label.
func (u *User) GetLabel() string {
	return u.Label
}

// NewUser creates a user model based on the provided label and session token.
func NewUser(label, token string) (*User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	label = strings.TrimSpace(label)
	if label == "" {
		return nil, errors.New("label cannot be empty")
	}

	api.InitClient(token)

	yearMap := make(map[int][]*Puzzle)
	for i := utils.FIRST_YEAR; i <= time.Now().Year(); i++ {
		yearMap[i] = make([]*Puzzle, 26)
	}

	newUser := &User{
		Label:      label,
		SessionTok: token,
		Years:      yearMap,
	}

	name := newUser.LoadDisplayName()
	newUser.DisplayName = name

	return newUser, nil
}

func (u *User) Display() {
	p := tea.NewProgram(u.NewModel())
	if _, err := p.Run(); err != nil {
		output.Fatal("Couldn't run viewport!", "err", err)
	}
}

func (u *User) LoadUser() {
	maxYear, _ := utils.GetCurrentMaxYearAndDay()

	numStars := make(map[int]int)
	year := utils.FIRST_YEAR

	for year <= maxYear {
		numStars[year] = 0
		day := 1
		for day <= 25 {
			puzzle := LoadOrCreatePuzzle(year, day, u.SessionTok)
			u.Years[year][day] = puzzle

			if puzzle.AnswerOne != "" {
				u.NumStars++
				numStars[year]++
				if puzzle.AnswerTwo != "" {
					u.NumStars++
					numStars[year]++
				}
			}
			day++
		}

		// There's only 1 puzzle on Day 25, so if they've earned 49 stars, they get the 50th for free
		if numStars[year] == 49 {
			u.Years[year][25].AnswerTwo = "Merry Christmas!"
			u.NumStars++
			numStars[year]++
		}
		year++
	}
}

func (u *User) LoadDisplayName() string {
	resp, err := api.NewGetReq("https://adventofcode.com/", u.SessionTok)
	if err != nil {
		output.Fatal("Unable to load user's information", "err", err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		output.Fatal("Error constructing new PageData.", "error", err)
	}

	nameDiv := doc.Find("div.user")
	nameClone := nameDiv.Clone()
	nameClone.Find("span").Remove()

	return strings.TrimSpace(nameClone.Text())
}
