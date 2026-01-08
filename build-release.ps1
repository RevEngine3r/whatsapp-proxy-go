<#
.SYNOPSIS
    Build WhatsApp Proxy Go release binaries with automatic datetime-based versioning.

.DESCRIPTION
    This script automates the Go release build process:
    - Generates version from current datetime (YYYY.MM.DD.HHmm)
    - Builds binaries for multiple OS/architectures
    - Creates compressed archives (tar.gz for Unix, zip for Windows)
    - Generates SHA256 checksums
    - Outputs to dist/ directory

.PARAMETER CustomVersion
    Optional custom version string

.PARAMETER SkipTests
    Skip running tests before building

.EXAMPLE
    .\build-release.ps1
    Builds release with automatic datetime version

.EXAMPLE
    .\build-release.ps1 -CustomVersion "1.2.3"
    Builds release with custom version

.NOTES
    Requires: Go 1.21+, PowerShell 5.1+
    Author: RevEngine3r
    Project: WhatsApp Proxy Go
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $false)]
    [string]$CustomVersion = "",
    
    [Parameter(Mandatory = $false)]
    [switch]$SkipTests
)

# Script configuration
$ErrorActionPreference = "Stop"
$Script:ProjectName = "whatsapp-proxy"
$Script:MainPackage = "./cmd/whatsapp-proxy"
$Script:DistDir = "dist"

# Build targets (GOOS/GOARCH/GOARM)
$Script:BuildTargets = @(
    @{OS="linux"; Arch="amd64"; Arm=""},
    @{OS="linux"; Arch="arm64"; Arm=""},
    @{OS="linux"; Arch="386"; Arm=""},
    @{OS="linux"; Arch="arm"; Arm="7"},
    @{OS="darwin"; Arch="amd64"; Arm=""},
    @{OS="darwin"; Arch="arm64"; Arm=""},
    @{OS="windows"; Arch="amd64"; Arm=""},
    @{OS="windows"; Arch="386"; Arm=""},
    @{OS="windows"; Arch="arm64"; Arm=""}
)

#region Helper Functions

function Write-Step {
    param([string]$Message)
    Write-Host "`n[STEP] $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "  ✓ $Message" -ForegroundColor Green
}

function Write-Info {
    param([string]$Message)
    Write-Host "  → $Message" -ForegroundColor Gray
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "  ✗ $Message" -ForegroundColor Red
}

function Get-DateTimeVersion {
    $now = Get-Date
    return $now.ToString("yyyy.MM.dd.HHmm")
}

function Get-GitCommitHash {
    try {
        $hash = git rev-parse --short HEAD 2>$null
        return $hash
    }
    catch {
        return "unknown"
    }
}

function Test-GoInstalled {
    try {
        $goVersion = go version 2>$null
        return $true
    }
    catch {
        return $false
    }
}

