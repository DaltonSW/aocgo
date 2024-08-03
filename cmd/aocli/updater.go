package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
)

const currentVersion = "v0.0.0"
const repoURL = "https://api.github.com/repos/DaltonSW/aocGo/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func checkForUpdate() bool {
	latestVersion, err := getLatestRelease()
	if err != nil {
		log.Error("Error checking for updates!", "error", err)
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

// Performs the update. Downloads the file and replaces the executable in place
func update() {
	release, err := getLatestRelease()
	if err != nil {
		log.Error("Error checking for updates.", "error", err)
	}

	var assetURL string
	for _, asset := range release.Assets {
		if asset.Name == fmt.Sprintf("aocli-%v-%v", runtime.GOOS, runtime.GOARCH) {
			assetURL = asset.DownloadURL
		}
	}

	if assetURL == "" {
		log.Error("Error obtaining a valid download URL.", "error", err)
	}

	resp, err := http.Get(assetURL)
	if err != nil {
		log.Error("Error downloading the new version.", "error", err)
	}
	defer resp.Body.Close()

	curExec, err := os.Executable()
	if err != nil {
		log.Error("Error obtaining current executable info.", "error", err)
	}

	tmpFile, err := os.CreateTemp("", "aocli-update-")
	if err != nil {
		log.Error("Error creating temp file.", "error", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the downloaded content to the temporary file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}

	// Close the file to flush the content
	if err := tmpFile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return
	}

	// Replace the current executable with the new one
	if err := os.Rename(tmpFile.Name(), curExec); err != nil {
		fmt.Println("Error replacing the executable:", err)
		return
	}

	fmt.Println("Updated successfully to version", release.TagName)

}
