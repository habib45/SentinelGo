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
	SupabaseURL = "https://hlbilthxcozyqolbkkok.supabase.co"
	SupabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImhsYmlsdGh4Y296eXFvbGJra29rIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ1MDcyNDIsImV4cCI6MjA4MDA4MzI0Mn0.G3WG-lKvB6K--kZLZzdeila-CG9DdhZna5jnjZS84B4"
)

type Payload struct {
	DeviceID string `json:"device_id"`
	Alive    string `json:"alive"`
	BSID     string `json:"employee_id"`
	OS       string `json:"os"`
}

func init() {
	// No .env file loading - using hardcoded credentials
}

func Send(ctx context.Context, cfg *config.Config, sysInfo *osinfo.SystemInfo) error {

	payload := Payload{
		DeviceID: cfg.DeviceID,
		Alive:    "true",
		BSID:     sysInfo.EmployeeId,
		OS:       sysInfo.OS,
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
