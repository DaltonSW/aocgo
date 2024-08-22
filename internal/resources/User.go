package resources

import (
	"errors"
	"strings"
	"time"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/session"

	// "dalton.dog/aocgo/internal/styles"
	"dalton.dog/aocgo/internal/utils"
	"github.com/PuerkitoBio/goquery"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// TODO: Load user's name from a page
// Attempt and do it bespokely on user creation
// so we don't need to rely on other requests for it

// User represents a session token and accompanying puzzles.
type User struct {
	DisplayName string
	NumStars    int
	Years       map[int][]*Puzzle
	SessionTok  string
}

// GetToken returns the user's session token.
func (u *User) GetToken() string {
	return u.SessionTok
}

// Creates a new user based on a provided session token.
// If none is provided, it'll be loaded from environment
// variable or from config file.
func NewUser(token string) (*User, error) {
	var err error
	if token == "" {
		token, err = session.GetSessionToken()
		if err != nil {
			return nil, err
		}
	}

	if token == "" {
		return nil, errors.New("Token was still empty after load attempts.")
	}

	token = strings.TrimSpace(token)
	api.InitClient(token)

	yearMap := make(map[int][]*Puzzle)
	for i := utils.FIRST_YEAR; i <= time.Now().Year(); i++ {
		yearMap[i] = make([]*Puzzle, 26)
	}

	return &User{
		SessionTok: token,
		Years:      yearMap,
	}, nil
}

func (u *User) Display() {
	p := tea.NewProgram(u.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Couldn't run viewport!", "err", err)
	}
}

func (u *User) LoadUser() {
	u.LoadDisplayName()
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

func (u *User) LoadDisplayName() {
	resp, err := api.NewGetReq("https://adventofcode.com/", u.SessionTok)
	if err != nil {
		log.Fatal("Unable to load user's information", "err", err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error constructing new PageData.", "error", err)
	}

	nameDiv := doc.Find("div.user")

	// log.Info(nameDiv.Text())

	nameClone := nameDiv.Clone()
	// log.Info(nameClone.Text())
	nameClone.Find("span").Remove()

	u.DisplayName = strings.TrimSpace(nameClone.Text())
}
