package types

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/OkanUysal/go-logger"
)

// GitHubTag represents GitHub API tag response
type GitHubTag struct {
	Name string `json:"name"`
}

// FetchLatestVersions fetches latest versions from GitHub for all libraries
func FetchLatestVersions() map[string]string {
	versions := make(map[string]string)
	client := &http.Client{Timeout: 5 * time.Second}

	libraries := []string{
		"go-auth",
		"go-migration",
		"go-logger",
		"go-cache",
		"go-swagger",
		"go-response",
		"go-validator",
		"go-pagination",
		"go-websocket",
		"go-metrics",
	}

	for _, lib := range libraries {
		version := fetchGitHubVersion(client, lib)
		if version != "" {
			versions[lib] = version
		}
	}

	logger.Info("Fetched versions from GitHub", logger.Int("count", len(versions)))
	return versions
}

// fetchGitHubVersion fetches the latest tag version from GitHub
func fetchGitHubVersion(client *http.Client, repoName string) string {
	// Use tags API instead of releases API
	url := fmt.Sprintf("https://api.github.com/repos/OkanUysal/%s/tags", repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Debug("Failed to create request", logger.String("repo", repoName), logger.Err(err))
		return ""
	}

	// GitHub API requires User-Agent
	req.Header.Set("User-Agent", "go-starter-api")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Debug("Failed to fetch version", logger.String("repo", repoName), logger.Err(err))
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Debug("GitHub API error", logger.String("repo", repoName), logger.Int("status", resp.StatusCode))
		return ""
	}

	var tags []GitHubTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		logger.Debug("Failed to parse response", logger.String("repo", repoName), logger.Err(err))
		return ""
	}

	// Return first tag (latest)
	if len(tags) > 0 {
		logger.Debug("Fetched version", logger.String("repo", repoName), logger.String("version", tags[0].Name))
		return tags[0].Name
	}

	logger.Debug("No tags found", logger.String("repo", repoName))
	return ""
}
