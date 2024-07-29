package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewportModel struct {
	Width  int
	Height int
	Keymap Keymap

	MouseWheelEnabled bool
	MouseWheelDelta   int

	YOffset   int
	YPosition int

	Style           lipgloss.Style
	HighPerformance bool

	initialized bool
	lines       []string
}

func NewViewport(width, height int) viewportModel {
	newViewport := viewportModel{
		Width:             width,
		Height:            height,
		Keymap:            defaultKeymap(),
		MouseWheelEnabled: true,
		MouseWheelDelta:   3,
		initialized:       true,
	}

	return newViewport
}

func (m *viewportModel) Init() tea.Cmd {
	return nil
}

func (m *viewportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.PageDown):

		}

	case tea.MouseMsg:
		if !m.MouseWheelEnabled || msg.Action != tea.MouseActionPress {
			break
		}

		switch msg.Button {
		case tea.MouseButtonWheelDown:

		case tea.MouseButtonWheelUp:
		}
	}

	return m, nil
}

func (m *viewportModel) View() string {
	return m.View()
}

// Keymap Util

type Keymap struct {
	PageDown key.Binding
	PageUp   key.Binding
	Down     key.Binding
	Up       key.Binding
}

func defaultKeymap() Keymap {
	return Keymap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "f"),
			key.WithHelp("f/PgDn", "Page Down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("b/PgUp", "Page Up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "Down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "Up"),
		),
	}
}
