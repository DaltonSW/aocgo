package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
)

const currentVersion = "v0.0.0"
const repoURL = "https://api.github.com/repos/DaltonSW/aocGo/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func checkForUpdate() bool {
	latestVersion, err := getLatestRelease()
	if err != nil {
		log.Error("Error checking for updates!", "error", err)
	}

	latestSemVer := semver.Canonical(latestVersion)
	currentSemVer := semver.Canonical(currentVersion)

	fmt.Printf("Current version: %v -- Latest version: %v\n", currentSemVer, latestSemVer)

	return semver.Compare(latestSemVer, currentSemVer) > 0
}

func getLatestRelease() (string, error) {
	resp, err := http.Get(repoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get the latest release: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", nil
	}

	return release.TagName, nil
}
