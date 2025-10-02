package main

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/utils"
)

var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "Shows a puzzle's daily leaderboard, or a yearly leaderboard.",
	Args:  cobra.NoArgs,
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
	var year int
	var day int
	var err error
	var lb resources.ViewableLB

	if yearIn == "0" {
		year, day, err := utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal("Error loading leaderboard based on current directory!", "err", err)
		}
		lb = resources.LoadOrCreateLeaderboard(year, day)
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal("Error parsing year!", "err", err)
		}
	}

	if lb == nil {
		if dayIn != "0" {
			day, err = utils.ParseDay(dayIn)
			if err != nil {
				log.Fatal("Error parsing day from args.", "err", err)
			}
			lb = resources.LoadOrCreateLeaderboard(year, day)
		} else {
			lb = resources.LoadOrCreateLeaderboard(year, 0)
		}
	}

	if lb == nil {
		log.Fatal("Unable to load/create leaderboard!")
		return
	}

	resources.NewLeaderboardViewport(lb.GetContent(), lb.GetTitle())
}
