package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// TODO: `help` - Print help and info about the CLI tool

// TODO: `get [day] [year]`  - if no inputs, checks today's date

// TODO: `submit [day] [year]` - if no inputs, tries today's date
//	Not sure if we're going to need to specifically accept submitting for Part 1 or 2

// TODO: `run` - Will benchmark and run files in current and subdirectory

// TODO: `view [day] [year]` - surface the description and examples in a pretty lipgloss text view or whatever

var helpStyle = lipgloss.NewStyle()

func main() {
	args := os.Args
	if len(args) == 1 {
		nothing()
		os.Exit(0)
	}

	switch args[1] {
	case "help":
		help(args)
	case "get":
		get(args)
	case "submit":
		submit(args)
	case "run":
		run(args)
	case "view":
		view(args)
	case "health":
		health(args)
	}
	return
}

// Runs when no command was given. Suggests to run "aocli help"
func nothing() {
	s := "No command was given! Try running `aocli help`."
	fmt.Println(helpStyle.Render(s))
}

// `help [command]` command
// Desc: Prints info and a list of commands
// Params:
//
//	[command] - command name to print help for
func help(args []string) {

}

// `get [year] [day]` command
// Desc: Gets input data for a specific day
// Params:
//
//	[year] - 2 or 4 digit year (16 or 2016)
//	[day]  - 1 or 2 digit day (1, 01, 21)
//
// TODO: Accept an -o flag to write the input to a file vs forcing a redirection
func get(args []string) {

}

// `submit [year] [day] [part] [answer]` command
func submit(args []string) {

}

func run(args []string) {

}

// `view [year] [day]` command
// Desc: Pretty-prints the puzzle's page data
func view(args []string) {

}

// `health` command
// Desc: Checks if a session key is available
func health(args []string) {

}
