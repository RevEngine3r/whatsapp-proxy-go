# WhatsApp Proxy Go - Release Tag Cleanup Script (PowerShell)
# This script removes the 'release' tag locally and remotely

# Requires PowerShell 5.1 or later
#Requires -Version 5.1

# Set error action preference
$ErrorActionPreference = "Continue"

function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Header {
    param([string]$Title)
    Write-Host ""
    Write-ColorOutput "================================================" "Cyan"
    Write-ColorOutput "  $Title" "Cyan"
    Write-ColorOutput "================================================" "Cyan"
    Write-Host ""
}

cls
Write-Header "WhatsApp Proxy Go - Release Tag Cleanup"

# Check if git is available
try {
    $gitVersion = git.exe --version 2>&1
    Write-ColorOutput "Git found: $gitVersion" "Gray"
}
catch {
    Write-ColorOutput "ERROR: git.exe not found!" "Red"
    Write-ColorOutput "Please ensure Git is installed and in your PATH." "Yellow"
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host ""
Write-ColorOutput "This will delete the 'release' tag both locally and remotely." "Yellow"
Write-Host ""

$confirm = Read-Host "Are you sure? (y/N)"
if ($confirm -ne "y" -and $confirm -ne "Y") {
    Write-ColorOutput "Cleanup cancelled." "Yellow"
    Read-Host "Press Enter to exit"
    exit 0
}

Write-Host ""

# Delete local tag
Write-ColorOutput "Deleting local 'release' tag..." "Green"
try {
    git.exe tag -d release 2>&1 | Out-Null
    Write-ColorOutput "  Local tag deleted successfully." "Green"
}
catch {
    Write-ColorOutput "  Warning: Local tag may not exist." "Yellow"
}

# Delete remote tag
Write-ColorOutput "Deleting remote 'release' tag..." "Green"
try {
    git.exe push origin :refs/tags/release 2>&1 | Out-Null
    Write-ColorOutput "  Remote tag deleted successfully." "Green"
}
catch {
    Write-ColorOutput "  Warning: Remote tag may not exist." "Yellow"
}

Write-Header "Cleanup Complete!"

Write-ColorOutput "The 'release' tag has been removed." "Green"
Write-ColorOutput "You can now create a new release whenever needed." "Gray"
Write-Host ""

Read-Host "Press Enter to exit"
