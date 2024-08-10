package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = true

const ViewportWidth = 80

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).Foreground(lipgloss.Color("#FFFF00")).Underline(true)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

type initMsg int

type helpKeymap struct {
	Up      key.Binding
	Down    key.Binding
	Browser key.Binding
	// Refresh key.Binding
	// Submit  key.Binding
	Input key.Binding
	Quit  key.Binding
}

func (k helpKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Input, k.Browser, k.Quit}
}

func (k helpKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit},
		{k.Input, k.Browser},
	}
}

var helpKeys = helpKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Browser: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "[B]rowser"),
	),
	// Refresh: key.NewBinding(
	// 	key.WithKeys("r"),
	// 	key.WithHelp("r", "[R]efresh Page"),
	// ),
	// Submit: key.NewBinding(
	// 	key.WithKeys("a"),
	// 	key.WithHelp("a", "[A]nswer Puzzle"),
	// ),
	Input: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "[S]ave Input"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "[Q]uit"),
	),
}
