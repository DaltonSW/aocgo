package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"
)

var submitCmd = &cobra.Command{
	Use:   "submit [-p {1|2}] <answer>",
	Short: "Submits the given answer to a puzzle.",
	Args:  cobra.ExactArgs(1),
	Annotations: map[string]string{
		requiresAuthAnnotation: "true",
	},
	Run: func(cmd *cobra.Command, args []string) {
		Submit(UserRsrc, Year, Day, args[0], AnswerPart)
	},
}

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

	switch answerResp {
	case resources.CorrectAnswer:
		fmt.Println(styles.CorrectAnswerStyle.Render("Correct answer!"))
		fmt.Println(styles.CorrectAnswerStyle.Render(message))
		user.NumStars++
	case resources.IncorrectAnswer:
		fmt.Println(styles.IncorrectAnswerStyle.Render("Incorrect answer!"))
		fmt.Println(styles.IncorrectAnswerStyle.Render(message))
	case resources.WarningAnswer:
		fmt.Println(styles.WarningAnswerStyle.Render("Answer not submitted!"))
		fmt.Println(styles.WarningAnswerStyle.Render(message))
	case resources.NeutralAnswer:
		fmt.Println(styles.NeutralAnswerStyle.Render("Answer not submitted!"))
		fmt.Println(styles.NeutralAnswerStyle.Render(message))
	}
}
