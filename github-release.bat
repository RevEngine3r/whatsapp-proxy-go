@echo off
REM ============================================================
REM GitHub Release Publisher
REM ============================================================
REM This script pushes a 'release' tag to trigger automated
REM GitHub Actions workflow for creating releases.
REM
REM Requirements:
REM   - git.exe in PATH
REM   - Committed changes (clean working tree recommended)
REM ============================================================

setlocal enabledelayedexpansion

REM Configuration
set TAG_NAME=release
set REMOTE=origin

REM Colors for output (if terminal supports it)
echo.
echo ============================================================
echo   GitHub Release Publisher
echo ============================================================
echo.

REM Check if git is installed
git --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] git.exe not found in PATH
    echo Please ensure Git is installed and added to your PATH
    pause
    exit /b 1
)

echo [OK] Git found: 
for /f "tokens=*" %%i in ('git --version') do echo     %%i
echo.

REM Check if we're in a git repository
git rev-parse --git-dir >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Not in a git repository
    echo Please run this script from the repository root
    pause
    exit /b 1
)

echo [OK] Git repository detected
echo.

REM Check for uncommitted changes (warning only)
for /f %%i in ('git status --porcelain ^| find /c /v ""') do set CHANGES=%%i
if !CHANGES! gtr 0 (
    echo [WARNING] You have uncommitted changes:
    git status --short
    echo.
    echo It's recommended to commit changes before creating a release.
    echo.
    set /p CONTINUE="Continue anyway? (y/N): "
    if /i not "!CONTINUE!"=="y" (
        echo Release cancelled.
        pause
        exit /b 0
    )
    echo.
)

REM Get current branch
for /f "tokens=*" %%i in ('git branch --show-current') do set CURRENT_BRANCH=%%i
echo [INFO] Current branch: !CURRENT_BRANCH!
echo.

REM Confirm release
echo This will:
echo   1. Delete local '%TAG_NAME%' tag (if exists)
echo   2. Create new '%TAG_NAME%' tag at current HEAD
echo   3. Force push tag to %REMOTE%
echo   4. Trigger GitHub Actions release workflow
echo.
set /p CONFIRM="Proceed with release? (y/N): "

if /i not "!CONFIRM!"=="y" (
    echo Release cancelled.
    pause
    exit /b 0
)

echo.
echo ============================================================
echo   Creating Release
echo ============================================================
echo.

REM Delete local tag if exists
echo [STEP 1/3] Deleting local tag (if exists)...
git tag -d %TAG_NAME% >nul 2>&1
if errorlevel 1 (
    echo     No existing local tag found
) else (
    echo     Local tag deleted
)
echo.

REM Create new tag
echo [STEP 2/3] Creating new tag at HEAD...
for /f "tokens=*" %%i in ('git rev-parse --short HEAD') do set COMMIT=%%i
echo     Commit: !COMMIT!

git tag %TAG_NAME%
if errorlevel 1 (
    echo [ERROR] Failed to create tag
    pause
    exit /b 1
)
echo     Tag created successfully
echo.

REM Push tag to remote (force)
echo [STEP 3/3] Pushing tag to %REMOTE%...
git push %REMOTE% %TAG_NAME% --force
if errorlevel 1 (
    echo [ERROR] Failed to push tag to remote
    echo     You may need to authenticate or check network connection
    pause
    exit /b 1
)

echo.
echo ============================================================
echo   SUCCESS!
echo ============================================================
echo.
echo Tag '%TAG_NAME%' pushed successfully!
echo.
echo GitHub Actions workflow will now:
echo   1. Run build-release.ps1 to create release binaries
echo   2. Auto-generate version from datetime
echo   3. Build for Linux, macOS, Windows (multiple architectures)
echo   4. Create GitHub release with all binaries and checksums
echo.
echo Check workflow progress at:
for /f "tokens=*" %%i in ('git remote get-url %REMOTE%') do set REPO_URL=%%i
REM Convert git URL to https if needed
set REPO_URL=!REPO_URL:git@github.com:=https://github.com/!
set REPO_URL=!REPO_URL:.git=!
echo   !REPO_URL!/actions
echo.

pause
exit /b 0
