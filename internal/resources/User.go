package resources

import (
	"errors"
	"strings"
	"time"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/session"
	"dalton.dog/aocgo/internal/styles"
	"dalton.dog/aocgo/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// TODO: Load user's name from a page
// Attempt and do it bespokely on user creation
// so we don't need to rely on other requests for it

// User represents a session token and accompanying puzzles.
type User struct {
	NumStars   int
	Years      map[int][]*Puzzle
	SessionTok string
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
	logger := styles.GetStdoutLogger()

	maxYear, _ := utils.GetCurrentMaxYearAndDay()

	numStars := make(map[int]int)
	year := utils.FIRST_YEAR

	for year <= maxYear {
		logger.Info("Loading year", "year", year)
		numStars[year] = 0
		day := 1
		for day <= 25 {
			logger.Info("Loading day", "day", day)
			puzzle := LoadOrCreatePuzzle(year, day, u.SessionTok)

			u.Years[year][day] = puzzle

			if puzzle.AnswerOne != "" {
				logger.Info("Answer one found!", "year", year, "day", day, "answer", puzzle.AnswerOne)
				u.NumStars++
				numStars[year]++
				if puzzle.AnswerTwo != "" {
					logger.Info("Answer two found!", "year", year, "day", day, "answer", puzzle.AnswerTwo)
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
		logger.Info("Ending year", "Stars found", numStars[year])

		year++
	}
}
