package heartbeat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"sentinelgo/internal/config"
	"sentinelgo/internal/osinfo"

	"github.com/joho/godotenv"
)

var (
	SupabaseURL = ""
	SupabaseKey = ""
)
var apiToken = ""

func init() {
	// Try to load .env file (if it exists)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load from environment variables
	//SupabaseURL = os.Getenv("SUPABASE_URL")
	//SupabaseKey = os.Getenv("SUPABASE_KEY")
	//apiToken = os.Getenv("API_TOKEN")
	apiToken = "Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6IkYxWmRESkNhZ3YvRGZrZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2hsYmlsdGh4Y296eXFvbGJra29rLnN1cGFiYXNlLmNvL2F1dGgvdjEiLCJzdWIiOiJiZGJhZThmYy0xOWU4LTQ5YjUtOWRjYi04YTMwYTk1MzQ0YjUiLCJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxNzcwNTUxMjQ1LCJpYXQiOjE3NzA1NDc2NDUsImVtYWlsIjoiaGFiaWIuY3NlcEBnbWFpbC5jb20iLCJwaG9uZSI6IiIsImFwcF9tZXRhZGF0YSI6eyJwcm92aWRlciI6ImVtYWlsIiwicHJvdmlkZXJzIjpbImVtYWlsIl19LCJ1c2VyX21ldGFkYXRhIjp7ImVtYWlsIjoiaGFiaWIuY3NlcEBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGhvbmVfdmVyaWZpZWQiOmZhbHNlLCJzdWIiOiJiZGJhZThmYy0xOWU4LTQ5YjUtOWRjYi04YTMwYTk1MzQ0YjUifSwicm9sZSI6ImF1dGhlbnRpY2F0ZWQiLCJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJwYXNzd29yZCIsInRpbWVzdGFtcCI6MTc3MDU0NzY0NX1dLCJzZXNzaW9uX2lkIjoiMzY3Mjk3ODUtMDEyYy00MGZjLWJlMzEtNjczNTUzMmE1ZTA5IiwiaXNfYW5vbnltb3VzIjpmYWxzZX0.DcCWFzQ1kBs3DUPHGeSeEQ9mFtSHEEI3HlyaX64zZog"

	// Fallback to build-time embedded values for backward compatibility
	if SupabaseURL == "" {
		SupabaseURL = "https://hlbilthxcozyqolbkkok.supabase.co"
	}
	// if SupabaseKey == "" {
	// 	SupabaseKey = ""
	// }
}

type Payload struct {
	DeviceID string `json:"device_id"`
	Alive    string `json:"alive"`
}

func Send(ctx context.Context, cfg *config.Config, sysInfo *osinfo.SystemInfo) error {
	// if SupabaseURL == "" || SupabaseKey == "" {
	// 	return fmt.Errorf("SUPABASE_URL and SUPABASE_KEY environment variables must be set")
	// }

	payload := Payload{
		DeviceID: cfg.DeviceID,
		Alive:    "true",
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImhsYmlsdGh4Y296eXFvbGJra29rIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ1MDcyNDIsImV4cCI6MjA4MDA4MzI0Mn0.G3WG-lKvB6K--kZLZzdeila-CG9DdhZna5jnjZS84B4")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6IkYxWmRESkNhZ3YvRGZrZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2hsYmlsdGh4Y296eXFvbGJra29rLnN1cGFiYXNlLmNvL2F1dGgvdjEiLCJzdWIiOiJiZGJhZThmYy0xOWU4LTQ5YjUtOWRjYi04YTMwYTk1MzQ0YjUiLCJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxNzcwNTYxMDM1LCJpYXQiOjE3NzA1NTc0MzUsImVtYWlsIjoiaGFiaWIuY3NlcEBnbWFpbC5jb20iLCJwaG9uZSI6IiIsImFwcF9tZXRhZGF0YSI6eyJwcm92aWRlciI6ImVtYWlsIiwicHJvdmlkZXJzIjpbImVtYWlsIl19LCJ1c2VyX21ldGFkYXRhIjp7ImVtYWlsIjoiaGFiaWIuY3NlcEBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGhvbmVfdmVyaWZpZWQiOmZhbHNlLCJzdWIiOiJiZGJhZThmYy0xOWU4LTQ5YjUtOWRjYi04YTMwYTk1MzQ0YjUifSwicm9sZSI6ImF1dGhlbnRpY2F0ZWQiLCJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJwYXNzd29yZCIsInRpbWVzdGFtcCI6MTc3MDU1NzQzNX1dLCJzZXNzaW9uX2lkIjoiN2JkYTRjNzUtZTI0My00NDhlLWJjZTAtMzVlZjI4NjdhZWIzIiwiaXNfYW5vbnltb3VzIjpmYWxzZX0.ViiZFbbdYo9qn7XHswgNRv11SS6UfE64li3XCb8kbE8")
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
