package main

import (
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/utils"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reloads the page data for a given puzzle.",
	Args:  cobra.NoArgs,
	Annotations: map[string]string{
		requiresAuthAnnotation: "true",
	},
	Run: func(cmd *cobra.Command, args []string) {
		Reload(UserRsrc, Year, Day)
	},
}

// Reload will force reload the puzzle data for a specific day
// Command: `reload [-y yyyy -d dd]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func Reload(user *resources.User, yearIn, dayIn string) {
	var year int
	var day int
	var err error

	if yearIn == "0" || dayIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			output.Fatal(err)
		}

	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			output.Fatal(err)
		}

		day, err = utils.ParseDay(dayIn)
		if err != nil {
			output.Fatal(err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	puzzle.ReloadPuzzleData()
}
