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
If you want to override defaults, create:
```
C:\ProgramData\sentinelgo\config.json
```
Example:
```json
{
  "heartbeat_interval": "5m",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v0.1.0"
}
```

## 5. Install as a Windows Service
Run PowerShell as Administrator:
```powershell
cd C:\ProgramFiles\SentinelGo
.\sentinelgo.exe -install
```

The service will be named **SentinelGo** and set to start automatically.

## 6. Start the Service
```powershell
Start-Service SentinelGo
```

## 7. Verify
```powershell
Get-Service SentinelGo
```
Status should be **Running**.

## 8. Uninstall (if needed)
```powershell
Stop-Service SentinelGo
.\sentinelgo.exe -uninstall
```

## Logs
Windows Event Viewer → Windows Logs → Application → Source “SentinelGo”.

## Notes
- Supabase connection is embedded at build time; users do not configure it.
- The agent checks for updates once every 24 hours.
- Heartbeat is sent every 5 minutes (configurable).
- The service runs under the LocalSystem account.
