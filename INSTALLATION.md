# SCANOSS Code Compare - Comprehensive Installation Guide

This guide provides detailed installation instructions for all supported platforms and distribution methods.

## Quick Install

The fastest way to install SCANOSS Code Compare is with our installation scripts:

### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install.sh | bash
```

### Windows (PowerShell as Administrator)
```powershell
irm https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install-windows.ps1 | iex
```

### Security-Conscious Installation

If you prefer to review the installation script before running:

```bash
# Download the script
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install.sh -o install.sh

# Review it
cat install.sh

# Run it
bash install.sh
```

The installation scripts automatically:
- ✅ Download the latest version
- ✅ Verify SHA256 checksums
- ✅ Install dependencies (Linux)
- ✅ Set up PATH configuration
- ✅ Create shortcuts/desktop entries

---

## Table of Contents

- [Quick Install](#quick-install) ⬆️
- [macOS](#macos)
  - [Automated Installation Script](#automated-installation-script)
  - [Homebrew](#homebrew-recommended)
  - [DMG](#dmg-manual-installation)
- [Windows](#windows)
  - [Automated Installation Script](#automated-installation-script-1)
  - [Manual Installation](#manual-installation)
- [Linux](#linux)
  - [Automated Installation Script](#automated-installation-script-2)
  - [Manual Installation](#manual-installation-1)
- [Troubleshooting](#troubleshooting)

---

## macOS

### Automated Installation Script

The installation script provides an interactive menu to choose between Homebrew (recommended) or direct DMG installation:

```bash
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install-macos.sh | bash
```

**Options:**
- **Homebrew installation** - Installs via Homebrew (will install Homebrew first if needed)
- **Direct installation** - Downloads DMG, installs to /Applications, creates CLI symlink

The script handles everything automatically, including:
- Installing Homebrew (if chosen and not present)
- Downloading and verifying the DMG
- Installing the full .app bundle to /Applications
- Creating `/usr/local/bin/scanoss-cc` symlink for CLI access
- Clearing macOS quarantine attributes

### Homebrew (Recommended)

Homebrew is the easiest way to install SCANOSS Code Compare on macOS. It automatically handles updates and PATH configuration.

```bash
# Add the SCANOSS tap (first time only)
brew tap scanoss/dist

# Install SCANOSS Code Compare
brew install scanoss-code-compare
```

**Benefits:**
- ✅ Automatic PATH configuration
- ✅ Easy updates with `brew upgrade`
- ✅ Clean uninstallation with `brew uninstall`

**Verify Installation:**
```bash
scanoss-cc --version
```

### DMG (Manual Installation)

The DMG provides a traditional drag-and-drop installation for macOS.

1. Download `scanoss-cc-mac.zip` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP to get the DMG file
3. Open the DMG file
4. Drag "scanoss-cc.app" to the Applications folder
5. Eject the DMG

**Setting Up CLI Access:**

After installing via DMG, you need to manually set up CLI access. Choose one of these methods:

**Option 1: Symlink (Recommended)**
```bash
sudo ln -s "/Applications/scanoss-cc.app/Contents/MacOS/scanoss-cc" /usr/local/bin/scanoss-cc
```

**Option 2: Add to PATH**
```bash
echo 'export PATH="/Applications/scanoss-cc.app/Contents/MacOS:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Note:** The full .app bundle must be in /Applications for GUI features (like folder picker dialogs) to work properly due to macOS bundle path resolution.

For more details, see [INSTALL_MACOS.md](INSTALL_MACOS.md).

---

## Windows

### Automated Installation Script

The PowerShell installation script handles everything automatically:

```powershell
# Run PowerShell as Administrator, then:
irm https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install-windows.ps1 | iex
```

The script automatically:
- Downloads the latest version
- Verifies SHA256 checksum
- Installs to `C:\Program Files\SCANOSS`
- Adds to system PATH
- Creates Start Menu shortcut
- Checks/installs WebView2 Runtime if needed

### Manual Installation

If you prefer manual installation:

1. Download `scanoss-cc-win.zip` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP file
3. Run `scanoss-cc.exe`

**Installation Options:**

**Option 1: Portable Use**
- Run directly from Downloads or any location
- No installation required
- Ideal for USB drives or temporary usage

**Option 2: Install to Program Files (Recommended)**
```powershell
# In PowerShell (Run as Administrator)
New-Item -Path "C:\Program Files\SCANOSS" -ItemType Directory -Force
Move-Item -Path ".\scanoss-cc.exe" -Destination "C:\Program Files\SCANOSS\" -Force
```

