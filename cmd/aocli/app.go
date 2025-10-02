package main

import (
	"context"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/resources"
)

var Year string
var Day string

var AnswerPart int
var OutFilename string
var BaseFilename string
var ClearUser bool

var UserRsrc *resources.User

func main() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "aocli [-y year] [-d day] [command]",
	Short: "A CLI tool for interacting with Advent of Code puzzles.",
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// log.SetLevel(log.DebugLevel)
		var err error

		UserRsrc, err = resources.NewUser("")

		if err != nil {
			log.Fatalf("Unable to create user to run requests as. %v", err)
		} else {
			log.Debug("User loaded", "token", UserRsrc.SessionTok)
		}

		err = cache.StartupDBM(UserRsrc.GetToken())
		if err != nil {
			log.Fatal(err)
		}

	},

	Run: func(cmd *cobra.Command, args []string) {
		RunLandingPage()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cache.ShutdownDBM()
	},
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	rootCmd.PersistentFlags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")

	submitCmd.Flags().IntVarP(&AnswerPart, "part", "p", 0, "--part [1|2]")

	getCmd.Flags().StringVarP(&OutFilename, "out", "o", "input.txt", "--out filename")

	newCmd.Flags().StringVarP(&BaseFilename, "base", "b", "base.go", "--base filename")
	newCmd.Flags().StringVarP(&OutFilename, "out", "o", "main.go", "--out filename")

	userCmd.Flags().BoolVar(&ClearUser, "clear", false, "Clears the stored puzzle data for a user.")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(leaderboardCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(viewCmd)
}
