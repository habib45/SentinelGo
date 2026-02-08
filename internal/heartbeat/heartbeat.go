package heartbeat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sentinelgo/internal/config"
	"sentinelgo/internal/osinfo"
)

var (
	// Build-time embedded via ldflags
	SupabaseURL = "https://hlbilthxcozyqolbkkok.supabase.co"
	SupabaseKey = ""
)

type Payload struct {
	DeviceID   string             `json:"device_id"`
	Version    string             `json:"version"`
	Timestamp  time.Time          `json:"timestamp"`
	Alive      bool               `json:"alive"`
	SystemInfo *osinfo.SystemInfo `json:"system_info"`
}

func Send(ctx context.Context, cfg *config.Config, sysInfo *osinfo.SystemInfo) error {
	if SupabaseURL == "" || SupabaseKey == "" || SupabaseURL == "https://placeholder.supabase.co" || SupabaseKey == "placeholder-key" {
		return fmt.Errorf("supabase URL/key not embedded at build time")
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
