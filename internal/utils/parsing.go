package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	FIRST_YEAR = 2015
)

func GetYearAndDayFromCWD() (int, int, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, 0, err
	}
	splitCWD := strings.Split(cwd, string(os.PathSeparator))

	cwdLen := len(splitCWD)
	curDir, parentDir := splitCWD[cwdLen-1], splitCWD[cwdLen-2]

	return GetYearAndDayFromDirInput(curDir, parentDir)
}

func GetYearAndDayFromDirInput(curDir, parentDir string) (int, int, error) {
	day, err := ParseDay(curDir)
	if err != nil {
		return 0, 0, err
	}

	year, err := ParseYear(parentDir)
	if err != nil {
		return day, 0, err
	}

	return year, day, nil
}

func ParseYear(yearStr string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	yearStr = strings.TrimSpace(yearStr)
	match := re.FindString(yearStr)
	var outYear, year int
	var err error
	if len(match) == 2 {
		year, err = strconv.Atoi(match)
		if err != nil {
			return 0, err
		}
		outYear = 2000 + year
	} else {
		outYear, err = strconv.Atoi(match)
		if err != nil {
			return 0, err
		}
	}

	if outYear < FIRST_YEAR {
		return 0, errors.New(fmt.Sprintf("Year parsed to be earlier than %v.", FIRST_YEAR))
	}

	var maxYear int
	if time.Now().Month() == time.December {
		maxYear = time.Now().Year()
	} else {
		maxYear = time.Now().Year() - 1
	}

	if outYear > maxYear {
		return 0, errors.New(fmt.Sprintf("Year parsed to be later than %v.", maxYear))
	}

	return outYear, nil
}

func ParseDay(dayStr string) (int, error) {
	dayStr = strings.TrimSpace(dayStr)
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(dayStr)
	outInt, err := strconv.Atoi(match)
	if err != nil {
		return 0, nil
	}

	if outInt < 1 {
		return 0, errors.New("Day parsed to be less than 1.")
	} else if outInt > 25 {
		return 0, errors.New("Day parsed to be greater than 25.")
	}

	return outInt, nil
}

// ParseDuration extracts and returns the duration from the input string
func ParseDuration(input string) (time.Duration, error) {
	// Regex for capturing "X minutes" and "X seconds" separately
	digitMinuteRegex := regexp.MustCompile(`(\d+)\s?minute[s]?`)
	wordMinuteRegex := regexp.MustCompile(`(?i)(\w+)\s*minute[s]?`)
	digitSecondRegex := regexp.MustCompile(`(\d+)\s*second[s]?`)
	wordSecondRegex := regexp.MustCompile(`(?i)(\w+)\s*second[s]?`)
	shortMinuteSecondRegex := regexp.MustCompile(`(\d+)m(\d+)s`)
	shortMinuteRegex := regexp.MustCompile(`(\d+)m`)
	shortSecondRegex := regexp.MustCompile(`(\d+)s`)

	var totalDuration time.Duration

	// Capture and sum all digit minutes
	digitMinuteMatches := digitMinuteRegex.FindAllStringSubmatch(input, -1)
	for _, match := range digitMinuteMatches {
		minutes, _ := strconv.Atoi(match[1])
		totalDuration += time.Duration(minutes) * time.Minute
	}

	// Capture and sum all word minutes
	wordMinuteMatches := wordMinuteRegex.FindAllStringSubmatch(input, -1)
	for _, match := range wordMinuteMatches {
		minutes := parseWordNumber(match[1])
		totalDuration += time.Duration(minutes) * time.Minute
	}

	// Capture and sum all digit seconds
	digitSecondMatches := digitSecondRegex.FindAllStringSubmatch(input, -1)
	for _, match := range digitSecondMatches {
		seconds, _ := strconv.Atoi(match[1])
		totalDuration += time.Duration(seconds) * time.Second
	}

	// Capture and sum all word seconds
	wordSecondMatches := wordSecondRegex.FindAllStringSubmatch(input, -1)
	for _, match := range wordSecondMatches {
		seconds := parseWordNumber(match[1])
		totalDuration += time.Duration(seconds) * time.Second
	}

	// Capture and sum all short format minutes and seconds together (e.g., "1m30s")
	shortMinuteSecondMatches := shortMinuteSecondRegex.FindAllStringSubmatch(input, -1)
	for _, match := range shortMinuteSecondMatches {
		minutes, _ := strconv.Atoi(match[1])
		seconds, _ := strconv.Atoi(match[2])
		totalDuration += time.Duration(minutes) * time.Minute
		totalDuration += time.Duration(seconds) * time.Second
	}

	// Capture and sum all short format minutes
	shortMinuteMatches := shortMinuteRegex.FindAllStringSubmatch(input, -1)
	for _, match := range shortMinuteMatches {
		minutes, _ := strconv.Atoi(match[1])
		totalDuration += time.Duration(minutes) * time.Minute
	}

	// Capture and sum all short format seconds
	shortSecondMatches := shortSecondRegex.FindAllStringSubmatch(input, -1)
	for _, match := range shortSecondMatches {
		seconds, _ := strconv.Atoi(match[1])
		totalDuration += time.Duration(seconds) * time.Second
	}

	if totalDuration == 0 {
		return 0, fmt.Errorf("no duration found in input")
	}

	return totalDuration, nil
}

// parseWordNumber converts words like "one" or "a" to their numeric equivalents
func parseWordNumber(word string) int {
	switch strings.ToLower(word) {
	case "one", "a":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	case "ten":
		return 10
	default:
		return 0
	}
}
