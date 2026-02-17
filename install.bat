@echo off
REM SentinelGo Windows Installation Script
REM Enhanced version with better error handling and auto-start configuration
REM Usage: install.bat [install|uninstall|update|status|help|enable-autostart|disable-autostart]

setlocal enabledelayedexpansion

REM Configuration
set SERVICE_NAME=sentinelgo
set BINARY_NAME=sentinelgo.exe
set INSTALL_DIR=C:\opt\sentinelgo
set CONFIG_DIR=%INSTALL_DIR%\.sentinelgo
set LOG_FILE=%TEMP%\sentinelgo-install.log

REM Logging function
echo [INFO] %date% %time% - Starting SentinelGo installation script >> %LOG_FILE%

REM Check if running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] This script must be run as Administrator
    echo Right-click on install.bat and select "Run as administrator"
    echo [ERROR] Installation failed - insufficient privileges
    pause
    exit /b 1
)

echo [INFO] SentinelGo Windows Installation Script v2.0
echo [INFO] Command: %1%

REM Get command from parameters
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=install

REM Enhanced command handling
if "%COMMAND%"=="install" goto install
if "%COMMAND%"=="uninstall" goto uninstall
if "%COMMAND%"=="update" goto update
if "%COMMAND%"=="status" goto status
if "%COMMAND%"=="help" goto help
if "%COMMAND%"=="enable-autostart" goto enable-autostart
if "%COMMAND%"=="disable-autostart" goto disable-autostart
if "%COMMAND%"=="test" goto test
if "%COMMAND%"=="logs" goto logs
goto unknown

:install
echo [INFO] Installing SentinelGo...
echo [INFO] Stopping any running instances...
taskkill /F /IM sentinelgo.exe 2>nul

echo [INFO] Creating directories...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"

echo [INFO] Detecting architecture...
if exist "%ProgramFiles(x86)%" (
    set ARCH=amd64
    set BINARY_SOURCE=sentinelgo-windows-amd64.exe
) else (
    set ARCH=amd64
    set BINARY_SOURCE=sentinelgo-windows-amd64.exe
)

echo [INFO] Installing Windows %ARCH% binary...
if exist "%BINARY_SOURCE%" (
    copy "%BINARY_SOURCE%" "%INSTALL_DIR%\%BINARY_NAME%"
    echo [SUCCESS] Binary installed to %INSTALL_DIR%
    
    echo [INFO] Testing binary...
    "%INSTALL_DIR%\%BINARY_NAME%" --version >nul 2>&1
    if %errorLevel% equ 0 (
        echo [SUCCESS] Binary test passed
    ) else (
        echo [WARNING] Binary test failed, but continuing...
    )
) else (
    echo [ERROR] Binary not found: %BINARY_SOURCE%
    echo [INFO] Please download from: https://github.com/habib45/SentinelGo/releases
    pause
    exit /b 1
)

echo [INFO] Creating default config...
if not exist "%CONFIG_DIR%\config.json" (
    echo {"heartbeat_interval":"5m0s","auto_update":false} > "%CONFIG_DIR%\config.json"
    echo [SUCCESS] Default config created
) else (
    echo [INFO] Config file already exists
)

echo [INFO] Installing Windows service with auto-start...
REM Enhanced service creation with better configuration
sc.exe create "%SERVICE_NAME%" binPath= "\"%INSTALL_DIR%\%BINARY_NAME%\" -run" start= auto DisplayName= "SentinelGo Agent" type= own error= ignore obj= ""%INSTALL_DIR%\%BINARY_NAME%\""
if %errorLevel% equ 0 (
    echo [SUCCESS] Service created with auto-start enabled
    echo [INFO] Service will automatically start on Windows boot
) else (
    echo [ERROR] Failed to create service
    echo [INFO] Trying alternative method...
    goto fallback-install
)

