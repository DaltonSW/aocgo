package models

import (
	"errors"
	"strings"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/session"
	"github.com/charmbracelet/log"
)

type User struct {
	numStars   int
	calendar   Calendar
	sessionTok string
}

func (u *User) GetToken() string {
	return u.sessionTok
}

func NewUser(token string) (*User, error) {
	var err error
	log.Debug("Trying to create user", "tokenParam", token)
	if token == "" {
		token, err = session.GetSessionToken()
		if err != nil {
			return nil, err
		}
	}
	// Try to load User object to store stuff like numStars and calendar info

	if token == "" {
		return nil, errors.New("Token was still empty after load attempts.")
	}
	token = strings.TrimSpace(token)

	api.InitClient(token)

	return &User{sessionTok: token}, nil
}
