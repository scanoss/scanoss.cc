
<img src="assets/appicon.png" alt="SCANOSS Code Compare Logo" width="150"/>

# SCANOSS Code Compare

SCANOSS Code Compare is a streamlined desktop application for managing open source findings with a clean, distraction-free interface. It features vim-style navigation (j/k), side-by-side code comparison for both snippet and 100% matches, and robust decision management that persists across scans. Users can quickly filter results, mark components as included/omitted/replaced with single keystrokes, and access their previous decisions in future scans.


## Features

- 🔍 Advanced code scanning and component identification
- 📊 Detailed dependency analysis and visualization
- 📝 License compliance checking and management
- 🔄 Real-time scanning results
- ⚡ Fast and efficient local processing
- 🎯 Accurate component matching
- 🖥️ Cross-platform support

## Prerequisites

- Go 1.x or higher
- Node.js and npm
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

## Installation

### Quick Install (Recommended)

Install SCANOSS Code Compare with a single command:

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install.sh | bash
```

**Windows (PowerShell as Administrator):**
```powershell
irm https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install-windows.ps1 | iex
```

### Security-Conscious Installation

Download and inspect the script before running:

```bash
curl -fsSL https://raw.githubusercontent.com/scanoss/scanoss.cc/main/scripts/install.sh -o install.sh
cat install.sh  # Review the script
bash install.sh
```

### Alternative Installation Methods

<details>
<summary><b>macOS</b></summary>

#### Homebrew
```bash
brew install scanoss/dist/scanoss-code-compare
```

#### Direct Download
1. Download `scanoss-cc-mac.zip` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Extract and open the DMG
3. Drag the app to Applications

See [INSTALLATION.md](INSTALLATION.md) for detailed instructions.
</details>

<details>
<summary><b>Windows</b></summary>

1. Download `scanoss-cc-win.zip` from the [releases page](https://github.com/scanoss/scanoss.cc/releases)
2. Extract the ZIP file
3. Run `scanoss-cc.exe`

See [INSTALLATION.md](INSTALLATION.md) for PATH setup and installation to Program Files.
</details>

<details>
<summary><b>Linux</b></summary>

Choose the appropriate binary for your system from the [releases page](https://github.com/scanoss/scanoss.cc/releases):
- **Ubuntu 22.04 and older / Debian 12 and older**: [scanoss-cc-linux-amd64.zip](https://github.com/scanoss/scanoss.cc/releases/latest/download/scanoss-cc-linux-amd64.zip)
- **Ubuntu 24.04+ / Debian 13+**: [scanoss-cc-linux-amd64-webkit41.zip](https://github.com/scanoss/scanoss.cc/releases/latest/download/scanoss-cc-linux-amd64-webkit41.zip)

Extract and run:
```bash
# For Ubuntu 22.04 and older
unzip scanoss-cc-linux-amd64.zip
sudo mv scanoss-cc-linux-amd64 /usr/local/bin/scanoss-cc

# For Ubuntu 24.04+ (webkit 4.1)
unzip scanoss-cc-linux-amd64-webkit41.zip
sudo mv scanoss-cc-linux-amd64-webkit41 /usr/local/bin/scanoss-cc
```

See [INSTALLATION.md](INSTALLATION.md) for detailed instructions.
</details>

### Verify Installation

After installation, verify that the CLI is available:
```bash
scanoss-cc --version
```

### Uninstalling

For complete uninstall instructions for all platforms, see [INSTALLATION.md](INSTALLATION.md#uninstalling).

## Usage

### CLI Parameters

| Parameter      | Description                                                                 | Default Value |
|----------------|-----------------------------------------------------------------------------|---------------|
| **scan-root**  | Scanned folder                                                              | $WORKDIR |
| **input**      | Path to results.json file of the scanned project                            | $WORKDIR/.scanoss/results.json |
| **config**     | Path to configuration file                                                  | $HOME/.scanoss/scanoss-cc-settings.json |
| **apiUrl**     | SCANOSS API URL                                                             | https://api.osskb.org |
| **key**        | SCANOSS API Key token (not required for default OSSKB URL)                  | - |
| **debug**      | Enable debug mode                                                           | false |

### Example Commands

```bash
# Open the GUI application
scanoss-cc

# Open the GUI application with custom parameters (you can also change these from the GUI)
scanoss-cc --scan-root /path/to/scanned/project --input /path/to/results.json

# Basic scan with default settings
scanoss-cc scan /path/to/project

# Scan with custom results path
scanoss-cc scan --input /path/to/results.json

# Scan current directory with multiple parameters
scanoss-cc scan . --key $SCANOSS_API_KEY --apiurl $SCANOSS_API_URL --debug
```

## Development

### Dependencies

Before setting up the development environment, ensure you have the following dependencies installed:

**Prerequisites:**
- Go 1.23+
- Node.js and npm

**System Dependencies (Debian/Ubuntu):**

For Ubuntu 22.04 and older:
```bash
sudo apt-get update
sudo apt-get install -y build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
```

For Ubuntu 24.04+ / Debian 13+:
```bash
sudo apt-get update
sudo apt-get install -y build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Install Wails CLI:**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
### Setting Up the Development Environment

1. Clone the repository:
```bash
git clone https://github.com/scanoss/scanoss.cc.git
cd scanoss.cc
```

2. Run in development mode:
```bash
make run
```

### Pre-commit Setup
This project uses pre-commit hooks to ensure code quality and consistency. To set up pre-commit, run:
```bash
pip3 install pre-commit
pre-commit install
```

This will install the pre-commit tool and set up the git hooks defined in the `.pre-commit-config.yaml` file to run automatically on each commit.

### Development with Custom Parameters

```bash
# Using make command
make run APPARGS="--scan-root <scanRootPath> --input <resultPath>"

# Using wails command
wails dev -appargs "--input <resultPath>"
```

### Building

```bash
# Build for the current platform
make build
```

### Building for Different Linux Versions

SCANOSS Code Compare provides two Linux build variants to support different WebKit versions:

- Use `make build` for Ubuntu 22.04 and older (WebKit 4.0)
- Use `make build_webkit41` for Ubuntu 24.04+/Debian 13+ (WebKit 4.1)

## Contributing

We welcome contributions! Please read our [Contributing Guidelines](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) before submitting pull requests.

### Reporting Bugs

When submitting bug reports, please include:
- SCANOSS Code Compare version
- Your system information (OS, Go version, etc.)
- Steps to reproduce the issue
- Expected vs actual behavior

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please:
1. Check our [documentation](https://scanoss.readthedocs.io)
2. Open an issue in this repository
