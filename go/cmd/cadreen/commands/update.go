package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateRepo      = "timothy-billingrails/cadreen-sdks"
	maxDownloadSize = 100 * 1024 * 1024
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Cadreen to the latest version",
	Long: `Check for a new version of Cadreen and update if available.

Downloads the correct binary for your OS and architecture
from GitHub Releases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		current := Version
		if current == "dev" {
			fmt.Println("You're running a development build.")
			fmt.Println("Update is only available for release builds.")
			return nil
		}

		fmt.Printf("Current version: %s\n", current)
		fmt.Println("Checking for updates...")

		latest, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("checking for updates: %w", err)
		}

		if latest.TagName == "cli-v"+current || latest.TagName == current {
			fmt.Println("You're already on the latest version.")
			return nil
		}

		fmt.Printf("New version available: %s\n", latest.TagName)
		fmt.Println("Downloading...")

		asset := findAsset(latest.Assets)
		if asset == nil {
			return fmt.Errorf("no binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
		}

		if err := downloadAndReplace(asset.BrowserDownloadURL); err != nil {
			return fmt.Errorf("downloading update: %w", err)
		}

		fmt.Printf("Updated to %s successfully.\n", latest.TagName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Name    string        `json:"name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int    `json:"size"`
}

func getLatestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", updateRepo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func findAsset(assets []githubAsset) *githubAsset {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	suffix := ""
	if goos == "windows" {
		suffix = ".exe"
	}

	patterns := []string{
		fmt.Sprintf("cadreen_%s_%s%s", goos, goarch, suffix),
	}

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		for _, p := range patterns {
			if strings.HasSuffix(name, strings.ToLower(p)) {
				return &asset
			}
		}
	}

	return nil
}

func downloadAndReplace(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable: %w", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("resolving symlink: %w", err)
	}

	tmpPath := execPath + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	limitedReader := io.LimitReader(resp.Body, int64(maxDownloadSize)+1)
	n, err := io.Copy(tmpFile, limitedReader)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("writing binary: %w", err)
	}

	if n == 0 {
		os.Remove(tmpPath)
		return fmt.Errorf("downloaded file is empty")
	}

	if n > int64(maxDownloadSize) {
		os.Remove(tmpPath)
		return fmt.Errorf("downloaded file exceeds maximum size (%dMB)", maxDownloadSize/1024/1024)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpPath, 0755); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("setting permissions: %w", err)
		}
	}

	if err := os.Rename(tmpPath, execPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("replacing binary: %w", err)
	}

	return nil
}
