package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func GetTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func GetCurrentMaxYearAndDay() (int, int) {
	loc, err := time.LoadLocation("America/New_York")

	if err != nil {
		log.Fatal("Error loading location:", err)
	}

	nowYear, nowMonth, nowDay := time.Now().In(loc).Date()
	if nowMonth != time.December {
		return nowYear - 1, 25
	} else {
		if nowDay > 24 {
			return nowYear, 25
		} else {
			return nowYear, nowDay
		}
	}
}

func LaunchURL(url string) error {
	log.Printf("Attempting to launch URL %v\n", url)
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}

	if len(args) > 1 {
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	data, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func GetResouceBucketID(year, day int) string {
	return strconv.Itoa(year) + strconv.Itoa(day)
}
