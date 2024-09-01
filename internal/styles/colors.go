package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Normal text colors that should be used on a non-colored background
	NormalTextColor = lipgloss.AdaptiveColor{Light: "#202124", Dark: "#E4E4E4"}
	RedTextColor    = lipgloss.AdaptiveColor{Light: "#C82828", Dark: "#FF5454"}
	GreenTextColor  = lipgloss.AdaptiveColor{Light: "#718C00", Dark: "#98E024"}
	YellowTextColor = lipgloss.AdaptiveColor{Light: "#EAB700", Dark: "#E0D561"}
	BlueTextColor   = lipgloss.AdaptiveColor{Light: "#4171AE", Dark: "#74B2FF"}
	PurpleTextColor = lipgloss.AdaptiveColor{Light: "#8959A8", Dark: "#AE81FF"}
	CyanTextColor   = lipgloss.AdaptiveColor{Light: "#3E999F", Dark: "#58E1DB"}
	BrownTextColor  = lipgloss.AdaptiveColor{Light: "#3D251E", Dark: "#4A2B22"}
	SubtitleColor   = lipgloss.AdaptiveColor{Light: "#D0D0D0", Dark: "#444444"}

	// Answer colors
	CorrectAnswerStyle   = lipgloss.NewStyle().Foreground(GreenTextColor)
	IncorrectAnswerStyle = lipgloss.NewStyle().Foreground(RedTextColor)
	NeutralAnswerStyle   = lipgloss.NewStyle().Foreground(PurpleTextColor)
	WarningAnswerStyle   = lipgloss.NewStyle().Foreground(YellowTextColor)

	// Leaderboard colors
	GoldColor   = lipgloss.AdaptiveColor{Light: "#D4AF37", Dark: "#D4AF37"}
	SilverColor = lipgloss.AdaptiveColor{Light: "#C0C0C0", Dark: "#C0C0C0"}
	BronzeColor = lipgloss.AdaptiveColor{Light: "#CD7F32", Dark: "#CD7F32"}

	// Table colors
	TableBorderColor = PurpleTextColor

	// Puzzle view colors
	ItalColor = lipgloss.AdaptiveColor{Light: "#FF3374", Dark: "#FF3374"}
	StarColor = lipgloss.AdaptiveColor{Light: "#F1FA8C", Dark: "#F1FA8C"}
	LinkColor = lipgloss.AdaptiveColor{Light: "#8BE9FD", Dark: "#8BE9FD"}
	CodeColor = lipgloss.AdaptiveColor{Light: "#FAC3D5", Dark: "#FAC3D5"}

	// User display colors
	BothStarsColor = lipgloss.Color("#FFFF66")
	FirstStarColor = lipgloss.Color("#9999CC")
	NoStarsColor   = lipgloss.Color("#0F0F23")

	// Misc colors
	UpdateSpinnerColor = lipgloss.Color("#FB25A0")
)
