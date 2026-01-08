# WhatsApp Proxy Go - Automated Release Tag Script (PowerShell)
# This script creates and pushes the 'release' tag to trigger CI/CD

# Requires PowerShell 5.1 or later
#Requires -Version 5.1

# Set error action preference
$ErrorActionPreference = "Stop"

# Function to write colored output
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

function Write-Step {
    param(
        [int]$Step,
        [int]$Total,
        [string]$Message
    )
    Write-ColorOutput "[$Step/$Total] $Message..." "Green"
}

# Main script
cls
Write-Header "WhatsApp Proxy Go - Release Publisher"

try {
    # Check if git is available
    Write-Step 1 6 "Checking Git installation"
    try {
        $gitVersion = git.exe --version 2>&1
        Write-ColorOutput "  Git found: $gitVersion" "Gray"
    }
    catch {
        Write-ColorOutput "ERROR: git.exe not found!" "Red"
        Write-ColorOutput "Please ensure Git is installed and in your PATH." "Yellow"
        Write-ColorOutput "Download from: https://git-scm.com/download/win" "Yellow"
        Read-Host "Press Enter to exit"
        exit 1
    }

    # Check if we're in a git repository
    Write-Step 2 6 "Verifying Git repository"
    try {
        git.exe status | Out-Null
    }
    catch {
        Write-ColorOutput "ERROR: Not a git repository!" "Red"
        Write-ColorOutput "Please run this script from the repository root." "Yellow"
        Read-Host "Press Enter to exit"
        exit 1
    }

    # Fetch latest changes
    Write-Step 3 6 "Fetching latest changes"
    try {
        git.exe fetch origin 2>&1 | Out-Null
    }
    catch {
        Write-ColorOutput "WARNING: Could not fetch from origin" "Yellow"
    }

    # Check current branch
    $currentBranch = git.exe rev-parse --abbrev-ref HEAD
    Write-ColorOutput "  Current branch: " "Gray" -NoNewline
    Write-ColorOutput $currentBranch "Cyan"

    if ($currentBranch -ne "main") {
        Write-Host ""
        Write-ColorOutput "WARNING: You are not on the 'main' branch!" "Yellow"
        Write-ColorOutput "Current branch: $currentBranch" "Yellow"
        Write-Host ""
        $continue = Read-Host "Do you want to continue anyway? (y/N)"
        if ($continue -ne "y" -and $continue -ne "Y") {
            Write-ColorOutput "Release cancelled." "Yellow"
            Read-Host "Press Enter to exit"
            exit 0
        }
    }
    else {
        Write-Step 4 6 "Updating main branch"
        try {
            git.exe pull origin main
        }
        catch {
            Write-ColorOutput "ERROR: Could not pull latest changes!" "Red"
            Write-ColorOutput "Please resolve conflicts manually." "Yellow"
            Read-Host "Press Enter to exit"
            exit 1
        }
    }

    # Check for uncommitted changes
    Write-Step 5 6 "Checking for uncommitted changes"
    $status = git.exe status --porcelain
    if ($status) {
        Write-ColorOutput "WARNING: You have uncommitted changes!" "Yellow"
        git.exe status --short
        Write-Host ""
        $continue = Read-Host "Do you want to continue anyway? (y/N)"
        if ($continue -ne "y" -and $continue -ne "Y") {
            Write-ColorOutput "Release cancelled." "Yellow"
            Write-ColorOutput "Please commit or stash your changes first." "Yellow"
            Read-Host "Press Enter to exit"
            exit 0
        }
    }

    # Clean up old release tags
    Write-Step 6 6 "Cleaning up old release tags"
    git.exe tag -d release 2>&1 | Out-Null
    git.exe push origin :refs/tags/release 2>&1 | Out-Null

    # Show release summary
    Write-Header "Release Summary"
    
    Write-Host "Repository: " -NoNewline
    Write-ColorOutput "RevEngine3r/whatsapp-proxy-go" "Cyan"
    
    Write-Host "Branch: " -NoNewline
    Write-ColorOutput $currentBranch "Cyan"
    
    $commitHash = git.exe log -1 --pretty=format:"%h"
    $commitMsg = git.exe log -1 --pretty=format:"%s"
    Write-Host "Latest commit: " -NoNewline
    Write-ColorOutput "$commitHash" "Cyan" -NoNewline
    Write-Host " - $commitMsg"
    
    # Calculate expected version
    $version = (Get-Date).ToUniversalTime().ToString("yyyy.MM.dd.HHmm")
    Write-Host "Expected version: " -NoNewline
    Write-ColorOutput "v$version" "Cyan"
    Write-Host ""

    # Confirm release
    $confirm = Read-Host "Are you sure you want to create a release? (y/N)"
    if ($confirm -ne "y" -and $confirm -ne "Y") {
        Write-ColorOutput "Release cancelled by user." "Yellow"
        Read-Host "Press Enter to exit"
        exit 0
    }

    # Create and push the release tag
    Write-Host ""
    Write-ColorOutput "Creating release tag..." "Green"
    git.exe tag release
    
    Write-ColorOutput "Pushing release tag to origin..." "Green"
    git.exe push origin release

    # Success message
    Write-Header "Release Tag Pushed Successfully!"
    
    Write-ColorOutput "What happens next:" "Cyan"
    Write-Host "1. GitHub Actions will start building binaries"
    Write-Host "2. Tests will run automatically"
    Write-Host "3. Multi-platform binaries will be created"
    Write-Host "4. A new release will be published"
    Write-Host ""
    
    Write-ColorOutput "Monitor progress at:" "Cyan"
    $actionsUrl = "https://github.com/RevEngine3r/whatsapp-proxy-go/actions"
    Write-ColorOutput $actionsUrl "Blue"
    Write-Host ""
    
    Write-ColorOutput "View releases at:" "Cyan"
    $releasesUrl = "https://github.com/RevEngine3r/whatsapp-proxy-go/releases"
    Write-ColorOutput $releasesUrl "Blue"
    Write-Host ""
    
    Write-ColorOutput "Note: After the workflow completes, cleanup with:" "Yellow"
    Write-Host "  .\scripts\cleanup-release-tag.ps1"
    Write-Host "  or"
    Write-Host "  git tag -d release"
    Write-Host "  git push origin :refs/tags/release"
    Write-Host ""

    # Ask if user wants to open browser
    $openBrowser = Read-Host "Open GitHub Actions in browser? (y/N)"
    if ($openBrowser -eq "y" -or $openBrowser -eq "Y") {
        Start-Process $actionsUrl
    }
}
catch {
    Write-ColorOutput "ERROR: An unexpected error occurred!" "Red"
    Write-ColorOutput $_.Exception.Message "Red"
    Write-ColorOutput $_.ScriptStackTrace "Gray"
}

Read-Host "Press Enter to exit"
