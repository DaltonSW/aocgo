package main

import (
	"fmt"
	"os"
	"strconv"

	// "dalton.dog/aocgo/internal/dirparse"
	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/models"
	"dalton.dog/aocgo/internal/session"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// TODO: `submit [day] [year]` - if no inputs, tries today's date
//	Not sure if we're going to need to specifically accept submitting for Part 1 or 2

// TODO: `run` - Will benchmark and run files in current and subdirectory

// TODO: `clear session [year] [day]` - Clears the stored information for a given session

var helpBodyStyle = lipgloss.NewStyle().Width(70)
var helpTitleStyle = lipgloss.NewStyle().Width(70)
var testStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

var User *models.User

func main() {
	log.SetLevel(log.DebugLevel)

	args := os.Args
	if len(args) == 1 {
		fmt.Println(helpBodyStyle.Render("Welcome to aocli! Try running `aocli help` for a list of available commands."))
		os.Exit(0)
	}

	User, err := models.NewUser("")
	if err != nil {
		log.Error("Unable to create user to run requests as. Try running `aocli health`.")
	}

	cache.Startup(User.GetToken())
	defer cache.Shutdown()

	switch args[1] {
	case "help":
		help(args)
	case "get":
		get(args)
	case "submit":
		submit(args)
	case "lb":
	case "leaderboard":
		leaderboard(args)
	// case "run":
	// 	run(args)
	case "view":
		view(args)
	case "health":
		health()
	// case "test":
	// 	test()
	default:
		fmt.Println("Not a valid command! Run `aocli help` to see valid commands.")
	}
	return
}

// `help [command]` command
// Desc: Prints info and a list of commands
// Params:
//
//	[command] - command name to print help for
func help(args []string) {
	// Clear terminal
	fmt.Print("\033[H\033[2J")

	// Too many args
	if len(args) > 3 {
		fmt.Println(helpBodyStyle.Render("Too many arguments passed!"))
		return
	}

	// They requested help for a specific command
	if len(args) == 3 {
		commandName := args[2]
		helptext, ok := HelpTextMap[commandName]
		if ok {
			helptext.Print()
		} else {
			fmt.Println(helpBodyStyle.Render("Not a valid command!"))
		}
		return
	}

	// Otherwise they just open-endedly requested help
	ht, ok := HelpTextMap["aocli"]
	if ok {
		outS := "\n"
		outS += NameStyle.Render("NAME:  ")
		outS += ht.name + "\n\n"

		outS += UseStyle.Render("USAGE: ")
		outS += ht.use + "\n\n"

		outS += DescStyle.Render("DESCRIPTION")
		outS += ht.desc + "\n\n"

		outS += ArgsStyle.Render("COMMANDS")
		outS += ht.args

		fmt.Println(outS)
	}

}

// `get [year] [day]` command
// Desc: Gets input data for a specific day
// Params:
//
//	[year] - 2 or 4 digit year (16 or 2016)
//	[day]  - 1 or 2 digit day (1, 01, 21)
func get(args []string) {
	// TODO: Validate input
	if len(args) < 4 {
		return
		// TODO: Try loading with today
		// TODO: Print `get` help message
	}
	user, err := models.NewUser("")
	if err != nil {
		log.Fatal("Unable to load/create user!", "err", err)
	}

	year, _ := strconv.Atoi(args[2])
	day, _ := strconv.Atoi(args[3])

	puzzle := models.NewPuzzle(year, day)
	if err != nil {
		log.Fatal("Unable to load puzzle data!", "year", year, "day", day, "err", err)
	}
	userInput, err := puzzle.GetUserPuzzleInput(user.GetToken())
	fmt.Print(string(userInput))
	return
}

// TODO: Handle days as well

// `leaderboard [year] [day]` command
func leaderboard(args []string) {
	if len(args) != 3 {
		fmt.Println("Only works with yearly leaderboards for right now. Run `aocli help leaderboard`")
		return
	}

	// Ensure that the API is initialized. We don't actually need the user
	year, _ := strconv.Atoi(args[2])
	var day int
	// if len(args) == 4 {
	// 	day, _ = strconv.Atoi(args[3])
	// }

	lb := models.NewLeaderboard(year, day)

	if lb == nil {
		log.Fatal("Unable to load/create leaderboard!")
		return
	}

	lb.Display()
}

// `submit [year] [day] [part] [answer]` command
func submit(args []string) {

}

func run(args []string) {

}

// `view [year] [day]` command
// Desc: Pretty-prints the puzzle's page data
func view(args []string) {
	// TODO: Load user since it needs the token to know whether or not to include Part B
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
	pageData := puzzle.GetPuzzlePageData(user.GetToken())
	pageData.PrintPageData()
}

// `health` command
// Desc: Checks if a session key is available
func health() {
	sessionKey, err := session.GetSessionToken()
	if err != nil {
		log.Error("Test failed! Couldn't properly load a session key.", "err", err)
	}

	log.Info("Test succeeded! Properly loaded session key", "key", sessionKey)
}

// Command:	`test`
// Desc:	Does whatever I need to test at the time :)
func test() {
	lb := models.NewLeaderboard(2020, 0)
	lb.Display()
	//dirparse.GetDayAndYearFromCWD()
}
