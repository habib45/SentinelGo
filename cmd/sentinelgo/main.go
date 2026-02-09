package main

import (
	"context"
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
	"sentinelgo/internal/osinfo"
	"sentinelgo/internal/updater"

	"github.com/kardianos/service"
)

var (
	logger service.Logger
)

type program struct {
	cfg *config.Config
}

func (p *program) Start(s service.Service) error {
	logger.Info("Starting SentinelGo service")
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping SentinelGo service")
	return nil
}

func (p *program) run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(p.cfg.HeartbeatInterval)
	defer ticker.Stop()

	// Initial heartbeat
	if err := heartbeat.Send(ctx, p.cfg, osinfo.Collect()); err != nil {
		logger.Errorf("Initial heartbeat failed: %v", err)
	}

	// Daily update check (once per day)
	updateTicker := time.NewTicker(24 * time.Hour)
	defer updateTicker.Stop()

	// Run update check on start (once)
	if err := updater.CheckAndApply(ctx, p.cfg); err != nil {
		logger.Errorf("Update check failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := heartbeat.Send(ctx, p.cfg, osinfo.Collect()); err != nil {
				logger.Errorf("Heartbeat failed: %v", err)
			}
		case <-updateTicker.C:
			if err := updater.CheckAndApply(ctx, p.cfg); err != nil {
				logger.Errorf("Update check failed: %v", err)
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

					// Try to extract version from command line
					if strings.Contains(info.CmdLine, "-version=") {
						parts := strings.Split(info.CmdLine, "-version=")
						if len(parts) > 1 {
							info.Version = strings.Split(parts[1], " ")[0]
						}
					} else {
						info.Version = "unknown"
					}
				}
			}
		case "linux", "darwin":
			if strings.Contains(line, "sentinelgo") && !strings.Contains(line, "grep") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					pid, _ := strconv.Atoi(fields[1])
					info.PID = pid
					info.CmdLine = strings.Join(fields[10:], " ")
					info.Status = "Running"

					// Try to extract version from command line
					if strings.Contains(info.CmdLine, "-version=") {
						parts := strings.Split(info.CmdLine, "-version=")
						if len(parts) > 1 {
							info.Version = strings.Split(parts[1], " ")[0]
						}
					} else {
						info.Version = "unknown"
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
		return nil
	}

	fmt.Printf("Found %d running SentinelGo process(es):\n\n", len(processes))
	for i, proc := range processes {
		fmt.Printf("Process %d:\n", i+1)
		fmt.Printf("  PID:     %d\n", proc.PID)
		fmt.Printf("  Version: %s\n", proc.Version)
		fmt.Printf("  Status:  %s\n", proc.Status)
		fmt.Printf("  Command: %s\n", proc.CmdLine)
		fmt.Println()
	}

	// Check for multiple versions
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

	return nil
}

func main() {
	cfgPath := flag.String("config", "", "Path to config file (optional)")
	install := flag.Bool("install", false, "Install service")
	uninstall := flag.Bool("uninstall", false, "Uninstall service")
	run := flag.Bool("run", false, "Run in foreground (console mode)")
	status := flag.Bool("status", false, "Show running SentinelGo processes and versions")
	stop := flag.Bool("stop", false, "Stop all running SentinelGo processes")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
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
		if err := svc.Install(); err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		logger.Info("Service installed")
		return
	}

	if *uninstall {
		if err := svc.Uninstall(); err != nil {
			log.Fatalf("Failed to uninstall service: %v", err)
		}
		logger.Info("Service uninstalled")
		return
	}

	if *run {
		// Run in console/foreground mode
		prg.run()
		return
	}

	// Run as a service
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := svc.Run(); err != nil {
		logger.Errorf("Service failed: %v", err)
	}
}
