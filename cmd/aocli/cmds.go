package main

// TODO: Update godocs

import (
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
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(viewCmd)
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
			log.Fatal("Unable to create user to run requests as. Run `aocli health`.")
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

		if cmd.Name() != "update" {
			CheckForUpdate()
		}
	},
}

var newCmd = &cobra.Command{
	Use:   "new [-b base.go]",
	Short: "Copies the given file into the ./<year>/<day>/main.go",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		New(Year, Day, BaseFilename, OutFilename)
	},
}

var getCmd = &cobra.Command{
	Use:   "get [-o filename]",
	Short: "Gets the puzzle input and saves it to disk.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Get(UserRsrc, Year, Day, OutFilename)
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Checks if aocli and aocgo have proper config to run.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Health()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints out the current version of the program.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Version()
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Checks for an update. Will download if there's a new version.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		RunUpdateModel()
	},
}

var submitCmd = &cobra.Command{
	Use:   "submit [-p {1|2}] <answer>",
	Short: "Submits the given answer to a puzzle.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Submit(UserRsrc, Year, Day, args[0], AnswerPart)
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reloads the page data for a given puzzle.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Reload(UserRsrc, Year, Day)
	},
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Views the puzzle's page inside of the terminal.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		View(UserRsrc, Year, Day)
	},
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Shows a visual representation of the user's puzzle progress.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		User(UserRsrc, ClearUser)
	},
}

var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "Shows a puzzle's daily leaderboard, or a yearly leaderboard.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Leaderboard(Year, Day)
	},
}
