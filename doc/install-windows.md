# SentinelGo – Windows Installation Guide

## Prerequisites
- Windows 10/11 or Windows Server 2016+
- Administrator privileges
- PowerShell 5.1+ (built‑in)

## 1. Download the Binary
Download the latest Windows binary from GitHub Releases:
```
https://github.com/habib45/SentinelGo/releases/latest
```
Choose the file named `sentinelgo-windows-amd64.exe`.

## 2. Create a Directory
```powershell
mkdir C:\ProgramFiles\SentinelGo
```

## 3. Copy the Binary
```powershell
Copy-Item -Path .\sentinelgo-windows-amd64.exe -Destination C:\ProgramFiles\SentinelGo\sentinelgo.exe
```

## 4. (Optional) Create a Configuration File
Create configuration directory and file:
```powershell
mkdir C:\ProgramData\sentinelgo
```

Create `C:\ProgramData\sentinelgo\config.json`:
```json
{
  "device_id": "windows-pc-001",
  "heartbeat_interval": "5m",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.8.4"
}
```

## 5. Install as a Windows Service (Recommended)
Run PowerShell as Administrator:
```powershell
cd C:\ProgramFiles\SentinelGo
.\sentinelgo.exe -install
```

This will automatically:
- Install SentinelGo as a Windows service named "SentinelGo"
- Configure the service to start automatically on system boot
- Set up proper logging to Windows Event Viewer
- Configure automatic updates

## 6. Start the Service
```powershell
Start-Service SentinelGo
```

## 7. Verify Installation
```powershell
# Check service status
Get-Service SentinelGo

# Check version
.\sentinelgo.exe -version

# Check running processes
.\sentinelgo.exe -status
```

Status should be **Running**.

## 8. Service Management Commands

### View Service Status
```powershell
.\sentinelgo.exe -status
```
Shows all running SentinelGo processes with versions.

### Stop All Processes
```powershell
.\sentinelgo.exe -stop
```
Safely stops all running SentinelGo processes.

### Run in Foreground (for testing)
```powershell
.\sentinelgo.exe -run
```
Runs the agent in console mode (not as service).

### Check Version
```powershell
.\sentinelgo.exe -version
```
Shows current build version and platform info.

## 9. Automatic Updates
SentinelGo now includes intelligent automatic update management:

- **Daily Update Checks**: Automatically checks for new versions every 24 hours
- **Safe Process Management**: Stops old processes before applying updates
- **Seamless Restart**: Automatically restarts with new version
- **Version Tracking**: Identifies which version each process is running
- **Fallback Protection**: Continues running if update fails

### Manual Update Check
```powershell
# Force update check (requires internet)
.\sentinelgo.exe -run
# The updater runs automatically in the background
```

## 10. Uninstall (if needed)
```powershell
Stop-Service SentinelGo
.\sentinelgo.exe -uninstall
```

## 11. Troubleshooting

### Check Windows Event Logs
Windows Event Viewer → Windows Logs → Application → Source "SentinelGo"

### Common Issues

#### Service Won't Start
```powershell
# Check permissions
Get-Acl C:\ProgramFiles\SentinelGo\sentinelgo.exe

# Test manually
.\sentinelgo.exe -run
```

#### Multiple Versions Running
```powershell
# Check for conflicts
.\sentinelgo.exe -status

# Stop all old versions
.\sentinelgo.exe -stop

# Restart service
Stop-Service SentinelGo
Start-Service SentinelGo
```

#### Update Issues
```powershell
# Stop all processes
.\sentinelgo.exe -stop

# Download latest version manually
# Replace binary in C:\ProgramFiles\SentinelGo\

# Restart service
Start-Service SentinelGo
```

## 12. Heartbeat Data
SentinelGo sends the following data to your Supabase database:

```json
{
  "device_id": "windows-pc-001",
  "alive": "true",
  "employee_id": "WINDOWS-PC",
  "os": "windows",
  "uptime": 86400,
  "uptime_formatted": "1 day 0 hours 0 minutes",
  "timestamp": "2024-02-10T17:30:00Z",
  "hostname": "WINDOWS-PC",
  "platform": "Microsoft Windows 10 Pro",
  "platform_version": "10.0.19045",
  "arch": "amd64",
  "cpu": {
    "model_name": "Intel(R) Core(TM) i7-9700K",
    "cores": 8,
    "usage": 15.5
  },
  "memory": {
    "total": 17179869184,
    "used": 8589934592,
    "free": 8589934592,
    "usage": 50.0
  },
  "disk": {
    "total": 107374182400,
    "used": 53687091200,
    "free": 53687091200
  },
  "network": [
    {
      "name": "Ethernet",
      "bytes_sent": 1048576000,
      "bytes_recv": 2097152000
    }
  ]
}
```

## Key Features

### New in This Version
- **Automatic Updates**: Safe, intelligent update management
- **Version Tracking**: Clear identification of running versions
- **Process Management**: Safe stopping of old processes
- **Enhanced Logging**: Better error reporting and diagnostics
- **Uptime Formatting**: Human-readable uptime display
- **Cross-Platform**: Consistent behavior across Windows, Linux, and macOS

### Configuration Options
- **Device ID**: Unique identifier for the machine
- **Heartbeat Interval**: How often to send data (default: 5m)
- **Update Frequency**: How often to check for updates (default: 24h)
- **GitHub Repository**: Where to download updates from

### Monitoring Capabilities
- **System Uptime**: Both raw seconds and formatted string
- **Resource Usage**: CPU, memory, disk, and network statistics
- **Service Health**: Automatic restart on failures
- **Update Status**: Track update success/failure

## Security Notes
- The service runs under the LocalSystem account
- All network communications use HTTPS
- Configuration files should be protected with appropriate permissions
- Supabase credentials are stored in configuration files (not environment variables)
