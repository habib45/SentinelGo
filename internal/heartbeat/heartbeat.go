package heartbeat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"sentinelgo/internal/config"
	"sentinelgo/internal/osinfo"

	"github.com/joho/godotenv"
)

var (
	// Build-time embedded via ldflags
	SupabaseURL = ""
	SupabaseKey = ""
)

func init() {
	// Try to load .env file (if it exists)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load from environment variables
	SupabaseURL = os.Getenv("SUPABASE_URL")
	SupabaseKey = os.Getenv("SUPABASE_KEY")

	// Fallback to build-time embedded values for backward compatibility
	if SupabaseURL == "" {
		SupabaseURL = "https://hlbilthxcozyqolbkkok.supabase.co"
	}
	if SupabaseKey == "" {
		SupabaseKey = ""
	}
}

type Payload struct {
	DeviceID   string             `json:"device_id"`
	Version    string             `json:"version"`
	Timestamp  time.Time          `json:"timestamp"`
	Alive      bool               `json:"alive"`
	SystemInfo *osinfo.SystemInfo `json:"system_info"`
}

func Send(ctx context.Context, cfg *config.Config, sysInfo *osinfo.SystemInfo) error {
	if SupabaseURL == "" || SupabaseKey == "" {
		return fmt.Errorf("SUPABASE_URL and SUPABASE_KEY environment variables must be set")
	}

	payload := Payload{
		DeviceID:   cfg.DeviceID,
		Version:    cfg.CurrentVersion,
		Timestamp:  time.Now(),
		Alive:      true,
		SystemInfo: sysInfo,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", SupabaseURL+"/rest/v1/heartbeat", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+SupabaseKey)
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send heartbeat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("heartbeat failed with status %d", resp.StatusCode)
	}

	return nil
}
