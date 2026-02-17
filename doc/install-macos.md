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
sudo cp sentinelgo-darwin-amd64 /opt/sentinelgo/sentinelgo
OR 
sudo cp sentinelgo-darwin-armd64 /opt/sentinelgo/sentinelgo

sudo chmod +x /opt/sentinelgo/sentinelgo
sudo chown -R $(whoami) /opt/sentinelgo

```

## 4. Configuration Options

### Option 1: Default Configuration (Recommended)
By default, SentinelGo will create a config file at `/opt/sentinelgo/.sentinelgo/config.json` with these settings:
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
# Config will be created automatically at /opt/sentinelgo/.sentinelgo/config.json
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

#### Method C: System-Wide Config
```bash
# Create system-wide config
sudo mkdir -p /etc/sentinelgo
sudo tee /etc/sentinelgo/config.json > /dev/null <<'EOF'
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

## 5. Install as Service (Recommended)
The SentinelGo binary now includes automated macOS service management:

```bash
# Navigate to the binary directory
cd /opt/sentinelgo

# Install as launchd service (requires sudo)
sudo ./sentinelgo -install
sudo chown -R $(whoami) ~/.sentinelgo
```

This will automatically:
- Create the launchd plist file at `/Library/LaunchDaemons/com.sentinelgo.agent.plist`
- Load and start the service
- Configure the service to start automatically on system boot
- Set up logging to `/var/log/sentinelgo.log` and `/var/log/sentinelgo.err`

### Important Notes MAC Alert “sentinelgo” Not Opened

**Allow from System Settings (recommended)**
1. Try to run the command again so the warning appears.
2. Open System Settings → Privacy & Security
3. Scroll down to the bottom.
4. You’ll see:“sentinelgo was blocked because it is not from an identified developer”
5. Click Allow Anyway
6. Run the command again.
7. This time click Open


## 6. Verify Installation
```bash
# Check service status
./sentinelgo -status

# Check launchd service specifically
sudo launchctl list | grep sentinelgo
```

## 7. Troubleshooting

### Config Loading Error
If you get this error:
```
Failed to load config: json: cannot unmarshal string into Go struct field Config.heartbeat_interval of type time.Duration
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
sudo launchctl list | grep sentinelgo

# Check logs
tail -f /var/log/sentinelgo.log

# Reinstall service
sudo ./sentinelgo -uninstall
sudo ./sentinelgo -install
```

### Auto-Start Not Working
```bash
# Use the updated install script
sudo ./install-macos-simple.sh

# Or manually fix
sudo launchctl unload /Library/LaunchDaemons/com.sentinelgo.plist
sudo launchctl load /Library/LaunchDaemons/com.sentinelgo.plist
sudo launchctl start com.sentinelgo
```

## 8. Service Management Commands

### Check Status
```bash
./sentinelgo -status
```
Shows all running SentinelGo processes and launchd service status.

### Stop All Processes
```bash
./sentinelgo -stop
```
Stops all running SentinelGo processes safely.

### Run in Foreground 
```bash
./sentinelgo -run
```
Runs the agent in console mode for debugging.

### Uninstall Service
```bash
sudo ./sentinelgo -uninstall
```
Stops and removes the launchd service completely.

## 8. Manual Installation (Alternative)
If you prefer manual setup instead of the automated installation:

```bash
# Create launchd plist manually
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
  <key>WorkingDirectory</key>
  <string>/opt/sentinelgo</string>
</dict>
</plist>
EOF

# Load and start the service
sudo launchctl load -w /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo launchctl start com.sentinelgo.agent
```

## 9. Manual Uninstall (Alternative)
```bash
sudo launchctl unload -w /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo rm -f /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo rm -rf /opt/sentinelgo
sudo rm -rf /etc/sentinelgo
```

## Logs
```bash
# View application logs
tail -f /var/log/sentinelgo.log

# View error logs
tail -f /var/log/sentinelgo.err

# View launchd service logs
log show --predicate 'process == "sentinelgo"' --last 1h
```

## Version Management
The SentinelGo agent includes automatic version detection and management:

```bash
# Check which versions are running
./sentinelgo -status

# Before updating, stop old versions
./sentinelgo -stop

# Install new version
sudo ./sentinelgo -install
```

The system will warn you if multiple versions are detected and help you clean up old instances.

## Troubleshooting

### Service Not Starting
```bash
# Check launchd service status
sudo launchctl list com.sentinelgo.agent

# Check for errors in logs
tail -f /var/log/sentinelgo.err

# Check plist file permissions
ls -la /Library/LaunchDaemons/com.sentinelgo.agent.plist
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

# Update plist if binary is in different location
sudo nano /Library/LaunchDaemons/com.sentinelgo.agent.plist
```

## Notes
- The agent runs as a system daemon (root) to collect full system metrics
- Heartbeat is sent every 5 minutes (configurable)
- Updates are checked once every 24 hours
- Service automatically starts on system boot
- All service management is handled through the SentinelGo binary commands
- Cross-platform compatible: same commands work on Linux, Windows, and macOS
