package dirparse

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	// "github.com/charmbracelet/log"
)

const ()

func GetDayAndYearFromCWD() (int, int, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, 0, err
	}
	splitCWD := strings.Split(cwd, string(os.PathSeparator))

	cwdLen := len(splitCWD)
	curDir, parentDir := splitCWD[cwdLen-1], splitCWD[cwdLen-2]

	return GetDayAndYearFromDirInput(curDir, parentDir)
}

func GetDayAndYearFromDirInput(curDir, parentDir string) (int, int, error) {
	day, err := parseDay(curDir)
	if err != nil {
		return 0, 0, err
	}
	year, err := parseYear(parentDir)
	if err != nil {
		return 0, 0, err
	}

	return day, year, nil
}

func parseYear(yearStr string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(yearStr)
	if len(match) == 2 {
		year, err := strconv.Atoi(match)
		if err != nil {
			return 0, nil
		}
		return 2000 + year, nil
	}
	return strconv.Atoi(match)
}

func parseDay(dayStr string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(dayStr)
	return strconv.Atoi(match)
}
