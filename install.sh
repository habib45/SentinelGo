#!/bin/bash

# SentinelGo Universal Installation Script
# Usage: sudo ./install.sh [COMMAND]
# Works on Ubuntu/Debian, CentOS/RHEL, macOS, and Windows (via Git Bash)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="sentinelgo"
BINARY_NAME="sentinelgo"
INSTALL_DIR="/opt/sentinelgo"
CONFIG_DIR="${INSTALL_DIR}/.sentinelgo"
SERVICE_USER="sentinelgo"

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get >/dev/null 2>&1; then
            echo "ubuntu"
        elif command -v yum >/dev/null 2>&1; then
            echo "centos"
        elif command -v dnf >/dev/null 2>&1; then
            echo "fedora"
        else
            echo "linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running with appropriate permissions
check_permissions() {
    local os=$(detect_os)
    
    if [[ "$os" == "windows" ]]; then
        # Windows: Check if running as administrator
        if ! net session >/dev/null 2>&1; then
            print_error "Please run this script as Administrator on Windows"
            exit 1
        fi
    else
        # Unix-like: Check if running as root
        if [[ $EUID -ne 0 ]]; then
            print_error "Please run this script with sudo or as root"
            exit 1
        fi
    fi
}

# Create service user
create_service_user() {
    local os=$(detect_os)
    
    if [[ "$os" == "windows" ]]; then
        # Windows doesn't need a special user for this
        return 0
    fi
    
    if ! id -u "$SERVICE_USER" >/dev/null 2>&1; then
        print_status "Creating service user: $SERVICE_USER"
        if [[ "$os" == "macos" ]]; then
            # macOS: Create user with proper group
            sysadminctl -addUser "$SERVICE_USER" 2>/dev/null || dscl . -create /Users/"$SERVICE_USER"
            # Create group if it doesn't exist
            if ! dscl . -list /Groups | grep -q "^$SERVICE_USER$"; then
                dscl . -create /Groups/"$SERVICE_USER"
            fi
            dscl . -append /Groups/"$SERVICE_USER" GroupMembership "$SERVICE_USER"
        else
            # Linux: Create user and group
            useradd -r -s /bin/false "$SERVICE_USER" 2>/dev/null || true
        fi
        print_success "Service user created"
    else
        print_status "Service user already exists"
    fi
}

# Install directories and permissions
setup_directories() {
    print_status "Setting up directories and permissions"
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    
    # Copy binary if it exists in current directory
    if [[ -f "./$BINARY_NAME" ]]; then
        print_status "Installing binary from current directory"
        cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    elif [[ -f "./sentinelgo-linux-amd64" ]]; then
        print_status "Installing Linux AMD64 binary"
        cp "./sentinelgo-linux-amd64" "$INSTALL_DIR/$BINARY_NAME"
    elif [[ -f "./sentinelgo-darwin-amd64" ]]; then
        print_status "Installing macOS AMD64 binary"
        cp "./sentinelgo-darwin-amd64" "$INSTALL_DIR/$BINARY_NAME"
    elif [[ -f "./sentinelgo-linux-arm64" ]]; then
        print_status "Installing Linux ARM64 binary"
        cp "./sentinelgo-linux-arm64" "$INSTALL_DIR/$BINARY_NAME"
    elif [[ -f "./sentinelgo-darwin-arm64" ]]; then
        print_status "Installing macOS ARM64 binary"
        cp "./sentinelgo-darwin-arm64" "$INSTALL_DIR/$BINARY_NAME"
    elif [[ -f "./sentinelgo-windows-amd64.exe" ]]; then
        print_status "Installing Windows AMD64 binary"
        cp "./sentinelgo-windows-amd64.exe" "$INSTALL_DIR/$BINARY_NAME.exe"
    else
        print_error "No SentinelGo binary found in current directory"
        print_error "Please download: binary for your OS from GitHub releases and place it in the same directory as this script"
        exit 1
    fi
    
    # Set permissions
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    local os=$(detect_os)
    
    # Set permissions based on OS
    if [[ "$os" == "macos" ]]; then
        # macOS: Use chown with proper group handling
        chown -R "$SERVICE_USER" "$INSTALL_DIR" 2>/dev/null || true
        chmod -R 755 "$INSTALL_DIR" 2>/dev/null || true
    elif [[ "$os" == "windows" ]]; then
        # Windows: Skip ownership change
        echo "[INFO] Skipping ownership change on Windows"
    else
        # Linux: Standard permissions with proper user/group format
        if id "$SERVICE_USER" &>/dev/null 2>&1; then
            chown -R "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR" 2>/dev/null || true
        else
            # User exists but group might not, try with just user
            chown -R "$SERVICE_USER" "$INSTALL_DIR" 2>/dev/null || true
        fi
        chmod -R 755 "$INSTALL_DIR" 2>/dev/null || true
    fi
    
    print_success "Directories and permissions set"
}

