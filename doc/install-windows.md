# SentinelGo – Windows Installation Guide

## What You Need
- Windows 10 or Windows 11 computer
- Administrator access to your computer
- PowerShell (comes with Windows)

## 🚀 Installation Methods

### Method 1: Automated Setup (Recommended)
**Use `setup-windows-service.bat` for automatic installation**

#### Step 1: Download SentinelGo
1. Go to: https://github.com/habib45/SentinelGo/releases/latest
2. Download: `sentinelgo-windows-amd64.exe`
3. Also download: `setup-windows-service.bat`
4. Save both files to a folder (e.g., `C:\temp\sentinelgo`)

#### Step 2: Run Automated Setup
1. **Right-click** on `setup-windows-service.bat`
2. **Select "Run as administrator"**
3. **Click "Yes"** on UAC prompt
4. **Follow** the on-screen instructions

#### Step 3: Verify Installation
The script will automatically:
- ✅ Create service with auto-start
- ✅ Configure background operation
- ✅ Set recovery actions
- ✅ Start the service
- ✅ Verify it's running

### Method 2: Manual Installation
**Step-by-step manual setup**

#### Step 1: Download SentinelGo
1. Go to: https://github.com/habib45/SentinelGo/releases/latest
2. Download: `sentinelgo-windows-amd64.exe`
3. Save it to your Downloads folder

#### Step 2: Create SentinelGo Folder
1. Open PowerShell as Administrator:
   - Click Start button
   - Type "PowerShell"
   - Right-click "Windows PowerShell"
   - Select "Run as administrator"
2. Create installation directory:
   ```powershell
   mkdir 'C:\Program Files\SentinelGo'
   ```

#### Step 3: Copy SentinelGo to Folder
1. In PowerShell, go to your Downloads folder:
   ```powershell
   cd Downloads
   ```
2. Copy the file:
   ```powershell
    Copy-Item .\sentinelgo-windows-amd64.exe 'C:\Program Files\SentinelGo\sentinelgo.exe'
   ```

#### Step 4: Configuration Options

### Option 1: Default Configuration (Recommended)
By default, SentinelGo will create a config file at `C:\Program Files\SentinelGo\.sentinelgo\config.json` with these settings:
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
```powershell
# Config will be created automatically at C:\Program Files\SentinelGo\.sentinelgo\config.json
C:\Program Files\SentinelGo\sentinelgo.exe -run
```

#### Method B: Specify Custom Path
```powershell
# Create config at custom location
New-Item -ItemType Directory -Path "C:\MySentinelGoConfig" -Force

$ConfigContent = @"
{
  "heartbeat_interval": "10m0s",
  "github_owner": "your-username",
  "github_repo": "your-repo",
  "current_version": "v1.9.9.0",
  "auto_update": true
}
"@

$ConfigContent | Out-File -FilePath "C:\MySentinelGoConfig\config.json" -Encoding UTF8 -Force

# Run with custom config
C:\Program Files\SentinelGo\sentinelgo.exe -run -config "C:\MySentinelGoConfig\config.json"
```

#### Method C: System-Wide Config
```powershell
# Create system-wide config
New-Item -ItemType Directory -Path "C:\sentinelgo" -Force

$ConfigContent = @"
{
  "heartbeat_interval": "5m0s",
  "github_owner": "habib45",
  "github_repo": "SentinelGo",
  "current_version": "v1.9.9.0",
  "auto_update": false
}
"@

$ConfigContent | Out-File -FilePath "C:\sentinelgo\config.json" -Encoding UTF8 -Force
```

### Option 3: Environment Variable
```powershell
# Set config path via environment variable
$env:SENTINELGO_CONFIG = "C:\Path\To\Your\config.json"
C:\Program Files\SentinelGo\sentinelgo.exe -run
```

### Option 4: Manual Config Creation
If you prefer to create config manually:

1. Create config directory:
   ```powershell
   mkdir "C:\Program Files\SentinelGo\.sentinelgo"
   ```
2. Create config file:
   ```powershell
   New-Item -Path "C:\Program Files\SentinelGo\.sentinelgo\config.json"
   ```
3. Edit the file:
   ```powershell
   notepad "C:\Program Files\SentinelGo\.sentinelgo\config.json"
   ```
4. Copy and paste this content:
   ```json
   {
     "heartbeat_interval": "5m0s",
     "github_owner": "habib45",
     "github_repo": "SentinelGo",
     "current_version": "v1.9.9.0",
     "auto_update": false
   }
   ```
5. Save the file (Ctrl+S) and close Notepad

**Important:** The `heartbeat_interval` must be a **string** in quotes (e.g., `"5m0s"`) not a number (e.g., `300`).

#### Step 5: Install as Windows Service
1. In PowerShell, go to SentinelGo folder:
   ```powershell
   cd "C:\Program Files\SentinelGo"
   ```
2. Install the service:
   ```powershell
   .\sentinelgo.exe -install
   ```

#### Step 6: Start SentinelGo
1. Start the service:
   ```powershell
   Start-Service SentinelGo
   ```
2. Check if it's running:
   ```powershell
   Get-Service SentinelGo
   ```
   You should see "Status: Running"

### Method 3: Quick Install Script
**Use `install.bat` for quick setup**

#### Step 1: Download Files
1. Download `sentinelgo-windows-amd64.exe`
2. Download `install.bat`
3. Save both to the same folder

#### Step 2: Run Install Script
1. **Right-click** on `install.bat`
2. **Select "Run as administrator"**
3. **Choose** installation option when prompted

## 🎯 What SentinelGo Does Now
- ✅ Automatically starts when your computer turns on
- ✅ Sends computer health information every 5 minutes
- ✅ Updates itself automatically (if enabled)
- ✅ Runs silently in the background
- ✅ Restarts automatically if it crashes

