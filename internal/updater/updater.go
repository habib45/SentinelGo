package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"sentinelgo/internal/config"
)

type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

func CheckAndApply(ctx context.Context, cfg *config.Config) error {
	latest, err := fetchLatestRelease(ctx, cfg)
	if err != nil {
		return fmt.Errorf("fetch latest release: %w", err)
	}

	if latest.TagName == cfg.CurrentVersion {
		fmt.Printf("Already up to date: %s\n", latest.TagName)
		return nil // already up-to-date
	}

	assetURL, err := selectAsset(latest, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("select asset: %w", err)
	}

	fmt.Printf("Found update: %s -> %s\n", cfg.CurrentVersion, latest.TagName)
	newPath, err := downloadAndReplace(ctx, assetURL, latest.TagName)
	if err != nil {
		return fmt.Errorf("download and replace: %w", err)
	}

	// Update config version
	cfg.CurrentVersion = latest.TagName
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Restart using new binary
	return restart(newPath)
}

func fetchLatestRelease(ctx context.Context, cfg *config.Config) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", cfg.GitHubOwner, cfg.GitHubRepo)
	fmt.Printf("Fetching release from: %s\n", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API status %d", resp.StatusCode)
	}

	var rel GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}

	fmt.Printf("Fetched release: %s with %d assets\n", rel.TagName, len(rel.Assets))
	return &rel, nil
}

func selectAsset(rel *GitHubRelease, goos, goarch string) (string, error) {
	var suffix string
	switch goos {
	case "windows":
		suffix = ".exe"
	case "linux", "darwin":
		suffix = ""
	default:
		return "", fmt.Errorf("unsupported OS %s", goos)
	}

	pattern := fmt.Sprintf("sentinelgo-%s-%s%s", goos, goarch, suffix)
	fmt.Printf("Looking for asset: %s\n", pattern)
	fmt.Printf("Available assets: %v\n", func() (names []string) {
		for _, asset := range rel.Assets {
			names = append(names, asset.Name)
		}
		return
	}())

	for _, asset := range rel.Assets {
		if asset.Name == pattern {
			fmt.Printf("Found matching asset: %s\n", asset.Name)
			return asset.URL, nil
		}
	}
	return "", fmt.Errorf("no matching asset for %s-%s", goos, goarch)
}

func downloadAndReplace(ctx context.Context, url, version string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed status %d", resp.StatusCode)
	}

	selfPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	newPath := selfPath + ".new"
	f, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}

	return newPath, nil
}

func restart(newPath string) error {
	selfPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Replace current binary with new one
	if runtime.GOOS != "windows" {
		if err := os.Rename(newPath, selfPath); err != nil {
			return err
		}
		cmd := exec.Command(selfPath)
		return cmd.Start()
	} else {
		// Windows: use batch script to replace after exit
		bat := selfPath + ".bat"
		script := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
move /Y "%s" "%s"
"%s"
del "%s"`, newPath, selfPath, selfPath, bat)
		if err := os.WriteFile(bat, []byte(script), 0644); err != nil {
			return err
		}
		cmd := exec.Command(bat)
		return cmd.Start()
	}
}
