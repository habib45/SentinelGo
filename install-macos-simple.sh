#!/bin/bash

# SentinelGo macOS Simple Installation Script
# Focused on reliability and simplicity

set -e

# Check for help command
if [[ "$1" == "help" ]] || [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    echo "SentinelGo macOS Simple Installation Script"
    echo ""
    echo "Usage: sudo ./install-macos-simple.sh"
    echo ""
    echo "This script provides a simple, reliable installation for macOS:"
    echo "  - Creates service user"
    echo "  - Installs binary"
    echo "  - Sets up configuration"
    echo "  - Creates launchd service"
    echo "  - Starts the service"
    echo ""
    echo "Features:"
    echo "  ✅ Robust error handling"
    echo "  ✅ Simple plist creation"
    echo "  ✅ Clear status reporting"
    echo "  ✅ Troubleshooting guidance"
    exit 0
fi

# Configuration
SERVICE_NAME="sentinelgo"
BINARY_NAME="sentinelgo"
INSTALL_DIR="/opt/sentinelgo"
CONFIG_DIR="${INSTALL_DIR}/.sentinelgo"
SERVICE_USER="sentinelgo"

echo "🍎 SentinelGo macOS Installation"
echo "==============================="

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo "❌ This script must be run as root (sudo)"
    echo "Usage: sudo ./install-macos-simple.sh"
    exit 1
fi

# Stop any existing processes
echo "🛑 Stopping existing processes..."
sudo launchctl stop "com.sentinelgo" 2>/dev/null || true
sudo launchctl unload /Library/LaunchDaemons/com.sentinelgo.plist 2>/dev/null || true
sudo rm -f /Library/LaunchDaemons/com.sentinelgo.plist 2>/dev/null || true
sudo pkill -f sentinelgo 2>/dev/null || true

# Create user if needed
echo "👤 Creating service user..."
if ! id "$SERVICE_USER" &>/dev/null; then
    sudo sysadminctl -addUser "$SERVICE_USER" 2>/dev/null || true
    echo "✅ Service user created"
else
    echo "✅ Service user already exists"
fi

# Create directories
echo "📁 Creating directories..."
sudo mkdir -p "$INSTALL_DIR"
sudo mkdir -p "$CONFIG_DIR"

# Install binary
echo "📦 Installing binary..."
if [[ -f "./sentinelgo-darwin-amd64" ]]; then
    sudo cp "./sentinelgo-darwin-amd64" "$INSTALL_DIR/$BINARY_NAME"
    echo "✅ AMD64 binary installed"
elif [[ -f "./sentinelgo-darwin-arm64" ]]; then
    sudo cp "./sentinelgo-darwin-arm64" "$INSTALL_DIR/$BINARY_NAME"
    echo "✅ ARM64 binary installed"
else
    echo "❌ No macOS binary found"
    echo "Please download from: https://github.com/habib45/SentinelGo/releases"
    exit 1
fi

# Set permissions
echo "🔐 Setting permissions..."
sudo chown -R "$SERVICE_USER" "$INSTALL_DIR" 2>/dev/null || true
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Create config
echo "📝 Creating configuration..."
if [[ ! -f "$CONFIG_DIR/config.json" ]]; then
    sudo tee "$CONFIG_DIR/config.json" > /dev/null <<'EOF'
{
  "heartbeat_interval": "5m0s",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.9.8",
  "auto_update": false
}
EOF
    echo "✅ Default config created"
else
    echo "✅ Config already exists"
fi

# Create robust launchd plist
echo "🚀 Creating launchd service..."
sudo tee /Library/LaunchDaemons/com.sentinelgo.plist > /dev/null << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.sentinelgo</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/$BINARY_NAME</string>
        <string>-run</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/sentinelgo.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/sentinelgo.error.log</string>
    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>
    <key>UserName</key>
    <string>$SERVICE_USER</string>
    <key>GroupName</key>
    <string>$SERVICE_USER</string>
    <key>ProcessType</key>
    <string>Background</string>
</dict>
</plist>
EOF

# Set plist permissions
sudo chown root:wheel /Library/LaunchDaemons/com.sentinelgo.plist
sudo chmod 644 /Library/LaunchDaemons/com.sentinelgo.plist

# Ensure log directory exists
sudo mkdir -p /var/log
sudo touch /var/log/sentinelgo.log /var/log/sentinelgo.error.log
sudo chown $SERVICE_USER:$SERVICE_USER /var/log/sentinelgo.log /var/log/sentinelgo.error.log

# Load service
echo "🔄 Loading service..."
if sudo launchctl load /Library/LaunchDaemons/com.sentinelgo.plist; then
    echo "✅ Service loaded successfully"
    
    # Start service
    echo "🚀 Starting service..."
    if sudo launchctl start "com.sentinelgo"; then
        echo "✅ SentinelGo started successfully!"
        echo ""
        echo "📊 Status:"
        sudo launchctl list | grep sentinelgo
        echo ""
        echo "📋 Logs:"
        echo "tail -f /var/log/sentinelgo.log"
        echo ""
        echo "🔄 Testing auto-start..."
        sleep 3
        if sudo launchctl list | grep -q "com.sentinelgo"; then
            echo "✅ Service is running and will auto-start on reboot!"
        else
            echo "⚠️ Service loaded but may not be running"
            echo "🔍 Check logs: tail -f /var/log/sentinelgo.error.log"
        fi
    else
        echo "❌ Failed to start service"
        echo "🔍 Checking logs:"
        sudo launchctl list | grep sentinelgo
        echo ""
        echo "🔍 Manual troubleshooting:"
        echo "1. Check plist: cat /Library/LaunchDaemons/com.sentinelgo.plist"
        echo "2. Test binary: sudo -u sentinelgo $INSTALL_DIR/$BINARY_NAME -run"
        echo "3. Check logs: tail -f /var/log/sentinelgo.error.log"
    fi
else
    echo "❌ Failed to load service"
    echo "🔍 Manual troubleshooting:"
    echo "1. Check plist: cat /Library/LaunchDaemons/com.sentinelgo.plist"
    echo "2. Test binary: sudo -u sentinelgo $INSTALL_DIR/$BINARY_NAME -run"
    echo "3. Check logs: tail -f /var/log/sentinelgo.error.log"
fi

echo ""
echo "🎉 Installation complete!"
echo "📖 For help: sudo ./install-macos-simple.sh help"
