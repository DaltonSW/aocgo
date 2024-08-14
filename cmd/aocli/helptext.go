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

var HelpTextMap map[string]*helpText = map[string]*helpText{
	// NOTE: `aocli` program
	"aocli": {
		name: "aocli - CLI tool for interacting with Advent of Code",
		use:  "aocli <command> [params...]",
		desc: `
	aocli is a CLI tool written in Go for interacting with Advent of Code information.
	It is able to view puzzle information and access user puzzle input, with more on the way!
	Requests are limited to 10/s, and responses are locally cached.`,
		args: `
	get    - Gets user input for a certain day and saves it to 'input.txt'
	view   - Pretty print puzzle page data to screen
	health - Verifies that aocli and aocgo have information required to run properly`,
		// submit <answer> <part> [year] [day] - Attempts to submit an answer for a puzzle
	},

	// NOTE: `view` command
	"view": {
		name: "view - Pretty prints the puzzle's page data to the screen",
		use:  "aocli view [year] [day]",
		desc: `
	Pretty-prints the puzzle's page data.
	Will attempt to parse the current and parent directory names for date information.
	If a user session token is set, the page will properly include Part B when unlocked.`,
		args: `
	year - Year to load input for. If not provided, will try to derive from parent dir name.
	day  - Day to load input for. If not provided, will try to derive from current dir names.`,
	},

	// NOTE: `get` command
	"get": {
		name: "get - Saves the user's puzzle input to a file",
		use:  "aocli get [year] [day]",
		desc: `
	Loads the user's puzzle input for a given year and day.
	The input will be saved to the current directory, named 'input.txt'.`,
		args: `
	year - Year to load input for. If not provided, will try to derive from parent dir name.
	day  - Day to load input for. If not provided, will try to derive from current dir names.`,
	},

	// NOTE: `health` command
	"health": {
		name: "health - Checks validity of aocutil module",
		use:  "aocli health",
		desc: `
	Will check the current session token for validity.`,
		args: "",
	},
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

const submitHelpText = `
Submits the provided answer. Will display the response from the server.
`

const getHelpText = `
Gets the input for the given year and day.
Will attempt to parse the current and parent directory names for date information.
By default, will save the input to a file in the current directory called 'input.txt'.
`

const healthHelpText = `
Checks if there's a session token that's able to be used.
`

const runHelpText = `
Runs the tests in the current directory, and any subdirectories.
`
