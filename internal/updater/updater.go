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
	"strconv"
	"strings"
	"time"

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

// ProcessInfo contains information about a running SentinelGo process
type ProcessInfo struct {
	PID     int
	Version string
	CmdLine string
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

	// Stop all old processes before applying update
	fmt.Println("Stopping old SentinelGo processes before update...")
	if err := stopOldProcesses(); err != nil {
		fmt.Printf("Warning: Failed to stop some old processes: %v\n", err)
	}

	// Give processes time to stop
	time.Sleep(2 * time.Second)

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

// findOldProcesses finds all running SentinelGo processes except the current one
func findOldProcesses() ([]ProcessInfo, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/fi", "imagename eq sentinelgo.exe", "/fo", "csv", "/v")
	case "linux", "darwin":
		cmd = exec.Command("ps", "aux")
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseProcessOutput(string(output)), nil
}

// parseProcessOutput parses the output of process listing commands
func parseProcessOutput(output string) []ProcessInfo {
	var processes []ProcessInfo
	lines := strings.Split(output, "\n")

	currentPID := os.Getpid()

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var info ProcessInfo

		switch runtime.GOOS {
		case "windows":
			if strings.Contains(line, "sentinelgo.exe") {
				fields := strings.Split(line, ",")
				if len(fields) >= 5 {
					pid, _ := strconv.Atoi(strings.Trim(fields[1], `"`))
					if pid != currentPID { // Skip current process
						info.PID = pid
						info.CmdLine = strings.Trim(fields[8], `"`)
						info.Version = extractVersionFromCmd(info.CmdLine)
						processes = append(processes, info)
					}
				}
			}
		case "linux", "darwin":
			if strings.Contains(line, "sentinelgo") && !strings.Contains(line, "grep") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pid, _ := strconv.Atoi(fields[1])
					if pid != currentPID { // Skip current process
						info.PID = pid
						info.CmdLine = strings.Join(fields[10:], " ")
						info.Version = extractVersionFromCmd(info.CmdLine)
						processes = append(processes, info)
					}
				}
			}
		}
	}

	return processes
}

// extractVersionFromCmd tries to extract version from command line arguments
func extractVersionFromCmd(cmdLine string) string {
	if strings.Contains(cmdLine, "-version=") {
		parts := strings.Split(cmdLine, "-version=")
		if len(parts) > 1 {
			version := strings.Split(parts[1], " ")[0]
			return strings.Trim(version, `"`)
		}
	}
	return "unknown"
}

// stopOldProcesses stops all running SentinelGo processes except the current one
func stopOldProcesses() error {
	processes, err := findOldProcesses()
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		fmt.Println("No old SentinelGo processes found")
		return nil
	}

	fmt.Printf("Found %d old SentinelGo process(es) to stop:\n", len(processes))
	for _, proc := range processes {
		fmt.Printf("  PID: %d, Version: %s\n", proc.PID, proc.Version)
	}

	fmt.Println("Stopping old processes...")
	for _, proc := range processes {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(proc.PID))
		case "linux", "darwin":
			cmd = exec.Command("kill", "-TERM", strconv.Itoa(proc.PID))
		}

		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to stop PID %d: %v\n", proc.PID, err)
		} else {
			fmt.Printf("Stopped PID %d\n", proc.PID)
		}
	}

	// Wait a moment for processes to stop
	time.Sleep(1 * time.Second)

	// Force kill any remaining processes
	processes, _ = findOldProcesses()
	for _, proc := range processes {
		fmt.Printf("Force killing PID %d\n", proc.PID)
		switch runtime.GOOS {
		case "windows":
			exec.Command("taskkill", "/F", "/PID", strconv.Itoa(proc.PID)).Run()
		case "linux", "darwin":
			exec.Command("kill", "-KILL", strconv.Itoa(proc.PID)).Run()
		}
	}

	return nil
}

// stopLaunchdService stops the launchd service on macOS
func stopLaunchdService() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	// Check if launchd service is running
	cmd := exec.Command("launchctl", "list", "com.sentinelgo.agent")
	if err := cmd.Run(); err != nil {
		// Service not found or not running
		return nil
	}

	fmt.Println("Stopping launchd service...")

	// Unload the service
	cmd = exec.Command("launchctl", "unload", "-w", "/Library/LaunchDaemons/com.sentinelgo.agent.plist")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Failed to unload launchd service: %v\n", err)
	}

	return nil
}

// startLaunchdService starts the launchd service on macOS
func startLaunchdService() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	fmt.Println("Starting launchd service...")

	// Load and start the service
	cmd := exec.Command("launchctl", "load", "-w", "/Library/LaunchDaemons/com.sentinelgo.agent.plist")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load launchd service: %w", err)
	}

	cmd = exec.Command("launchctl", "start", "com.sentinelgo.agent")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start launchd service: %w", err)
	}

	return nil
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

	// For macOS, handle launchd service specially
	if runtime.GOOS == "darwin" {
		// Stop launchd service before replacing binary
		if err := stopLaunchdService(); err != nil {
			fmt.Printf("Warning: Failed to stop launchd service: %v\n", err)
		}

		// Give service time to stop
		time.Sleep(1 * time.Second)

		// Replace current binary with new one
		if err := os.Rename(newPath, selfPath); err != nil {
			return fmt.Errorf("failed to replace binary: %w", err)
		}

		// Start launchd service with new binary
		if err := startLaunchdService(); err != nil {
			fmt.Printf("Warning: Failed to start launchd service: %v\n", err)
			// Fallback to direct execution
			cmd := exec.Command(selfPath, "-run")
			return cmd.Start()
		}

		fmt.Println("Successfully updated and restarted launchd service")
		return nil
	}

	// For Linux and Windows
	if runtime.GOOS != "windows" {
		// Replace current binary with new one
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
