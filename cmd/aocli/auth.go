package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const requiresAuthAnnotation = "requiresAuth"

var authCmd = &cobra.Command{
	Use:   "auth [token]",
	Short: "Authenticate your Advent of Code session token",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		Auth(args)
	},
}

type authModel struct {
	input     textinput.Model
	message   string
	token     string
	cancelled bool
}

func newAuthModel() authModel {
	ti := textinput.New()
	ti.Prompt = "» "
	ti.Placeholder = "Paste your Advent of Code session token"
	ti.Focus()
	ti.CharLimit = 128
	ti.SetWidth(60)
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	styles := ti.Styles
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#74B2FF"))
	styles.Focused.Text = lipgloss.NewStyle()
	styles.Blurred.Prompt = styles.Focused.Prompt
	styles.Blurred.Text = styles.Focused.Text
	styles.Cursor.Color = lipgloss.Color("#F1FA8C")
	ti.Styles = styles

	return authModel{
		input: ti,
	}
}

func (m authModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			value := strings.TrimSpace(m.input.Value())
			if value == "" {
				m.message = "Token cannot be empty."
				return m, nil
			}
			m.token = value
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m authModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Enter Advent of Code Session Token")
	instructions := lipgloss.NewStyle().Faint(true).Render("Press Enter to save, Esc to cancel.")
	view := lipgloss.JoinVertical(lipgloss.Left, header, instructions, "", m.input.View())
	if m.message != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F"))
		view = lipgloss.JoinVertical(lipgloss.Left, view, "", errorStyle.Render(m.message))
	}
	return view
}

func Auth(args []string) {
	var token string
	var err error

	if len(args) > 0 {
		token = strings.TrimSpace(args[0])
	} else {
		var model tea.Model
		model, err = tea.NewProgram(newAuthModel()).Run()
		if err != nil {
			log.Fatal("Unable to capture session token", "err", err)
		}

		m, ok := model.(authModel)
		if !ok {
			log.Fatal("Unexpected auth model response")
		}

		if m.cancelled {
			log.Info("Authentication cancelled; no token stored.")
			return
		}

		token = strings.TrimSpace(m.token)
	}

	if token == "" {
		log.Fatal("Session token cannot be empty.")
	}

	tokenPath, err := writeSessionToken(token)
	if err != nil {
		log.Fatal("Unable to store session token", "err", err)
	}

	log.Info("Session token saved.", "path", tokenPath)
}

func writeSessionToken(token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", errors.New("session token cannot be empty")
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(userHomeDir, ".config", "aocgo")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return "", err
	}

	tokenPath := filepath.Join(configDir, "session.token")
	if err := os.WriteFile(tokenPath, []byte(strings.TrimSpace(token)), 0o600); err != nil {
		return "", err
	}

	return tokenPath, nil
}
