# LazyTodo Windows Installation Script
# This script downloads and installs LazyTodo for Windows

param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:LOCALAPPDATA\LazyTodo"
)

# Colors for output
$ErrorColor = "Red"
$SuccessColor = "Green"
$InfoColor = "Cyan"
$WarningColor = "Yellow"

function Write-ColorText {
    param(
        [string]$Text,
        [string]$Color = "White"
    )
    Write-Host $Text -ForegroundColor $Color
}

function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Header
Write-ColorText "[INSTALL] LazyTodo Windows Installer" $InfoColor
Write-ColorText "=====================================" $InfoColor
Write-Host ""

# Check PowerShell version
if ($PSVersionTable.PSVersion.Major -lt 5) {
    Write-ColorText "[ERROR] PowerShell 5.0 or higher is required. Current version: $($PSVersionTable.PSVersion)" $ErrorColor
    exit 1
}

# Create installation directory
Write-ColorText "[INFO] Creating installation directory..." $InfoColor
try {
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-ColorText "[SUCCESS] Created directory: $InstallDir" $SuccessColor
    } else {
        Write-ColorText "[INFO] Directory already exists: $InstallDir" $WarningColor
    }
}
catch {
    Write-ColorText "[ERROR] Failed to create installation directory: $($_.Exception.Message)" $ErrorColor
    exit 1
}

# Download latest release
Write-ColorText "[INFO] Downloading LazyTodo..." $InfoColor
try {
    $apiUrl = "https://api.github.com/repos/DhirajZope/lazytodo/releases/latest"
    if ($Version -ne "latest") {
        $apiUrl = "https://api.github.com/repos/DhirajZope/lazytodo/releases/tags/$Version"
    }
    
    $release = Invoke-RestMethod -Uri $apiUrl -ErrorAction Stop
    $asset = $release.assets | Where-Object { $_.name -match "lazytodo.*windows.*\.zip" } | Select-Object -First 1
    
    if (!$asset) {
        Write-ColorText "[ERROR] No Windows binary found in release $($release.tag_name)" $ErrorColor
        exit 1
    }
    
    $downloadUrl = $asset.browser_download_url
    $zipPath = Join-Path $env:TEMP "lazytodo-windows.zip"
    
    Write-ColorText "[INFO] Downloading: $($asset.name)" $InfoColor
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -ErrorAction Stop
    Write-ColorText "[SUCCESS] Downloaded successfully" $SuccessColor
}
catch {
    Write-ColorText "[ERROR] Failed to download: $($_.Exception.Message)" $ErrorColor
    exit 1
}

# Extract and install
Write-ColorText "[INFO] Extracting files..." $InfoColor
try {
    Expand-Archive -Path $zipPath -DestinationPath $InstallDir -Force
    Remove-Item $zipPath -Force
    Write-ColorText "[SUCCESS] Extraction completed" $SuccessColor
}
catch {
    Write-ColorText "[ERROR] Failed to extract: $($_.Exception.Message)" $ErrorColor
    exit 1
}

# Add to PATH
Write-ColorText "[INFO] Configuring PATH..." $InfoColor
try {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        if ($userPath) {
            $newPath = "$userPath;$InstallDir"
        } else {
            $newPath = $InstallDir
        }
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-ColorText "[SUCCESS] Added $InstallDir to user PATH" $SuccessColor
        Write-ColorText "[WARNING] Please restart your terminal to use 'lazytodo' command" $WarningColor
    } else {
        Write-ColorText "[INFO] Installation directory already in PATH" $InfoColor
    }
}
catch {
    Write-ColorText "[WARNING] Could not add to PATH. You may need to add manually: $InstallDir" $WarningColor
}

# Create desktop shortcut (optional)
$createShortcut = Read-Host "[PROMPT] Create desktop shortcut? (y/N)"
if (($createShortcut -eq "y") -or ($createShortcut -eq "Y")) {
    try {
        $WshShell = New-Object -comObject WScript.Shell
        $Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\LazyTodo.lnk")
        $Shortcut.TargetPath = Join-Path $InstallDir "lazytodo.exe"
        $Shortcut.WorkingDirectory = $InstallDir
        $Shortcut.Description = "LazyTodo - Beautiful Terminal Todo Manager"
        $Shortcut.Save()
        Write-ColorText "[SUCCESS] Desktop shortcut created" $SuccessColor
    }
    catch {
        Write-ColorText "[WARNING] Could not create desktop shortcut" $WarningColor
    }
}

# Installation complete
Write-Host ""
Write-ColorText "[SUCCESS] Installation Complete!" $SuccessColor
Write-ColorText "=========================" $SuccessColor
Write-ColorText "Installation Directory: $InstallDir" $InfoColor
$exePath = Join-Path $InstallDir 'lazytodo.exe'
Write-ColorText "Executable: $exePath" $InfoColor
Write-Host ""
Write-ColorText "Getting Started:" $InfoColor
Write-ColorText "  * Open a new terminal and run: lazytodo" $InfoColor
Write-ColorText "  * Or run directly: `"$InstallDir\lazytodo.exe`"" $InfoColor
Write-ColorText "  * Press 'n' to create your first todo list" $InfoColor
Write-ColorText "  * Press '?' for help" $InfoColor
Write-Host ""
Write-ColorText "Documentation: https://github.com/DhirajZope/lazytodo" $InfoColor
Write-Host ""

# Test installation
$testRun = Read-Host "[PROMPT] Test installation now? (y/N)"
if (($testRun -eq "y") -or ($testRun -eq "Y")) {
    Write-ColorText "[INFO] Testing installation..." $InfoColor
    try {
        $lazyTodoPath = Join-Path $InstallDir "lazytodo.exe"
        if (Test-Path $lazyTodoPath) {
            & $lazyTodoPath --version
            Write-ColorText "[SUCCESS] Installation test successful!" $SuccessColor
        } else {
            Write-ColorText "[ERROR] Binary not found at expected location" $ErrorColor
        }
    }
    catch {
        Write-ColorText "[WARNING] Could not test installation: $($_.Exception.Message)" $WarningColor
    }
} 