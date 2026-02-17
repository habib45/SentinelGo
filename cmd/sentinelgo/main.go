package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"sentinelgo/internal/config"
	"sentinelgo/internal/heartbeat"
	"sentinelgo/internal/lockfile"
	"sentinelgo/internal/osinfo"
	"sentinelgo/internal/updater"

	"github.com/kardianos/service"
)

var (
	// Build version injected at build time
	Version = config.Version
)

var (
	logger service.Logger
)

type program struct {
	cfg      *config.Config
	lockFile *lockfile.LockFile
}

func (p *program) Start(s service.Service) error {
	if err := logger.Info("Starting SentinelGo service"); err != nil {
		// Log error but continue - service can still start
		fmt.Printf("Warning: failed to log start message: %v\n", err)
	}

	// Acquire process lock to prevent multiple instances
	version := getCurrentVersion()
	p.lockFile = lockfile.NewLockFile(fmt.Sprintf("sentinelgo-%s", version))

	// Check for existing lock
	locked, err := p.lockFile.CheckExistingLock()
	if err != nil {
		if err := logger.Errorf("Failed to check existing lock: %v", err); err != nil {
			fmt.Printf("Warning: failed to log error: %v\n", err)
		}
		return err
	}
	if locked {
		errMsg := fmt.Sprintf("Another instance of SentinelGo v%s is already running", version)
		if err := logger.Error(errMsg); err != nil {
			fmt.Printf("Warning: failed to log error: %v\n", err)
		}
		return errors.New(errMsg)
	}

	// Try to acquire lock
	if err := p.lockFile.TryAcquire(); err != nil {
		errMsg := fmt.Sprintf("Failed to acquire process lock: %v", err)
		if err := logger.Error(errMsg); err != nil {
			fmt.Printf("Warning: failed to log error: %v\n", err)
		}
		return errors.New(errMsg)
	}

	if err := logger.Infof("Acquired process lock for SentinelGo v%s", version); err != nil {
		fmt.Printf("Warning: failed to log info: %v\n", err)
	}

	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	if err := logger.Info("Stopping SentinelGo service"); err != nil {
		fmt.Printf("Warning: failed to log stop message: %v\n", err)
	}

	// Release process lock
	if p.lockFile != nil {
		if err := p.lockFile.Release(); err != nil {
			if err := logger.Errorf("Failed to release process lock: %v", err); err != nil {
				fmt.Printf("Warning: failed to log error: %v\n", err)
			}
		} else {
			if err := logger.Info("Released process lock"); err != nil {
				fmt.Printf("Warning: failed to log info: %v\n", err)
			}
		}
	}

	return nil
}

func (p *program) run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start auto-updater in background if enabled
	if p.cfg.AutoUpdate {
		go updater.AutoUpdateChecker(ctx, p.cfg)
	}

	ticker := time.NewTicker(p.cfg.HeartbeatInterval)
	defer ticker.Stop()

	// Initial heartbeat
	if err := heartbeat.Send(ctx, p.cfg, osinfo.Collect()); err != nil {
		if err := logger.Errorf("Initial heartbeat failed: %v", err); err != nil {
			fmt.Printf("Warning: failed to log error: %v\n", err)
		}
	}

	// Daily update check (once per day)
	updateTicker := time.NewTicker(24 * time.Hour)
	defer updateTicker.Stop()

	// Run update check on start (once)
	if err := updater.CheckAndApply(ctx, p.cfg); err != nil {
		if err := logger.Errorf("Update check failed: %v", err); err != nil {
			fmt.Printf("Warning: failed to log error: %v\n", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := heartbeat.Send(ctx, p.cfg, osinfo.Collect()); err != nil {
				if err := logger.Errorf("Heartbeat failed: %v", err); err != nil {
					fmt.Printf("Warning: failed to log error: %v\n", err)
				}
			}
		case <-updateTicker.C:
			if err := updater.CheckAndApply(ctx, p.cfg); err != nil {
				if err := logger.Errorf("Update check failed: %v", err); err != nil {
					fmt.Printf("Warning: failed to log error: %v\n", err)
				}
			}
		}
	}
}

