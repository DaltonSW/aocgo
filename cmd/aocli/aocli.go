package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/resources"
	"dalton.dog/aocgo/internal/session"
	"dalton.dog/aocgo/internal/styles"
	"dalton.dog/aocgo/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func main() {
	// HACK: Triggering this immediately so it doesn't run into
	// a possible race condition with BubbleTea competing for stdout
	// https://github.com/charmbracelet/bubbletea/issues/1071
	_ = lipgloss.DefaultRenderer().HasDarkBackground()

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

		logStyles.Levels[log.FatalLevel] = styles.LoggerFatalStyle
		logStyles.Levels[log.ErrorLevel] = styles.LoggerErrorStyle
		logStyles.Levels[log.InfoLevel] = styles.LoggerInfoStyle

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
		fmt.Println(styles.HelpTextStyle.Render("Welcome to aocli! Try running `aocli help` for a list of available commands."))
		os.Exit(0)
	}

	log.Debug("Args parsed", "args", args[1:])

	// User agnostic commands
	if args[1] == "health" {
		health()
		return
	} else if args[1] == "help" {
		Help(args)
		return
	} else if args[1] == "update" {
		// Update()
		RunUpdateModel()
		return
	} else if args[1] == "leaderboard" {
		Leaderboard(args)
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

	// User dependent functions. I recognize that I used if/else above and switch statement here, oh well
	switch args[1] {
	case "version":
		update := Version()
		if update {
			fmt.Printf("New version available! Run `aocli update` to get the new version (or `sudo aocli update` if your executable is in a protected location)")
		}
	case "get":
		Get(args, user)
	case "submit":
		submit(args, user)
	// case "load-user":
	//	 loadUser(args, user)
	case "reload":
		Reload(args, user)
	// case "run":
	// 	run(args)
	case "view":
		view(args, user)
	// case "test":
	// 	test(user)
	case "user":
		User(args, user)
	case "clear-user":
		clearUser(user)
	default:
		fmt.Println("Not a valid command! Run `aocli help` to see valid commands.")
	}
	return
}

func clearUser(user *resources.User) {
	cache.ClearUserDatabase(user.SessionTok)
}

// Help prints info and a list of commands
// Associated command: `help`
// Params:
//
//	[command] - command name to print Help for
func Help(args []string) {
	// Clear terminal
	fmt.Print("\033[H\033[2J")

	// Too many args
	if len(args) > 3 {
		fmt.Println(styles.HelpTextStyle.Render("Too many arguments passed!"))
		return
	}

	// They requested help for a specific command
	if len(args) == 3 {
		commandName := args[2]
		helptext, ok := HelpTextMap[commandName]
		if ok {
			helptext.Print()
		} else {
			fmt.Println(styles.HelpTextStyle.Render("Not a valid command!"))
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

// Get obtains input data for a specific day, outputting it to the current directory as `input.txt`
// Command: `get [year] [day]`
// Params:
//
//	[year] - 2 or 4 digit year (16 or 2016)
//	[day]  - 1 or 2 digit day (1, 01, 21)
func Get(args []string, user *resources.User) {
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

// Leaderboard obtains and displays Leaderboard information for a specific year or day
// Command: `Leaderboard year [day]`
// Params:
//
//	year  - 2 or 4 digit year (16 or 2016)
//	[day] - 1 or 2 digit day (1, 01, 21)
func Leaderboard(args []string) {
	// TODO: Validation and help message
	year, err := utils.ParseYear(args[2])
	if err != nil {
		log.Fatal("Error parsing year!", "err", err)
	}

	var lb resources.ViewableLB
	if len(args) == 4 {
		day, _ := utils.ParseDay(args[3])
		lb = resources.NewDayLB(year, day)
	} else {
		lb = resources.NewYearLB(year)
	}

	if lb == nil {
		log.Fatal("Unable to load/create leaderboard!")
		return
	}

	resources.NewLeaderboardViewport(lb.GetContent(), lb.GetTitle())
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

// TODO: Document

// `submit [answer] -y <yyyy> -d <dd>` command
func submit(args []string, user *resources.User) {
	year, day, err := utils.GetYearAndDayFromCWD()
	if err != nil {
		log.Fatal(err)
	}

	answer := args[2]

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.SessionTok)

	// TODO: Allow this to take a part as an argument
	answerResp, message := puzzle.SubmitAnswer(answer, 0)

	// TODO: Move these styles to the 'styles' package
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

// Reload will force reload the puzzle data for a specific day
// Command: `reload [year] [day]`
// Params:
//
//	[year] - 2 or 4 digit year (16 or 2016)
//	[day]  - 1 or 2 digit day (1, 01, 21)
func Reload(args []string, user *resources.User) {
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
	puzzle.ReloadPuzzleData()
}

// TODO: Implement
func run(args []string) {

}

func User(args []string, user *resources.User) {
	user.LoadUser()
	user.Display()
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
			log.Fatal("Unable to parse year/day from current directory.", "err", err)
		}
	} else {
		year, err = utils.ParseYear(args[2])
		if err != nil {
			log.Fatal("Unable to parse year from current directory.", "err", err)
		}

		day, err = utils.ParseDay(args[3])
		if err != nil {
			log.Fatal("Unable to parse day from current directory.", "err", err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	puzzle.Display()
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
// func test(user *resources.User) {
// 	lb := resources.NewDayLB(2016, 1)
// 	resources.NewLeaderboardViewport(lb.GetContent(), lb.GetTitle())
// }
