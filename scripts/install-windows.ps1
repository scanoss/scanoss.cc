# SCANOSS Code Compare - Windows Installer
# This script installs SCANOSS Code Compare on Windows
#
# Parameters:
#   -InstallPath  : Installation directory (default: C:\Program Files\SCANOSS)
#   -NoPath       : Skip adding to system PATH
#   -NoShortcuts  : Skip creating Start Menu shortcuts
#   -NoVerify     : Skip SHA256 checksum verification (not recommended)
#   -Version      : Specific version to install (default: latest)
#
# Examples:
#   # Standard installation
#   .\install-windows.ps1
#
#   # Install without PATH modification
#   .\install-windows.ps1 -NoPath
#
#   # Install specific version
#   .\install-windows.ps1 -Version "0.9.0"
#
#   # Skip checksum verification (not recommended)
#   .\install-windows.ps1 -NoVerify

#Requires -Version 5.1

[CmdletBinding()]
param(
    [string]$InstallPath = "$env:ProgramFiles\SCANOSS",
    [switch]$NoPath,
    [switch]$NoShortcuts,
    [switch]$NoVerify,
    [string]$Version
)

# Configuration
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$Script:Repo = "scanoss/scanoss.cc"
$Script:AppName = "scanoss-cc"
$Script:ExeName = "$AppName.exe"

