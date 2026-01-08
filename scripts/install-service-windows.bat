@echo off
REM WhatsApp Proxy Go - Windows Service Installation Script
REM Requires administrator privileges

setlocal enabledelayedexpansion

REM Check for admin privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] This script requires administrator privileges
    echo Please run as Administrator
    pause
    exit /b 1
)

REM Configuration
set SERVICE_NAME=WhatsAppProxy
set DISPLAY_NAME=WhatsApp Proxy Server
set DESCRIPTION=WhatsApp Proxy with SOCKS5 upstream support
set BINARY=whatsapp-proxy.exe
set CONFIG_FILE=config.yaml

echo ========================================
echo   WhatsApp Proxy - Service Installer
echo ========================================
echo.

REM Check if binary exists
if not exist "%BINARY%" (
    if exist "dist\%BINARY%" (
        set BINARY=dist\%BINARY%
    ) else (
        echo [ERROR] Binary not found: %BINARY%
        echo Please build the project first: make build
        pause
        exit /b 1
    )
)

REM Get full path to binary
for %%I in ("%BINARY%") do set BINARY_PATH=%%~fI

REM Check if config exists
if not exist "%CONFIG_FILE%" (
    if exist "configs\config.example.yaml" (
        echo [INFO] Copying example config...
        copy "configs\config.example.yaml" "%CONFIG_FILE%"
        echo [WARN] Please edit config.yaml before starting the service
    ) else (
        echo [ERROR] Config file not found: %CONFIG_FILE%
        pause
        exit /b 1
    )
)

REM Get full path to config
for %%I in ("%CONFIG_FILE%") do set CONFIG_PATH=%%~fI

REM Check if service already exists
sc query "%SERVICE_NAME%" >nul 2>&1
if %errorLevel% equ 0 (
    echo [WARN] Service already exists. Removing old service...
    sc stop "%SERVICE_NAME%" >nul 2>&1
    timeout /t 2 /nobreak >nul
    sc delete "%SERVICE_NAME%"
    timeout /t 2 /nobreak >nul
)

REM Create service
echo [INFO] Creating service...
sc create "%SERVICE_NAME%" binPath= "\"%BINARY_PATH%\" --config \"%CONFIG_PATH%\"" start= auto DisplayName= "%DISPLAY_NAME%"

if %errorLevel% neq 0 (
    echo [ERROR] Failed to create service
    pause
    exit /b 1
)

REM Set service description
sc description "%SERVICE_NAME%" "%DESCRIPTION%"

REM Set service recovery options (restart on failure)
sc failure "%SERVICE_NAME%" reset= 86400 actions= restart/5000/restart/10000/restart/30000

echo.
echo ========================================
echo   Installation Complete!
echo ========================================
echo.
echo Service Name: %SERVICE_NAME%
echo Binary: %BINARY_PATH%
echo Config: %CONFIG_PATH%
echo.
echo Service commands:
echo   Start:   sc start %SERVICE_NAME%
echo   Stop:    sc stop %SERVICE_NAME%
echo   Status:  sc query %SERVICE_NAME%
echo   Remove:  sc delete %SERVICE_NAME%
echo.
echo Or use Services Manager (services.msc)
echo.
echo [NOTE] Service is created but not started.
echo        Please edit config.yaml and start manually.
echo.
pause