// findSentinelGoProcesses finds all running SentinelGo processes
func findSentinelGoProcesses() ([]ProcessInfo, error) {
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

// ProcessInfo contains information about a running SentinelGo process
type ProcessInfo struct {
	PID     int
	Version string
	CmdLine string
	Status  string
}

// parseProcessOutput parses the output of process listing commands
func parseProcessOutput(output string) []ProcessInfo {
	var processes []ProcessInfo
	lines := strings.Split(output, "\n")

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
					info.PID = pid
					info.CmdLine = strings.Trim(fields[8], `"`)
					info.Status = "Running"

					// Extract version from command line or use getBinaryVersion
					info.Version = extractVersionFromCmd(info.CmdLine)
					if info.Version == "unknown" {
						info.Version = getBinaryVersion(info.CmdLine)
					}
				}
			}
		case "linux", "darwin":
			if strings.Contains(line, "sentinelgo") && !strings.Contains(line, "grep") && !strings.Contains(line, "systemctl") && !strings.Contains(line, "journalctl") && !strings.Contains(line, "editor") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pid, _ := strconv.Atoi(fields[1])
					info.PID = pid
					info.CmdLine = strings.Join(fields[10:], " ")
					info.Status = "Running"

					// Extract version from command line or use getBinaryVersion
					info.Version = extractVersionFromCmd(info.CmdLine)
					if info.Version == "unknown" {
						info.Version = getBinaryVersion(info.CmdLine)
					}
				}
			}
		}

		if info.PID > 0 {
			processes = append(processes, info)
		}
	}

	return processes
}

// extractVersionFromCmd tries to extract version from command line arguments
func extractVersionFromCmd(cmdLine string) string {
	// Look for -version flag in command line
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

// stopSentinelGoProcesses stops all running SentinelGo processes
func stopSentinelGoProcesses() error {
	processes, err := findSentinelGoProcesses()
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		fmt.Println("No running SentinelGo processes found")
		return nil
	}

	fmt.Printf("Found %d SentinelGo process(es):\n", len(processes))
	for _, proc := range processes {
		fmt.Printf("  PID: %d, Version: %s, Status: %s\n", proc.PID, proc.Version, proc.Status)
	}

	fmt.Println("\nStopping processes...")
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

	return nil
}

// showSentinelGoStatus shows all running SentinelGo processes
func showSentinelGoStatus() error {
	processes, err := findSentinelGoProcesses()
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		fmt.Println("No running SentinelGo processes found")
	} else {
		fmt.Printf("Found %d running SentinelGo process(es):\n\n", len(processes))
		for i, proc := range processes {
			fmt.Printf("Process %d:\n", i+1)
			fmt.Printf("  PID:     %d\n", proc.PID)
			fmt.Printf("  Version: %s\n", proc.Version)
			fmt.Printf("  Status:  %s\n", proc.Status)
			fmt.Printf("  Command: %s\n", proc.CmdLine)
			fmt.Println()
		}
	}

	// Check for multiple versions
	if len(processes) > 0 {
		versions := make(map[string]int)
		for _, proc := range processes {
			versions[proc.Version]++
		}

		if len(versions) > 1 {
			fmt.Println("WARNING: Multiple versions are running!")
			for version, count := range versions {
				fmt.Printf("  Version %s: %d process(es)\n", version, count)
			}
			fmt.Println("Consider stopping old versions before running the new one.")
		} else {
			fmt.Println("All processes are running the same version.")
		}
	}

	// Check launchd service status on macOS
	if runtime.GOOS == "darwin" {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("macOS launchd Service Status:")
		fmt.Println(strings.Repeat("=", 50))
		if err := checkLaunchdService(); err != nil {
			fmt.Printf("Failed to check launchd service: %v\n", err)
		}
	}

	return nil
}

