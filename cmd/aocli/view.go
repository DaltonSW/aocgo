package main

import (
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/utils"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Views the puzzle's page inside of the terminal.",
	Args:  cobra.NoArgs,
	Annotations: map[string]string{
		requiresAuthAnnotation: "true",
	},
	Run: func(cmd *cobra.Command, args []string) {
		View(UserRsrc, Year, Day)
	},
}

// View will pretty print the puzzle's page data.
// Command: `aocli view [-y yyyy -d dd]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func View(user *resources.User, yearIn, dayIn string) {
	var year int
	var day int
	var err error

	if yearIn == "0" || dayIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			output.Fatal("Unable to parse year/day from current directory.", "err", err)
		}
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			output.Fatal("Unable to parse year from current directory.", "err", err)
		}

		day, err = utils.ParseDay(dayIn)
		if err != nil {
			output.Fatal("Unable to parse day from current directory.", "err", err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	puzzle.Display()
}
