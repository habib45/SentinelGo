@echo off
REM SentinelGo Windows Service Setup Script
REM Configures SentinelGo to run as a proper Windows service in background

setlocal enabledelayedexpansion

REM Configuration
set SERVICE_NAME=sentinelgo
set DISPLAY_NAME=SentinelGo Agent
REM Try multiple possible binary locations
if exist "C:\Program Files\SentinelGo\sentinelgo.exe" (
    set BINARY_PATH=C:\Program Files\SentinelGo\sentinelgo.exe
) else if exist "C:\Program Files (x86)\SentinelGo\sentinelgo.exe" (
    set BINARY_PATH=C:\Program Files (x86)\SentinelGo\sentinelgo.exe
) else if exist "%~dp0sentinelgo.exe" (
    set BINARY_PATH=%~dp0sentinelgo.exe
) else (
    echo [ERROR] SentinelGo binary not found in standard locations
    echo [INFO] Please install SentinelGo first using install.bat
    pause
    exit /b 1
)
set DESCRIPTION=SentinelGo monitoring agent for system health and performance

echo [INFO] SentinelGo Windows Service Setup
echo =====================================

REM Check if running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] This script must be run as Administrator
    echo Right-click on setup-windows-service.bat and select "Run as administrator"
    pause
    exit /b 1
)

REM Check if binary exists
if not exist "%BINARY_PATH%" (
    echo [ERROR] SentinelGo binary not found: %BINARY_PATH%
    echo [INFO] Please install SentinelGo first using install.bat
    pause
    exit /b 1
)

echo [INFO] Found SentinelGo binary at %BINARY_PATH%

REM Test binary functionality
echo [INFO] Testing binary functionality...
"%BINARY_PATH%" --version >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] Binary test failed - checking if binary is valid...
    echo [INFO] Binary size: 
    dir "%BINARY_PATH%" | find "sentinelgo.exe"
    pause
    exit /b 1
)
echo [SUCCESS] Binary test passed

REM Stop and remove existing service
echo [INFO] Removing any existing SentinelGo service...
sc.exe stop "%SERVICE_NAME%" 2>nul
sc.exe delete "%SERVICE_NAME%" 2>nul
timeout /t 2 /nobreak >nul

REM Create Windows Service with proper background configuration
echo [INFO] Creating Windows Service...
sc.exe create "%SERVICE_NAME%" ^
    binPath= "\"%BINARY_PATH%\" -run" ^
    DisplayName= "%DISPLAY_NAME%" ^
    Description= "%DESCRIPTION%" ^
    start= auto ^
    type= own ^
    obj= LocalSystem ^
    error= ignore

if %errorLevel% neq 0 (
    echo [ERROR] Failed to create Windows Service
    echo [INFO] Trying alternative configuration...
    
    REM Alternative service creation
    sc.exe create "%SERVICE_NAME%" ^
        binPath= "\"%BINARY_PATH%\" -run" ^
        DisplayName= "%DISPLAY_NAME%" ^
        start= auto ^
        type= own
        
    if %errorLevel% neq 0 (
        echo [ERROR] Service creation failed completely
        pause
        exit /b 1
    )
)

echo [SUCCESS] Windows Service created

REM Configure service for background operation
echo [INFO] Configuring service for background operation...

REM Set service to run in background (no desktop interaction)
sc.exe config "%SERVICE_NAME%" type= own
sc.exe config "%SERVICE_NAME%" start= auto
sc.exe config "%SERVICE_NAME%" error= ignore

REM Configure service recovery (restart on failure)
sc.exe failure "%SERVICE_NAME%" reset= 86400 actions= restart/5000/restart/5000/restart/5000

REM Set service dependencies (start after network)
sc.exe config "%SERVICE_NAME%" depend= Tcpip/Dnscache

REM Configure service to run as SYSTEM with proper privileges
sc.exe sidtype "%SERVICE_NAME%" unrestricted

echo [SUCCESS] Service configuration completed

REM Set service permissions
echo [INFO] Setting service permissions...
sc.exe sdset "%SERVICE_NAME%" "D:(A;;CC;;;AU)(A;;CC;;;SY)(A;;CC;;;BA)"

REM Start the service
echo [INFO] Starting SentinelGo service...
sc.exe start "%SERVICE_NAME%"

REM Wait for service to start
timeout /t 5 /nobreak >nul

REM Check service status more reliably
echo [INFO] Checking service status...
sc.exe query "%SERVICE_NAME%" >nul 2>&1
if %errorLevel% equ 0 (
    REM Service exists, check if running
    for /f "tokens=3 delims=: " %%i in ('sc.exe query "%SERVICE_NAME%" ^| find "STATE"') do (
        if "%%i"=="RUNNING" (
            echo [SUCCESS] SentinelGo service is running in background
            echo [INFO] Service will automatically start on Windows boot
            echo [INFO] Service runs without desktop interaction
            set SERVICE_RUNNING=1
        ) else (
            echo [WARNING] Service exists but not running (STATE: %%i)
            set SERVICE_RUNNING=0
        )
    )
) else (
    echo [ERROR] Service not found or query failed
    set SERVICE_RUNNING=0
    
    REM Try Task Scheduler fallback
    echo [INFO] Trying Task Scheduler fallback...
    schtasks /create /tn "SentinelGo" /tr "\"%BINARY_PATH%\" -run" /sc onstart /ru SYSTEM /rl highest /f 2>nul
    if !errorLevel! equ 0 (
        echo [SUCCESS] Task Scheduler fallback created
        echo [INFO] Starting Task Scheduler task...
        schtasks /run /tn "SentinelGo" 2>nul
        timeout /t 3 /nobreak >nul
        
        REM Check if process is running
        tasklist /FI "IMAGENAME eq sentinelgo.exe" 2>nul | find "sentinelgo.exe" >nul
        if !errorLevel! equ 0 (
            echo [SUCCESS] SentinelGo is running via Task Scheduler
            set SERVICE_RUNNING=1
        ) else (
            echo [ERROR] Task Scheduler task failed to start
        )
    ) else (
        echo [ERROR] Task Scheduler fallback failed
    )
)

echo.
echo [INFO] Service Management Commands:
echo   Start:    sc.exe start %SERVICE_NAME%
echo   Stop:     sc.exe stop %SERVICE_NAME%
echo   Status:   sc.exe query %SERVICE_NAME%
echo   Config:   sc.exe qc %SERVICE_NAME%
echo   Delete:   sc.exe delete %SERVICE_NAME%
echo   Task Start: schtasks /run /tn "SentinelGo"
echo   Task Stop:  schtasks /end /tn "SentinelGo"
echo   Task Delete: schtasks /delete /tn "SentinelGo" /f
echo   Logs:     eventvwr.msc (Windows Event Viewer)
echo.
echo [INFO] Background Operation Features:
echo   - Runs as Windows Service (primary)
echo   - Task Scheduler fallback (backup)
echo   - No desktop interaction required
echo   - Automatic start on boot
echo   - Auto-restart on failure
echo   - Runs with SYSTEM privileges
echo   - Logs to Windows Event Log
echo   - Multiple installation paths supported
echo   - Comprehensive error handling

pause
