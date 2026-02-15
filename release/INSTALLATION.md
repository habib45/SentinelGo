# SentinelGo Installation Guide

## 🚀 Quick Installation (One Command)

After downloading the release from GitHub, simply run:

```bash
# For Linux/macOS
sudo ./install.sh

# For Windows (run as Administrator)
./install.sh
```

That's it! SentinelGo will be installed and configured to start automatically on system boot.

---

## 📦 What You Need to Download

From the GitHub releases page, download:
1. **The binary** for your operating system:
   - `sentinelgo-linux-amd64` (for Ubuntu/Debian/CentOS)
   - `sentinelgo-darwin-amd64` (for macOS)
   - `sentinelgo-windows-amd64.exe` (for Windows)
2. **The installation script**: `install.sh`

Place both files in the same directory.

---

## 🛠️ Installation Commands

The `install.sh` script supports multiple commands:

### Install (Default)
```bash
sudo ./install.sh install
```
- Installs SentinelGo as a system service
- Creates necessary directories and permissions
- Sets up auto-start on boot
- Starts the service immediately

### Update
```bash
sudo ./install.sh update
```
- Updates the binary while keeping configuration
- Restarts the service with new version

### Status Check
```bash
sudo ./install.sh status
```
- Shows current service status
- Displays log locations

### Uninstall
```bash
sudo ./install.sh uninstall
```
- Removes the service
- Optionally removes all data and user

### Help
```bash
./install.sh help
```
- Shows all available commands

---

## 🖥️ Platform-Specific Details

### Ubuntu/Debian/CentOS
- Uses **systemd** for service management
- Binary location: `/opt/sentinelgo/sentinelgo`
- Config location: `/opt/sentinelgo/.sentinelgo/`
- Service name: `sentinelgo`

#### Service Management:
```bash
sudo systemctl start sentinelgo      # Start
sudo systemctl stop sentinelgo       # Stop
sudo systemctl restart sentinelgo    # Restart
sudo systemctl status sentinelgo     # Status
sudo systemctl enable sentinelgo     # Enable on boot
sudo systemctl disable sentinelgo    # Disable on boot
journalctl -u sentinelgo -f         # View logs
```

### macOS
- Uses **launchd** for service management
- Binary location: `/opt/sentinelgo/sentinelgo`
- Config location: `/opt/sentinelgo/.sentinelgo/`
- Service name: `com.sentinelgo.agent`

#### Service Management:
```bash
sudo launchctl load /Library/LaunchDaemons/com.sentinelgo.agent.plist    # Load
sudo launchctl unload /Library/LaunchDaemons/com.sentinelgo.agent.plist  # Unload
sudo launchctl start com.sentinelgo.agent                               # Start
sudo launchctl stop com.sentinelgo.agent                                # Stop
launchctl list | grep sentinelgo                                        # Status
tail -f /var/log/sentinelgo.log                                         # Logs
```

### Windows
- Uses **Windows Service** for service management
- Binary location: `C:\opt\sentinelgo\sentinelgo.exe`
- Config location: `C:\opt\sentinelgo\.sentinelgo\`
- Service name: `sentinelgo`

#### Service Management (Run as Administrator):
```cmd
sc.exe start sentinelgo          # Start
sc.exe stop sentinelgo           # Stop
sc.exe query sentinelgo          # Status
sc.exe delete sentinelgo         # Delete service
```

---

## 🔧 Troubleshooting

### Service Not Starting
1. Check permissions: Ensure running as root/administrator
2. Check logs: Use platform-specific log commands above
3. Verify binary: Ensure correct binary for your OS

### Permission Errors
- **Linux/macOS**: Always use `sudo`
- **Windows**: Right-click and "Run as Administrator"

### Binary Not Found
- Ensure the binary is in the same directory as `install.sh`
- Verify you downloaded the correct binary for your OS

### Port Conflicts
- SentinelGo uses default ports for communication
- Check if another service is using the same ports

---

## 📁 File Locations After Installation

### Linux/macOS
```
/opt/sentinelgo/
├── sentinelgo              # Main binary
└── .sentinelgo/
    └── config.json         # Configuration file
## 🗑️ Complete Removal

To completely remove SentinelGo:
```bash
sudo ./install.sh uninstall
# Answer 'y' when asked to remove all data
```

This removes:
- Service configuration
- Binary files
- Configuration files
- Service user account
- Log files

---

## 🆘 Support

If you encounter issues:

1. Check service status: `sudo ./install.sh status`
2. View logs for your platform
3. Ensure correct binary for your OS
4. Verify running with proper permissions

For additional support, check the GitHub repository or contact your system administrator.
