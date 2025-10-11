package resources

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"go.dalton.dog/aocgo/internal/output"
)

type LeaderboardModel struct {
	left, right string
	viewport    viewport.Model
	title       string
	ready       bool
}

type ViewableLB interface {
	GetTitle() string
	GetContent() string
}

func NewLeaderboardViewport(leftTable, rightTable, title string) {
	m := LeaderboardModel{
		left:  leftTable,
		right: rightTable,
		title: title,
		ready: false,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Couldn't run viewport:", err)
		os.Exit(1)
	}
}

func (m LeaderboardModel) Init() tea.Cmd {
	output.Debug("'Init' function")

	return nil
}

func (m LeaderboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())

		width := msg.Width
		height := max(msg.Height-headerHeight, 0)

		if !m.ready {
			m.viewport = viewport.New(
				viewport.WithWidth(width),
				viewport.WithHeight(height),
			)
			m.ready = true
			m.viewport.YPosition = headerHeight
		} else {
			m.viewport.SetWidth(width)
			m.viewport.SetHeight(height)
			m.viewport.YPosition = headerHeight
		}
		var content string
		if m.viewport.Width() < lipgloss.Width(m.left)+lipgloss.Width(m.right) {
			content = lipgloss.JoinVertical(lipgloss.Left, m.left, m.right)
		} else {
			content = lipgloss.JoinHorizontal(lipgloss.Top, m.left, m.right)
		}
		m.viewport.SetContent(content)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m LeaderboardModel) View() string {
	return fmt.Sprintf("%s\n%s", m.headerView(), m.viewport.View())
}

func (m LeaderboardModel) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}
