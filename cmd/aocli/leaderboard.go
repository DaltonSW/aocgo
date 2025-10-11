package main

import (
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/utils"
)

var leaderboardCmd = &cobra.Command{
	Use:     "leaderboard",
	Aliases: []string{"lb", "board"},
	Short:   "Shows a puzzle's daily leaderboard, or a yearly leaderboard.",
	Args:    cobra.NoArgs,
	Annotations: map[string]string{
		requiresAuthAnnotation: "true",
	},
	Run: func(cmd *cobra.Command, args []string) {
		Leaderboard(Year, Day)
	},
}

// Leaderboard obtains and displays Leaderboard information for a specific year or day
// Command: `aocli leaderboard -y yyyy [-d dd]`
// Params:
//
//	(Req) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func Leaderboard(yearIn, dayIn string) {
	var (
		year int
		day  int
		err  error
	)

	if yearIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			output.Fatal("Error loading leaderboard based on current directory!", "err", err)
		}
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			output.Fatal("Error parsing year!", "err", err)
		}

		if dayIn != "0" {
			day, err = utils.ParseDay(dayIn)
			if err != nil {
				output.Fatal("Error parsing day from args.", "err", err)
			}
		}
	}

	lb, err := resources.LoadOrCreateLeaderboard(year, day)
	if err != nil {
		output.Fatal("Unable to load leaderboard data.", "err", err)
	}

	tOne, tTwo := lb.GetContent()

	resources.NewLeaderboardViewport(tOne, tTwo, lb.GetTitle())
}