**Add to PATH (Optional - for CLI access):**
```powershell
# In PowerShell (Run as Administrator)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\SCANOSS", "Machine")
```

After adding to PATH, open a new terminal and verify:
```powershell
scanoss-cc --version
```

**Create Desktop Shortcut (Optional):**
- Right-click `scanoss-cc-windows.exe` → Send to → Desktop (create shortcut)

**Uninstallation:**
- Simply delete the executable and any shortcuts you created

---

## Linux

### Automated Installation Script

The Linux installation script automatically detects your distribution and installs dependencies:

```bash
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install-linux.sh | bash
```

The script automatically:
- Detects your Linux distribution
- Installs required dependencies (GTK3, WebKit2GTK)
- Downloads and verifies the latest version
- Installs to `/usr/local/bin/scanoss-cc`
- Optionally creates a desktop application entry

**Supported distributions:**
- Debian/Ubuntu (apt)
- Fedora/RHEL/CentOS (dnf)
- Arch/Manjaro (pacman)
- openSUSE (zypper)

### Manual Installation

#### All Distributions

1. Download `scanoss-cc-linux.zip` from [releases](https://github.com/scanoss/scanoss.cc/releases)
2. Extract and install:

```bash
unzip scanoss-cc-linux.zip
chmod +x scanoss-cc-linux
sudo mv scanoss-cc-linux /usr/local/bin/scanoss-cc
```

3. Install dependencies based on your distribution:

**Debian/Ubuntu:**
```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
```

**Fedora/RHEL:**
```bash
sudo dnf install gtk3 webkit2gtk4.0
```

**Arch Linux:**
```bash
sudo pacman -S gtk3 webkit2gtk
```

**openSUSE:**
```bash
sudo zypper install gtk3 webkit2gtk3
```

#### Optional: Desktop Entry

Create a desktop application entry for launching from your application menu:

```bash
cat > ~/.local/share/applications/scanoss-code-compare.desktop <<EOF
[Desktop Entry]
Name=SCANOSS Code Compare
Comment=Open source code comparison and analysis tool
Exec=/usr/local/bin/scanoss-cc
Icon=code
Terminal=false
Type=Application
Categories=Development;Utility;
EOF
```

---

## Troubleshooting

### macOS

**Issue:** "SCANOSS Code Compare.app is damaged and can't be opened"

**Solution:** This is a Gatekeeper warning. The app is signed and notarized, but you may need to:
```bash
xattr -cr "/Applications/SCANOSS Code Compare.app"
```

**Issue:** `scanoss-cc: command not found`

**Solution:** Make sure `/usr/local/bin` is in your PATH:
```bash
echo $PATH | grep /usr/local/bin
```

If not present, add it:
```bash
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Windows

**Issue:** `scanoss-cc` not recognized in terminal

**Solution:**
1. Make sure you added the directory containing `scanoss-cc-windows.exe` to your PATH
2. Open a new terminal window (PATH changes don't affect existing windows)
3. Verify the installation directory is in your PATH:
```powershell
$env:Path -split ';' | Select-String SCANOSS
```

**Issue:** Windows Defender SmartScreen warning

**Solution:** The executable is code-signed, but SmartScreen may show a warning for new downloads. Click "More info" and then "Run anyway".

**Issue:** Application won't start or crashes immediately

**Solution:** SCANOSS Code Compare requires WebView2 runtime. It's typically pre-installed on Windows 11, but Windows 10 users may need to install it:
- Download from: https://developer.microsoft.com/en-us/microsoft-edge/webview2/
- Or the app will prompt to download it automatically on first launch

### Linux

**Issue:** Missing GTK/WebKit dependencies

**Solution:** Install required libraries:
```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37

# Fedora
sudo dnf install gtk3 webkit2gtk4.0

# Arch
sudo pacman -S gtk3 webkit2gtk
```

---

## Getting Help

If you encounter issues not covered in this guide:

1. Check the [GitHub Issues](https://github.com/scanoss/scanoss.cc/issues) for similar problems
2. Review the [main README](README.md) for additional information
3. Open a new issue with:
   - Your operating system and version
   - Installation method used
   - Complete error messages
   - Steps to reproduce the problem

## Next Steps

After installation, check out the [Usage Guide](README.md#usage) to learn how to use SCANOSS Code Compare effectively.