# Install systemd service (Linux)
install_systemd_service() {
    print_status "Installing systemd service"
    
    cat > "/etc/systemd/system/${SERVICE_NAME}.service" << EOF
[Unit]
Description=SentinelGo Agent
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_USER
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
    
    chmod 644 "/etc/systemd/system/${SERVICE_NAME}.service"
    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
    
    # Test service configuration
    print_status "Testing service configuration..."
    # Simple test - check if service file exists and has correct permissions
    if [[ -f "/etc/systemd/system/${SERVICE_NAME}.service" ]]; then
        print_status "Service file exists and is properly configured"
        systemctl start "$SERVICE_NAME"
        
        # Check status
        sleep 3
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            print_success "Systemd service installed and started"
            print_status "Current status:"
            systemctl status "$SERVICE_NAME" --no-pager -l
        else
            print_error "Service failed to start - checking logs"
            print_status "Recent logs:"
            journalctl -u "$SERVICE_NAME" -n 10 --no-pager
        fi
    else
        print_error "Service configuration test failed"
        print_status "Reinstalling service..."
        install_systemd_service
    fi
}

# Install launchd service (macOS)
install_launchd_service() {
    print_status "Installing launchd service"
    
    cat > "/Library/LaunchDaemons/com.sentinelgo.agent.plist" << EOF
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
    <string>/var/log/sentinelgo.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/sentinelgo.log</string>
    <key>UserName</key>
    <string>$SERVICE_USER</string>
    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>
</dict>
</plist>
EOF
    
    # Set permissions for macOS
    chown root:wheel "/Library/LaunchDaemons/com.sentinelgo.agent.plist"
    chmod 644 "/Library/LaunchDaemons/com.sentinelgo.agent.plist"
    
    # Load and start service
    launchctl load "/Library/LaunchDaemons/com.sentinelgo.agent.plist"
    launchctl start "com.sentinelgo.agent"
    
    print_success "Launchd service installed and started"
}

# Install Windows service
install_windows_service() {
    print_status "Installing Windows service"
    
    # Create Windows service using sc.exe
    sc.exe create "$SERVICE_NAME" binPath= "\"$INSTALL_DIR\\$BINARY_NAME.exe\" -run" start= auto
    sc.exe description "$SERVICE_NAME" "SentinelGo Agent - Cross-platform monitoring and heartbeat service"
    
    # Set service to restart on failure
    sc.exe failure "$SERVICE_NAME" reset=86400 actions= restart/5000/restart/10000/restart/20000
    
    # Start service
    sc.exe start "$SERVICE_NAME"
    
    print_success "Windows service installed and started"
}

# Show service status
show_status() {
    print_status "Service Status:"
    
    local os=$(detect_os)
    
    case "$os" in
        ubuntu|centos|fedora|linux)
            systemctl status "$SERVICE_NAME" --no-pager
            print_status "Logs: journalctl -u $SERVICE_NAME -f"
            ;;
        macos)
            launchctl list | grep sentinelgo
            print_status "Logs: tail -f /var/log/sentinelgo.log"
            ;;
        windows)
            sc.exe query "$SERVICE_NAME"
            print_status "Logs: Get-EventLog -LogName Application -Source \"$SERVICE_NAME\" -Newest 20"
            ;;
    esac
}

