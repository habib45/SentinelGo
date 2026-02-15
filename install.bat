@echo off
REM SentinelGo Windows Installation Script
REM Usage: install.bat [install|uninstall|update|status|help]

setlocal enabledelayedexpansion

REM Configuration
set SERVICE_NAME=sentinelgo
set BINARY_NAME=sentinelgo.exe
set INSTALL_DIR=C:\opt\sentinelgo
set CONFIG_DIR=%INSTALL_DIR%\.sentinelgo

REM Check if running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] This script must be run as Administrator
    echo Right-click on install.bat and select "Run as administrator"
    pause
    exit /b 1
)

REM Get command from parameters
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=install

echo [INFO] SentinelGo Windows Installation Script
echo [INFO] Command: %COMMAND%

if "%COMMAND%"=="install" goto install
if "%COMMAND%"=="uninstall" goto uninstall
if "%COMMAND%"=="update" goto update
if "%COMMAND%"=="status" goto status
if "%COMMAND%"=="help" goto help
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
) else (
    echo [ERROR] Binary not found: %BINARY_SOURCE%
    echo [INFO] Please download from: https://github.com/habib45/SentinelGo/releases
    pause
    exit /b 1
)

echo [INFO] Creating default config...
if not exist "%CONFIG_DIR%\config.json" (
    echo {"heartbeat_interval":"5m","auto_update":false} > "%CONFIG_DIR%\config.json"
    echo [SUCCESS] Default config created
)

echo [INFO] Installing Windows service...
sc.exe create "%SERVICE_NAME%" binPath= "%INSTALL_DIR%\%BINARY_NAME%" -run start= auto
if %errorLevel% equ 0 (
    echo [SUCCESS] Service created
) else (
    echo [ERROR] Failed to create service
    pause
    exit /b 1
)

echo [INFO] Starting service...
sc.exe start "%SERVICE_NAME%"
if %errorLevel% equ 0 (
    echo [SUCCESS] SentinelGo service started successfully!
    echo [INFO] Use 'sc.exe query sentinelgo' to check status
) else (
    echo [ERROR] Failed to start service
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

echo [INFO] Stopping processes...
taskkill /F /IM sentinelgo.exe 2>nul

echo [INFO] Removing directories...
rmdir /S /Q "%INSTALL_DIR%" 2>nul

echo [SUCCESS] SentinelGo uninstalled successfully!
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
echo [INFO] SentinelGo Service Status:
sc.exe query "%SERVICE_NAME%"
echo.
echo [INFO] Running processes:
tasklist | findstr sentinelgo.exe
goto end

:help
echo SentinelGo Windows Installation Script
echo.
echo Usage: install.bat [COMMAND]
echo.
echo Commands:
echo   install     Install SentinelGo as Windows service ^(default^)
echo   uninstall   Remove SentinelGo service and data
echo   update      Update SentinelGo binary
echo   status      Show service status
echo   help        Show this help message
echo.
echo Examples:
echo   install.bat install
echo   install.bat uninstall
echo   install.bat status
echo.
echo Note: Must be run as Administrator
goto end

:unknown
echo [ERROR] Unknown command: %COMMAND%
echo.
echo Use 'install.bat help' to see available commands
pause
exit /b 1

:end
echo [INFO] Operation completed!
pause
