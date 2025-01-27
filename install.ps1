# Ensure we're running with elevated privileges
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "Please run this script as Administrator!"
    exit 1
}

# Configuration
$owner = "voioo"
$repo = "td"
$installDir = "$env:LOCALAPPDATA\Programs\td"
$pathType = "User" # Can be "User" or "Machine"

# Determine architecture
$arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ([System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Error "32-bit systems are not supported"
    exit 1
}

Write-Host "Installing td for Windows ($arch)..." -ForegroundColor Cyan

# Cleanup function for failed installations
function Cleanup {
    if (Test-Path $zipPath) {
        Remove-Item $zipPath -Force
    }
}

# Check for existing installation
if (Get-Command td -ErrorAction SilentlyContinue) {
    $currentVersion = td --version
    Write-Host "Found existing installation: $currentVersion"
    $confirmation = Read-Host "Do you want to continue with the installation? (y/N)"
    if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
        Write-Host "Installation cancelled" -ForegroundColor Yellow
        exit 0
    }
}

# Create install directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Get latest release info from GitHub
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$owner/$repo/releases/latest"
    $asset = $release.assets | Where-Object { $_.name -like "*windows_${arch}.zip" }
    if (-not $asset) {
        Write-Error "No release found for Windows $arch"
        exit 1
    }
    $version = $release.tag_name
} catch {
    Write-Error "Failed to get latest release information: $_"
    exit 1
}

# Download the latest release
$zipPath = Join-Path $installDir "td.zip"
try {
    Write-Host "Downloading version $version..."
    Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $zipPath
} catch {
    Write-Error "Failed to download release: $_"
    Cleanup
    exit 1
}

# Extract the archive
try {
    Write-Host "Extracting files..."
    Expand-Archive -Path $zipPath -DestinationPath $installDir -Force
    Remove-Item $zipPath
} catch {
    Write-Error "Failed to extract archive: $_"
    Cleanup
    exit 1
}

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", $pathType)
if (-not $currentPath.Contains($installDir)) {
    try {
        Write-Host "Adding to PATH..."
        $newPath = "$currentPath;$installDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, $pathType)
        $env:Path = $newPath
    } catch {
        Write-Error "Failed to add to PATH: $_"
        exit 1
    }
}

Write-Host "`nSuccessfully installed td $version!" -ForegroundColor Green
Write-Host "Please restart your terminal to use 'td' command." -ForegroundColor Yellow 