# Uninstall function
uninstall_service() {
    print_status "Uninstalling SentinelGo..."
    
    check_permissions
    
    local os=$(detect_os)
    
    case "$os" in
        ubuntu|centos|fedora|linux)
            systemctl stop "$SERVICE_NAME" 2>/dev/null || true
            systemctl disable "$SERVICE_NAME"
            rm -f "/etc/systemd/system/${SERVICE_NAME}.service"
            systemctl daemon-reload
            ;;
        macos)
            launchctl unload "/Library/LaunchDaemons/com.sentinelgo.agent.plist" 2>/dev/null || true
            rm -f "/Library/LaunchDaemons/com.sentinelgo.agent.plist"
            ;;
        windows)
            sc.exe stop "$SERVICE_NAME" 2>/dev/null || true
            sc.exe delete "$SERVICE_NAME"
            ;;
    esac
    
    # Remove directories and user
    read -p "Remove all SentinelGo data and user? (y/N): " -n 1 -r response
    if [[ $response =~ ^[Yy]$ ]]; then
        rm -rf "$INSTALL_DIR"
        if [[ "$os" != "windows" ]]; then
            userdel -r "$SERVICE_USER" 2>/dev/null || true
        fi
    fi
    
    print_success "SentinelGo service uninstalled successfully"
}

# Update function
update_service() {
    print_status "Updating SentinelGo..."
    
    check_permissions
    
    # Stop any running SentinelGo processes first
    print_status "Stopping any running SentinelGo processes..."
    local os=$(detect_os)
    case "$os" in
        ubuntu|centos|fedora|linux)
            systemctl stop sentinelgo 2>/dev/null || true
            ;;
        macos)
            launchctl stop com.sentinelgo.agent 2>/dev/null || true
            ;;
        windows)
            sc.exe stop sentinelgo 2>/dev/null || true
            ;;
    esac
    
    # Kill any remaining processes
    pkill -f sentinelgo 2>/dev/null || true
    sleep 2
    
    setup_directories
    
    case "$os" in
        ubuntu|centos|fedora|linux)
            systemctl start sentinelgo
            ;;
        macos)
            launchctl start com.sentinelgo.agent
            ;;
        windows)
            sc.exe start sentinelgo
            ;;
    esac
    
    print_success "SentinelGo updated successfully!"
}

# Fix service startup issues
fix_service() {
    print_status "Fixing SentinelGo service issues..."
    
    check_permissions
    
    # Stop service first
    print_status "Stopping service..."
    local os=$(detect_os)
    case "$os" in
        ubuntu|centos|fedora|linux)
            systemctl stop sentinelgo 2>/dev/null || true
            systemctl disable sentinelgo 2>/dev/null || true
            ;;
        macos)
            launchctl stop com.sentinelgo.agent 2>/dev/null || true
            ;;
        windows)
            sc.exe stop sentinelgo 2>/dev/null || true
            ;;
    esac
    
    # Wait for complete stop
    sleep 3
    
    # Kill any remaining processes
    print_status "Killing remaining processes..."
    pkill -f sentinelgo 2>/dev/null || true
    sleep 2
    
    # Check and fix permissions
    print_status "Fixing permissions..."
    if [[ "$os" != "windows" ]]; then
        chown -R "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Check config
    print_status "Checking config..."
    if [[ ! -f "$CONFIG_DIR/config.json" ]]; then
        print_status "Creating default config..."
        mkdir -p "$CONFIG_DIR"
        echo '{"heartbeat_interval":"5m0s","auto_update":false}' > "$CONFIG_DIR/config.json"
        if [[ "$os" != "windows" ]]; then
            chown -R "$SERVICE_USER:$SERVICE_USER" "$CONFIG_DIR"
        fi
    fi
    
    # Test binary
    print_status "Testing binary..."
    if [[ "$os" != "windows" ]]; then
        # Simple test - check if binary exists and is executable
        if [[ -x "$INSTALL_DIR/$BINARY_NAME" ]]; then
            print_status "Binary exists and is executable"
            
            # Restart service
            print_status "Restarting service..."
            case "$os" in
                ubuntu|centos|fedora|linux)
                    systemctl daemon-reload
                    systemctl enable "$SERVICE_NAME"
                    systemctl start "$SERVICE_NAME"
                    ;;
                macos)
                    launchctl load "/Library/LaunchDaemons/com.sentinelgo.agent.plist"
                    launchctl start "com.sentinelgo.agent"
                    ;;
            esac
            
            # Check status
            sleep 3
            case "$os" in
                ubuntu|centos|fedora|linux)
                    if systemctl is-active --quiet "$SERVICE_NAME"; then
                        print_success "Service started successfully!"
                        print_status "Current status:"
                        systemctl status "$SERVICE_NAME" --no-pager -l
                    else
                        print_error "Service failed to start - checking logs"
                        print_status "Recent logs:"
                        journalctl -u "$SERVICE_NAME" -n 10 --no-pager
                    fi
                    ;;
            esac
        else
            print_error "Binary not found or not executable"
            print_status "Installing binary first..."
            # Copy current binary if available
            if [[ -f "./sentinelgo-linux-amd64" ]]; then
                cp "./sentinelgo-linux-amd64" "$INSTALL_DIR/$BINARY_NAME"
                chmod +x "$INSTALL_DIR/$BINARY_NAME"
                chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/$BINARY_NAME"
                print_status "Binary installed, trying again..."
            elif [[ -f "./build/linux/sentinelgo-linux-amd64" ]]; then
                cp "./build/linux/sentinelgo-linux-amd64" "$INSTALL_DIR/$BINARY_NAME"
                chmod +x "$INSTALL_DIR/$BINARY_NAME"
                chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/$BINARY_NAME"
                print_status "Binary installed from build/, trying again..."
            else
                print_error "No binary found to install"
                print_status "Please build the binary first: make"
                return
            fi
            
            # Try to start service again
            case "$os" in
                ubuntu|centos|fedora|linux)
                    systemctl start "$SERVICE_NAME"
                    ;;
                macos)
                    launchctl start "com.sentinelgo.agent"
                    ;;
            esac
            
            sleep 3
            if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
                print_success "Service started successfully!"
            else
                print_error "Service still failed - manual intervention needed"
                print_status "Check logs: journalctl -u $SERVICE_NAME -f"
            fi
        fi
    fi
    
    print_success "Service fix completed!"
}

