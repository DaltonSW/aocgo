package aocli

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

// `help` command - Prints info and a list of commands
func help(args []string) {

}

func get(args []string) {

}

func submit(args []string) {

}

func run(args []string) {

}

func view(args []string) {

}

func health(args []string) {

}
