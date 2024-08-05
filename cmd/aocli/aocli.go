package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/models"
	"dalton.dog/aocgo/internal/session"
	"dalton.dog/aocgo/internal/tui"
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
	debugFile, err := os.Create("./debug.log")
	if err != nil {
		log.Fatal("Unable to create debug file.", "error", err)
	}
	defer debugFile.Close()

	log.SetOutput(debugFile)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetTimeFormat(time.StampMicro)

	args := os.Args
	if len(args) == 1 {
		fmt.Println(helpBodyStyle.Render("Welcome to aocli! Try running `aocli help` for a list of available commands."))
		os.Exit(0)
	}

	log.Debug("Trying to create user")
	user, err := models.NewUser("")
	if err != nil {
		log.Error("Unable to create user to run requests as. Try running `aocli health`.")
	}
	log.Debug("Created user")

	// log.Debug("Trying to startup database")
	// cache.StartupDBM(user.GetToken())
	// defer cache.ShutdownDBM()
	// log.Debug("Database started")

	log.Debug("Args parsed", "args", args[1:])

	switch args[1] {
	case "check-update":
		update := checkForUpdate()
		if update {
			fmt.Printf("New version available! Run `aocli update` to get the new version.")
		}
	case "get":
		get(args, user)
	case "health":
		health()
	case "help":
		help(args)
	// case "submit":
	// 	submit(args)
	case "leaderboard":
		leaderboard(args)
	// case "run":
	// 	run(args)
	case "view":
		view(args, user)
	case "test":
		test(user)
	case "clear-data":
		clearData(user)
	case "update":
		update()
	default:
		fmt.Println("Not a valid command! Run `aocli help` to see valid commands.")
	}
	return
}

func clearData(user *models.User) {
	cache.ClearUserDatabase(user.SessionTok)
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
// Desc: Gets input data for a specific day, outputting it to the current directory as `input.txt`
// Params:
//
//	[year] - 2 or 4 digit year (16 or 2016)
//	[day]  - 1 or 2 digit day (1, 01, 21)
func get(args []string, user *models.User) {
	// TODO: Validate input
	if len(args) < 4 {
		return
		// TODO: Try loading with today
		// TODO: Print `get` help message
	}

	year, _ := strconv.Atoi(args[2])
	day, _ := strconv.Atoi(args[3])

	puzzle := models.LoadOrCreatePuzzle(year, day, user.GetToken())
	userInput, _ := puzzle.GetUserInput()

	out, _ := os.Create("./input.txt")
	out.Write(userInput)
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
func view(args []string, user *models.User) {
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

	puzzle := models.LoadOrCreatePuzzle(year, day, user.GetToken())
	tui.StartViewportWithArr(puzzle.GetPrettyPageData(), puzzle.Title, true)
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
func test(user *models.User) {
	puzzle := models.LoadOrCreatePuzzle(2023, 1, user.GetToken())
	tui.StartViewportWithArr(puzzle.GetPrettyPageData(), puzzle.Title, true)
}
