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
Write-ColorText "🚀 LazyTodo Windows Installer" $InfoColor
Write-ColorText "==============================" $InfoColor
Write-Host ""

# Check PowerShell version
if ($PSVersionTable.PSVersion.Major -lt 5) {
    Write-ColorText "❌ PowerShell 5.0 or higher is required. Current version: $($PSVersionTable.PSVersion)" $ErrorColor
    exit 1
}

# Create installation directory
Write-ColorText "📁 Creating installation directory..." $InfoColor
try {
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-ColorText "✅ Created directory: $InstallDir" $SuccessColor
    } else {
        Write-ColorText "📂 Directory already exists: $InstallDir" $WarningColor
    }
} catch {
    Write-ColorText "❌ Failed to create installation directory: $($_.Exception.Message)" $ErrorColor
    exit 1
}

# Download latest release
Write-ColorText "⬇️  Downloading LazyTodo..." $InfoColor
try {
    $apiUrl = "https://api.github.com/repos/yourusername/lazytodo/releases/latest"
    if ($Version -ne "latest") {
        $apiUrl = "https://api.github.com/repos/yourusername/lazytodo/releases/tags/$Version"
    }
    
    $release = Invoke-RestMethod -Uri $apiUrl -ErrorAction Stop
    $asset = $release.assets | Where-Object { $_.name -match "lazytodo.*windows.*\.zip" } | Select-Object -First 1
    
    if (!$asset) {
        Write-ColorText "❌ No Windows binary found in release $($release.tag_name)" $ErrorColor
        exit 1
    }
    
    $downloadUrl = $asset.browser_download_url
    $zipPath = Join-Path $env:TEMP "lazytodo-windows.zip"
    
    Write-ColorText "📦 Downloading: $($asset.name)" $InfoColor
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -ErrorAction Stop
    Write-ColorText "✅ Downloaded successfully" $SuccessColor
    
} catch {
    Write-ColorText "❌ Failed to download: $($_.Exception.Message)" $ErrorColor
    exit 1
}

# Extract and install
Write-ColorText "📦 Extracting files..." $InfoColor
try {
    Expand-Archive -Path $zipPath -DestinationPath $InstallDir -Force
    Remove-Item $zipPath -Force
    Write-ColorText "✅ Extraction completed" $SuccessColor
} catch {
    Write-ColorText "❌ Failed to extract: $($_.Exception.Message)" $ErrorColor
    exit 1
}

# Add to PATH
Write-ColorText "🔧 Configuring PATH..." $InfoColor
try {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        $newPath = if ($userPath) { "$userPath;$InstallDir" } else { $InstallDir }
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-ColorText "✅ Added $InstallDir to user PATH" $SuccessColor
        Write-ColorText "⚠️  Please restart your terminal to use 'lazytodo' command" $WarningColor
    } else {
        Write-ColorText "📂 Installation directory already in PATH" $InfoColor
    }
} catch {
    Write-ColorText "⚠️  Could not add to PATH. You may need to add manually: $InstallDir" $WarningColor
}

# Create desktop shortcut (optional)
$createShortcut = Read-Host "🖥️  Create desktop shortcut? (y/N)"
if ($createShortcut -eq "y" -or $createShortcut -eq "Y") {
    try {
        $WshShell = New-Object -comObject WScript.Shell
        $Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\LazyTodo.lnk")
        $Shortcut.TargetPath = Join-Path $InstallDir "lazytodo.exe"
        $Shortcut.WorkingDirectory = $InstallDir
        $Shortcut.Description = "LazyTodo - Beautiful Terminal Todo Manager"
        $Shortcut.Save()
        Write-ColorText "✅ Desktop shortcut created" $SuccessColor
    } catch {
        Write-ColorText "⚠️  Could not create desktop shortcut" $WarningColor
    }
}

# Installation complete
Write-Host ""
Write-ColorText "🎉 Installation Complete!" $SuccessColor
Write-ColorText "========================" $SuccessColor
Write-ColorText "Installation Directory: $InstallDir" $InfoColor
Write-ColorText "Executable: $(Join-Path $InstallDir 'lazytodo.exe')" $InfoColor
Write-Host ""
Write-ColorText "📚 Getting Started:" $InfoColor
Write-ColorText "  • Open a new terminal and run: lazytodo" $InfoColor
Write-ColorText "  • Or run directly: `"$InstallDir\lazytodo.exe`"" $InfoColor
Write-ColorText "  • Press 'n' to create your first todo list" $InfoColor
Write-ColorText "  • Press '?' for help" $InfoColor
Write-Host ""
Write-ColorText "🔗 Documentation: https://github.com/yourusername/lazytodo" $InfoColor
Write-Host ""

# Test installation
$testRun = Read-Host "🧪 Test installation now? (y/N)"
if ($testRun -eq "y" -or $testRun -eq "Y") {
    Write-ColorText "🧪 Testing installation..." $InfoColor
    try {
        $lazyTodoPath = Join-Path $InstallDir "lazytodo.exe"
        if (Test-Path $lazyTodoPath) {
            & $lazyTodoPath --version
            Write-ColorText "✅ Installation test successful!" $SuccessColor
        } else {
            Write-ColorText "❌ Binary not found at expected location" $ErrorColor
        }
    } catch {
        Write-ColorText "⚠️  Could not test installation: $($_.Exception.Message)" $WarningColor
    }
} 