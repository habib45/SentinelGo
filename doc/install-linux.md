# SentinelGo – Linux Installation Guide

## Prerequisites
- Linux (systemd-based distro: Ubuntu 18.04+, Debian 10+, CentOS 8+, RHEL 8+, etc.)
- `sudo` or root access

## 1. Download the Binary
Download the latest Linux binary from GitHub Releases:
```
https://github.com/habib45/SentinelGo/releases/latest
```
Choose the file named `sentinelgo-linux-amd64` (or `-arm64` for ARM).

## 2. Create a Directory
```bash
sudo mkdir -p /opt/sentinelgo
```

## 3. Copy the Binary
```bash
sudo cp sentinelgo-linux-amd64 /opt/sentinelgo/sentinelgo
sudo chmod +x /opt/sentinelgo/sentinelgo
```

## 4. Set Environment Variables
Create a `.env` file with your Supabase credentials:
```bash
sudo tee /opt/sentinelgo/.env > /dev/null <<'EOF'
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
API_TOKEN=your-api-token
EOF
sudo chmod 600 /opt/sentinelgo/.env
```

## 5. Configuration Options

### Option 1: Default Configuration (Recommended)
By default, SentinelGo will create a config file at `/etc/sentinelgo/config.json` with these settings:
```json
{
  "heartbeat_interval": "5m0s",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.9.9.0",
  "auto_update": false
}
```

### Option 2: Custom Configuration
Create a custom config file at your preferred location:

#### Method A: Use Default Location
```bash
# Config will be created automatically at /etc/sentinelgo/config.json
/opt/sentinelgo/sentinelgo -run
```

#### Method B: Specify Custom Path
```bash
# Create config at custom location
mkdir -p ~/my-sentinelgo-config
tee ~/my-sentinelgo-config/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": "10m0s",
  "github_owner": "your-username",
  "github_repo": "your-repo",
  "current_version": "v1.9.9.0",
  "auto_update": true
}
EOF

# Run with custom config
/opt/sentinelgo/sentinelgo -run -config ~/my-sentinelgo-config/config.json
```

#### Method C: User-Level Config
```bash
# Create user-level config
mkdir -p ~/.sentinelgo
tee ~/.sentinelgo/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": "5m0s",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.9.9.0",
  "auto_update": false
}
EOF
```

### Option 3: Environment Variable
```bash
# Set config path via environment variable
export SENTINELGO_CONFIG="/path/to/your/config.json"
/opt/sentinelgo/sentinelgo -run
```

**Important:** The `heartbeat_interval` must be a **string** in quotes (e.g., `"5m0s"`) not a number (e.g., `300`).

## 6. Create a systemd Service
```bash
sudo tee /etc/systemd/system/sentinelgo.service > /dev/null <<'EOF'
[Unit]
Description=SentinelGo Agent
After=network.target

[Service]
WorkingDirectory=/opt/sentinelgo
Type=simple
ExecStart=/opt/sentinelgo/sentinelgo -run
Restart=on-failure
RestartSec=10
User=root
Group=root

[Install]
WantedBy=multi-user.target
EOF
```

## 7. Enable and Start the Service
```bash
sudo systemctl daemon-reload
sudo systemctl enable sentinelgo
sudo systemctl start sentinelgo
```

## 8. Verify
```bash
sudo systemctl status sentinelgo
```
Should show `active (running)`.

## 9. Troubleshooting

### Config Loading Error
If you get this error:
```
Failed to load config: json: cannot unmarshal number into Go struct field .heartbeat_interval of type string
```

**Solution:**
```bash
# Fix config file
sudo rm -f /etc/sentinelgo/config.json
sudo tee /etc/sentinelgo/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": "5m0s",
  "auto_update": false
}
EOF

# Ensure binary is updated (v1.9.9+)
/opt/sentinelgo/sentinelgo -version
```

### Service Not Starting
```bash
# Check service status
sudo systemctl status sentinelgo

# Check logs
sudo journalctl -u sentinelgo -f

# Reinstall service
sudo systemctl stop sentinelgo
sudo systemctl disable sentinelgo
sudo systemctl enable sentinelgo
sudo systemctl start sentinelgo
```

### Permission Issues
```bash
# Fix permissions
sudo chown -R sentinelgo:sentinelgo /opt/sentinelgo
sudo chmod +x /opt/sentinelgo/sentinelgo
sudo chmod -R 755 /etc/sentinelgo
```

### Auto-Start Not Working
```bash
# Enable service properly
sudo systemctl daemon-reload
sudo systemctl enable sentinelgo
sudo systemctl start sentinelgo

# Check if enabled
sudo systemctl is-enabled sentinelgo
```

## 9. Uninstall (if needed)
```bash
sudo systemctl stop sentinelgo
sudo systemctl disable sentinelgo
sudo rm -f /etc/systemd/system/sentinelgo.service
sudo systemctl daemon-reload
sudo rm -rf /opt/sentinelgo
sudo rm -rf /etc/sentinelgo
```

## Logs
```bash
sudo journalctl -u sentinelgo -f
```

## Notes
- Supabase connection is configured via environment variables in `/opt/sentinelgo/.env`.
- The agent checks for updates once every 24 hours from GitHub Releases.
- Heartbeat is sent every 5 minutes (configurable).
- Runs as root to collect full system metrics. If you prefer non-root, change `User`/`Group` in the service file and adjust permissions.
