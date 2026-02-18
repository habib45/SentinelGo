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
sudo mkdir -p /opt/sentinelgo/.sentinelgo
```

## 3. Copy the Binary
```bash
sudo cp sentinelgo-darwin-amd64 /opt/sentinelgo/sentinelgo
OR 
sudo cp sentinelgo-darwin-armd64 /opt/sentinelgo/sentinelgo

sudo chmod +x /opt/sentinelgo/sentinelgo
sudo chown -R $(whoami) /opt/sentinelgo

```

## 4. Configuration Options

### Option 1: Default Configuration (Recommended)
By default, SentinelGo will create a config file at `/opt/sentinelgo/.sentinelgo/config.json` with these settings:
```bash
sudo tee /opt/sentinelgo/.sentinelgo/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": 5,
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.0.0",
  "auto_update": false
}
EOF
```

```bash

sudo tee /opt/sentinelgo/.sentinelgo/config.json > /dev/null <<'EOF'
{
  "heartbeat_interval": "5m0s",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.9.9.0",
  "auto_update": false
}
EOF
```

```bash
# Navigate to the binary directory
cd /opt/sentinelgo

# Install as launchd service (requires sudo)
sudo ./sentinelgo -install
sudo chown -R $(whoami) ~/.sentinelgo
```

```bash
# Config will be created automatically at /opt/sentinelgo/.sentinelgo/config.json
sudo /opt/sentinelgo/sentinelgo -run
```

## 6. Verify Installation
```bash
# Check service install
./sentinelgo -install

# Check service status
./sentinelgo -status

# Check service version
./sentinelgo -version

# Check service run
./sentinelgo -run

# Check service stop
./sentinelgo -stop

# Check service uninstall
./sentinelgo -uninstall

# Check service status
sudo launchctl list | grep sentinelgo

#start application
sudo launchctl start com.sentinelgo

#stop application
sudo launchctl stop com.sentinelgo

# Check launchd service specifically
sudo launchctl list | grep sentinelgo

# Check launchd service status
sudo launchctl list com.sentinelgo.agent

# Check logs
tail -f /var/log/sentinelgo.log

# View error logs
tail -f /var/log/sentinelgo.err

# View launchd service logs
log show --predicate 'process == "sentinelgo"' --last 1h



# Reinstall service
sudo ./sentinelgo -uninstall
sudo ./sentinelgo -install
```

### Permission Issues
Ensure the binary has proper permissions:
```bash
sudo chmod +x /opt/sentinelgo/sentinelgo
sudo chown root:wheel /opt/sentinelgo/sentinelgo
```

### Binary Not Found
Make sure the binary path in the plist matches your installation:
```bash
# Verify binary exists
ls -la /opt/sentinelgo/sentinelgo

## Notes
- The agent runs as a system daemon (root) to collect full system metrics
- Heartbeat is sent every 5 minutes (configurable)
- Updates are checked once every 24 hours
- Service automatically starts on system boot
- All service management is handled through the SentinelGo binary commands
- Cross-platform compatible: same commands work on Linux, Windows, and macOS