# Color output functions
function Write-Info {
    param([string]$Message)
    Write-Host "==> " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warn {
    param([string]$Message)
    Write-Host "Warning: " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-Err {
    param([string]$Message)
    Write-Host "Error: " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

function Write-Debug-Info {
    param([string]$Message)
    if ($DebugPreference -ne "SilentlyContinue") {
        Write-Host "Debug: " -ForegroundColor Blue -NoNewline
        Write-Host $Message
    }
}

# Check if running as Administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Get latest version from GitHub
function Get-LatestVersion {
    Write-Info "Fetching latest version..."

    try {
        $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$Script:Repo/releases/latest"
        $version = $releases.tag_name -replace '^v', ''
        Write-Info "Latest version: v$version"
        return $version
    }
    catch {
        Write-Err "Failed to fetch latest version: $_"
        throw
    }
}

# Download file with progress
function Get-FileWithProgress {
    param(
        [string]$Url,
        [string]$OutputPath
    )

    Write-Info "Downloading: $Url"

    try {
        $webClient = New-Object System.Net.WebClient
        $webClient.DownloadFile($Url, $OutputPath)
        Write-Info "Download complete"
    }
    catch {
        Write-Err "Download failed: $_"
        throw
    }
}

# Verify SHA256 checksum
function Test-Checksum {
    param(
        [string]$FilePath,
        [string]$ExpectedHash
    )

    Write-Info "Verifying checksum..."

    $actualHash = (Get-FileHash -Path $FilePath -Algorithm SHA256).Hash

    if ($actualHash -eq $ExpectedHash) {
        Write-Info "Checksum verified: ✓"
        return $true
    }
    else {
        Write-Err "Checksum mismatch!"
        Write-Err "Expected: $ExpectedHash"
        Write-Err "Got:      $actualHash"
        return $false
    }
}

# Get asset download URL
function Get-AssetUrl {
    param([string]$Version)

    $releaseUrl = "https://api.github.com/repos/$Script:Repo/releases/tags/v$Version"
    $assetPattern = "$Script:AppName-win.zip"

    try {
        $release = Invoke-RestMethod -Uri $releaseUrl
        $asset = $release.assets | Where-Object { $_.name -like "*$assetPattern*" } | Select-Object -First 1

        if (-not $asset) {
            throw "Asset not found: $assetPattern"
        }

        return $asset.browser_download_url
    }
    catch {
        Write-Err "Failed to get asset URL: $_"
        throw
    }
}

# Get checksum for a specific asset from SHA256SUMS file
function Get-AssetChecksum {
    param(
        [string]$Version,
        [string]$AssetName
    )

    Write-Debug-Info "Fetching checksum for: $AssetName"

    $checksumsUrl = "https://github.com/$Script:Repo/releases/download/v$Version/SHA256SUMS"

    try {
        $checksums = (Invoke-WebRequest -Uri $checksumsUrl -UseBasicParsing -ErrorAction SilentlyContinue).Content

        if ($checksums) {
            # Parse the SHA256SUMS file (format: "hash  filename")
            $lines = $checksums -split "`n"
            foreach ($line in $lines) {
                if ($line -match "^\s*([a-fA-F0-9]+)\s+.*$AssetName") {
                    $checksum = $matches[1]
                    Write-Debug-Info "Found checksum: $checksum"
                    return $checksum
                }
            }
        }

        Write-Warn "No checksum file found for this release"
        return $null
    }
    catch {
        Write-Debug-Info "Failed to fetch checksums: $_"
        return $null
    }
}

# Check if WebView2 is installed
function Test-WebView2 {
    $webView2Key = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"

    if (Test-Path $webView2Key) {
        $version = (Get-ItemProperty -Path $webView2Key -Name "pv" -ErrorAction SilentlyContinue).pv
        if ($version) {
            Write-Info "WebView2 is installed (version: $version)"
            return $true
        }
    }

    return $false
}

# Install WebView2
function Install-WebView2 {
    Write-Info "WebView2 Runtime is not installed."
    Write-Info "SCANOSS Code Compare requires WebView2 to function."
    Write-Host ""

    $response = Read-Host "Download and install WebView2 Runtime? (Y/n)"
    if ($response -match '^[Nn]') {
        Write-Warn "Skipping WebView2 installation"
        Write-Warn "The application may not work without it"
        return
    }

    Write-Info "Downloading WebView2 Runtime..."
    $webView2Url = "https://go.microsoft.com/fwlink/p/?LinkId=2124703"
    $webView2Installer = "$env:TEMP\MicrosoftEdgeWebview2Setup.exe"

    try {
        (New-Object System.Net.WebClient).DownloadFile($webView2Url, $webView2Installer)

        Write-Info "Installing WebView2 Runtime..."
        Start-Process -FilePath $webView2Installer -ArgumentList "/silent", "/install" -Wait

        Remove-Item $webView2Installer -Force
        Write-Info "✓ WebView2 Runtime installed"
    }
    catch {
        Write-Err "Failed to install WebView2: $_"
        Write-Info "You can manually download it from: https://developer.microsoft.com/en-us/microsoft-edge/webview2/"
    }
}

# Add to system PATH
function Add-ToPath {
    param([string]$Path)

    Write-Info "Adding to system PATH..."

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")

    if ($currentPath -like "*$Path*") {
        Write-Info "Already in PATH"
        return
    }

    $newPath = "$currentPath;$Path"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")

    # Also update current session
    $env:Path = "$env:Path;$Path"

    Write-Info "✓ Added to PATH"
}

# Create Start Menu shortcut
function New-StartMenuShortcut {
    param(
        [string]$TargetPath,
        [string]$ShortcutName
    )

    Write-Info "Creating Start Menu shortcut..."

    $startMenuPath = "$env:ProgramData\Microsoft\Windows\Start Menu\Programs"
    $shortcutPath = "$startMenuPath\$ShortcutName.lnk"

    $WScriptShell = New-Object -ComObject WScript.Shell
    $shortcut = $WScriptShell.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = $TargetPath
    $shortcut.WorkingDirectory = [System.IO.Path]::GetDirectoryName($TargetPath)
    $shortcut.Description = "SCANOSS Code Compare - Open source code analysis"
    $shortcut.Save()

    Write-Info "✓ Start Menu shortcut created"
}

# Main installation function
function Install-ScanossCodeCompare {
    Write-Host ""
    Write-Info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Info "  SCANOSS Code Compare - Windows Installation"
    Write-Info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    Write-Host ""

    # Check administrator rights
    if (-not (Test-Administrator)) {
        Write-Err "This installation requires administrator privileges."
        Write-Host ""
        Write-Host "Please run PowerShell as Administrator and try again:"
        Write-Host "  1. Right-click PowerShell"
        Write-Host "  2. Select 'Run as Administrator'"
        Write-Host "  3. Run the installation command again"
        Write-Host ""
        exit 1
    }

    # Create installation directory
    if (-not (Test-Path $InstallPath)) {
        Write-Info "Creating installation directory: $InstallPath"
        New-Item -Path $InstallPath -ItemType Directory -Force | Out-Null
    }

    # Create temporary directory
    $tempDir = New-Item -Path "$env:TEMP\scanoss-install-$(Get-Random)" -ItemType Directory -Force

    try {
        # Get version
        if (-not $Version) {
            $Version = Get-LatestVersion
        }

        # Download
        $downloadUrl = Get-AssetUrl -Version $Version
        $zipPath = "$tempDir\$Script:AppName.zip"
        $assetName = "$Script:AppName-win.zip"

        Get-FileWithProgress -Url $downloadUrl -OutputPath $zipPath

        # Verify checksum (unless -NoVerify is set)
        if (-not $NoVerify) {
            Write-Host ""
            $expectedChecksum = Get-AssetChecksum -Version $Version -AssetName $assetName

            if ($expectedChecksum) {
                if (-not (Test-Checksum -FilePath $zipPath -ExpectedHash $expectedChecksum)) {
                    # Checksum failed - clean up and abort
                    Remove-Item -Path $zipPath -Force -ErrorAction SilentlyContinue
                    throw "Checksum verification failed! The downloaded file may be corrupted or tampered with."
                }
            }
            else {
                Write-Warn "No checksum available for verification"
                Write-Warn "Installation will proceed without checksum verification"
            }
        }
        else {
            Write-Warn "Checksum verification skipped (--NoVerify flag set)"
        }

        # Extract
        Write-Host ""
        Write-Info "Extracting..."
        Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force

        # Find the executable
        $exePath = Get-ChildItem -Path $tempDir -Filter $Script:ExeName -Recurse | Select-Object -First 1

        if (-not $exePath) {
            throw "Executable not found in archive: $Script:ExeName"
        }

        # Copy to installation directory
        Write-Info "Installing to $InstallPath..."
        Copy-Item -Path $exePath.FullName -Destination "$InstallPath\$Script:ExeName" -Force

        Write-Host ""
        Write-Info "✓ Binary installed successfully"

        # Add to PATH
        if (-not $NoPath) {
            Write-Host ""
            Add-ToPath -Path $InstallPath
        }

        # Create shortcuts
        if (-not $NoShortcuts) {
            Write-Host ""
            New-StartMenuShortcut -TargetPath "$InstallPath\$Script:ExeName" -ShortcutName "SCANOSS Code Compare"
        }

        # Check WebView2
        Write-Host ""
        if (-not (Test-WebView2)) {
            Install-WebView2
        }

        # Show completion message
        Write-Host ""
        Write-Info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        Write-Info "✨ SCANOSS Code Compare has been successfully installed!"
        Write-Info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        Write-Host ""
        Write-Info "Verify installation:"
        Write-Host "    $Script:AppName --version"
        Write-Host ""
        Write-Info "Get started:"
        Write-Host "    $Script:AppName --help"
        Write-Host ""
        Write-Info "Or launch from Start Menu: SCANOSS Code Compare"
        Write-Host ""
        Write-Info "Documentation:"
        Write-Host "    https://github.com/$Script:Repo"
        Write-Host ""

        # Verify installation
        Write-Info "Verifying installation..."
        $installedPath = "$InstallPath\$Script:ExeName"
        if (Test-Path $installedPath) {
            try {
                & $installedPath --version
                Write-Host ""
                Write-Info "Installation verified successfully!"
            }
            catch {
                Write-Warn "Could not get version (but binary is installed)"
            }
        }
    }
    catch {
        Write-Err "Installation failed: $_"
        Write-Host ""
        exit 1
    }
    finally {
        # Cleanup
        if (Test-Path $tempDir) {
            Remove-Item -Path $tempDir -Recurse -Force
        }
    }
}

# Run installation
Install-ScanossCodeCompare
