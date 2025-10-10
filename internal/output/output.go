package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	infoLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightCyan).Bold(true)
	infoTextStyle  = lipgloss.NewStyle()

	successLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Bold(true)
	successTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Green)

	warnLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow).Bold(true)
	warnTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Yellow)

	errorLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true)
	errorTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Red)

	debugLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightMagenta).Bold(true)
	debugTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Magenta).Faint(true)
)

var debugEnabled = os.Getenv("AOCGO_DEBUG") != ""

// EnableDebug toggles debug output on or off.
func EnableDebug(enabled bool) {
	debugEnabled = enabled
}

// Info prints a styled informational message.
func Info(msg any, fields ...any) {
	printMessage(os.Stdout, infoLabelStyle, infoTextStyle, "INFO", msg, fields...)
}

// Infof prints a formatted informational message.
func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

// Success prints a styled success message.
func Success(msg any, fields ...any) {
	printMessage(os.Stdout, successLabelStyle, successTextStyle, " OK ", msg, fields...)
}

// Successf prints a formatted success message.
func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}

// Warn prints a styled warning message.
func Warn(msg any, fields ...any) {
	printMessage(os.Stderr, warnLabelStyle, warnTextStyle, "WARN", msg, fields...)
}

// Warnf prints a formatted warning message.
func Warnf(format string, args ...any) {
	Warn(fmt.Sprintf(format, args...))
}

// Error prints a styled error message.
func Error(msg any, fields ...any) {
	printMessage(os.Stderr, errorLabelStyle, errorTextStyle, "ERROR", msg, fields...)
}

// Errorf prints a formatted error message.
func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

// Debug prints a styled debug message when debug output is enabled.
func Debug(msg string, fields ...any) {
	if !debugEnabled {
		return
	}
	printMessage(os.Stdout, debugLabelStyle, debugTextStyle, "DEBUG", msg, fields...)
}

// Debugf prints a formatted debug message when debug output is enabled.
func Debugf(format string, args ...any) {
	Debug(fmt.Sprintf(format, args...))
}

// Fatal prints a styled fatal message and exits the process.
func Fatal(msg any, fields ...any) {
	printMessage(os.Stderr, errorLabelStyle, errorTextStyle, "FATAL", msg, fields...)
	os.Exit(1)
}

// Fatalf prints a formatted fatal message and exits the process.
func Fatalf(format string, args ...any) {
	Fatal(fmt.Sprintf(format, args...))
}

func printMessage(w io.Writer, labelStyle, textStyle lipgloss.Style, label string, msg any, fields ...any) {
	message := strings.TrimSpace(fmt.Sprint(msg))
	fieldStr := formatFields(fields...)

	switch {
	case message != "" && fieldStr != "":
		message = fmt.Sprintf("%s (%s)", message, fieldStr)
	case message == "" && fieldStr != "":
		message = fieldStr
	}

	labelText := labelStyle.PaddingLeft(1).Render(fmt.Sprintf("[%s]", label))
	if message == "" {
		fmt.Fprintln(w, labelText)
		return
	}

	fmt.Fprintln(w, labelText+" "+textStyle.Render(message))
}

func formatFields(fields ...any) string {
	if len(fields) == 0 {
		return ""
	}

	if len(fields) == 1 {
		return fmt.Sprint(fields[0])
	}

	pairs := make([]string, 0, len(fields)/2+1)
	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprint(fields[i])
		if i+1 < len(fields) {
			pairs = append(pairs, fmt.Sprintf("%s=%v", key, fields[i+1]))
		} else {
			pairs = append(pairs, key)
		}
	}

	return strings.Join(pairs, " ")
}
