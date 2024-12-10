// Package aocgo provides functions to get your Advent of Code puzzle inputs in a non-intrusive way.
package aocgo

import (
	"fmt"
	"strings"
	"time"

	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/session"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var correctTestColor = lipgloss.Color("#1d8509")
var incorrectTestColor = lipgloss.Color("#c40e4f")
var puzzleSolveColor = lipgloss.Color("#674dd9")

// InputData interface is a Generic wrapper around the possible forms you can get puzzle input data in
type InputData interface {
	string | []string | [][]string | []byte
}

// AnswerData interface is for any form that the puzzle answer can be output as
type AnswerData interface {
	int | string
}

// Func is a function that will take in
type Solver[input InputData, answer AnswerData] func(input) answer

// RunTest will run a given function with the given input, and compare it against a known output.
// The result will print with a color based on if your function returns the same result as was expected.
func RunTest[In InputData, Out AnswerData](title string, solver Solver[In, Out], inputData In, expected Out) {
	start := time.Now()

	answer := solver(inputData)

	timeTaken := time.Since(start)

	// Convert the answer to a string
	answerStr := fmt.Sprintf("%v", answer)

	expectedStr := fmt.Sprintf("%v", expected)

	var outColor lipgloss.Color

	if answerStr == expectedStr {
		outColor = correctTestColor
	} else {
		outColor = incorrectTestColor
	}

	// Create styles using lipgloss
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(outColor).
		Padding(0, 1).Margin(1, 0, 0).
		AlignHorizontal(lipgloss.Center)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(outColor).
		Padding(0, 2)

	// Pretty info
	prettyAnswer := fmt.Sprintf("Answer  : %v\n", answerStr)
	prettyExpected := fmt.Sprintf("Expected: %v\n", expectedStr)
	prettyTime := fmt.Sprintf("Runtime : %v", timeTaken)

	// Render the answer inside the border
	wrappedInfo := borderStyle.Render(prettyAnswer + prettyExpected + prettyTime)

	// Render the title
	titleBox := titleStyle.Width(lipgloss.Width(wrappedInfo)).Render(title)

	// Combine the title and wrapped answer
	output := lipgloss.JoinVertical(lipgloss.Top, titleBox, wrappedInfo)

	// Display the output
	fmt.Println(output)
}

// RunSolve will attempt to run the input function with the input data.
// It will print out information about the function run in a pretty table.
func RunSolve[In InputData, Out AnswerData](title string, solver Solver[In, Out], inputData In) {
	start := time.Now()

	answer := solver(inputData)

	timeTaken := time.Since(start)

	// Convert the answer to a string
	answerStr := fmt.Sprintf("%v", answer)

	// Create styles using lipgloss
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(puzzleSolveColor).
		Padding(0, 1).Margin(1, 0, 0).
		AlignHorizontal(lipgloss.Center)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(puzzleSolveColor).
		Padding(0, 2)

	// Pretty info
	prettyAnswer := fmt.Sprintf("Answer  : %v\n", answerStr)
	prettyTime := fmt.Sprintf("Runtime : %v", timeTaken)

	// Render the answer inside the border
	wrappedInfo := borderStyle.Render(prettyAnswer + prettyTime)

	// Render the title
	titleBox := titleStyle.Width(lipgloss.Width(wrappedInfo)).Render(title)

	// Combine the title and wrapped answer
	output := lipgloss.JoinVertical(lipgloss.Top, titleBox, wrappedInfo)

	// Display the output
	fmt.Println(output)
}

// GetInputAsByteArray will return the user's puzzle input, as determined by the file's working directory, as an array of bytes.
func GetInputAsByteArray() []byte {
	year, day, err := utils.GetYearAndDayFromCWD()
	if err != nil {
		log.Fatal(err)
	}

	return getData(year, day)
}

// GetInputAsString will return the user's puzzle input, as determined by the file's working directory, as a single string.
func GetInputAsString() string {
	return string(GetInputAsByteArray())
}

// GetInputAsLineArray will return the user's puzzle input, as determined by the file's working directory, as an array of strings, split on newline.
func GetInputAsLineArray() []string {
	return strings.Split(GetInputAsString(), "\n")
}

// GetInputAsCharMatrix will return the user's puzzle input, as determined by the file's working directory, as a 2D matrix, split on newlines and then by every character
func GetInputAsCharMatrix() [][]string {
	var out [][]string
	for _, line := range GetInputAsLineArray() {
		out = append(out, strings.Split(line, ""))
	}

	return out
}

func getData(year int, day int) []byte {
	userToken, err := session.GetSessionToken(false)
	if err != nil {
		log.Fatal(err)
	}

	_, err = resources.NewUser(userToken)
	if err != nil {
		log.Fatal(err)
	}

	err = cache.StartupDBM(userToken)
	if err != nil {
		log.Fatal(err)
	}

	puzzle := resources.LoadOrCreatePuzzle(year, day, userToken)
	input, err := puzzle.GetUserInput()
	if err != nil {
		log.Fatal(err)
	}
	return input
}
