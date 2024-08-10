package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	// completeColor = lipgloss.CompleteAdaptiveColor{
	// 	Light: lipgloss.CompleteColor{TrueColor: "", ANSI256: "", ANSI: ""},
	// 	Dark:  lipgloss.CompleteColor{TrueColor: "", ANSI256: "", ANSI: ""},
	// },

	goldColor   = lipgloss.CompleteColor{TrueColor: "#D4AF37", ANSI256: "178", ANSI: "11"}
	silverColor = lipgloss.CompleteColor{TrueColor: "#C0C0C0", ANSI256: "145", ANSI: "7"}
	bronzeColor = lipgloss.CompleteColor{TrueColor: "#CD7F32", ANSI256: "94", ANSI: "3"}

	lbBorderColor  = lipgloss.CompleteColor{TrueColor: "#8787FF", ANSI256: "99", ANSI: "13"}
	viewTitleColor = lipgloss.CompleteColor{TrueColor: "#FFFF00", ANSI256: "184", ANSI: "11"}

	italColor = lipgloss.CompleteColor{TrueColor: "#FF3374", ANSI256: "197", ANSI: "13"}
	starColor = lipgloss.CompleteColor{TrueColor: "#F1FA8C", ANSI256: "228", ANSI: "11"}
	linkColor = lipgloss.CompleteColor{TrueColor: "#8BE9FD", ANSI256: "117", ANSI: "14"}
	codeColor = lipgloss.CompleteColor{TrueColor: "#FAC3D5", ANSI256: "102", ANSI: "7"}
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
)
