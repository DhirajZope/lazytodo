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
        
        # Refresh PATH for current session
        $env:PATH = [Environment]::GetEnvironmentVariable("Path", "User")
        Write-ColorText "[INFO] PATH refreshed for current session" $InfoColor
    } else {
        Write-ColorText "[INFO] Installation directory already in PATH" $InfoColor
    }
}
catch {
    Write-ColorText "[WARNING] Could not add to PATH. You may need to add manually: $InstallDir" $WarningColor
}

# Verify PATH configuration
Write-ColorText "[INFO] Verifying PATH configuration..." $InfoColor
try {
    $pathCheck = Get-Command "lazytodo" -ErrorAction SilentlyContinue
    if ($pathCheck) {
        Write-ColorText "[SUCCESS] lazytodo command is available in PATH" $SuccessColor
    } else {
        Write-ColorText "[WARNING] lazytodo not found in PATH. You may need to restart your terminal" $WarningColor
    }
}
catch {
    Write-ColorText "[WARNING] Could not verify PATH configuration" $WarningColor
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
            # Test version command
            Write-ColorText "[INFO] Testing version command..." $InfoColor
            & $lazyTodoPath --version
            
            # Test database initialization
            Write-ColorText "[INFO] Testing database initialization..." $InfoColor
            $infoOutput = & $lazyTodoPath --info 2>&1
            
            if ($LASTEXITCODE -eq 0) {
                Write-ColorText "[SUCCESS] Database initialization successful!" $SuccessColor
                Write-ColorText "[INFO] Database info:" $InfoColor
                $infoOutput | ForEach-Object { Write-ColorText "  $($_)" $InfoColor }
            } else {
                Write-ColorText "[WARNING] Database initialization failed. This might be a first-run issue." $WarningColor
                Write-ColorText "[INFO] Common solutions:" $InfoColor
                Write-ColorText "  * If you see migration errors, delete the database file:" $InfoColor
                Write-ColorText "    Remove-Item `"$env:USERPROFILE\.lazytodo\lazytodo.db`" -Force" $InfoColor
                Write-ColorText "  * Then run: lazytodo --info" $InfoColor
            }
            
            Write-ColorText "[SUCCESS] Installation test completed!" $SuccessColor
        } else {
            Write-ColorText "[ERROR] Binary not found at expected location" $ErrorColor
        }
    }
    catch {
        Write-ColorText "[WARNING] Could not complete installation test: $($_.Exception.Message)" $WarningColor
        Write-ColorText "[INFO] This is usually not critical - the installation may still work" $WarningColor
    }
}

# Final instructions
Write-Host ""
Write-ColorText "[INFO] Installation Notes:" $InfoColor
Write-ColorText "  * If you encounter database migration errors on first run:" $WarningColor
Write-ColorText "    Remove-Item `"$env:USERPROFILE\.lazytodo\lazytodo.db`" -Force" $InfoColor
Write-ColorText "  * If 'lazytodo' command not found, restart your terminal" $WarningColor
Write-ColorText "  * Or run directly: `"$InstallDir\lazytodo.exe`"" $InfoColor 