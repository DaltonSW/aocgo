package resources

import (
	"fmt"
	"os"
	"strings"

	"dalton.dog/aocgo/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

type PuzzleModel struct {
	puzzle   *Puzzle
	content  string
	viewport viewport.Model
	help     help.Model
	keys     helpKeymap
	status   string
}

func NewPuzzleViewport(puzzle *Puzzle) {
	content := puzzle.GetPrettyPageData()
	contentStr := strings.Join(content, "")
	m := PuzzleModel{
		puzzle:  puzzle,
		content: contentStr,
		keys:    helpKeys,
		help:    help.New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Couldn't run viewport:", err)
		os.Exit(1)
	}
}
func (m PuzzleModel) Init() tea.Cmd {
	log.Debug("'Init' function")

	return func() tea.Msg { return initMsg(0) }
}
func (m PuzzleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case initMsg:
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return m, nil
		}

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport = viewport.New(min(ViewportWidth, width), height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = UseHighPerformanceRenderer
		m.viewport.SetContent(m.content)
		m.viewport.YPosition = headerHeight + 1

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			m.status = "Quitting!"
			return m, tea.Quit
		case "b":
			utils.LaunchURL(m.puzzle.URL)
			m.status = "Page launched in browser!"
			return m, nil

		// BUG: Refreshing isn't working quite right. Stuff has to scroll before visualizing
		case "r":
			err := m.puzzle.ReloadPuzzleData()
			if err != nil {
				log.Fatal(err)
			}
			m.content = strings.Join(m.puzzle.GetPrettyPageData(), "\n")
			m.status = "Page refreshed!"
			// Clear terminal
			fmt.Print("\033[H\033[2J")

			return m, func() tea.Msg { return initMsg(1) }
		case "s":
			out, err := os.Create("./input.txt")
			if err != nil {
				log.Fatal(err)
			}
			userInput, err := m.puzzle.GetUserInput()
			if err != nil {
				log.Fatal(err)
			}
			out.Write(userInput)
			out.Close()
			m.status = "Input saved to 'input.txt'"
			return m, nil
		case "a":
			// TODO: Answer question
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport.Width = min(ViewportWidth, msg.Width)
		m.viewport.Height = msg.Height - verticalMarginHeight

		if UseHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m PuzzleModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m PuzzleModel) headerView() string {
	title := puzzleTitleStyle.Render(m.puzzle.Title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m PuzzleModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	sOut := lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	sOut += "\n" + lipgloss.JoinHorizontal(lipgloss.Center, m.help.View(m.keys))
	if m.status != "" {
		sOut += " -- " + m.status
	}

	return sOut
}

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
