package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	// Version can be injected at build time via ldflags
	Version = "dev"
)

type Config struct {
	Path              string `json:"-"` // Path to the config file
	HeartbeatInterval string `json:"heartbeat_interval"`
	GitHubOwner       string `json:"github_owner"`
	GitHubRepo        string `json:"github_repo"`
	CurrentVersion    string `json:"current_version"`
	DeviceID          string `json:"device_id"`   // persistent unique identifier
	AutoUpdate        bool   `json:"auto_update"` // Enable automatic updates
}

// GetHeartbeatInterval returns the heartbeat interval as time.Duration
func (c *Config) GetHeartbeatInterval() time.Duration {
	duration, _ := time.ParseDuration(c.HeartbeatInterval)
	return duration
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Path:              path,
		HeartbeatInterval: "5m0s",
		GitHubOwner:       "habib45",
		GitHubRepo:        "SentinelGo",
		CurrentVersion:    Version, // Use injected version
		AutoUpdate:        false,   // Disabled by default for safety
	}

	if path == "" {
		// Default location
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback to /opt/sentinelgo for service environment
			home = "/opt/sentinelgo"
		}
		configDir := filepath.Join(home, ".sentinelgo")
		// Ensure config directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}
		cfg.Path = filepath.Join(configDir, "config.json")
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
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp if random fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
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
