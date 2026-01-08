@echo off
REM WhatsApp Proxy Go - Release Tag Cleanup Script
REM This script removes the 'release' tag locally and remotely

setlocal

echo.
echo ================================================
echo   WhatsApp Proxy Go - Release Tag Cleanup
echo ================================================
echo.

REM Check if git.exe is available
git.exe --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: git.exe not found!
    echo Please ensure Git is installed and in your PATH.
    pause
    exit /b 1
)

echo This will delete the 'release' tag both locally and remotely.
echo.
set /p CONFIRM="Are you sure? (y/N): "
if /i not "%CONFIRM%"=="y" (
    echo Cleanup cancelled.
    pause
    exit /b 0
)

echo.
echo Deleting local 'release' tag...
git.exe tag -d release
if errorlevel 1 (
    echo Warning: Local tag may not exist.
) else (
    echo Local tag deleted successfully.
)

echo.
echo Deleting remote 'release' tag...
git.exe push origin :refs/tags/release
if errorlevel 1 (
    echo Warning: Remote tag may not exist.
) else (
    echo Remote tag deleted successfully.
)

echo.
echo ================================================
echo   Cleanup Complete!
echo ================================================
echo.
echo The 'release' tag has been removed.
echo You can now create a new release whenever needed.
echo.

pause
exit /b 0
