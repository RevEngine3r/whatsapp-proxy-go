@echo off
REM WhatsApp Proxy Go - Automated Release Tag Script
REM This script creates and pushes the 'release' tag to trigger CI/CD

setlocal enabledelayedexpansion

REM Colors for output (if supported)
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "BLUE=[94m"
set "NC=[0m"

echo.
echo %BLUE%================================================%NC%
echo %BLUE%  WhatsApp Proxy Go - Release Publisher%NC%
echo %BLUE%================================================%NC%
echo.

REM Check if git.exe is available
git.exe --version >nul 2>&1
if errorlevel 1 (
    echo %RED%ERROR: git.exe not found!%NC%
    echo Please ensure Git is installed and in your PATH.
    echo Download from: https://git-scm.com/download/win
    pause
    exit /b 1
)

echo %GREEN%[1/6] Checking Git status...%NC%
git.exe status >nul 2>&1
if errorlevel 1 (
    echo %RED%ERROR: Not a git repository!%NC%
    echo Please run this script from the repository root.
    pause
    exit /b 1
)

echo %GREEN%[2/6] Fetching latest changes...%NC%
git.exe fetch origin
if errorlevel 1 (
    echo %YELLOW%WARNING: Could not fetch from origin%NC%
)

REM Check current branch
for /f "tokens=*" %%i in ('git.exe rev-parse --abbrev-ref HEAD') do set CURRENT_BRANCH=%%i
echo Current branch: %BLUE%%CURRENT_BRANCH%%NC%

if not "%CURRENT_BRANCH%"=="main" (
    echo.
    echo %YELLOW%WARNING: You are not on the 'main' branch!%NC%
    echo Current branch: %CURRENT_BRANCH%
    echo.
    set /p CONTINUE="Do you want to continue anyway? (y/N): "
    if /i not "!CONTINUE!"=="y" (
        echo %YELLOW%Release cancelled.%NC%
        pause
        exit /b 0
    )
) else (
    echo %GREEN%[3/6] Updating main branch...%NC%
    git.exe pull origin main
    if errorlevel 1 (
        echo %RED%ERROR: Could not pull latest changes!%NC%
        echo Please resolve conflicts manually.
        pause
        exit /b 1
    )
)

REM Check for uncommitted changes
echo %GREEN%[4/6] Checking for uncommitted changes...%NC%
git.exe diff --quiet
if errorlevel 1 (
    echo %YELLOW%WARNING: You have uncommitted changes!%NC%
    git.exe status --short
    echo.
    set /p CONTINUE="Do you want to continue anyway? (y/N): "
    if /i not "!CONTINUE!"=="y" (
        echo %YELLOW%Release cancelled.%NC%
        echo Please commit or stash your changes first.
        pause
        exit /b 0
    )
)

REM Delete existing release tag if it exists
echo %GREEN%[5/6] Cleaning up old release tags...%NC%
git.exe tag -d release >nul 2>&1
git.exe push origin :refs/tags/release >nul 2>&1

REM Show release summary
echo.
echo %BLUE%================================================%NC%
echo %BLUE%  Release Summary%NC%
echo %BLUE%================================================%NC%
echo.
echo Repository: %BLUE%RevEngine3r/whatsapp-proxy-go%NC%
echo Branch: %BLUE%%CURRENT_BRANCH%%NC%

REM Get latest commit info
for /f "tokens=*" %%i in ('git.exe log -1 --pretty^=format:"%%h"') do set COMMIT_HASH=%%i
for /f "tokens=*" %%i in ('git.exe log -1 --pretty^=format:"%%s"') do set COMMIT_MSG=%%i
echo Latest commit: %BLUE%%COMMIT_HASH%%NC% - %COMMIT_MSG%

REM Calculate expected version
for /f "tokens=1-6 delims=/:. " %%a in ('%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe -Command "(Get-Date).ToUniversalTime().ToString('yyyy.MM.dd.HHmm')"') do (
    set VERSION=%%a.%%b.%%c.%%d%%e
)
echo Expected version: %BLUE%v%VERSION%%NC%
echo.

REM Confirm release
set /p CONFIRM="Are you sure you want to create a release? (y/N): "
if /i not "%CONFIRM%"=="y" (
    echo %YELLOW%Release cancelled by user.%NC%
    pause
    exit /b 0
)

REM Create and push the release tag
echo.
echo %GREEN%[6/6] Creating and pushing release tag...%NC%
git.exe tag release
if errorlevel 1 (
    echo %RED%ERROR: Could not create tag!%NC%
    pause
    exit /b 1
)

echo Pushing release tag to origin...
git.exe push origin release
if errorlevel 1 (
    echo %RED%ERROR: Could not push tag!%NC%
    echo Cleaning up local tag...
    git.exe tag -d release
    pause
    exit /b 1
)

echo.
echo %GREEN%================================================%NC%
echo %GREEN%  Release Tag Pushed Successfully!%NC%
echo %GREEN%================================================%NC%
echo.
echo %BLUE%What happens next:%NC%
echo 1. GitHub Actions will start building binaries
echo 2. Tests will run automatically
echo 3. Multi-platform binaries will be created
echo 4. A new release will be published
echo.
echo %BLUE%Monitor progress at:%NC%
echo https://github.com/RevEngine3r/whatsapp-proxy-go/actions
echo.
echo %BLUE%View releases at:%NC%
echo https://github.com/RevEngine3r/whatsapp-proxy-go/releases
echo.
echo %YELLOW%Note: After the workflow completes, cleanup with:%NC%
echo   git tag -d release
echo   git push origin :refs/tags/release
echo.

pause
exit /b 0
