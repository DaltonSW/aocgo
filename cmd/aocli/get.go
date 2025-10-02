package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/utils"
)

var getCmd = &cobra.Command{
	Use:   "get [-o filename]",
	Short: "Gets the puzzle input and saves it to disk.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Get(UserRsrc, Year, Day, OutFilename)
	},
}

// Get obtains input data for a specific day, outputting it to the current directory `input.txt`.
// Command: `aocli get [-y yyyy -d dd -o output_name.txt]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
//	(Opt) filename  - overrides "input.txt" name if it's provided
func Get(user *resources.User, yearIn, dayIn string, filename string) {
	var year int
	var day int
	var err error

	if yearIn == "0" || dayIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal("Unable to parse year/day from current directory.", "err", err)
		}
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal("Unable to parse year from current directory.", "err", err)
		}

		day, err = utils.ParseDay(dayIn)
		if err != nil {
			log.Fatal("Unable to parse day from current directory.", "err", err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	userInput, _ := puzzle.GetUserInput()

	out, _ := os.Create(filename)
	defer out.Close()
	out.Write(userInput)

	log.Infof("Input saved to %v!", filename)
}
