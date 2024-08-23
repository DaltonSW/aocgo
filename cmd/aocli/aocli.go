package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"dalton.dog/aocgo/internal/cache"
	"dalton.dog/aocgo/internal/resources"
	"dalton.dog/aocgo/internal/session"
	"dalton.dog/aocgo/internal/styles"
	"dalton.dog/aocgo/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
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
		RunLandingPage()
		return
		// fmt.Println("Welcome to aocli! Try running `aocli help` for a list of available commands.")
		// os.Exit(0)
	}

	log.Debug("Args parsed", "args", args[1:])

	// User agnostic commands
	if args[1] == "health" {
		Health()
		return
	} else if args[1] == "help" {
		Help(args)
		return
	} else if args[1] == "update" {
		RunUpdateModel() // Runs the request such that masterAPI doesn't need to be initialized
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

	// User dependent functions. I recognize that I used if/else above and switch statement here, oh well
	switch args[1] {
	case "version":
		Version()
	case "get":
		Get(args, user)
	case "submit":
		Submit(args, user)
	case "reload":
		Reload(args, user)
	// case "run":
	// 	run(args)
	case "view":
		View(args, user)
	// case "test":
	// 	test(user)
	case "user":
		User(args, user)
	case "clear-user":
		ClearUser(user)
	case "leaderboard":
		Leaderboard(args)
	default:
		fmt.Println("Not a valid command! Run `aocli help` to see valid commands.")
	}

	latestVersion, err := getLatestRelease()
	if err != nil {
		log.Fatal("Error checking for updates!", "error", err)
	}

	if !strings.Contains(latestVersion.TagName, "aocli-") {
		return
	} else {
		latestVersion.TagName = strings.Replace(latestVersion.TagName, "aocli-", "", 1)
	}

	latestSemVer := semver.Canonical(latestVersion.TagName)
	currentSemVer := semver.Canonical(currentVersion)

	if semver.Compare(latestSemVer, currentSemVer) > 0 {
		fmt.Println(styles.GlobalSpacingStyle.Render(styles.NormalTextStyle.Render(updateMessage)))
	}

	return
}

// ClearUser will delete the database file associated with the current session token.
// Command: `aocli clear-user`
func ClearUser(user *resources.User) {
	cache.ClearUserDatabase(user.SessionTok)
}

// Help prints info and a list of commands
// Command: `aocli help`
// Params:
//
//	[command] - command name to print Help for
func Help(args []string) {
	utils.ClearTerminal()

	// They requested help for a specific command
	if len(args) == 3 {
		commandName := args[2]
		ht, ok := HelpTextMap[commandName]
		if ok {
			ht.Print()
		} else {
			log.Error("Not a valid command!")
		}
		return
	}

	// Otherwise they just open-endedly requested help
	ht, ok := HelpTextMap["aocli"]
	if ok {
		outS := NameStyle.Render("NAME:  ")
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
// Command: `aocli get [year] [day]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
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
// Command: `aocli leaderboard year [day]`
// Params:
//
//	(Req) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func Leaderboard(args []string) {
	// TODO: Validation and help message
	year, err := utils.ParseYear(args[2])
	if err != nil {
		log.Fatal("Error parsing year!", "err", err)
	}

	var lb resources.ViewableLB
	if len(args) == 4 {
		day, err := utils.ParseDay(args[3])
		if err != nil {
			log.Fatal("Error parsing day from args.", "err", err)
		}
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

// Submit will submit the answer provided.
// If date arguments aren't provided, they will be parsed from the current directory.
// Command: `aocli submit <answer> [year] [day]`
// Params:
//
//	(Req) answer - Answer to submit to the server
//	(Opt) year   - 2 or 4 digit year (16 or 2016)
//	(Opt) day    - 1 or 2 digit day (1, 01, 21)
func Submit(args []string, user *resources.User) {
	var year, day int
	var err error

	if len(args) > 3 {
		year, err = utils.ParseYear(args[3])
		log.Fatal("Couldn't parse provided year argument.", "err", err)
	}

	if len(args) > 4 {
		day, err = utils.ParseDay(args[4])
		log.Fatal("Couldn't parse provided day argument.", "err", err)
	}

	parseYear, parseDay, err := utils.GetYearAndDayFromCWD()
	if err != nil {
		log.Fatal(err)
	}

	if year == 0 {
		year = parseYear
	}

	if day == 0 {
		day = parseDay
	}

	answer := args[2]

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.SessionTok)

	// TODO: Allow this to take a part as an argument
	answerResp, message := puzzle.SubmitAnswer(answer, 0)

	if answerResp == resources.CorrectAnswer {
		fmt.Println(styles.CorrectAnswerStyle.Render("Correct answer!"))
		fmt.Println(styles.CorrectAnswerStyle.Render(message))
		user.NumStars++
	} else if answerResp == resources.IncorrectAnswer {
		fmt.Println(styles.IncorrectAnswerStyle.Render("Incorrect answer!"))
		fmt.Println(styles.IncorrectAnswerStyle.Render(message))
	} else if answerResp == resources.WarningAnswer {
		fmt.Println(styles.WarningAnswerStyle.Render("Answer not submitted!"))
		fmt.Println(styles.WarningAnswerStyle.Render(message))
	} else if answerResp == resources.NeutralAnswer {
		fmt.Println(styles.NeutralAnswerStyle.Render("Answer not submitted!"))
		fmt.Println(styles.NeutralAnswerStyle.Render(message))
	}
}

// Reload will force reload the puzzle data for a specific day
// Command: `reload [year] [day]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
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

// User will print out a table visualization of the user's star progress.
// Command: `aocli user`
func User(args []string, user *resources.User) {
	user.Display()
}

// View will pretty print the puzzle's page data.
// Command: `aocli view [year] [day]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func View(args []string, user *resources.User) {
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

// Health will check if a session key is available so that the program can run.
// Command: `aocli health`
func Health() {
	sessionKey, err := session.GetSessionToken()
	if err != nil {
		log.Fatal("Test failed! Couldn't properly load a session key.", "err", err)
	}

	log.Info("Test succeeded! Properly loaded session key", "key", sessionKey)
}

// Test does whatever I need to test at the time :)
// Command:	`aocli test`
// func test(user *resources.User) {
// 	dispName := user.LoadDisplayName()
// 	log.Info(user.DisplayName)
// 	log.Info(dispName)
// }
