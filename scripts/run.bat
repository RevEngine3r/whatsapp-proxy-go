@echo off
REM WhatsApp Proxy Go - Windows Runner Script
REM Usage: run.bat [config_file]

setlocal enabledelayedexpansion

REM Configuration
if "%~1"=="" (
    set CONFIG_FILE=config.yaml
) else (
    set CONFIG_FILE=%~1
)

set BINARY=whatsapp-proxy.exe

REM Check if binary exists
if not exist "%BINARY%" (
    if exist "dist\whatsapp-proxy.exe" (
        set BINARY=dist\whatsapp-proxy.exe
    ) else (
        echo [ERROR] Binary not found: %BINARY%
        echo Please build the project first: make build
        exit /b 1
    )
)

REM Check if config file exists
if not exist "%CONFIG_FILE%" (
    echo [ERROR] Config file not found: %CONFIG_FILE%
    echo Usage: %0 [config_file]
    echo.
    echo Copy the example config to get started:
    echo   copy configs\config.example.yaml config.yaml
    exit /b 1
)

echo Starting WhatsApp Proxy Go
echo Config: %CONFIG_FILE%
echo.

REM Run the proxy
"%BINARY%" --config "%CONFIG_FILE%"
