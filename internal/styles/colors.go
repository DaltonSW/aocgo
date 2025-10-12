package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

var (
	// Normal text colors that should be used on a non-colored background
	NormalTextColor = compat.AdaptiveColor{Light: lipgloss.Color("#202124"), Dark: lipgloss.Color("#E4E4E4")}
	RedTextColor    = compat.AdaptiveColor{Light: lipgloss.Color("#C82828"), Dark: lipgloss.Color("#FF5454")}
	GreenTextColor  = compat.AdaptiveColor{Light: lipgloss.Color("#718C00"), Dark: lipgloss.Color("#98E024")}
	YellowTextColor = compat.AdaptiveColor{Light: lipgloss.Color("#EAB700"), Dark: lipgloss.Color("#E0D561")}
	BlueTextColor   = compat.AdaptiveColor{Light: lipgloss.Color("#4171AE"), Dark: lipgloss.Color("#74B2FF")}
	PurpleTextColor = compat.AdaptiveColor{Light: lipgloss.Color("#8959A8"), Dark: lipgloss.Color("#AE81FF")}
	CyanTextColor   = compat.AdaptiveColor{Light: lipgloss.Color("#3E999F"), Dark: lipgloss.Color("#58E1DB")}
	BrownTextColor  = compat.AdaptiveColor{Light: lipgloss.Color("#3D251E"), Dark: lipgloss.Color("#6A4A3A")}
	SubtitleColor   = compat.AdaptiveColor{Light: lipgloss.Color("#D0D0D0"), Dark: lipgloss.Color("#777B7E")}

	// Answer colors
	CorrectAnswerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Green)
	IncorrectAnswerStyle = lipgloss.NewStyle().Foreground(lipgloss.Red)
	NeutralAnswerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Magenta)
	WarningAnswerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Yellow)

	// Leaderboard colors
	GoldColor   = compat.AdaptiveColor{Light: lipgloss.Color("#D4AF37"), Dark: lipgloss.Color("#D4AF37")}
	SilverColor = compat.AdaptiveColor{Light: lipgloss.Color("#C0C0C0"), Dark: lipgloss.Color("#C0C0C0")}
	BronzeColor = compat.AdaptiveColor{Light: lipgloss.Color("#CD7F32"), Dark: lipgloss.Color("#CD7F32")}

	// Table colors
	TableHeaderColor = lipgloss.Magenta
	TableBorderColor = lipgloss.BrightRed

	// User display colors
	BothStarsColor = lipgloss.BrightYellow
	FirstStarColor = lipgloss.Cyan
	NoStarsColor   = lipgloss.Black

	// Misc colors
	UpdateSpinnerColor = lipgloss.BrightMagenta
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

	ItalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Italic(true).Bold(true)
	StarStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow)
	LinkStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightCyan).Underline(true)
	CodeStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Bold(true).Italic(true)

	GlobalSpacingStyle = lipgloss.NewStyle().Margin(1).AlignHorizontal(lipgloss.Center)

	UserTableStyle = lipgloss.NewStyle().Foreground(NormalTextColor).
			BorderForeground(TableHeaderColor).Align(lipgloss.Center)
)
