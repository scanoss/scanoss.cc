# Installing and Using SCANOSS Code Compare as a CLI Tool on macOS

`SCANOSS Code Compare` is primarily a GUI application, but it can also be run directly from the terminal as a command-line tool. Follow the steps below to set up and use it as a CLI tool.

## Quick Install (Recommended)

Use the automated installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install-macos.sh | bash
```

The script will:
- Offer to install via Homebrew (recommended) or direct download
- Install the app to /Applications
- Create the `scanoss-cc` CLI command automatically
- Handle all PATH configuration

## Manual Installation

### Step 1: Install the Application
1. Download `scanoss-cc-mac.zip` from [releases](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP to get the DMG
3. Open the `.dmg` file
4. Drag and drop `scanoss-cc.app` into the `/Applications` folder

### Step 2: Add the CLI to Your PATH

You have two options to make the CLI accessible from your terminal:

#### Option A: Create a Symlink (Recommended)
Run the following command in your terminal:
```bash
sudo ln -s "/Applications/scanoss-cc.app/Contents/MacOS/scanoss-cc" /usr/local/bin/scanoss-cc
```

#### Option B: Add to PATH
Add the following line to your shell configuration file (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):
```bash
export PATH="/Applications/scanoss-cc.app/Contents/MacOS:$PATH"
```

Then reload your shell configuration:
```bash
source ~/.zshrc  # or source ~/.bash_profile
```

**Important:** The full .app bundle must be in `/Applications` for GUI features (like folder picker dialogs) to work properly. This is due to how macOS resolves bundle resources.

## Verification

To verify the installation, open a terminal and run:
```bash
scanoss-cc --help
```

You should see a list of available commands and options.

## Uninstallation

To remove the CLI tool:

### If you used Option A (Symlink):
```bash
sudo rm /usr/local/bin/scanoss-cc
```

### If you used Option B (PATH):
Remove the export PATH line and the alias line from your shell configuration file (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`).
