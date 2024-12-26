# Installing and Using SCANOSS Lui as a CLI Tool on macOS

`SCANOSS Lui` is primarily a GUI application, but it can also be run directly from the terminal as a command-line tool. Follow the steps below to set up and use it as a CLI tool.

## Installation

### Step 1: Install the Application
1. Open the `.dmg` file
2. Drag and drop `scanoss-lui` into the `/Applications` folder

### Step 2: Add the CLI to Your PATH

You have two options to make the CLI accessible from your terminal:

#### Option A: Create a Symlink (Recommended)
Run the following command in your terminal:
```bash
sudo ln -s /Applications/scanoss-lui.app/Contents/MacOS/scanoss-lui /usr/local/bin/scanoss-lui
```

#### Option B: Add to PATH
Add the following line to your shell configuration file (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):
```bash
export PATH="/Applications/scanoss-lui.app/Contents/MacOS:$PATH"
```

Then reload your shell configuration:
```bash
source ~/.zshrc  # or source ~/.bash_profile
```

## Verification

To verify the installation, open a terminal and run:
```bash
scanoss-lui --help
```

You should see a list of available commands and options.

## Uninstallation

To remove the CLI tool:

### If you used Option A (Symlink):
```bash
sudo rm /usr/local/bin/scanoss-lui
```

### If you used Option B (PATH):
Remove the export PATH line from your shell configuration file (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`).
