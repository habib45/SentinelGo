#!/bin/bash

# SentinelGo Linux Simple Installation Script
# Focused on reliability and simplicity

set -e

# Check for help command
if [[ "$1" == "help" ]] || [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    echo "SentinelGo Linux Simple Installation Script"
    echo ""
    echo "Usage: sudo ./install-linux-simple.sh"
    echo ""
    echo "This script provides a simple, reliable installation for Linux:"
    echo "  - Creates service user"
    echo "  - Installs binary"
    echo "  - Sets up configuration"
    echo "  - Creates systemd service"
    echo "  - Starts the service"
    echo ""
    echo "Features:"
    echo "  ✅ Robust error handling"
    echo "  ✅ Simple user management"
    echo "  ✅ Clean systemd setup"
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

echo "🐧 SentinelGo Linux Installation"
echo "============================="

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo "❌ This script must be run as root (sudo)"
    echo "Usage: sudo ./install-linux-simple.sh"
    exit 1
fi

# Stop any existing processes
echo "🛑 Stopping existing processes..."
sudo systemctl stop "$SERVICE_NAME" 2>/dev/null || true
sudo pkill -f sentinelgo 2>/dev/null || true

# Create user if needed
echo "👤 Creating service user..."
if ! id "$SERVICE_USER" &>/dev/null; then
    sudo useradd -r -s /bin/false "$SERVICE_USER" 2>/dev/null || true
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
if [[ -f "./sentinelgo-linux-amd64" ]]; then
    sudo cp "./sentinelgo-linux-amd64" "$INSTALL_DIR/$BINARY_NAME"
    echo "✅ AMD64 binary installed"
elif [[ -f "./sentinelgo-linux-arm64" ]]; then
    sudo cp "./sentinelgo-linux-arm64" "$INSTALL_DIR/$BINARY_NAME"
    echo "✅ ARM64 binary installed"
else
    echo "❌ No Linux binary found"
    echo "Please download from: https://github.com/habib45/SentinelGo/releases"
    exit 1
fi

# Set permissions (simplified)
echo "🔐 Setting permissions..."
sudo chown -R "$SERVICE_USER" "$INSTALL_DIR" 2>/dev/null || true
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
sudo chmod -R 755 "$INSTALL_DIR" 2>/dev/null || true

# Create config
echo "📝 Creating configuration..."
if [[ ! -f "$CONFIG_DIR/config.json" ]]; then
    echo '{"heartbeat_interval":"5m0s","auto_update":false}' | sudo tee "$CONFIG_DIR/config.json"
    echo "✅ Default config created"
fi

# Create simple systemd service
echo "🚀 Creating systemd service..."
cat > /etc/systemd/system/"$SERVICE_NAME".service << EOF
[Unit]
Description=SentinelGo Agent
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/$BINARY_NAME -run
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sentinelgo
Environment=HOME=$INSTALL_DIR
Environment=XDG_CONFIG_HOME=$CONFIG_DIR

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$INSTALL_DIR
ReadWritePaths=$CONFIG_DIR

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
echo "🔄 Enabling service..."
sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"

# Start service
echo "🚀 Starting service..."
if sudo systemctl start "$SERVICE_NAME"; then
    echo "✅ SentinelGo started successfully!"
    echo ""
    echo "📊 Status:"
    sudo systemctl status "$SERVICE_NAME" --no-pager -l
    echo ""
    echo "📋 Logs:"
    echo "sudo journalctl -u $SERVICE_NAME -f"
else
    echo "❌ Failed to start service"
    echo "🔍 Checking logs:"
    sudo journalctl -u "$SERVICE_NAME" -n 10 --no-pager
fi

echo ""
echo "🎉 Installation complete!"
echo "📖 For help: sudo ./install-linux-simple.sh help"
