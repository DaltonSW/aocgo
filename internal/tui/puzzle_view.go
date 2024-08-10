package tui

import (
	"fmt"
	"os"
	"strings"

	"dalton.dog/aocgo/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

type PuzzleModel struct {
	userInput []byte
	title     string
	content   string
	url       string
	viewport  viewport.Model
	help      help.Model
	keys      helpKeymap
}

func NewPuzzleViewport(content []string, title, url string, userInput []byte) {

	contentStr := strings.Join(content, "")
	m := PuzzleModel{
		content:   contentStr,
		title:     title,
		keys:      helpKeys,
		help:      help.New(),
		userInput: userInput,
		url:       url,
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
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
		m.viewport.SetContent(m.content)
		m.viewport.YPosition = headerHeight + 1

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "b":
			utils.LaunchURL(m.url)
			return m, nil
		case "s":
			out, _ := os.Create("./input.txt")
			out.Write(m.userInput)
			out.Close()
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

		if useHighPerformanceRenderer {
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
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m PuzzleModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	sOut := lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	sOut += "\n" + lipgloss.JoinHorizontal(lipgloss.Center, m.help.View(m.keys))

	return sOut
}