function Build-Binary {
    param(
        [string]$OS,
        [string]$Arch,
        [string]$Arm,
        [string]$Version,
        [string]$CommitHash
    )
    
    $outputName = "$Script:ProjectName-$Version-$OS-$Arch"
    if ($Arm) {
        $outputName += "v$Arm"
    }
    
    $ext = ""
    if ($OS -eq "windows") {
        $ext = ".exe"
    }
    
    $outputPath = Join-Path $Script:DistDir "$outputName$ext"
    
    Write-Info "Building $OS/$Arch$(if ($Arm) {"v$Arm"})..."
    
    # Set environment variables
    $env:GOOS = $OS
    $env:GOARCH = $Arch
    $env:CGO_ENABLED = "0"
    
    if ($Arm) {
        $env:GOARM = $Arm
    }
    else {
        Remove-Item Env:\GOARM -ErrorAction SilentlyContinue
    }
    
    # Build flags
    $ldflags = "-s -w"
    $ldflags += " -X main.version=$Version"
    $ldflags += " -X main.buildTime=$((Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ'))"
    $ldflags += " -X main.gitCommit=$CommitHash"
    
    # Build command
    $buildArgs = @(
        "build",
        "-trimpath",
        "-ldflags=`"$ldflags`"",
        "-o", $outputPath,
        $Script:MainPackage
    )
    
    try {
        & go @buildArgs
        if ($LASTEXITCODE -ne 0) {
            throw "Build failed with exit code $LASTEXITCODE"
        }
    }
    catch {
        throw "Failed to build $OS/$Arch: $_"
    }
    
    # Create archive
    if ($OS -eq "windows") {
        # Create ZIP for Windows
        $zipName = "$outputName.zip"
        $zipPath = Join-Path $Script:DistDir $zipName
        
        Compress-Archive -Path $outputPath -DestinationPath $zipPath -Force
        Remove-Item $outputPath
        
        $size = [math]::Round((Get-Item $zipPath).Length / 1MB, 2)
        Write-Success "Created $zipName ($size MB)"
    }
    else {
        # Create tar.gz for Unix
        $tarName = "$outputName.tar.gz"
        $tarPath = Join-Path $Script:DistDir $tarName
        
        # Use tar if available (Git for Windows includes it)
        try {
            $null = & tar --version 2>$null
            & tar -czf $tarPath -C $Script:DistDir (Split-Path $outputPath -Leaf)
            Remove-Item $outputPath
        }
        catch {
            # Fallback to keeping uncompressed
            Write-Info "tar not available, keeping uncompressed binary"
            $tarPath = $outputPath
        }
        
        $size = [math]::Round((Get-Item $tarPath).Length / 1MB, 2)
        Write-Success "Created $(Split-Path $tarPath -Leaf) ($size MB)"
    }
}

function New-Checksums {
    Write-Info "Generating checksums..."
    
    $checksumFile = Join-Path $Script:DistDir "checksums.txt"
    $checksums = @()
    
    Get-ChildItem -Path $Script:DistDir -File | Where-Object { $_.Name -ne "checksums.txt" } | ForEach-Object {
        $hash = (Get-FileHash -Path $_.FullName -Algorithm SHA256).Hash.ToLower()
        $checksums += "$hash  $($_.Name)"
    }
    
    $checksums | Out-File -FilePath $checksumFile -Encoding UTF8
    Write-Success "Checksums saved to checksums.txt"
}

#endregion

#region Main Script

function Main {
    Write-Host "`n" + "=" * 60 -ForegroundColor Cyan
    Write-Host "  WhatsApp Proxy Go - Release Builder" -ForegroundColor Cyan
    Write-Host "=" * 60 -ForegroundColor Cyan
    
    try {
        # Step 1: Validate environment
        Write-Step "Validating environment"
        
        if (-not (Test-GoInstalled)) {
            throw "Go is not installed or not in PATH"
        }
        
        $goVersion = go version
        Write-Success "Go found: $goVersion"
        
        # Check if we're in the right directory
        if (-not (Test-Path "go.mod")) {
            throw "go.mod not found - please run from project root"
        }
        Write-Success "Found go.mod"
        
        # Step 2: Determine version
        Write-Step "Determining version"
        
        if ($CustomVersion) {
            $version = $CustomVersion
            Write-Info "Using custom version: $version"
        }
        else {
            $version = Get-DateTimeVersion
            Write-Info "Generated datetime version: $version"
        }
        
        $commitHash = Get-GitCommitHash
        Write-Info "Git commit: $commitHash"
        
        # Step 3: Run tests
        if (-not $SkipTests) {
            Write-Step "Running tests"
            
            & go test -race -coverprofile=coverage.txt -covermode=atomic ./...
            if ($LASTEXITCODE -ne 0) {
                throw "Tests failed"
            }
            Write-Success "All tests passed"
        }
        else {
            Write-Step "Skipping tests"
        }
        
        # Step 4: Clean and create dist directory
        Write-Step "Preparing build directory"
        
        if (Test-Path $Script:DistDir) {
            Remove-Item -Path $Script:DistDir -Recurse -Force
            Write-Info "Cleaned existing dist directory"
        }
        
        New-Item -ItemType Directory -Path $Script:DistDir -Force | Out-Null
        Write-Success "Created dist directory"
        
        # Step 5: Build binaries
        Write-Step "Building binaries for $($Script:BuildTargets.Count) platforms"
        
        foreach ($target in $Script:BuildTargets) {
            Build-Binary -OS $target.OS -Arch $target.Arch -Arm $target.Arm -Version $version -CommitHash $commitHash
        }
        
        # Step 6: Generate checksums
        Write-Step "Generating checksums"
        New-Checksums
        
        # Success summary
        Write-Host "
" + "=" * 60 -ForegroundColor Green
        Write-Host "  BUILD SUCCESSFUL" -ForegroundColor Green
        Write-Host "=" * 60 -ForegroundColor Green
        Write-Host ""
        Write-Host "  Version:  $version" -ForegroundColor White
        Write-Host "  Commit:   $commitHash" -ForegroundColor White
        Write-Host "  Output:   $Script:DistDir/" -ForegroundColor White
        Write-Host "  Files:    $($Script:BuildTargets.Count) binaries + checksums" -ForegroundColor White
        Write-Host ""
        
        # List generated files
        Write-Host "  Generated files:" -ForegroundColor Yellow
        Get-ChildItem -Path $Script:DistDir -File | ForEach-Object {
            $size = [math]::Round($_.Length / 1MB, 2)
            Write-Host "    - $($_.Name) ($size MB)" -ForegroundColor Gray
        }
        Write-Host ""
        
        return 0
    }
    catch {
        Write-Host ""
        Write-Host "=" * 60 -ForegroundColor Red
        Write-Host "  BUILD FAILED" -ForegroundColor Red
        Write-Host "=" * 60 -ForegroundColor Red
        Write-ErrorMsg $_.Exception.Message
        Write-Host ""
        return 1
    }
}

# Execute main function
exit (Main)

#endregion
