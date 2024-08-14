package styles

import "github.com/charmbracelet/lipgloss"

// TODO: Actually differentiate between light and dark. I just don't wanna delve into this right now
var (
	// Colors
	goldColor   = lipgloss.AdaptiveColor{Light: "#D4AF37", Dark: "#D4AF37"}
	silverColor = lipgloss.AdaptiveColor{Light: "#C0C0C0", Dark: "#C0C0C0"}
	bronzeColor = lipgloss.AdaptiveColor{Light: "#CD7F32", Dark: "#CD7F32"}

	lbBorderColor  = lipgloss.AdaptiveColor{Light: "#8787FF", Dark: "#8787FF"}
	viewTitleColor = lipgloss.AdaptiveColor{Light: "#FFFF00", Dark: "#FFFF00"}

	italColor = lipgloss.AdaptiveColor{Light: "#FF3374", Dark: "#FF3374"}
	starColor = lipgloss.AdaptiveColor{Light: "#F1FA8C", Dark: "#F1FA8C"}
	linkColor = lipgloss.AdaptiveColor{Light: "#8BE9FD", Dark: "#8BE9FD"}
	codeColor = lipgloss.AdaptiveColor{Light: "#FAC3D5", Dark: "#FAC3D5"}
)

var (
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

	italStyle = lipgloss.NewStyle().Foreground(italColor)
	starStyle = lipgloss.NewStyle().Foreground(starColor)
	linkStyle = lipgloss.NewStyle().Foreground(linkColor).Underline(true)
	codeStyle = lipgloss.NewStyle().Foreground(codeColor).Bold(true)

	HelpTextStyle = lipgloss.NewStyle().MaxWidth(70)
)

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
		return style.Foreground(goldColor)
	} else if row == 2 {
		return style.Foreground(silverColor)
	} else if row == 3 {
		return style.Foreground(bronzeColor)
	} else {
		return style
	}
}
