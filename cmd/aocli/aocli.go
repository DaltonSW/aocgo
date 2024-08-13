package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/resources"
	"dalton.dog/aocgo/internal/session"
	"dalton.dog/aocgo/internal/utils"
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

var User *resources.User

func main() {

	// Flag Parsing
	debugFlag := flag.Bool("debug", false, "Use to enable debug logging")
	profFlag := flag.Bool("prof", false, "Use to enable performance profiling")
	flag.Parse()

	// Debug Logging
	if *debugFlag {
		debugFile, err := os.Create("./debug.log")
		if err != nil {
			log.Fatal("Unable to create debug file.", "error", err)
		}
		defer debugFile.Close()
		log.SetOutput(debugFile)
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		log.SetTimeFormat(time.StampMicro)
	} else {
		logStyles := log.DefaultStyles()
		logStyles.Levels[log.FatalLevel] = lipgloss.NewStyle().
			SetString("FATAL").
			Padding(0, 1).
			Background(lipgloss.Color("#af5fd7")).
			Foreground(lipgloss.Color("0"))
		logStyles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
			SetString("ERROR").
			Padding(0, 1).
			Background(lipgloss.Color("204")).
			Foreground(lipgloss.Color("0"))
		logStyles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
			SetString("INFO").
			Padding(0, 1).
			Background(lipgloss.Color("#42ccd4")).
			Foreground(lipgloss.Color("0"))
		log.SetTimeFormat(time.Kitchen)
		log.SetStyles(logStyles)
	}

	if *profFlag {
		profFile, err := os.Create("./aoc.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file.", "error", err)
		}
		defer profFile.Close()

		pprof.StartCPUProfile(profFile)
		defer pprof.StopCPUProfile()
	}

	args := os.Args
	if len(args) == 1 {
		fmt.Println(helpBodyStyle.Render("Welcome to aocli! Try running `aocli help` for a list of available commands."))
		os.Exit(0)
	}

	if args[1] == "health" {
		health()
		return
	} else if args[1] == "help" {
		help(args)
		return
	}

	log.Debug("Trying to create user")
	user, err := resources.NewUser("")
	if err != nil {
		log.Fatal("Unable to create user to run requests as. Check the README to ensure you have the proper setup. Then try running `aocli health`.")
	}
	log.Debug("Created user")

	log.Debug("Trying to startup database")
	cache.StartupDBM(user.GetToken())
	defer cache.ShutdownDBM()
	log.Debug("Database started")

	log.Debug("Args parsed", "args", args[1:])

	switch args[1] {
	case "check-update":
		update := checkForUpdate()
		if update {
			fmt.Printf("New version available! Run `aocli update` to get the new version.")
		}
	case "get":
		get(args, user)
	case "submit":
		submit(args, user)
	case "leaderboard":
		leaderboard(args)
	case "load-user":
		loadUser(args, user)
	// case "run":
	// 	run(args)
	case "view":
		view(args, user)
	case "test":
		test(user)
	case "clear-user":
		clearUser(user)
	case "update":
		update()
	default:
		fmt.Println("Not a valid command! Run `aocli help` to see valid commands.")
	}
	return
}

