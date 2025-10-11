package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
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

	ItalStyle = lipgloss.NewStyle().Foreground(ItalColor)
	StarStyle = lipgloss.NewStyle().Foreground(StarColor)
	LinkStyle = lipgloss.NewStyle().Foreground(LinkColor).Underline(true)
	CodeStyle = lipgloss.NewStyle().Foreground(CodeColor).Bold(true).Italic(true)

	GlobalSpacingStyle = lipgloss.NewStyle().Padding(1, 1, 0)

	UserTableStyle = lipgloss.NewStyle().Foreground(NormalTextColor).
			BorderForeground(TableBorderColor).Align(lipgloss.Center)
)
