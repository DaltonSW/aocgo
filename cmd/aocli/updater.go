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
	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
)

// Internally tracked version to compare against GitHub releases
const currentVersion = "v0.9.4"
const repoURL = "https://api.github.com/repos/DaltonSW/aocGo/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckForUpdate will get the latest release and compare the program versions.
// Associated command: `check-update`
func CheckForUpdate() bool {
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

	logger.Info("Successfully downloaded")

	curExec, err := os.Executable()
	if err != nil {
		logger.Fatal("Error obtaining current executable info.", "error", err)
	}

	tmpFile, err := os.CreateTemp("", "aocli-update-")
	if err != nil {
		logger.Fatal("Error creating temp file.", "error", err)
	}
	defer os.Remove(tmpFile.Name())

	logger.Info("Writing downloaded content to temp file")

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

	// Replace the current executable with the new one
	if err := os.Rename(tmpFile.Name(), curExec); err != nil {
		logger.Fatal("Error replacing the executable, maybe try sudo: ", err)
		return
	}

	logger.Infof("Updated successfully to version %v", release.TagName)
}