REM Configure service recovery and dependencies
echo [INFO] Configuring service recovery...
sc.exe failure "%SERVICE_NAME%" reset= 86400 actions= restart/5000/restart/5000/restart/5000
sc.exe config "%SERVICE_NAME%" start= auto
sc.exe config "%SERVICE_NAME%" type= own

echo [SUCCESS] Service auto-start configured"

echo [INFO] Starting service...
sc.exe start "%SERVICE_NAME%"
if %errorLevel% equ 0 (
    echo [SUCCESS] SentinelGo service started successfully!
    echo [INFO] Use 'sc.exe query sentinelgo' to check status
    
) else (
    echo [ERROR] Failed to start service
:fallback-install
echo [INFO] Installing SentinelGo in background mode...
if not exist "%INSTALL_DIR%\%BINARY_NAME%" (
    echo [ERROR] Binary not found at %INSTALL_DIR%\%BINARY_NAME%
    pause
    exit /b 1
)

REM Create startup task for auto-start
echo [INFO] Creating Windows Task Scheduler task for auto-start...
schtasks /create /tn "SentinelGo" /tr "\"%INSTALL_DIR%\%BINARY_NAME%\" -run" /sc onstart /ru SYSTEM /rl highest /f
if %errorLevel% equ 0 (
    echo [SUCCESS] Task Scheduler task created for auto-start
) else (
    echo [WARNING] Failed to create Task Scheduler task
)

REM Start background process
echo [INFO] Starting SentinelGo in background mode...
start "SentinelGo" /D "%INSTALL_DIR%" "%INSTALL_DIR%\%BINARY_NAME%" -run
if %errorLevel% equ 0 (
    echo [SUCCESS] SentinelGo started in background mode
    echo [INFO] Check Task Manager for sentinelgo.exe process
    echo [INFO] Use 'install.bat status' to check status
) else (
    echo [ERROR] Failed to start background process
    pause
    exit /b 1
)
goto end

:uninstall
echo [INFO] Uninstalling SentinelGo...
echo [INFO] Stopping service...
sc.exe stop "%SERVICE_NAME%" 2>nul

echo [INFO] Deleting service...
sc.exe delete "%SERVICE_NAME%" 2>nul

echo [INFO] Removing Task Scheduler task...
schtasks /delete /tn "SentinelGo" /f 2>nul

echo [INFO] Stopping processes...
taskkill /F /IM sentinelgo.exe 2>nul

echo [INFO] Removing directories...
if exist "%INSTALL_DIR%" rmdir /S /Q "%INSTALL_DIR%" 2>nul
if exist "%CONFIG_DIR%" rmdir /S /Q "%CONFIG_DIR%" 2>nul

echo [SUCCESS] SentinelGo uninstalled successfully!
echo [INFO] Auto-start configurations removed
goto end

:update
echo [INFO] Updating SentinelGo...
echo [INFO] Stopping service...
sc.exe stop "%SERVICE_NAME%"

echo [INFO] Installing new binary...
if exist "%BINARY_SOURCE%" (
    copy "%BINARY_SOURCE%" "%INSTALL_DIR%\%BINARY_NAME%"
    echo [SUCCESS] Binary updated
) else (
    echo [ERROR] Binary not found: %BINARY_SOURCE%
    pause
    exit /b 1
)

echo [INFO] Starting service...
sc.exe start "%SERVICE_NAME%"
if %errorLevel% equ 0 (
    echo [SUCCESS] SentinelGo updated successfully!
) else (
    echo [ERROR] Failed to start service after update
    pause
    exit /b 1
)
goto end

:status
echo [INFO] Checking SentinelGo status...
echo.
echo [INFO] Service Status:
sc.exe query "%SERVICE_NAME%" 2>nul
if %errorLevel% equ 0 (
    echo [SUCCESS] Service is installed
) else (
    echo [WARNING] Service is not installed
)

echo.
echo [INFO] Running Processes:
tasklist /FI "IMAGENAME eq sentinelgo.exe" 2>nul
if %errorLevel% equ 0 (
    echo [SUCCESS] SentinelGo process is running
) else (
    echo [WARNING] No SentinelGo process found
)

echo.
echo [INFO] Task Scheduler Status:
schtasks /query /tn "SentinelGo" 2>nul
if %errorLevel% equ 0 (
    echo [SUCCESS] Auto-start task is configured
) else (
    echo [WARNING] Auto-start task is not configured
)

echo.
echo [INFO] Configuration Files:
if exist "%CONFIG_DIR%\config.json" (
    echo [SUCCESS] Config file exists: %CONFIG_DIR%\config.json
) else (
    echo [WARNING] Config file not found
)

goto end

:help
echo SentinelGo Windows Installation Script v2.0
echo.
echo Usage: install.bat [COMMAND]
echo.
echo Commands:
echo   install         Install SentinelGo as Windows service ^(default^)
echo   uninstall       Remove SentinelGo service and data
echo   update          Update SentinelGo binary
echo   status          Show service and process status
echo   enable-autostart Enable auto-start on Windows boot
echo   disable-autostart Disable auto-start on Windows boot
echo   test            Test binary functionality
echo   logs            Show installation logs
echo   help            Show this help message
echo.
echo Examples:
echo   install.bat install         # Install service
echo   install.bat status          # Check status
echo   install.bat enable-autostart # Enable auto-start
echo   install.bat uninstall        # Remove completely
echo.
echo Auto-start Features:
echo   - Windows Service with 'start= auto' configuration
echo   - Task Scheduler fallback for reliability
echo   - Automatic recovery on failure
echo   - Starts on system boot
echo.
echo Examples:
echo   install.bat install
echo   install.bat uninstall
echo   install.bat status
echo Note: Must be run as Administrator
goto end

:enable-autostart
echo [INFO] Enabling auto-start for SentinelGo...
echo [INFO] Configuring Windows Service for auto-start...
sc.exe config "%SERVICE_NAME%" start= auto
if %errorLevel% equ 0 (
    echo [SUCCESS] Service auto-start enabled
) else (
    echo [WARNING] Service not found, creating Task Scheduler task...
    schtasks /create /tn "SentinelGo" /tr "\"%INSTALL_DIR%\%BINARY_NAME%\" -run" /sc onstart /ru SYSTEM /rl highest /f
    if %errorLevel% equ 0 (
        echo [SUCCESS] Task Scheduler auto-start enabled
    ) else (
        echo [ERROR] Failed to enable auto-start
    )
)
goto end

:disable-autostart
echo [INFO] Disabling auto-start for SentinelGo...
echo [INFO] Disabling Windows Service auto-start...
sc.exe config "%SERVICE_NAME%" start= demand
if %errorLevel% equ 0 (
    echo [SUCCESS] Service auto-start disabled
) else (
    echo [WARNING] Service not found
)

echo [INFO] Removing Task Scheduler task...
schtasks /delete /tn "SentinelGo" /f 2>nul
if %errorLevel% equ 0 (
    echo [SUCCESS] Task Scheduler auto-start disabled
) else (
    echo [WARNING] Task Scheduler task not found
)
goto end

:test
echo [INFO] Testing SentinelGo binary...
if not exist "%INSTALL_DIR%\%BINARY_NAME%" (
    echo [ERROR] Binary not found: %INSTALL_DIR%\%BINARY_NAME%
    pause
    exit /b 1
)

echo [INFO] Running binary test...
"%INSTALL_DIR%\%BINARY_NAME%" --version
if %errorLevel% equ 0 (
    echo [SUCCESS] Binary test passed
) else (
    echo [ERROR] Binary test failed
)
goto end

:logs
echo [INFO] Showing installation logs...
if exist "%LOG_FILE%" (
    echo.
    echo [INFO] Installation log file: %LOG_FILE%
    echo.
    type "%LOG_FILE%"
) else (
    echo [WARNING] No log file found at: %LOG_FILE%
)
goto end

:end
echo [INFO] Operation completed successfully!

:end
echo [INFO] Operation completed!
pause