// macOS specific launchd service management
func createLaunchdPlist() error {
	// Get current version for the plist
	currentVersion := Version
	if currentVersion == "" {
		currentVersion = "unknown"
	}

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.sentinelgo.agent</string>
  <key>ProgramArguments</key>
  <array>
    <string>/opt/sentinelgo/sentinelgo</string>
    <string>-run</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>/var/log/sentinelgo.log</string>
  <key>StandardErrorPath</key>
  <string>/var/log/sentinelgo.err</string>
  <key>WorkingDirectory</key>
  <string>/opt/sentinelgo</string>
  <key>Comment</key>
  <string>SentinelGo Agent v%s - Cross-platform system monitoring</string>
</dict>
</plist>`, currentVersion)

	// Create directory if it doesn't exist
	if err := os.MkdirAll("/Library/LaunchDaemons", 0755); err != nil {
		return fmt.Errorf("create LaunchDaemons directory: %w", err)
	}

	// Write the plist file
	if err := os.WriteFile("/Library/LaunchDaemons/com.sentinelgo.agent.plist", []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("write plist file: %w", err)
	}

	fmt.Println("Created launchd plist: /Library/LaunchDaemons/com.sentinelgo.agent.plist")
	return nil
}

func loadLaunchdService() error {
	cmd := exec.Command("launchctl", "load", "-w", "/Library/LaunchDaemons/com.sentinelgo.agent.plist")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("load launchd service: %w", err)
	}
	fmt.Println("Loaded launchd service: com.sentinelgo.agent")
	return nil
}

func startLaunchdService() error {
	cmd := exec.Command("launchctl", "start", "com.sentinelgo.agent")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("start launchd service: %w", err)
	}
	fmt.Println("Started launchd service: com.sentinelgo.agent")
	return nil
}

func unloadLaunchdService() error {
	cmd := exec.Command("launchctl", "unload", "-w", "/Library/LaunchDaemons/com.sentinelgo.agent.plist")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unload launchd service: %w", err)
	}
	fmt.Println("Unloaded launchd service: com.sentinelgo.agent")
	return nil
}

func removeLaunchdPlist() error {
	if err := os.Remove("/Library/LaunchDaemons/com.sentinelgo.agent.plist"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove plist file: %w", err)
	}
	fmt.Println("Removed launchd plist: /Library/LaunchDaemons/com.sentinelgo.agent.plist")
	return nil
}

func checkLaunchdService() error {
	cmd := exec.Command("sh", "-c", "launchctl list | grep sentinelgo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to check launchd service: %v\n", err)
		return nil
	}
	if len(strings.TrimSpace(string(output))) == 0 {
		fmt.Println("No SentinelGo services found in launchctl")
	} else {
		fmt.Printf("Launchd service status:\n%s\n", string(output))
	}
	return nil
}

// getCurrentVersion returns the current version of the running process
func getCurrentVersion() string {
	// Try to get version from config or use build version
	cfg, err := config.Load("")
	if err == nil && cfg.CurrentVersion != "" {
		return cfg.CurrentVersion
	}
	return Version
}

func main() {
	cfgPath := flag.String("config", "", "Path to config file (optional)")
	install := flag.Bool("install", false, "Install service")
	uninstall := flag.Bool("uninstall", false, "Uninstall service")
	run := flag.Bool("run", false, "Run in foreground (console mode)")
	status := flag.Bool("status", false, "Show running SentinelGo processes and versions")
	stop := flag.Bool("stop", false, "Stop all running SentinelGo processes")
	enableAutoUpdate := flag.Bool("enable-auto-update", false, "Enable automatic updates")
	version := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("SentinelGo version: %s\n", Version)
		fmt.Printf("Build info: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	// Handle enable-auto-update flag
	if *enableAutoUpdate {
		// Enable auto-update in config
		var cfg *config.Config
		var err error
		if len(*cfgPath) == 0 {
			// Load default config
			cfg, err = config.Load("")
		} else {
			// Load specified config
			cfg, err = config.Load(*cfgPath)
		}
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// Enable auto-update
		cfg.AutoUpdate = true
		if err := cfg.Save(); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
		fmt.Println("Auto-update enabled in config")
		return
	}

	var cfg *config.Config
	var err error
	cfg, err = config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Handle status command
	if *status {
		if err := showSentinelGoStatus(); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
		return
	}

	// Handle stop command
	if *stop {
		if err := stopSentinelGoProcesses(); err != nil {
			log.Fatalf("Failed to stop processes: %v", err)
		}
		return
	}

	// Check for existing processes before starting new one
	if *run || (!*install && !*uninstall) {
		processes, err := findSentinelGoProcesses()
		if err != nil {
			log.Printf("Warning: Could not check for existing processes: %v", err)
		} else if len(processes) > 0 {
			fmt.Printf("WARNING: Found %d running SentinelGo process(es):\n", len(processes))
			for _, proc := range processes {
				fmt.Printf("  PID: %d, Version: %s\n", proc.PID, proc.Version)
			}
			fmt.Println("\nConsider running './sentinelgo -stop' to stop old versions first")
			if !*run {
				fmt.Println("Or use './sentinelgo -run' to run in foreground mode")
			}
		}
	}

	prg := &program{cfg: cfg}

	svcCfg := &service.Config{
		Name:        "SentinelGo",
		DisplayName: "SentinelGo Agent",
		Description: "Cross-platform agent to collect OS info and report heartbeat to Supabase",
		Arguments:   []string{"-config", cfg.Path},
	}

	svc, err := service.New(prg, svcCfg)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	logger, err = svc.Logger(nil)
	if err != nil {
		log.Fatalf("Failed to get service logger: %v", err)
	}

	if *install {
		if runtime.GOOS == "darwin" {
			// macOS: Use launchd
			fmt.Println("Installing SentinelGo as launchd service...")

			// Create launchd plist
			if err := createLaunchdPlist(); err != nil {
				log.Fatalf("Failed to create launchd plist: %v", err)
			}

			// Load the service
			if err := loadLaunchdService(); err != nil {
				log.Fatalf("Failed to load launchd service: %v", err)
			}

			// Start the service
			if err := startLaunchdService(); err != nil {
				log.Fatalf("Failed to start launchd service: %v", err)
			}

			fmt.Println("SentinelGo service installed and started successfully!")
			fmt.Println("Service will start automatically on system boot.")
			fmt.Println("Logs: /var/log/sentinelgo.log and /var/log/sentinelgo.err")
		} else {
			// Linux/Windows: Use kardianos/service
			if err := svc.Install(); err != nil {
				log.Fatalf("Failed to install service: %v", err)
			}
			if err := logger.Info("Service installed"); err != nil {
				log.Printf("Warning: failed to log service installation: %v", err)
			}
		}
		return
	}

	if *uninstall {
		if runtime.GOOS == "darwin" {
			// macOS: Use launchd
			fmt.Println("Uninstalling SentinelGo launchd service...")

			// Stop and unload the service
			if err := unloadLaunchdService(); err != nil {
				log.Printf("Warning: Failed to unload launchd service: %v", err)
			}

			// Remove the plist file
			if err := removeLaunchdPlist(); err != nil {
				log.Fatalf("Failed to remove launchd plist: %v", err)
			}

			fmt.Println("SentinelGo service uninstalled successfully!")
		} else {
			// Linux/Windows: Use kardianos/service
			if err := svc.Uninstall(); err != nil {
				log.Fatalf("Failed to uninstall service: %v", err)
			}
			if err := logger.Info("Service uninstalled"); err != nil {
				log.Fatalf("Failed to log service uninstallation: %v", err)
			}
		}
		return
	}

	if *run {
		// Run in console/foreground mode
		// Acquire process lock to prevent multiple instances
		version := getCurrentVersion()
		lockFile := lockfile.NewLockFile(fmt.Sprintf("sentinelgo-%s", version))

		// Check for existing lock
		locked, err := lockFile.CheckExistingLock()
		if err != nil {
			log.Printf("Warning: Failed to check existing lock: %v", err)
		} else if locked {
			log.Printf("Another instance of SentinelGo v%s is already running", version)
			fmt.Printf("Error: Another instance of SentinelGo v%s is already running\n", version)
			fmt.Println("Use './sentinelgo -stop' to stop the running instance first")
			return
		}

		// Try to acquire lock
		if err := lockFile.TryAcquire(); err != nil {
			log.Printf("Failed to acquire process lock: %v", err)
			fmt.Printf("Error: Failed to acquire process lock: %v\n", err)
			return
		}
		defer func() {
			if err := lockFile.Release(); err != nil {
				log.Printf("Warning: Failed to release lock: %v", err)
			}
		}()

		fmt.Printf("Started SentinelGo v%s in foreground mode\n", version)

		// Set lockFile in program struct for proper cleanup
		prg.lockFile = lockFile
		prg.run()
		return
	}

	// Run as a service
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := svc.Run(); err != nil {
		if err := logger.Errorf("Service failed: %v", err); err != nil {
			log.Printf("Warning: failed to log service error: %v", err)
		}
	}
}