## 📋 Service Management Commands

### Using PowerShell
```powershell
# Check service status
Get-Service SentinelGo

# Start service
Start-Service SentinelGo

# Stop service
Stop-Service SentinelGo

# Restart service
Restart-Service SentinelGo

# Remove service
Remove-Service SentinelGo
```

### Using Command Prompt (Admin)
```cmd
# Check service status
sc.exe query sentinelgo

# Start service
sc.exe start sentinelgo

# Stop service
sc.exe stop sentinelgo

# Delete service
sc.exe delete sentinelgo
```

### Using SentinelGo Binary
```powershell
cd "C:\Program Files\SentinelGo"

# Check status
.\sentinelgo.exe -status

# Stop all processes
.\sentinelgo.exe -stop

# Check version
.\sentinelgo.exe -version

# Test run (not as service)
.\sentinelgo.exe -run

# Enable auto-updates
.\sentinelgo.exe -enable-auto-update
```

## 🔧 Troubleshooting

### Problem: "Access is denied"
**What it means:** You need administrator rights

**Easy Fix:**
1. Close PowerShell/Command Prompt
2. Right-click and select "Run as administrator"
3. Try the command again

### Problem: "Cannot start service"
**What it means:** SentinelGo can't start properly

**Easy Fix:**
1. Test if SentinelGo works manually:
   ```powershell
   cd "C:\Program Files\SentinelGo"
   .\sentinelgo.exe -run
   ```
2. If it works, reinstall service:
   ```powershell
   .\sentinelgo.exe -uninstall
   .\sentinelgo.exe -install
   Start-Service SentinelGo
   ```
3. Or use automated setup:
   ```powershell
   .\setup-windows-service.bat
   ```

### Problem: "Application Control policy blocked"
**What it means:** Windows is blocking the program for security

**Easy Fix:**
1. Open Windows Security (search in Start menu)
2. Go to "App & Browser Control"
3. Click "Exploit protection settings"
4. Click "Controlled folder access"
5. Click "Allow an app through Controlled folder access"
6. Click "Add an allowed app" → "Recently blocked apps"
7. Find and select `sentinelgo.exe`
8. Click "Add"
9. Try the command again

### Problem: Service fails to start after restart
**What it means:** Configuration or permission issues

**Easy Fix:**
1. Check Event Viewer:
   - Press `Win + R`, type `eventvwr.msc`
   - Go to **Windows Logs → Application**
   - Look for "sentinelgo" errors
2. Verify config file exists:
   ```powershell
   Test-Path "C:\Program Files\SentinelGo\.sentinelgo\config.json"
   ```
3. Re-run setup script:
   ```powershell
   .\setup-windows-service.bat
   ```

### Problem: Binary not found
**What it means:** SentinelGo executable is missing

**Easy Fix:**
1. Verify binary exists:
   ```powershell
   Test-Path "C:\Program Files\SentinelGo\sentinelgo.exe"
   ```
2. If missing, copy it again:
   ```powershell
   Copy-Item .\sentinelgo-windows-amd64.exe "C:\Program Files\SentinelGo\sentinelgo.exe"
   ```

## 🗑️ How to Uninstall SentinelGo

### Method 1: Using SentinelGo Binary
```powershell
cd "C:\Program Files\SentinelGo"

# Stop and remove service
.\sentinelgo.exe -uninstall

# Delete folder
Remove-Item "C:\Program Files\SentinelGo" -Recurse -Force
```

### Method 2: Using PowerShell
```powershell
# Stop service
Stop-Service SentinelGo

# Remove service
Remove-Service SentinelGo

# Delete folder
Remove-Item "C:\Program Files\SentinelGo" -Recurse -Force
```

### Method 3: Using Command Prompt
```cmd
# Stop and delete service
sc.exe stop sentinelgo
sc.exe delete sentinelgo

# Delete folder
rmdir /s /q "C:\Program Files\SentinelGo"
```

## 📊 What Information SentinelGo Sends
SentinelGo sends this information to your dashboard:
- Computer name and unique device ID
- Operating system version (Windows)
- System uptime and performance
- CPU, memory, and disk usage
- Network information and connectivity
- Software version and update status

**Your data is secure and only sent to your own dashboard.**

## 🆘 Need Help?

### Quick Help Commands
```powershell
# Show all available options
.\sentinelgo.exe -help

# Check current status
.\sentinelgo.exe -status

# Show running processes
.\sentinelgo.exe -status
```

### Common Issues
1. **Always run as Administrator** - Required for service management
2. **Check Event Viewer** - For detailed error messages
3. **Verify file paths** - Ensure binary and config exist
4. **Use automated setup** - `setup-windows-service.bat` handles most issues

### Support Resources
- Check the "Common Problems" section above
- Verify you're running PowerShell/Command Prompt as administrator
- Check the GitHub repository for latest updates
- Use the automated setup script for reliable installation

## 📋 Installation Summary

### ✅ After Successful Installation:
- SentinelGo is installed as a Windows Service
- Service starts automatically when computer boots
- Sends regular health updates to your dashboard
- Updates itself automatically (if enabled)
- Runs silently in the background
- Restarts automatically if it crashes

### 🎯 Recommended Method:
**Use `setup-windows-service.bat` for the most reliable installation**
- Handles all configuration automatically
- Sets up proper service permissions
- Configures recovery actions
- Provides detailed feedback
- Includes troubleshooting features

### 🔄 Alternative Methods:
- **`install.bat`** - Quick installation with options
- **Manual setup** - Full control over installation process
- **PowerShell commands** - For advanced users

Choose the method that best fits your needs and technical comfort level!
