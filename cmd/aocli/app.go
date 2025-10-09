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
var cacheInitialized bool

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
		if !commandRequiresAuth(cmd) {
			log.Debug("Skipping session bootstrap for command", "command", cmd.CommandPath())
			return
		}

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
		cacheInitialized = true

	},

	Run: func(cmd *cobra.Command, args []string) {
		RunLandingPage()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cacheInitialized {
			cache.ShutdownDBM()
			cacheInitialized = false
		}
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

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(leaderboardCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(viewCmd)
}

func commandRequiresAuth(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		if current.Annotations != nil {
			if value, ok := current.Annotations[requiresAuthAnnotation]; ok {
				return value == "true"
			}
		}
	}
	return false
}
