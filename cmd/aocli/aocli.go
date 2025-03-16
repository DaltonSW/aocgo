package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/session"
	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/charmbracelet/log"
)

func main() {
	rootCmd.Execute()
}

// region: User-agnostic commands

// NOTE: Cobra handles help inherently now, but this is still
//	here in case I want to override it again in the future

// Help prints info and a list of commands
// Command: `aocli help`
// Params:
//
//	[command] - command name to print Help for
// func Help(args []string) {
// 	utils.ClearTerminal()
//
// 	// They requested help for a specific command
// 	if len(args) == 3 {
// 		commandName := args[2]
// 		ht, ok := HelpTextMap[commandName]
// 		if ok {
// 			ht.Print()
// 		} else {
// 			log.Error("Not a valid command!")
// 		}
// 		return
// 	}
//
// 	// Otherwise they just open-endedly requested help
// 	ht, ok := HelpTextMap["aocli"]
// 	if ok {
// 		outS := NameStyle.Render("NAME:  ")
// 		outS += ht.name + "\n\n"
//
// 		outS += UseStyle.Render("USAGE: ")
// 		outS += ht.use + "\n\n"
//
// 		outS += DescStyle.Render("DESCRIPTION")
// 		outS += ht.desc + "\n\n"
//
// 		outS += ArgsStyle.Render("COMMANDS")
// 		outS += ht.args
//
// 		fmt.Println(outS)
// 	}
//
// }

// Leaderboard obtains and displays Leaderboard information for a specific year or day
// Command: `aocli leaderboard -y yyyy [-d dd]`
// Params:
//
//	(Req) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func Leaderboard(yearIn, dayIn string) {
	var year int
	var day int
	var err error
	var lb resources.ViewableLB

	if yearIn == "0" {
		year, day, err := utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal("Error loading leaderboard based on current directory!", "err", err)
		}
		lb = resources.LoadOrCreateLeaderboard(year, day)
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal("Error parsing year!", "err", err)
		}
	}

	if lb == nil {
		if dayIn != "0" {
			day, err = utils.ParseDay(dayIn)
			if err != nil {
				log.Fatal("Error parsing day from args.", "err", err)
			}
			lb = resources.LoadOrCreateLeaderboard(year, day)
		} else {
			lb = resources.LoadOrCreateLeaderboard(year, 0)
		}
	}

	if lb == nil {
		log.Fatal("Unable to load/create leaderboard!")
		return
	}

	resources.NewLeaderboardViewport(lb.GetContent(), lb.GetTitle())
}

// Health will check if a session key is available so that the program can run.
// Command: `aocli health`
func Health() {
	sessionToken, err := session.GetSessionToken(true)

	if err != nil {
		log.Fatal("Test failed! Couldn't properly load a session token.", "err", err)
	}

	log.Info("Session token check success!")

	api.InitClient(sessionToken)

	log.Info("API Client initialization check success!")

	// user, err := resources.NewUser(sessionToken)

	resources.LoadOrCreatePuzzle(2016, 1, sessionToken)

	log.Info("Session token appears to be valid, happy solving!")
}

// User-specific functions

