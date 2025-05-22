package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// GitHubRelease holds the information we need from the GitHub API
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// PackageJSON to parse version from package.json
type PackageJSON struct {
	Version string `json:"version"`
}

const LATEST_RELEASE_LINK = "https://api.github.com/repos/Didiloy/spotwrap-next/releases/latest"

// GetCurrentAppVersion reads the version from package.json
func GetCurrentAppVersion() (string, error) {
	data, err := os.ReadFile("./frontend/package.json")
	if err != nil {
		data, err = os.ReadFile("package.json")
		if err != nil {
			return "", fmt.Errorf("failed to read package.json from ./frontend.package.json or package.json: %w", err)
		}
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", fmt.Errorf("failed to parse package.json: %w", err)
	}
	return pkg.Version, nil
}

// FetchLatestReleaseInfo fetches the latest release tag_name and html_url from GitHub.
func FetchLatestReleaseInfo() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", LATEST_RELEASE_LINK, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// GitHub API prefers a User-Agent
	req.Header.Set("User-Agent", "spotwrap-next-update-checker")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release from %s: %w", LATEST_RELEASE_LINK, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch latest release: status %s, body: %s", resp.Status, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release JSON (body: %s): %w", string(body), err)
	}
	// Ensure tag_name is cleaned up if it has 'v'
	release.TagName = strings.TrimPrefix(release.TagName, "v")
	release.HTMLURL = LATEST_RELEASE_LINK
	return &release, nil
}

// parseVersionString converts "1.2.3" to [1, 2, 3]
func parseVersionString(versionStr string) ([3]int, error) {
	cleanedVersion := strings.TrimPrefix(versionStr, "v") // Ensure "v" is removed
	parts := strings.Split(cleanedVersion, ".")
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("invalid version format: '%s'. Expected 'x.y.z'", versionStr)
	}

	var versionParts [3]int
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid version segment '%s' in version '%s': %w", part, versionStr, err)
		}
		versionParts[i] = num
	}
	return versionParts, nil
}

// IsNewerVersion checks if latestVersionStr is strictly newer than currentVersionStr.
func IsNewerVersion(currentVersionStr, latestVersionStr string) (bool, error) {
	currentV, err := parseVersionString(currentVersionStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse current version '%s': %w", currentVersionStr, err)
	}
	latestV, err := parseVersionString(latestVersionStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse latest version '%s': %w", latestVersionStr, err)
	}

	if latestV[0] > currentV[0] {
		return true, nil
	}
	if latestV[0] == currentV[0] && latestV[1] > currentV[1] {
		return true, nil
	}
	if latestV[0] == currentV[0] && latestV[1] == currentV[1] && latestV[2] > currentV[2] {
		return true, nil
	}
	return false, nil
}
