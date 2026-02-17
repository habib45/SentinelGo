#!/bin/bash

# SentinelGo macOS Reliable Installation Script
# Handles launchd issues with fallback to background process

set -e

# Configuration
SERVICE_NAME="sentinelgo"
BINARY_NAME="sentinelgo"
INSTALL_DIR="/opt/sentinelgo"
CONFIG_DIR="${INSTALL_DIR}/.sentinelgo"
SERVICE_USER="sentinelgo"

echo "🍎 SentinelGo macOS Reliable Installation"
echo "======================================"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo "❌ This script must be run as root (sudo)"
    echo "Usage: sudo ./install-macos-reliable.sh"
    exit 1
fi

# Stop any existing processes
echo "🛑 Stopping existing processes..."
sudo launchctl stop "com.sentinelgo.agent" 2>/dev/null || true
sudo launchctl unload "/Library/LaunchDaemons/com.sentinelgo.agent.plist" 2>/dev/null || true
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
    echo '{"heartbeat_interval":"5m0s","auto_update":false}' | sudo tee "$CONFIG_DIR/config.json"
    echo "✅ Default config created"
fi

# Test binary first
echo "🧪 Testing binary..."
if sudo -u "$SERVICE_USER" "$INSTALL_DIR/$BINARY_NAME" --version >/dev/null 2>&1; then
    echo "✅ Binary test passed"
else
    echo "⚠️  Binary test failed, but continuing..."
fi

# Try launchd service first
echo "🚀 Attempting launchd service..."
sudo tee /Library/LaunchDaemons/com.sentinelgo.agent.plist > /dev/null << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.sentinelgo.agent</string>
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
    <string>/tmp/sentinelgo.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/sentinelgo.log</string>
    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>
    <key>UserName</key>
    <string>$SERVICE_USER</string>
</dict>
</plist>
EOF

# Set plist permissions
sudo chown root:wheel /Library/LaunchDaemons/com.sentinelgo.agent.plist
sudo chmod 644 /Library/LaunchDaemons/com.sentinelgo.agent.plist

# Try to load launchd service
echo "🔄 Loading launchd service..."
LAUNCHD_SUCCESS=false
if sudo launchctl load /Library/LaunchDaemons/com.sentinelgo.agent.plist 2>/dev/null; then
    echo "✅ Launchd service loaded"
    
    # Try to start
    if sudo launchctl start "com.sentinelgo.agent" 2>/dev/null; then
        echo "✅ Launchd service started successfully!"
        LAUNCHD_SUCCESS=true
    else
        echo "⚠️  Launchd service failed to start"
    fi
else
    echo "⚠️  Launchd service failed to load"
fi

# Fallback to background process if launchd fails
if [[ "$LAUNCHD_SUCCESS" != "true" ]]; then
    echo "🔄 Falling back to background process mode..."
    
    # Create a simple startup script
    sudo tee "$INSTALL_DIR/start-sentinelgo.sh" > /dev/null << 'EOF'
#!/bin/bash
cd "$INSTALL_DIR"
exec "$INSTALL_DIR/$BINARY_NAME" -run >> /tmp/sentinelgo.log 2>&1
EOF
    
    sudo chmod +x "$INSTALL_DIR/start-sentinelgo.sh"
    sudo chown "$SERVICE_USER" "$INSTALL_DIR/start-sentinelgo.sh"
    
    # Start in background
    echo "🚀 Starting SentinelGo in background..."
    if sudo -u "$SERVICE_USER" nohup "$INSTALL_DIR/start-sentinelgo.sh" >/dev/null 2>&1 & then
        echo "✅ SentinelGo started in background mode"
        echo "📋 Process info:"
        ps aux | grep sentinelgo | grep -v grep
    else
        echo "❌ Failed to start background process"
    fi
fi

# Show status
echo ""
echo "📊 Current Status:"
if sudo launchctl list | grep -q "com.sentinelgo.agent"; then
    echo "✅ Launchd service: $(sudo launchctl list | grep sentinelgo)"
else
    echo "ℹ️  Launchd service: Not loaded"
fi

echo ""
echo "🔍 Running processes:"
ps aux | grep sentinelgo | grep -v grep || echo "No processes found"

echo ""
echo "📋 Logs:"
echo "tail -f /tmp/sentinelgo.log"

echo ""
echo "🎉 Installation complete!"
echo ""
echo "📖 Management Commands:"
echo "  Stop:  sudo pkill -f sentinelgo"
echo "  Start: sudo -u sentinelgo nohup $INSTALL_DIR/$BINARY_NAME -run &"
echo "  Status: ps aux | grep sentinelgo"
echo "  Logs:  tail -f /tmp/sentinelgo.log"
