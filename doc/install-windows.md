# SentinelGo – Windows Installation Guide

## What You Need
- Windows 10 or Windows 11 computer
- Administrator access to your computer
- PowerShell (comes with Windows)

## Step 1: Download SentinelGo
1. Go to this website: https://github.com/habib45/SentinelGo/releases/latest
2. Download the file named: `sentinelgo-windows-amd64.exe`
3. Save it to your Downloads folder

## Step 2: Create SentinelGo Folder
1. Open PowerShell as Administrator:
   - Click Start button
   - Type "PowerShell"
   - Right-click "Windows PowerShell"
   - Select "Run as administrator"
2. Type this command and press Enter:
   ```powershell
   mkdir C:\ProgramFiles\SentinelGo
   ```

## Step 3: Copy SentinelGo to Folder
1. In PowerShell, go to your Downloads folder:
   ```powershell
   cd Downloads
   ```
2. Copy the file:
   ```powershell
   Copy-Item .\sentinelgo-windows-amd64.exe C:\ProgramFiles\SentinelGo\sentinelgo.exe
   ```

## Step 4: Create Settings File (Optional)
1. Create settings folder:
   ```powershell
   mkdir C:\ProgramData\sentinelgo
   ```
2. Create settings file:
   ```powershell
   notepad C:\ProgramData\sentinelgo\config.json
   ```
3. Copy and paste this text into Notepad:
   ```json
   {
     "device_id": "my-computer-001",
     "heartbeat_interval": "5m0s",
     "github_owner": "habib45",
     "github_repo": "SentinelGo",
     "current_version": "v1.8.4"
   }
   ```
4. Save the file and close Notepad

## Step 5: Install SentinelGo as a Service
1. In PowerShell, go to SentinelGo folder:
   ```powershell
   cd C:\ProgramFiles\SentinelGo
   ```
2. Install the service:
   ```powershell
   .\sentinelgo.exe -install
   ```

## Step 6: Start SentinelGo
1. Start the service:
   ```powershell
   Start-Service SentinelGo
   ```
2. Check if it's running:
   ```powershell
   Get-Service SentinelGo
   ```
   You should see "Status: Running"

## What SentinelGo Does Now
- ✅ Automatically starts when your computer turns on
- ✅ Sends computer health information every 5 minutes
- ✅ Updates itself automatically
- ✅ Runs silently in the background

## Quick Commands You Can Use

### Check SentinelGo Status
```powershell
cd C:\ProgramFiles\SentinelGo
.\sentinelgo.exe -status
```

### Stop SentinelGo
```powershell
cd C:\ProgramFiles\SentinelGo
.\sentinelgo.exe -stop
```

### Check Version
```powershell
cd C:\ProgramFiles\SentinelGo
.\sentinelgo.exe -version
```

### Test SentinelGo (Not as Service)
```powershell
cd C:\ProgramFiles\SentinelGo
.\sentinelgo.exe -run
```

## Common Problems and Easy Fixes

### Problem: "Access is denied"
**What it means:** You need administrator rights

**Easy Fix:**
1. Close PowerShell
2. Right-click PowerShell icon
3. Select "Run as administrator"
4. Try the command again

### Problem: "Cannot start service"
**What it means:** SentinelGo can't start properly

**Easy Fix:**
1. Test if SentinelGo works:
   ```powershell
   cd C:\ProgramFiles\SentinelGo
   .\sentinelgo.exe -run
   ```
2. If it works, reinstall:
   ```powershell
   .\sentinelgo.exe -uninstall
   .\sentinelgo.exe -install
   Start-Service SentinelGo
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

### Problem: Service won't start after restart
**What it means:** Settings file is missing

**Easy Fix:**
1. Create settings file (see Step 4)
2. Restart the service:
   ```powershell
   Stop-Service SentinelGo
   Start-Service SentinelGo
   ```

## How to Uninstall SentinelGo
1. Stop the service:
   ```powershell
   Stop-Service SentinelGo
   ```
2. Remove the service:
   ```powershell
   cd C:\ProgramFiles\SentinelGo
   .\sentinelgo.exe -uninstall
   ```
3. Delete the folder:
   ```powershell
   Remove-Item C:\ProgramFiles\SentinelGo -Recurse -Force
   ```

## What Information SentinelGo Sends
SentinelGo sends this information to your dashboard:
- Computer name and ID
- Operating system (Windows)
- How long computer has been running
- CPU, memory, and disk usage
- Network information

**Your data is secure and only sent to your own dashboard.**

## Need Help?
If you have any problems:
1. Check the "Common Problems" section above
2. Make sure you're running PowerShell as administrator
3. Contact your IT support person

## Summary
After following these steps:
- ✅ SentinelGo is installed on your computer
- ✅ It starts automatically when computer turns on
- ✅ It sends regular health updates
- ✅ It updates itself automatically
- ✅ You can control it with simple commands
