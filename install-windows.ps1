# Windows PowerShell installation script for jot CLI
# Usage: Run as Administrator or ensure your user has Go properly configured
# Example: powershell -ExecutionPolicy Bypass -File install-windows.ps1

Write-Host "Installing jot CLI for Windows..." -ForegroundColor Green

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Found Go: $goVersion" -ForegroundColor Yellow
} catch {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

# Install jot
Write-Host "Installing jot from GitHub..." -ForegroundColor Yellow
try {
    go install github.com/sk25469/jot@latest
    Write-Host "Installation completed!" -ForegroundColor Green
} catch {
    Write-Host "Error: Failed to install jot" -ForegroundColor Red
    Write-Host "Make sure you have internet connection and Go is properly configured" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

# Check installation
$jotPath = "$env:USERPROFILE\go\bin\jot.exe"
if (Test-Path $jotPath) {
    Write-Host "✅ jot.exe found at: $jotPath" -ForegroundColor Green
} else {
    Write-Host "❌ jot.exe not found in expected location" -ForegroundColor Red
}

# Check if GOPATH/bin is in PATH
$goPath = "$env:USERPROFILE\go\bin"
if ($env:PATH -split ';' -contains $goPath) {
    Write-Host "✅ Go bin directory is in PATH" -ForegroundColor Green
} else {
    Write-Host "⚠️  Go bin directory not in PATH" -ForegroundColor Yellow
    Write-Host "Adding to PATH for current session..." -ForegroundColor Yellow
    $env:PATH += ";$goPath"
    
    Write-Host "To add permanently, run as Administrator:" -ForegroundColor Yellow
    Write-Host "[Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$goPath', 'User')" -ForegroundColor Cyan
}

# Test jot command
Write-Host "Testing jot command..." -ForegroundColor Yellow
try {
    jot --help | Out-Null
    Write-Host "✅ jot command works!" -ForegroundColor Green
    
    Write-Host ""
    Write-Host "🎉 Installation successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Try these commands:" -ForegroundColor Yellow
    Write-Host "  jot new `"My first note`" -t windows -t getting-started" -ForegroundColor Cyan
    Write-Host "  jot list" -ForegroundColor Cyan
    Write-Host "  jot stats" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Configuration: $env:APPDATA\jot\" -ForegroundColor Gray
    Write-Host "Notes location: $env:APPDATA\jot\notes\" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Default editor will be detected automatically (notepad.exe fallback)" -ForegroundColor Gray
    Write-Host "To customize: Edit $env:APPDATA\jot\config.yaml" -ForegroundColor Gray
    
} catch {
    Write-Host "❌ jot command not accessible" -ForegroundColor Red
    Write-Host "Try running: $jotPath --help" -ForegroundColor Yellow
    Write-Host "If that works, you need to restart your terminal or add Go bin to PATH" -ForegroundColor Yellow
}

Write-Host ""
Read-Host "Press Enter to exit"