func clearUser(user *resources.User) {
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
func get(args []string, user *resources.User) {
	var year int
	var day int
	var err error

	if len(args) < 4 {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal(err)
		}
	} else {

		year, err = utils.ParseYear(args[2])
		if err != nil {
			log.Fatal(err)
		}

		day, err = utils.ParseDay(args[3])
		if err != nil {
			log.Fatal(err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	userInput, _ := puzzle.GetUserInput()

	out, _ := os.Create("./input.txt")
	defer out.Close()
	out.Write(userInput)

	log.Info("Input saved to input.txt!")
	return
}

// TODO: Handle days as well

// `leaderboard [year] [day]` command
func leaderboard(args []string) {
	if len(args) != 3 {
		log.Error("Only works with yearly leaderboards for right now. Run `aocli help leaderboard`")
		return
	}

	// Ensure that the API is initialized. We don't actually need the user
	year, _ := strconv.Atoi(args[2])
	var day int
	// if len(args) == 4 {
	// 	day, _ = strconv.Atoi(args[3])
	// }

	lb := resources.NewLeaderboard(year, day)

	if lb == nil {
		log.Fatal("Unable to load/create leaderboard!")
		return
	}

	lb.Display()
}

func loadUser(args []string, user *resources.User) {
	logger := log.New(os.Stdout)
	year := 2015
	var maxYear int
	if time.Now().Month() == time.December {
		maxYear = time.Now().Year()
	} else {
		maxYear = time.Now().Year() - 1
	}

	if len(args) > 2 {
		parseYear, err := utils.ParseYear(args[2])
		if err != nil {
			log.Fatal(err)
		}

		year = parseYear
		maxYear = parseYear
	}

	numStars := make(map[int]int)

	for year <= maxYear {
		logger.Info("Loading year", "year", year)
		numStars[year] = 0
		day := 1
		for day <= 25 {
			logger.Info("Loading day", "day", day)
			puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())

			user.Years[year][day] = puzzle

			if puzzle.AnswerOne != "" {
				logger.Info("Answer one found!", "year", year, "day", day, "answer", puzzle.AnswerOne)
				user.NumStars++
				numStars[year]++
				if puzzle.AnswerTwo != "" {
					logger.Info("Answer two found!", "year", year, "day", day, "answer", puzzle.AnswerTwo)
					user.NumStars++
					numStars[year]++
				}
			}

			day++
		}

		// There's only 1 puzzle on Day 25, so if they've earned 49 stars, they get the 50th for free
		if numStars[year] == 49 {
			user.NumStars++
			numStars[year]++
		}
		logger.Info("Ending year", "Stars found", numStars[year])

		year++
	}

	for val, key := range numStars {
		fmt.Printf("%d -- %d\n", key, val)
	}
}

// `submit [answer] -y <yyyy> -d <dd>` command
func submit(args []string, user *resources.User) {
	year, day, err := utils.GetYearAndDayFromCWD()
	if err != nil {
		log.Fatal(err)
	}

	answer := args[2]

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.SessionTok)

	answerResp, message := puzzle.SubmitAnswer(answer)

	if answerResp == resources.CorrectAnswer {
		correctStyle := lipgloss.NewStyle().Background(lipgloss.Color("34"))
		fmt.Println(correctStyle.Render("Correct answer!"))
		fmt.Println(correctStyle.Render(message))
		user.NumStars++
	} else if answerResp == resources.IncorrectAnswer {
		incorrectStyle := lipgloss.NewStyle().Background(lipgloss.Color("124"))
		fmt.Println(incorrectStyle.Render("Incorrect answer!"))
		fmt.Println(incorrectStyle.Render(message))
	} else if answerResp == resources.WarningAnswer {
		incorrectStyle := lipgloss.NewStyle().Background(lipgloss.Color("184"))
		fmt.Println(incorrectStyle.Render("Answer not submitted!"))
		fmt.Println(incorrectStyle.Render(message))
	}
}

func run(args []string) {

}

// `view [year] [day]` command
// Desc: Pretty-prints the puzzle's page data
func view(args []string, user *resources.User) {
	var year int
	var day int
	var err error

	if len(args) < 4 {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal("Unable to parse year/day from current directory.")
		}
	} else {
		year, err = utils.ParseYear(args[2])
		if err != nil {
			log.Fatal("Unable to parse year from current directory.")
		}

		day, err = utils.ParseDay(args[3])
		if err != nil {
			log.Fatal("Unable to parse day from current directory.")
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	// userInput, _ := puzzle.GetUserInput()
	puzzle.Display()
	// tui.NewPuzzleViewport(puzzle.GetPrettyPageData(), puzzle.Title, puzzle.URL, userInput)
}

// `health` command
// Desc: Checks if a session key is available
func health() {
	sessionKey, err := session.GetSessionToken()
	if err != nil {
		log.Fatal("Test failed! Couldn't properly load a session key.", "err", err)
	}

	log.Info("Test succeeded! Properly loaded session key", "key", sessionKey)
}

// Command:	`test`
// Desc:	Does whatever I need to test at the time :)
func test(user *resources.User) {
	puzzle := resources.LoadOrCreatePuzzle(2023, 1, user.GetToken())
	puzzle.Display()
}
