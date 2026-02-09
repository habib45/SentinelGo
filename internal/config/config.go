package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

var (
	// Version can be injected at build time via ldflags
	Version = "dev"
)

type Config struct {
	Path              string        `json:"-"` // Path to the config file
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	GitHubOwner       string        `json:"github_owner"`
	GitHubRepo        string        `json:"github_repo"`
	CurrentVersion    string        `json:"current_version"`
	DeviceID          string        `json:"device_id"` // persistent unique identifier
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Path:              path,
		HeartbeatInterval: 5 * time.Minute,
		GitHubOwner:       "habib45",
		GitHubRepo:        "SentinelGo",
		CurrentVersion:    Version, // Use injected version
	}

	if path == "" {
		// Default location
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfg.Path = filepath.Join(home, ".sentinelgo", "config.json")
	}

	if _, err := os.Stat(cfg.Path); err == nil {
		data, err := os.ReadFile(cfg.Path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Ensure DeviceID exists
	if cfg.DeviceID == "" {
		cfg.DeviceID = generateDeviceID()
		if err := cfg.Save(); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func generateDeviceID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *Config) Save() error {
	dir := filepath.Dir(c.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.Path, data, 0644)
}
