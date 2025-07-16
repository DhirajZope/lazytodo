#!/usr/bin/env pwsh

# LazyTodo Build Script

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("build", "run", "clean", "install")]
    [string]$Action = "build",
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("windows", "linux", "darwin")]
    [string]$OS = "",
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("amd64", "386", "arm64")]
    [string]$Arch = "amd64"
)

Write-Host "🎯 LazyTodo Build Script" -ForegroundColor Magenta
Write-Host "========================" -ForegroundColor Magenta

switch ($Action) {
    "build" {
        Write-Host "🔨 Building LazyTodo..." -ForegroundColor Yellow
        
        if ($OS -eq "") {
            # Build for current platform
            Write-Host "Building for current platform..." -ForegroundColor Green
            go build -o lazytodo.exe cmd/main.go
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ Build successful! Executable: lazytodo.exe" -ForegroundColor Green
            } else {
                Write-Host "❌ Build failed!" -ForegroundColor Red
                exit 1
            }
        } else {
            # Cross-platform build
            $outputName = "lazytodo-$OS-$Arch"
            if ($OS -eq "windows") {
                $outputName += ".exe"
            }
            
            Write-Host "Building for $OS/$Arch..." -ForegroundColor Green
            $env:GOOS = $OS
            $env:GOARCH = $Arch
            go build -o $outputName cmd/main.go
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ Cross-platform build successful! Executable: $outputName" -ForegroundColor Green
            } else {
                Write-Host "❌ Cross-platform build failed!" -ForegroundColor Red
                exit 1
            }
        }
    }
    
    "run" {
        Write-Host "🚀 Running LazyTodo..." -ForegroundColor Yellow
        go run cmd/main.go
    }
    
    "clean" {
        Write-Host "🧹 Cleaning build artifacts..." -ForegroundColor Yellow
        Remove-Item -Path "lazytodo*.exe" -ErrorAction SilentlyContinue
        Remove-Item -Path "lazytodo-*" -ErrorAction SilentlyContinue
        Write-Host "✅ Clean completed!" -ForegroundColor Green
    }
    
    "install" {
        Write-Host "📦 Installing dependencies..." -ForegroundColor Yellow
        go mod tidy
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Dependencies installed successfully!" -ForegroundColor Green
        } else {
            Write-Host "❌ Failed to install dependencies!" -ForegroundColor Red
            exit 1
        }
    }
}

Write-Host ""
Write-Host "Usage examples:" -ForegroundColor Cyan
Write-Host "  .\build.ps1 build          # Build for current platform" -ForegroundColor White
Write-Host "  .\build.ps1 run            # Run the application" -ForegroundColor White
Write-Host "  .\build.ps1 clean          # Clean build artifacts" -ForegroundColor White
Write-Host "  .\build.ps1 install        # Install dependencies" -ForegroundColor White
Write-Host "  .\build.ps1 build linux    # Cross-compile for Linux" -ForegroundColor White
Write-Host "  .\build.ps1 build darwin   # Cross-compile for macOS" -ForegroundColor White 