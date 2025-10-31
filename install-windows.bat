@echo off
REM Windows installation script for jot CLI
REM Usage: Run this script as Administrator or in a directory in your PATH

echo Installing jot CLI for Windows...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Install jot
echo Installing jot from GitHub...
go install github.com/sk25469/jot@latest

if %errorlevel% neq 0 (
    echo Error: Failed to install jot
    echo Make sure you have internet connection and Go is properly configured
    pause
    exit /b 1
)

REM Check if GOPATH/bin is in PATH
echo Checking PATH configuration...
echo %PATH% | findstr /i gopath >nul
if %errorlevel% neq 0 (
    echo.
    echo WARNING: GOPATH\bin might not be in your PATH
    echo.
    echo To add it permanently:
    echo 1. Press Win+R, type sysdm.cpl, press Enter
    echo 2. Click "Environment Variables"
    echo 3. Under "User variables", find and edit "Path"
    echo 4. Add: %USERPROFILE%\go\bin
    echo.
    echo Or run this command in an elevated PowerShell:
    echo [Environment]::SetEnvironmentVariable("Path", $env:Path + ";%USERPROFILE%\go\bin", "User"^)
    echo.
)

REM Try to run jot to verify installation
echo Testing installation...
jot --help >nul 2>&1
if %errorlevel% equ 0 (
    echo.
    echo ✅ jot installed successfully!
    echo.
    echo Try these commands:
    echo   jot new "My first note" -t windows -t getting-started
    echo   jot list
    echo   jot stats
    echo.
    echo Default editor will be: notepad.exe
    echo To change editor: Edit %APPDATA%\jot\config.yaml
    echo.
) else (
    echo.
    echo ❌ Installation completed but jot command not found in PATH
    echo.
    echo Try running: %USERPROFILE%\go\bin\jot.exe --help
    echo.
    echo If that works, you need to add %USERPROFILE%\go\bin to your PATH
)

echo.
echo Configuration will be stored in: %APPDATA%\jot\
echo Notes will be stored in: %APPDATA%\jot\notes\
echo.

pause