package main

import (
	"fmt"
	"os"
	"strconv"

	"dalton.dog/aocutil/internal/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// TODO: `help` - Print help and info about the CLI tool

// TODO: `get [day] [year]`  - if no inputs, checks today's date

// TODO: `submit [day] [year]` - if no inputs, tries today's date
//	Not sure if we're going to need to specifically accept submitting for Part 1 or 2

// TODO: `run` - Will benchmark and run files in current and subdirectory

var helpStyle = lipgloss.NewStyle()
var testStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

// NOTE: Bold just changes to a different font for me. Not sure why. But it DOES do something
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
	case "test":
		test(args)
	default:
		fmt.Println("Not a valid command! Run `aocli help` to see valid commands.")
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
	// TODO: Validate input
	if len(args) < 4 {
		return
		// TODO: Try loading with today
		// TODO: Print `get` help message
	}
	user, err := models.NewUser("")
	if err != nil {
		log.Error("Unable to load/create user!", "err", err)
	}

	year, _ := strconv.Atoi(args[2])
	day, _ := strconv.Atoi(args[3])

	puzzle := models.NewPuzzle(year, day)
	userInput := puzzle.GetUserPuzzleInput(user.GetToken())
	fmt.Print(userInput)
}

// `submit [year] [day] [part] [answer]` command
func submit(args []string) {

}

func run(args []string) {

}

// `view [year] [day]` command
// Desc: Pretty-prints the puzzle's page data
func view(args []string) {
	// TODO: Validate input
	if len(args) < 4 {
		return
		// TODO: Try loading with today
		// TODO: Print `get` help message
	}

	year, _ := strconv.Atoi(args[2])
	day, _ := strconv.Atoi(args[3])

	puzzle := models.NewPuzzle(year, day)
	rawPage := puzzle.GetPuzzlePageData()
	pageData := models.NewPageData(rawPage)
	pageData.PrintPageData()
}

// `health` command
// Desc: Checks if a session key is available
func health(args []string) {

}

// Command:	`test`
// Desc:	Does whatever I need to test at the time :)
func test(args []string) {

}
