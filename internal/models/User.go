package models

import (
	"errors"
	"strings"
	"time"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/session"
	// "github.com/charmbracelet/log"
)

// User represents a session token and accompanying puzzles.
type User struct {
	NumStars   int
	Years      map[int][]Puzzle
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

	yearMap := make(map[int][]Puzzle)
	for i := 2015; i <= time.Now().Year(); i++ {
		yearMap[i] = make([]Puzzle, 25)
	}

	return &User{
		SessionTok: token,
		Years:      yearMap,
	}, nil
}
