# SCANOSS Code Compare - Comprehensive Installation Guide

This guide provides detailed installation instructions for all supported platforms and distribution methods.

## Table of Contents

- [macOS](#macos)
  - [Homebrew](#homebrew-recommended)
  - [DMG](#dmg-manual-installation)
- [Windows](#windows)
  - [NSIS Installer](#nsis-installer)
- [Linux](#linux)
  - [Distribution Packages](#distribution-packages)
- [Troubleshooting](#troubleshooting)

---

## macOS

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

1. Download `SCANOSS-Code-Compare-*.dmg` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Open the DMG file
3. Drag "SCANOSS Code Compare.app" to the Applications folder
4. Eject the DMG

**Setting Up CLI Access:**

After installing via DMG, you need to manually set up CLI access. Choose one of these methods:

**Option 1: Symlink (Recommended)**
```bash
sudo ln -s "/Applications/SCANOSS Code Compare.app/Contents/MacOS/SCANOSS Code Compare" /usr/local/bin/scanoss-cc
```

**Option 2: Add to PATH**
```bash
echo 'export PATH="/Applications/SCANOSS Code Compare.app/Contents/MacOS:$PATH"' >> ~/.zshrc
echo 'alias scanoss-cc="/Applications/SCANOSS Code Compare.app/Contents/MacOS/SCANOSS Code Compare"' >> ~/.zshrc
source ~/.zshrc
```

For more details, see [INSTALL_MACOS.md](INSTALL_MACOS.md).

---

## Windows

### Standalone Executable

SCANOSS Code Compare for Windows is distributed as a portable, standalone executable.

1. Download `scanoss-cc-win.zip` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP file
3. Run `scanoss-cc-windows.exe`

**Installation Options:**

**Option 1: Portable Use**
- Run directly from Downloads or any location
- No installation required
- Ideal for USB drives or temporary usage

**Option 2: Install to Program Files (Recommended)**
```powershell
# In PowerShell (Run as Administrator)
New-Item -Path "C:\Program Files\SCANOSS" -ItemType Directory -Force
Move-Item -Path ".\scanoss-cc-windows.exe" -Destination "C:\Program Files\SCANOSS\" -Force
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

### Distribution Packages

#### Debian/Ubuntu (.deb)

1. Download the Linux ZIP file from [releases](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP
3. Install dependencies:
```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
```
4. The binary is ready to use

Move to PATH:
```bash
sudo mv scanoss-cc /usr/local/bin/
```

#### Fedora/RHEL (.rpm)

1. Download the Linux ZIP file from [releases](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP
3. Install dependencies:
```bash
sudo dnf install gtk3 webkit2gtk4.0
```
4. Move to PATH:
```bash
sudo mv scanoss-cc /usr/local/bin/
```

#### Arch Linux (AUR)

Community-maintained AUR packages may be available. Check [AUR](https://aur.archlinux.org/) for `scanoss-code-compare`.

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
