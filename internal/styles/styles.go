package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// TODO: Actually differentiate between light and dark. I just don't wanna delve into this right now
var (
	// Colors
	lbBorderColor  = lipgloss.AdaptiveColor{Light: "#8787FF", Dark: "#8787FF"}
	viewTitleColor = lipgloss.AdaptiveColor{Light: "#FFFF00", Dark: "#FFFF00"}
)

var (
	// Text Styles
	NormalTextStyle = lipgloss.NewStyle().Foreground(NormalTextColor)
	RedTextStyle    = lipgloss.NewStyle().Foreground(RedTextColor)
	GreenTextStyle  = lipgloss.NewStyle().Foreground(GreenTextColor)
	YellowTextStyle = lipgloss.NewStyle().Foreground(YellowTextColor)
	BlueTextStyle   = lipgloss.NewStyle().Foreground(BlueTextColor)
	PurpleTextStyle = lipgloss.NewStyle().Foreground(PurpleTextColor)
	CyanTextStyle   = lipgloss.NewStyle().Foreground(CyanTextColor)
	BrownTextStyle  = lipgloss.NewStyle().Foreground(BrownTextColor)

	SubtitleStyle = lipgloss.NewStyle().Foreground(SubtitleColor).Italic(true)

	// Styles
	viewportTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).Foreground(viewTitleColor).Underline(true)
	}()

	viewportScrollStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	ItalStyle = lipgloss.NewStyle().Foreground(ItalColor)
	StarStyle = lipgloss.NewStyle().Foreground(StarColor)
	LinkStyle = lipgloss.NewStyle().Foreground(LinkColor).Underline(true)
	CodeStyle = lipgloss.NewStyle().Foreground(CodeColor).Bold(true).Italic(true)

	LoggerFatalStyle = lipgloss.NewStyle().
				SetString("FATAL").
				Padding(0, 1).
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("#FF5F5F"))

	LoggerInfoStyle = lipgloss.NewStyle().
			SetString("INFO").
			Padding(0, 1).
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("#5FFFD7"))

	LoggerErrorStyle = lipgloss.NewStyle().
				SetString("ERROR").
				Padding(0, 1).
				Background(lipgloss.Color("204")).
				Foreground(lipgloss.Color("0"))

	GlobalSpacingStyle = lipgloss.NewStyle().Padding(1, 1, 0)

	UserTableStyle = lipgloss.NewStyle().Foreground(NormalTextColor).
			BorderForeground(TableBorderColor).Align(lipgloss.Center)
)
