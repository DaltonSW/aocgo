package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"go.dalton.dog/aocgo/internal/styles"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
)

// Internally tracked version to compare against GitHub releases
const currentVersion = "v0.9.11"

// const currentVersion = "v0.0.0"

const repoURL = "https://api.github.com/repos/DaltonSW/aocGo/releases/latest"

const updateMessage = "New version available: %v\nRun `aocli update` to get the new version (or `sudo aocli update` if your executable is in a protected location)"

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Version will print the current version of the program.
// It will also check for any updates available.
// Associated command: `version`
func Version() {
	fmt.Println(styles.GlobalSpacingStyle.Render(styles.NormalTextStyle.Render(("Current version: " + currentVersion))))
}

// CheckForUpdate will run at the end of program executions to alert
// the user if there's a program update available.
func CheckForUpdate() {
	latestVersion, err := getLatestRelease()
	if err != nil {
		log.Fatal("Error checking for updates!", "error", err)
	}

	if !strings.Contains(latestVersion.TagName, "aocli-") {
		return
	} else {
		latestVersion.TagName = strings.Replace(latestVersion.TagName, "aocli-", "", 1)
	}

	latestSemVer := semver.Canonical(latestVersion.TagName)
	currentSemVer := semver.Canonical(currentVersion)

	if semver.Compare(latestSemVer, currentSemVer) > 0 {
		fmt.Println(styles.GlobalSpacingStyle.Render(styles.NormalTextStyle.Render(fmt.Sprintf(updateMessage, latestSemVer))))
	}
}

// Gets the latest GitHub release's tag name (version number) and asset info
func getLatestRelease() (*githubRelease, error) {
	resp, err := http.Get(repoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get the latest release: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

type initMsg int
type doneMsg int
type urlMsg struct {
	assetURL string
	version  string
}

type fileMsg struct {
	curFile string
	tmpFile *os.File
}

type errMsg error

type updateModel struct {
	spinner spinner.Model
	status  string
	done    bool
	err     error

	version  string
	assetURL string
	curFile  string
	tmpFile  *os.File
}

// RunUpdateModel downloads the newest version of aocli and
// replaces the executable in place. Wrapped in a prettifier.
// Associated command: `update`
func RunUpdateModel() {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.UpdateSpinnerColor))

	model := updateModel{
		spinner: s,
		status:  "Starting up!",
		done:    false,
	}

	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

func (m updateModel) Init() tea.Cmd {
	initCmd := func() tea.Msg { return initMsg(1) }

	return tea.Batch(initCmd, m.spinner.Tick)
}

func (m updateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.status = "Querying for latest release"
		cmds = append(cmds, func() tea.Msg {
			release, err := getLatestRelease()
			if err != nil {
				return err
			}

			var assetURL string
			for _, asset := range release.Assets {
				if asset.Name == fmt.Sprintf("aocli-%v-%v", runtime.GOOS, runtime.GOARCH) {
					assetURL = asset.DownloadURL
				} else if asset.Name == fmt.Sprintf("aocli-%v-%v.exe", runtime.GOOS, runtime.GOARCH) {
					assetURL = asset.DownloadURL
				}
			}

			if assetURL == "" {
				return err
			}

			return urlMsg{assetURL: assetURL, version: release.TagName}
		})

	case urlMsg:
		m.status = "Downloading latest release to temp file"
		url := msg.assetURL
		m.version = msg.version

		cmds = append(cmds, urlCmd(url))

	case fileMsg:
		m.status = "Replacing current file with new version, then cleaning up"
		cur := msg.curFile
		tmp := msg.tmpFile

		cmds = append(cmds, fileCmd(cur, tmp))
	case doneMsg:
		m.done = true
		m.status = fmt.Sprintf("Successfully updated to version %s", m.version)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m updateModel) View() string {
	var symbol string
	var status string

	if m.err != nil {
		symbol = lipgloss.NewStyle().Foreground(styles.RedTextColor).Render(styles.FailureX)
		status = m.err.Error()
	} else if m.done {
		symbol = lipgloss.NewStyle().Foreground(styles.GreenTextColor).Render(styles.Checkmark)
		status = m.status
	} else {
		symbol = m.spinner.View()
		status = m.status
	}

	return fmt.Sprintf("\n %s %s\n", symbol, status)
}

func urlCmd(assetURL string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(assetURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		curExec, err := os.Executable()
		if err != nil {
			return err
		}

		tmpFile, err := os.CreateTemp("", "aocli-update-")
		if err != nil {
			return err
		}

		// Write the downloaded content to the temporary file
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return err
		}

		return fileMsg{curFile: curExec, tmpFile: tmpFile}
	}
}

func fileCmd(curFile string, tmpFile *os.File) tea.Cmd {
	return func() tea.Msg {
		// Close the file to flush the content
		if err := tmpFile.Close(); err != nil {
			return err
		}

		// Make the temp file executable
		if err := os.Chmod(tmpFile.Name(), 0700); err != nil {
			return err
		}

		// Replace the current executable with the new one
		if err := os.Rename(tmpFile.Name(), curFile); err != nil {
			return err
		}

		os.Remove(tmpFile.Name())
		return doneMsg(1)
	}
}
