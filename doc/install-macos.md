# SentinelGo – macOS Installation Guide

## Prerequisites
- macOS 11+ (Intel or Apple Silicon)
- Admin privileges

## 1. Download the Binary
Download the latest macOS binary from GitHub Releases:
```
https://github.com/habib45/SentinelGo/releases/latest
```
Choose the file named `sentinelgo-darwin-amd64` (Intel) or `sentinelgo-darwin-arm64` (Apple Silicon).

## 2. Create a Directory
```bash
sudo mkdir -p /opt/sentinelgo
```

## 3. Copy the Binary
```bash
sudo cp sentinelgo-darwin-* /opt/sentinelgo/sentinelgo
sudo chmod +x /opt/sentinelgo/sentinelgo
```

## 4. (Optional) Create a Configuration File
```bash
sudo mkdir -p /etc/sentinelgo
sudo tee /etc/sentinelgo/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": "5m",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v0.1.0"
}
EOF
```

## 5. Create a launchd Agent
```bash
sudo tee /Library/LaunchDaemons/com.sentinelgo.agent.plist > /dev/null <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.sentinelgo.agent</string>
  <key>ProgramArguments</key>
  <array>
    <string>/opt/sentinelgo/sentinelgo</string>
    <string>-run</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>/var/log/sentinelgo.log</string>
  <key>StandardErrorPath</key>
  <string>/var/log/sentinelgo.err</string>
</dict>
</plist>
EOF
```

## 6. Load and Start the Service
```bash
sudo launchctl load -w /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo launchctl start com.sentinelgo.agent
```

## 7. Verify
```bash
sudo launchctl list | grep sentinelgo
```
You should see `com.sentinelgo.agent` with a PID.

## 8. Uninstall (if needed)
```bash
sudo launchctl unload -w /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo rm -f /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo rm -rf /opt/sentinelgo
sudo rm -rf /etc/sentinelgo
```

## Logs
```bash
tail -f /var/log/sentinelgo.log
tail -f /var/log/sentinelgo.err
```

## Notes
- Supabase connection is embedded at build time; users do not configure it.
- The agent checks for updates once every 24 hours.
- Heartbeat is sent every 5 minutes (configurable).
- Runs as root (launchd daemon) to collect full system metrics.
