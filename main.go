package aocgo

import (
	"strings"

	"dalton.dog/aocgo/internal/dirparse"
	"dalton.dog/aocgo/internal/models"
	"dalton.dog/aocgo/internal/session"
	"github.com/charmbracelet/log"
)

func GetInputAsByteArray() []byte {
	year, day, err := dirparse.GetYearAndDayFromCWD()
	if err != nil {
		log.Fatal(err)
	}

	return getData(year, day)
}

func GetInputAsString() string {
	return string(GetInputAsByteArray())
}

func GetInputAsLineArray() []string {
	return strings.Split(GetInputAsString(), "\n")
}

func getData(year int, day int) []byte {
	userToken, err := session.GetSessionToken()
	if err != nil {
		log.Fatal(err)
	}
	puzzle := models.NewPuzzle(year, day)
	input, err := puzzle.GetUserPuzzleInput(userToken)
	if err != nil {
		log.Fatal(err)
	}
	return input
}
