package styles

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func GetStdoutLogger() *log.Logger {
	logger := log.New(os.Stdout)

	logger.SetReportTimestamp(true)
	logger.SetTimeFormat(time.Stamp)

	logStyles := log.DefaultStyles()

	logStyles.Levels[log.FatalLevel] = LoggerFatalStyle
	logStyles.Levels[log.ErrorLevel] = LoggerErrorStyle
	logStyles.Levels[log.InfoLevel] = LoggerInfoStyle
	logger.SetStyles(logStyles)

	return logger
}

func GetLeaderboardStyle(row, col int) lipgloss.Style {
	if row == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).Align(lipgloss.Center)
	}

	var style lipgloss.Style
	if col == 2 {
		style = lipgloss.NewStyle().Width(40)
	} else if col == 0 {
		style = lipgloss.NewStyle().Width(5).Align(lipgloss.Center)
	} else {
		style = lipgloss.NewStyle().Width(17).Align(lipgloss.Center)
	}

	if row == 1 {
		return style.Foreground(GoldColor)
	} else if row == 2 {
		return style.Foreground(SilverColor)
	} else if row == 3 {
		return style.Foreground(BronzeColor)
	} else {
		return style
	}
}
