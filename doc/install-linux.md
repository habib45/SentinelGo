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
EOF
sudo chmod 600 /opt/sentinelgo/.env
```

## 5. (Optional) Create a Configuration File
You can override defaults like heartbeat interval. Only include GitHub fields if you use a different repo:
```bash
sudo mkdir -p /etc/sentinelgo
sudo tee /etc/sentinelgo/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": "5m"
}
EOF
```

If you use a custom GitHub repo, also set:
```json
{
  "heartbeat_interval": "5m",
  "github_owner": "habib45",
  "github_repo": "SentinelGo"
}
```

## 6. Create a systemd Service
```bash
sudo tee /etc/systemd/system/sentinelgo.service > /dev/null <<'EOF'
[Unit]
Description=SentinelGo Agent
After=network.target

[Service]
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