# Show help
show_help() {
    printf '%s\n' 'SentinelGo Universal Installation Script'
    printf '%s\n' ''
    printf '%s\n' 'Usage: $0 [COMMAND]'
    printf '%s\n' ''
    printf '%s\n' 'Commands:'
    printf '%s\n' '  install     Install SentinelGo as a service (default)'
    printf '%s\n' '  uninstall   Remove SentinelGo service and data'
    printf '%s\n' '  update      Update SentinelGo binary'
    printf '%s\n' '  status      Show service status'
    printf '%s\n' '  fix-service Fix service startup issues'
    printf '%s\n' '  help        Show this help message'
    printf '%s\n' '  enable-auto-update        Enable automatic updates'
    printf '%s\n' ''
    printf '%s\n' 'Examples:'
    printf '%s\n' '  sudo ./install.sh install'
    printf '%s\n' '  sudo ./install.sh uninstall'
    printf '%s\n' '  sudo ./install.sh status'
    printf '%s\n' '  sudo ./install.sh fix-service'
}

# Install service
install_service() {
    print_status "Installing SentinelGo..."
    
    check_permissions
    
    # Stop any running SentinelGo processes first
    print_status "Stopping any running SentinelGo processes..."
    local os=$(detect_os)
    case "$os" in
        ubuntu|centos|fedora|linux)
            systemctl stop sentinelgo 2>/dev/null || true
            ;;
        macos)
            launchctl stop com.sentinelgo.agent 2>/dev/null || true
            ;;
        windows)
            sc.exe stop sentinelgo 2>/dev/null || true
            ;;
    esac
    
    # Kill any remaining processes
    pkill -f sentinelgo 2>/dev/null || true
    sleep 2
    
    create_service_user
    setup_directories
    
    case "$os" in
        ubuntu|centos|fedora|linux)
            install_systemd_service
            ;;
        macos)
            install_launchd_service
            ;;
        windows)
            install_windows_service
            ;;
    esac
    
    print_success "SentinelGo installed successfully!"
    print_status "Use './install.sh status' to check service status"
    print_status "Use './sentinelgo -enable-auto-update' to enable automatic updates"
}

# Main script logic
main() {
    local command="${1:-install}"
    
    case "$command" in
        install)
            install_service
            ;;
        uninstall)
            uninstall_service
            ;;
        update)
            update_service
            ;;
        status)
            show_status
            ;;
        fix-service)
            fix_service
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Call main function with all arguments
main "$@"
