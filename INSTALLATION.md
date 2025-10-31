# SCANOSS Code Compare - Comprehensive Installation Guide

This guide provides detailed installation instructions for all supported platforms and distribution methods.

## Table of Contents

- [macOS](#macos)
  - [Homebrew](#homebrew-recommended)
  - [PKG Installer](#pkg-installer)
  - [DMG](#dmg-manual-installation)
- [Windows](#windows)
  - [NSIS Installer](#nsis-installer)
- [Linux](#linux)
  - [AppImage](#appimage-recommended)
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

### PKG Installer

The `.pkg` installer provides a native macOS installation experience with automatic CLI setup.

1. Download `SCANOSS-Code-Compare-*-Installer.pkg` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Double-click the downloaded file
3. Follow the installation wizard
4. The app will be installed to `/Applications`
5. A symlink will be automatically created at `/usr/local/bin/scanoss-cc`

**What's Installed:**
- Application: `/Applications/SCANOSS Code Compare.app`
- CLI Symlink: `/usr/local/bin/scanoss-cc`

**Verify Installation:**
```bash
scanoss-cc --version
```

**Uninstallation:**
The PKG installer does not include an automatic uninstaller. To remove:
```bash
sudo rm -rf "/Applications/SCANOSS Code Compare.app"
sudo rm /usr/local/bin/scanoss-cc
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

### NSIS Installer

The NSIS installer provides a traditional Windows installation experience with customizable options.

1. Download `SCANOSS-Code-Compare-Setup.exe` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Run the installer
3. During installation:
   - Choose installation directory (default: `C:\Program Files\SCANOSS`)
   - ✅ **Check "Add scanoss-cc to PATH"** to enable CLI access from any command prompt
4. Click Install

**What's Installed:**
- Application files in `C:\Program Files\SCANOSS\`
- Start Menu shortcuts
- Desktop shortcut
- Optional: Added to system PATH

**Verify Installation:**
Open a new Command Prompt or PowerShell window:
```powershell
scanoss-cc --version
```

**Note:** You must open a new terminal window after installation for PATH changes to take effect.

**Uninstallation:**
- Use "Add or Remove Programs" in Windows Settings
- Or run the uninstaller from the Start Menu

---

## Linux

### AppImage (Recommended)

AppImage is a portable format that runs on most Linux distributions without installation.

1. Download `SCANOSS-Code-Compare-*.AppImage` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)

2. Make it executable:
```bash
chmod +x SCANOSS-Code-Compare-*.AppImage
```

3. Run:
```bash
./SCANOSS-Code-Compare-*.AppImage
```

**Install to System** (Optional):
```bash
# Move to a location in PATH
sudo mv SCANOSS-Code-Compare-*.AppImage /usr/local/bin/scanoss-cc

# Now you can run from anywhere
scanoss-cc
```

**Desktop Integration** (Optional):
Use [AppImageLauncher](https://github.com/TheAssassin/AppImageLauncher) for automatic desktop integration.

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

**Issue:** `scanoss-cc: command not found` after PKG installation

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

**Issue:** `scanoss-cc` not recognized after installer

**Solution:**
1. Make sure you checked "Add to PATH" during installation
2. Open a new terminal window (PATH changes don't affect existing windows)
3. Verify the installation directory is in your PATH:
```powershell
$env:Path -split ';' | Select-String SCANOSS
```

**Issue:** Windows Defender SmartScreen warning

**Solution:** The installer is signed, but SmartScreen may show a warning for new downloads. Click "More info" and then "Run anyway".

### Linux

**Issue:** AppImage won't run on Ubuntu 24.04

**Solution:** Ubuntu 24.04 ships with WebKit 4.1, but the AppImage expects 4.0. Create symlinks:
```bash
sudo ln -sf /usr/lib/x86_64-linux-gnu/libwebkit2gtk-4.1.so.0 /usr/lib/x86_64-linux-gnu/libwebkit2gtk-4.0.so.37
sudo ln -sf /usr/lib/x86_64-linux-gnu/libjavascriptcoregtk-4.1.so.0 /usr/lib/x86_64-linux-gnu/libjavascriptcoregtk-4.0.so.18
```

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
