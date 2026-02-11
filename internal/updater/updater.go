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

	// Give processes more time to fully stop
	fmt.Println("Waiting for old processes to fully terminate...")
	time.Sleep(5 * time.Second)

	// Double-check no old processes remain
	processes, _ := findOldProcesses()
	if len(processes) > 0 {
		fmt.Printf("Warning: %d old process(es) still running, proceeding anyway...\n", len(processes))
		for _, proc := range processes {
			fmt.Printf("  PID: %d, Version: %s\n", proc.PID, proc.Version)
		}
		// Force kill remaining processes
		fmt.Println("Force killing remaining old processes...")
		for _, proc := range processes {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(proc.PID))
			case "linux", "darwin":
				cmd = exec.Command("kill", "-KILL", strconv.Itoa(proc.PID))
			}
			cmd.Run()
		}
		// Wait for force kill to take effect
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("All old processes stopped successfully")
	}

	newPath, err := downloadAndReplace(ctx, assetURL, latest.TagName)
	if err != nil {
		return fmt.Errorf("download and replace: %w", err)
	}

	// Verify binary replacement was successful
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return fmt.Errorf("new binary not found after replacement: %w", err)
	}

	fmt.Printf("Successfully updated to version %s\n", latest.TagName)

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
	currentVersion := getCurrentVersion()

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
						info.Version = getProcessVersion(info.CmdLine, pid)
						// Only include if it's a different version or unknown
						if info.Version != currentVersion {
							processes = append(processes, info)
						}
					}
				}
			}
		case "linux", "darwin":
			if strings.Contains(line, "sentinelgo") && !strings.Contains(line, "grep") && !strings.Contains(line, "systemctl") && !strings.Contains(line, "journalctl") && !strings.Contains(line, "editor") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pid, _ := strconv.Atoi(fields[1])
					if pid != currentPID { // Skip current process
						info.PID = pid
						info.CmdLine = strings.Join(fields[10:], " ")
						info.Version = getProcessVersion(info.CmdLine, pid)
						// Only include if it's a different version or unknown
						if info.Version != currentVersion {
							processes = append(processes, info)
						}
					}
				}
			}
		}
	}

	return processes
}

// getProcessVersion determines the version of a running process
func getProcessVersion(cmdLine string, pid int) string {
	// Try to extract version from command line first
	if version := extractVersionFromCmd(cmdLine); version != "unknown" {
		return version
	}

	// Try to get version from binary
	if version := getBinaryVersion(cmdLine); version != "unknown" {
		return version
	}

	// Try to get version from executable path
	if version := extractVersionFromPath(cmdLine); version != "unknown" {
		return version
	}

	return "unknown"
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

	// Look for version flag as separate argument
	if strings.Contains(cmdLine, "-version") || strings.Contains(cmdLine, "--version") {
		// Try to find version after the flag
		parts := strings.Fields(cmdLine)
		for i, part := range parts {
			if (part == "-version" || part == "--version") && i+1 < len(parts) {
				return strings.Trim(parts[i+1], `"`)
			}
		}
	}

	return "unknown"
}

// getBinaryVersion tries to get version from the binary executable
func getBinaryVersion(cmdLine string) string {
	// Extract binary path from command line
	var binaryPath string
	parts := strings.Fields(cmdLine)

	if len(parts) > 0 {
		binaryPath = parts[0]
		// Handle relative paths
		if !strings.Contains(binaryPath, "/") && runtime.GOOS != "windows" {
			// Try to find binary in PATH
			if path, err := exec.LookPath(binaryPath); err == nil {
				binaryPath = path
			}
		}
	}

	// Try to get version by running binary with -version flag
	if binaryPath != "" {
		cmd := exec.Command(binaryPath, "-version")
		output, err := cmd.Output()
		if err == nil {
			outputStr := string(output)
			// Parse version output
			lines := strings.Split(outputStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "version:") || strings.Contains(line, "version") {
					// Extract version from line like "SentinelGo version: v1.0.0"
					parts := strings.Fields(line)
					for i, part := range parts {
						if strings.Contains(part, "version") && i+1 < len(parts) {
							return strings.Trim(parts[i+1], ",")
						}
					}
				}
			}
		}
	}

	return "unknown"
}

