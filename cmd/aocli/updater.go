package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"dalton.dog/aocgo/internal/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
)

// Internally tracked version to compare against GitHub releases
const currentVersion = "v0.9.5"
const repoURL = "https://api.github.com/repos/DaltonSW/aocGo/releases/latest"

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
func Version() bool {
	// TODO: Reorder the printing of this to make it "current version" focused at first

	// TODO: Some styling too
	latestVersion, err := getLatestRelease()
	if err != nil {
		log.Fatal("Error checking for updates!", "error", err)
	}

	if !strings.Contains(latestVersion.TagName, "aocli-") {
		return false
	} else {
		latestVersion.TagName = strings.Replace(latestVersion.TagName, "aocli-", "", 1)
	}

	latestSemVer := semver.Canonical(latestVersion.TagName)
	currentSemVer := semver.Canonical(currentVersion)

	fmt.Printf("Current version: %v -- Latest version: %v\n", currentSemVer, latestSemVer)

	return semver.Compare(latestSemVer, currentSemVer) > 0
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

// Update downloads the newest version of aocli and replaces the executable in place
// Associated command: `update`
func Update() {
	logger := styles.GetStdoutLogger()

	logger.Info("Querying for latest release")

	release, err := getLatestRelease()
	if err != nil {
		logger.Fatal("Error checking for updates.", "error", err)
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
		logger.Fatal("Error obtaining a valid download URL.", "error", err)
	}

	logger.Info("Attempting to download latest asset")

	resp, err := http.Get(assetURL)
	if err != nil {
		logger.Fatal("Error downloading the new version.", "error", err)
	}
	defer resp.Body.Close()

	logger.Info("Request successful")

	curExec, err := os.Executable()
	if err != nil {
		logger.Fatal("Error obtaining current executable info.", "error", err)
	}

	tmpFile, err := os.CreateTemp("", "aocli-update-")
	if err != nil {
		logger.Fatal("Error creating temp file.", "error", err)
	}
	defer os.Remove(tmpFile.Name())

	logger.Info("Downloading content to temp file")

	// Write the downloaded content to the temporary file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		logger.Fatal("Error writing to temporary file:", err)
		return
	}

	// Close the file to flush the content
	if err := tmpFile.Close(); err != nil {
		logger.Fatal("Error closing temporary file:", err)
		return
	}

	// Make the temp file executable
	if err := os.Chmod(tmpFile.Name(), 0700); err != nil {
		logger.Fatal("Error changing mode to allow execution: ", err)
		return
	}

	// Replace the current executable with the new one
	if err := os.Rename(tmpFile.Name(), curExec); err != nil {
		logger.Fatal("Error replacing the executable, maybe try sudo: ", err)
		return
	}

	logger.Infof("Updated successfully to version %v", release.TagName)
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

type errMsg struct{ err error }

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
				return errMsg{err}
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
				return errMsg{err}
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
		symbol = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("")
		status = m.err.Error()
	} else if m.done {
		symbol = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("󰸞")
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
			return errMsg{err}
		}
		defer resp.Body.Close()

		curExec, err := os.Executable()
		if err != nil {
			return errMsg{err}
		}

		tmpFile, err := os.CreateTemp("", "aocli-update-")
		if err != nil {
			return errMsg{err}
		}

		// Write the downloaded content to the temporary file
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			return errMsg{err}
		}

		return fileMsg{curFile: curExec, tmpFile: tmpFile}
	}
}

func fileCmd(curFile string, tmpFile *os.File) tea.Cmd {
	return func() tea.Msg {
		// Close the file to flush the content
		if err := tmpFile.Close(); err != nil {
			return errMsg{err}
		}

		// Make the temp file executable
		if err := os.Chmod(tmpFile.Name(), 0700); err != nil {
			return errMsg{err}
		}

		// Replace the current executable with the new one
		if err := os.Rename(tmpFile.Name(), curFile); err != nil {
			return errMsg{err}
		}

		os.Remove(tmpFile.Name())
		return doneMsg(1)
	}
}
