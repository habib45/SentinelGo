package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
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

func main() {
	cfgPath := flag.String("config", "", "Path to config file (optional)")
	install := flag.Bool("install", false, "Install service")
	uninstall := flag.Bool("uninstall", false, "Uninstall service")
	run := flag.Bool("run", false, "Run in foreground (console mode)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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
