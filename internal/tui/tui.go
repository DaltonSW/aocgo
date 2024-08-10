package tui

import (
	"fmt"
	"os"
	"strings"

	// "dalton.dog/aocgo/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

const useHighPerformanceRenderer = false

const ViewportWidth = 80

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type keymap struct {
	Up      key.Binding
	Down    key.Binding
	Browser key.Binding
	Refresh key.Binding
	Submit  key.Binding
	Input   key.Binding
	Quit    key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Input, k.Browser, k.Refresh, k.Quit}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit},
		{k.Input, k.Browser, k.Refresh},
	}
}

var keys = keymap{
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
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "[R]efresh Page"),
	),
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

type model struct {
	content   string
	ready     bool
	showHelp  bool
	keys      keymap
	viewport  viewport.Model
	help      help.Model
	viewTitle string
	userInput []byte
}

func StartViewportWithArr(input []string, title string, userInput []byte, showHelp bool) {
	log.Debug("Trying to start viewport with string array")
	contentStr := strings.Join(input, "")
	StartViewportWithString(contentStr, title, userInput, showHelp)
}

func StartViewportWithString(input string, title string, userInput []byte, showHelp bool) {
	log.Debug("Starting viewport, now creating model")
	m := model{
		content:   input,
		viewTitle: title,
		keys:      keys,
		help:      help.New(),
		showHelp:  showHelp,
		userInput: userInput,
	}
	log.Debug("Model created, now creating program")
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	log.Debug("Created program, now starting run")
	if _, err := p.Run(); err != nil {
		fmt.Println("Couldn't run viewport:", err)
		os.Exit(1)
	}
}

// TODO: Add commands for submitting and downloading input

func (m model) Init() tea.Cmd {
	log.Debug("'Init' function")
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil
	}

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	m.viewport = viewport.New(min(ViewportWidth, width), height-verticalMarginHeight)
	m.viewport.YPosition = headerHeight
	m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
	m.viewport.SetContent(m.content)
	m.ready = true
	m.viewport.YPosition = headerHeight + 1
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "b":
			// utils.LaunchURL()
		case "s":
			out, _ := os.Create("./input.txt")
			out.Write(m.userInput)
		case "a":
			// TODO: Answer question
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			log.Debug("Begin 'Update' init")
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(min(ViewportWidth, msg.Width), msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
			log.Debug("Finished 'Update' init")

		} else {
			m.viewport.Width = min(ViewportWidth, msg.Width)
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render(m.viewTitle)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	sOut := lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	if m.showHelp {
		sOut += "\n" + lipgloss.JoinHorizontal(lipgloss.Center, m.help.View(m.keys))
	}

	return sOut
}
