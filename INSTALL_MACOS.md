# Installing and Using SCANOSS Code Compare as a CLI Tool on macOS

`SCANOSS Code Compare` is primarily a GUI application, but it can also be run directly from the terminal as a command-line tool. Follow the steps below to set up and use it as a CLI tool.

## Installation

### Step 1: Install the Application
1. Open the `.dmg` file
2. Drag and drop `SCANOSS Code Compare` into the `/Applications` folder

### Step 2: Add the CLI to Your PATH

You have two options to make the CLI accessible from your terminal:

#### Option A: Create a Symlink (Recommended)
Run the following command in your terminal:
```bash
sudo ln -s "/Applications/SCANOSS Code Compare.app/Contents/MacOS/SCANOSS Code Compare" /usr/local/bin/scanoss-cc
```

#### Option B: Add to PATH
Add the following lines to your shell configuration file (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):
```bash
export PATH="/Applications/SCANOSS Code Compare.app/Contents/MacOS:$PATH"
alias scanoss-cc='"/Applications/SCANOSS Code Compare.app/Contents/MacOS/SCANOSS Code Compare"'
```

Then reload your shell configuration:
```bash
source ~/.zshrc  # or source ~/.bash_profile
```

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