// Submit will submit the answer provided.
// If date arguments aren't provided, they will be parsed from the current directory.
// Command: `aocli submit <answer> [-y yyyy -d dd --part {1|2}]`
// Params:
//
//	(Req) answer - Answer to submit to the server
//	(Opt) year   - 2 or 4 digit year (16 or 2016)
//	(Opt) day    - 1 or 2 digit day (1, 01, 21)
func Submit(user *resources.User, yearIn, dayIn, answer string, partIn int) {
	var year, day int
	var err error

	if yearIn != "0" {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal("Couldn't parse provided year argument.", "err", err)
		}
	}

	if dayIn != "0" {
		day, err = utils.ParseDay(dayIn)
		if err != nil {
			log.Fatal("Couldn't parse provided day argument.", "err", err)
		}
	}

	if day == 0 || year == 0 {
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

	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.SessionTok)

	var part int
	if partIn < 0 || partIn > 2 {
		part = 0
		log.Error("Part provided by option is invalid. Using default part for submission.")
	} else {
		part = partIn
	}

	answerResp, message := puzzle.SubmitAnswer(answer, part)

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
// Command: `reload [-y yyyy -d dd]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func Reload(user *resources.User, yearIn, dayIn string) {
	var year int
	var day int
	var err error

	if yearIn == "0" || dayIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal(err)
		}

		day, err = utils.ParseDay(dayIn)
		if err != nil {
			log.Fatal(err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	puzzle.ReloadPuzzleData()
}

// User will print out a table visualization of the user's star progress.
// Command: `aocli user [--clear]`
func User(user *resources.User, clearUser bool) {
	if ClearUser {
		cache.ClearUserDatabase(UserRsrc.SessionTok)
	} else {
		UserRsrc.Display()
	}
}

// View will pretty print the puzzle's page data.
// Command: `aocli view [-y yyyy -d dd]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
func View(user *resources.User, yearIn, dayIn string) {
	var year int
	var day int
	var err error

	if yearIn == "0" || dayIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal("Unable to parse year/day from current directory.", "err", err)
		}
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal("Unable to parse year from current directory.", "err", err)
		}

		day, err = utils.ParseDay(dayIn)
		if err != nil {
			log.Fatal("Unable to parse day from current directory.", "err", err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	puzzle.Display()
}

// Get obtains input data for a specific day, outputting it to the current directory `input.txt`.
// Command: `aocli get [-y yyyy -d dd -o output_name.txt]`
// Params:
//
//	(Opt) year - 2 or 4 digit year (16 or 2016)
//	(Opt) day  - 1 or 2 digit day (1, 01, 21)
//	(Opt) filename  - overrides "input.txt" name if it's provided
func Get(user *resources.User, yearIn, dayIn string, filename string) {
	var year int
	var day int
	var err error

	if yearIn == "0" || dayIn == "0" {
		year, day, err = utils.GetYearAndDayFromCWD()
		if err != nil {
			log.Fatal("Unable to parse year/day from current directory.", "err", err)
		}
	} else {
		year, err = utils.ParseYear(yearIn)
		if err != nil {
			log.Fatal("Unable to parse year from current directory.", "err", err)
		}

		day, err = utils.ParseDay(dayIn)
		if err != nil {
			log.Fatal("Unable to parse day from current directory.", "err", err)
		}
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, user.GetToken())
	userInput, _ := puzzle.GetUserInput()

	out, _ := os.Create(filename)
	defer out.Close()
	out.Write(userInput)

	log.Infof("Input saved to %v!", filename)
}

// New will make a copy of the given base file into the given out file in a ./yearIn/dayIn/ directory.
// Command `aocli new [-y yyyy -d dd -b base.go -o main.go]`
func New(yearIn, dayIn, baseFile, outFile string) {
	if yearIn == "0" {
		log.Fatal("Must provide a year with the -y option.")
	}

	if dayIn == "0" {
		log.Fatal("Must provide a day with the -d option.")
	}

	// Build the directory path, e.g. "./2025/15"
	dirPath := filepath.Join(".", yearIn, dayIn)

	// Create the directory (with parents, if needed).
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dirPath, err)
	}

	// Open the source/base file.
	src, err := os.Open(baseFile)
	if err != nil {
		log.Fatalf("Failed to open base file %s: %v", baseFile, err)
	}
	defer src.Close()

	// Create the destination file in the new directory.
	dstPath := filepath.Join(dirPath, outFile)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", dstPath, err)
	}
	defer dst.Close()

	// Copy the entire contents from baseFile to outFile.
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("Failed to copy data from %s to %s: %v", baseFile, dstPath, err)
	}

	log.Infof("Successfully copied %s to %s", baseFile, dstPath)

}

// Test does whatever I need to test at the time :)
// Command:	`aocli test`
// func test(user *resources.User) {
// 	dispName := user.LoadDisplayName()
// 	log.Info(user.DisplayName)
// 	log.Info(dispName)
// }
