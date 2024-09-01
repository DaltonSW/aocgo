package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// TODO: Move styles to 'styles' package

// TODO: Make sure documentation here is accurate as of v1.0

var NameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
var UseStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
var DescStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00"))
var ArgsStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF00FF"))

var (
	aocliHelpText = helpText{
		name: "aocli",
		use:  "aocli <command> [params...]",
		desc: `
	aocli is a CLI tool written in Go for interacting with Advent of Code information.`,
		args: `
	get --------- Get the user input for a given year and day and save it to a local file
	health ------ Checks to see if the system has valid configuration in place to successfully run the program
	help -------- Shows the help information for a specific command
	leaderboard - Shows the leaderboard for the given year, or given year and day
	reload ------ Refresh the page data for the puzzle on a given year and day
	submit ------ Submit a puzzle answer for a given year and day
	user -------- View the stars obtained for the current user
	view -------- Pretty print the puzzle for a given day`,
	}

	getHelpText = helpText{
		name: "get - Saves the user's puzzle input to a file",
		use:  "aocli get [year] [day]",
		desc: `
	Loads the user's puzzle input for a given year and day.
	The input will be saved to the current directory, named 'input.txt'.`,
		args: `
	year - Year to load input for. If not provided, will try to derive from parent dir name.
	day  - Day to load input for. If not provided, will try to derive from current dir names.`,
	}

	healthHelpText = helpText{
		name: "health - Checks validity of aocutil module",
		use:  "aocli health",
		desc: `
	Will check the current session token for validity.`,
		args: "",
	}

	helpHelpText = helpText{
		name: "help",
		use:  "aocli help <command>",
		desc: "Shows the help information for the program as a whole, or for a specific command",
	}

	leaderboardHelpText = helpText{
		name: "leaderboard",
		use:  "aocli leaderboard [year] [day]",
		desc: "Shows the leaderboard for a given year, or a given year and day",
	}

	reloadHelpText = helpText{
		name: "reload",
		use:  "aocli reload [year] [day]",
		desc: "Reloads the page data for the puzzle on a given year and day",
	}

	submitHelpText = helpText{
		name: "submit",
		desc: `Submits an answer to the puzzle on the user's behalf and prints out the response.
Attempts to derive the year and day from current and parent directory name.`,
	}

	userHelpText = helpText{
		name: "user",
		use:  "aocli user",
		desc: "Loads all puzzle information for a user and will display their star count for each year and day in a table.",
	}

	viewHelpText = helpText{
		name: "view - Pretty prints the puzzle's page data to the screen",
		use:  "aocli view [year] [day]",
		desc: `
	Pretty-prints the puzzle's page data.
	Will attempt to parse the current and parent directory names for date information.
	If a user session token is set, the page will properly include Part B when unlocked.`,
		args: `
	year - Year to load input for. If not provided, will try to derive from parent dir name.
	day  - Day to load input for. If not provided, will try to derive from current dir names.`,
	}
)

var HelpTextMap map[string]helpText = map[string]helpText{
	// Parent program help text
	"aocli": aocliHelpText,

	// Commands help text
	"get":         getHelpText,
	"health":      healthHelpText,
	"help":        helpHelpText,
	"leaderboard": leaderboardHelpText,
	"reload":      reloadHelpText,
	"submit":      submitHelpText,
	"user":        userHelpText,
	"view":        viewHelpText,
}

type helpText struct {
	name string
	use  string
	desc string
	args string
}

func (ht *helpText) Print() {
	outS := NameStyle.Render("NAME:  ")
	outS += ht.name + "\n\n"

	outS += UseStyle.Render("USAGE: ")
	outS += ht.use + "\n\n"

	outS += DescStyle.Render("DESCRIPTION")
	outS += ht.desc + "\n\n"

	if ht.args != "" {
		outS += ArgsStyle.Render("ARGUMENTS")
		outS += ht.args + "\n"
	}

	fmt.Println(outS)
}
