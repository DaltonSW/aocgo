package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Normal text colors that should be used on a non-colored background
	NormalText = lipgloss.AdaptiveColor{Light: "#202124", Dark: "#E4E4E4"}
	RedText    = lipgloss.AdaptiveColor{Light: "#C82828", Dark: "#FF5454"}
	GreenText  = lipgloss.AdaptiveColor{Light: "#718C00", Dark: "#98E024"}
	YellowText = lipgloss.AdaptiveColor{Light: "#EAB700", Dark: "#E0D561"}
	BlueText   = lipgloss.AdaptiveColor{Light: "#4171AE", Dark: "#74B2FF"}
	PurpleText = lipgloss.AdaptiveColor{Light: "#8959A8", Dark: "#AE81FF"}
	CyanText   = lipgloss.AdaptiveColor{Light: "#3E999F", Dark: "#58E1DB"}

	// Answer colors
	CorrectAnswerStyle   = lipgloss.NewStyle().Foreground(GreenText)
	IncorrectAnswerStyle = lipgloss.NewStyle().Foreground(RedText)
	NeutralAnswerStyle   = lipgloss.NewStyle().Foreground(PurpleText)
	WarningAnswerStyle   = lipgloss.NewStyle().Foreground(YellowText)

	// Leaderboard colors
	GoldColor   = lipgloss.AdaptiveColor{Light: "#D4AF37", Dark: "#D4AF37"}
	SilverColor = lipgloss.AdaptiveColor{Light: "#C0C0C0", Dark: "#C0C0C0"}
	BronzeColor = lipgloss.AdaptiveColor{Light: "#CD7F32", Dark: "#CD7F32"}

	// Table colors
	TableBorderColor = PurpleText

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