// getCurrentVersion returns the current version of the running process
func getCurrentVersion() string {
	// Try to get version from config or use build version
	cfg, err := config.Load("")
	if err == nil && cfg.CurrentVersion != "" {
		return cfg.CurrentVersion
	}
	return config.Version
}

// extractVersionFromPath extracts version from executable path
func extractVersionFromPath(path string) string {
	// Extract version from filename like "sentinelgo-v1.8.4"
	parts := strings.Split(path, "-")
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.HasPrefix(parts[i], "v") {
			return strings.Trim(parts[i], `"`)
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
	time.Sleep(3 * time.Second)

	// Check if any processes are still running
	remaining, _ := findOldProcesses()
	if len(remaining) > 0 {
		fmt.Printf("Force killing %d remaining process(es)...\n", len(remaining))
		for _, proc := range remaining {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(proc.PID))
			case "linux", "darwin":
				cmd = exec.Command("kill", "-KILL", strconv.Itoa(proc.PID))
			}
			cmd.Run()
		}
		// Wait for force kill to take effect
		time.Sleep(2 * time.Second)
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

	// Check if plist file exists
	plistPath := "/Library/LaunchDaemons/com.sentinelgo.agent.plist"
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		fmt.Printf("Launchd plist not found at %s\n", plistPath)
		return fmt.Errorf("launchd plist file not found - service may not be installed")
	}

	// Load the service
	cmd := exec.Command("launchctl", "load", "-w", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to load launchd service: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return fmt.Errorf("failed to load launchd service: %w", err)
	}

	// Wait a moment before starting
	time.Sleep(500 * time.Millisecond)

	// Start the service
	cmd = exec.Command("launchctl", "start", "com.sentinelgo.agent")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to start launchd service: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return fmt.Errorf("failed to start launchd service: %w", err)
	}

	// Wait and verify service is running
	time.Sleep(1 * time.Second)
	cmd = exec.Command("launchctl", "list", "com.sentinelgo.agent")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Warning: Could not verify launchd service status: %v\n", err)
	} else {
		if strings.Contains(string(output), "com.sentinelgo.agent") {
			fmt.Println("Launchd service started successfully")
		} else {
			fmt.Printf("Warning: Launchd service may not be running properly. Output: %s\n", string(output))
		}
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
		time.Sleep(3 * time.Second)

		// Replace current binary with new one
		if err := os.Rename(newPath, selfPath); err != nil {
			return fmt.Errorf("failed to replace binary: %w", err)
		}

		// Verify binary replacement was successful
		if _, err := os.Stat(selfPath); os.IsNotExist(err) {
			return fmt.Errorf("new binary not found after replacement: %w", err)
		}

		fmt.Printf("Successfully updated to version %s\n", extractVersionFromPath(newPath))

		// Wait before starting to ensure old processes are fully terminated
		time.Sleep(2 * time.Second)

		// Start launchd service with new binary
		if err := startLaunchdService(); err != nil {
			fmt.Printf("Warning: Failed to start launchd service: %v\n", err)
			fmt.Println("Falling back to direct execution...")
			// Fallback to direct execution
			cmd := exec.Command(selfPath, "-run")
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start fallback execution: %w", err)
			}
			fmt.Println("Started SentinelGo in direct execution mode")
			return nil
		}

		// Final verification - ensure only new version is running
		time.Sleep(3 * time.Second)
		fmt.Println("Verifying only new version is running...")

		// Check for any remaining old processes
		finalCheck, _ := findOldProcesses()
		if len(finalCheck) > 0 {
			fmt.Printf("Warning: Found %d old process(es) still running after update:\n", len(finalCheck))
			for _, proc := range finalCheck {
				fmt.Printf("  PID: %d, Version: %s\n", proc.PID, proc.Version)
			}
			fmt.Println("Force stopping remaining old processes...")
			for _, proc := range finalCheck {
				var cmd *exec.Cmd
				switch runtime.GOOS {
				case "windows":
					cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(proc.PID))
				case "linux", "darwin":
					cmd = exec.Command("kill", "-KILL", strconv.Itoa(proc.PID))
				}
				cmd.Run()
			}
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Success: Only new version is running")
		}

		return nil
	}

	// For Linux and Windows
	if runtime.GOOS != "windows" {
		// Replace current binary with new one
		if err := os.Rename(newPath, selfPath); err != nil {
			return err
		}
		// Wait before starting new process
		time.Sleep(2 * time.Second)
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
