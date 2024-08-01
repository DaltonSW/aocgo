package dirparse

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	// "github.com/charmbracelet/log"
)

const ()

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

	if outYear < 2016 {
		return 0, errors.New("Year parsed to be earlier than 2016.")
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
