package styles

import "github.com/charmbracelet/lipgloss/v2"

func GetLeaderboardStyle(row, col int) lipgloss.Style {
	if row == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).Align(lipgloss.Center)
	}

	var style lipgloss.Style

	switch col {
	case 0:
		style = lipgloss.NewStyle().Width(5).Align(lipgloss.Center)

	case 2:
		style = lipgloss.NewStyle().Width(40)

	default:
		style = lipgloss.NewStyle().Width(17).Align(lipgloss.Center)
	}

	switch row {
	case 1:
		return style.Foreground(GoldColor)

	case 2:
		return style.Foreground(SilverColor)

	case 3:
		return style.Foreground(BronzeColor)
	default:
		return style
	}
}
