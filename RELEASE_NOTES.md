# SentinelGo v1.9.5 Release Notes

## 🚀 **Complete Enterprise Deployment Solution**

This release provides a comprehensive, production-ready SentinelGo deployment system with automatic updates, cross-platform support, and enterprise-grade installation tools.

---

## 📦 **What's Included**

### **Universal Installation Script**
- ✅ **Cross-platform** - Works on Ubuntu, macOS, and Windows
- ✅ **Auto-detection** - Automatically detects OS and architecture
- ✅ **Service integration** - Installs as proper system service
- ✅ **One-command install** - Simple for non-technical users

### **Automatic Updates**
- ✅ **GitHub integration** - Fetches latest releases automatically
- ✅ **Safe updates** - Stops old processes before updating
- ✅ **Version control** - Only updates if newer version available
- ✅ **Process isolation** - Prevents conflicts during updates

### **Duplicate Prevention**
- ✅ **Lock mechanism** - Only one instance per version can run
- ✅ **Process detection** - Properly identifies old vs new versions
- ✅ **Clean transitions** - No duplicate heartbeats

---

## 🎯 **Key Features**

### **1. Easy Installation**
```bash
# Download the package for your OS/architecture
# Extract and run:
sudo ./install.sh
```

### **2. Automatic Updates**
```bash
# Enable auto-updates:
sudo ./sentinelgo -enable-auto-update
```

### **3. Cross-Platform Support**
- ✅ **Linux** (AMD64 + ARM64) - Ubuntu, Debian, CentOS, RHEL
- ✅ **macOS** (AMD64 + ARM64) - Intel and Apple Silicon Macs
- ✅ **Windows** (AMD64) - Windows 10/11, Server

---

## 📋 **Download Options**

### **Linux**
- `sentinelgo-v1.9.5-3-g004af9b-dirty-linux-amd64.tar.gz` (5.8MB)
- `sentinelgo-v1.9.5-3-g004af9b-dirty-linux-arm64.tar.gz` (5.4MB)

### **macOS**
- `sentinelgo-v1.9.5-3-g004af9b-dirty-darwin-amd64.tar.gz` (4.7MB)
- `sentinelgo-v1.9.5-3-g004af9b-dirty-darwin-arm64.tar.gz` (4.4MB)

### **Windows**
- `sentinelgo-v1.9.5-3-g004af9b-dirty-windows.tar.gz` (4.8MB)

---

## 🛠️ **Installation Instructions**

### **Quick Install (Recommended)**
1. Download the package for your OS/architecture
2. Extract the archive
3. Run: `sudo ./install.sh`
4. Enable auto-updates: `sudo ./sentinelgo -enable-auto-update`

### **Manual Install**
1. Download the binary for your platform
2. Place in `/opt/sentinelgo/` (Linux/macOS) or `C:\opt\sentinelgo\` (Windows)
3. Run service installation commands (see INSTALLATION.md)

---

## 🔧 **Configuration**

### **Enable/Disable Auto-Updates**
```bash
# Enable
sudo ./sentinelgo -enable-auto-update

# Disable (edit config file)
vim ~/.sentinelgo/config.json
# Set "auto_update": false
```

### **Service Management**
```bash
# Check status
sudo ./install.sh status

# Update manually
sudo ./install.sh update

# Uninstall
sudo ./install.sh uninstall
```

---

## 🎛️ **Safety Features**

### **Process Management**
- **Lock mechanism** prevents multiple instances of same version
- **Graceful shutdown** ensures clean transitions
- **Process detection** properly identifies old vs new versions
- **Configuration backup** preserves settings during updates

### **Update Safety**
- **Version comparison** only updates if newer version available
- **Rollback capability** if update fails
- **Service restart** ensures continuous operation
- **Error handling** with comprehensive logging

---

## 📊 **Technical Specifications**

### **System Requirements**
- **Linux**: Ubuntu 18.04+, Debian 9+, CentOS 7+, RHEL 7+
- **macOS**: macOS 10.15+ (Catalina and newer)
- **Windows**: Windows 10/11, Windows Server 2016+

### **Dependencies**
- **Linux**: systemd (for service management)
- **macOS**: launchd (built-in)
- **Windows**: Windows Service API (built-in)

### **Security**
- **Dedicated user** `sentinelgo` for service isolation
- **Limited permissions** - only necessary system access
- **Secure updates** - downloads from official GitHub releases
- **Configuration protection** - user-owned config files

---

## 🔄 **Update Process**

When a new version is released:

1. **Auto-check** every hour (if enabled)
2. **Download** new binary from GitHub releases
3. **Verify** version is newer than current
4. **Stop** old processes gracefully
5. **Install** new binary in correct location
6. **Restart** service with updated version
7. **Log** update completion

---

## 🆘 **Troubleshooting**

### **Common Issues**
- **Permission denied**: Run with sudo/administrator
- **Service not starting**: Check logs for specific error
- **Update failed**: Verify internet connection to GitHub
- **Multiple instances**: Use `./sentinelgo -stop` to clean up

### **Support Commands**
```bash
# Show all running processes
sudo ./install.sh status

# View logs (Linux)
journalctl -u sentinelgo -f

# View logs (macOS)
tail -f /var/log/sentinelgo.log

# Force stop all instances
sudo ./sentinelgo -stop
```

---

## 📈 **What's Fixed**

### **Previous Issues Resolved**
- ✅ **Duplicate heartbeats** - Lock mechanism prevents multiple instances
- ✅ **Manual updates** - Automatic updates from GitHub releases
- ✅ **Complex installation** - One-command universal installer
- ✅ **Architecture support** - Full AMD64 + ARM64 support
- ✅ **Cross-platform** - Works on all target operating systems

### **Enterprise Features Added**
- ✅ **Service management** - Proper system service integration
- ✅ **Auto-start** - Automatic startup on system boot
- ✅ **Configuration management** - Persistent settings
- ✅ **Logging integration** - System-level logging
- ✅ **Security hardening** - Isolated service user

---

## 🎉 **Ready for Production**

This release provides a **complete, enterprise-ready deployment solution** that:

- ✅ **Solves the original problem** - No more duplicate heartbeats
- ✅ **Enables automatic updates** - Keeps systems current
- ✅ **Simplifies deployment** - Easy for non-technical users
- ✅ **Supports all platforms** - Cross-platform compatibility
- ✅ **Provides professional tools** - Enterprise-grade installation

Deploy with confidence! 🚀
