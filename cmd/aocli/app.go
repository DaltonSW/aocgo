package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"

	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/session"
	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"
)

const Version = "2.0.0b"

var Year string
var Day string

var AnswerPart int
var OutFilename string
var BaseFilename string

var UserRsrc *resources.User
var cacheInitialized bool

func main() {
	fang.Execute(context.Background(), rootCmd, fang.WithoutCompletions(), fang.WithVersion(Version))
}

var rootCmd = &cobra.Command{
	Use:   "aocli",
	Short: "A CLI tool for interacting with Advent of Code puzzles.",
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !commandRequiresAuth(cmd) {
			output.Debug("Skipping session bootstrap for command", "command", cmd.CommandPath())
			return
		}

		activeUser, err := session.GetActiveUser(false)
		if err != nil {
			output.Fatalf("Unable to determine active user. %v", err)
		}

		UserRsrc, err = resources.NewUser(activeUser.Label, activeUser.Token)

		if err != nil {
			output.Fatalf("Unable to create user to run requests as. %v", err)
		} else {
			output.Debug("User loaded", "label", UserRsrc.GetLabel())
		}

		err = cache.StartupDBM(UserRsrc.GetLabel())
		if err != nil {
			output.Fatal(err)
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
	rootCmd.AddGroup(&cobra.Group{ID: "puzzles", Title: "Puzzles"})

	getCmd.Flags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	getCmd.Flags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")
	getCmd.Flags().StringVarP(&OutFilename, "out", "o", "input.txt", "--out filename")
	rootCmd.AddCommand(getCmd)

	leaderboardCmd.Flags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	leaderboardCmd.Flags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")
	rootCmd.AddCommand(leaderboardCmd)

	newCmd.Flags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	newCmd.Flags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")
	newCmd.Flags().StringVarP(&BaseFilename, "base", "b", "base.go", "--base filename")
	newCmd.Flags().StringVarP(&OutFilename, "out", "o", "main.go", "--out filename")
	rootCmd.AddCommand(newCmd)

	reloadCmd.Flags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	reloadCmd.Flags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")
	rootCmd.AddCommand(reloadCmd)

	submitCmd.Flags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	submitCmd.Flags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")
	submitCmd.Flags().IntVarP(&AnswerPart, "part", "p", 0, "--part [1|2]")
	rootCmd.AddCommand(submitCmd)

	viewCmd.Flags().StringVarP(&Year, "year", "y", "0", "--year [2015...2024]")
	viewCmd.Flags().StringVarP(&Day, "day", "d", "0", "--day [1...25]")
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

const tree = `
*
/_\
/_o_\
/o___o\
/__o____\
/_o___o___\
/__o__o__o__\
/o_________o__\
/____o___o______\
/o______o____o____\
/__o___o_____o______\
/_o___o_______o___o___\
|||
|||
`

var Colors = []lipgloss.Style{
	styles.BlueTextStyle,
	styles.RedTextStyle,
	styles.PurpleTextStyle,
	styles.CyanTextStyle,
}

func RunLandingPage() {
	var sOut string
	for i := range tree {
		c := string(tree[i])
		switch c {
		case "*":
			sOut += styles.YellowTextStyle.Render(c)
		case "o":
			sOut += Colors[rand.Intn(len(Colors))].Render(c)
		case "|":
			sOut += styles.BrownTextStyle.Render(c)
		default:
			sOut += styles.GreenTextStyle.Render(c)
		}
	}

	outStr := lipgloss.JoinVertical(lipgloss.Center,
		sOut,
		"Welcome to aocli!",
		fmt.Sprintf("Run %s to see a list of available commands.", styles.CodeStyle.Render("aocli help")),
		styles.SubtitleStyle.Render("\naocli by Dalton Williams (https://dalton.dog)"),
		styles.SubtitleStyle.Render("Advent of Code by Eric Wastl (http://was.tl)"),
		styles.SubtitleStyle.Render(fmt.Sprintf("ver. %s", Version)),
	)

	utils.ClearTerminal()
	lipgloss.Println(lipgloss.NewStyle().PaddingTop(1).PaddingLeft(2).Render(outStr))
}
