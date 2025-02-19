# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Upcoming changes...

## [0.6.1] 2025-02-19

### Fixed
- Fix scan dialog animation not working when clicking "Scan With Options"
- Fix path building when selecting scan root
- Fix scan dialog resetting to default values when opening

### Changed
- Use a single function to get all scan root derived paths in ScanDialog (results path, settings path, etc)

## [0.6.0] 2024-02-03

### Added
- Add sort options for results:
  - By path
  - By match percentage
- Add empty state for results
- Add scan root selector and info in top bar


### Fixed
- Improved diff viewer scrolling performance
- Show help command without running the process
- Check for several python default install locations when running scan command
- Set proper scan root when executing the app from symlink

### Changed
- Use same color for left and right code viewers
- "Scan With Options" improvements
  - Add advanced options input
  - Do not allow to manually select output path
  - Show boolean options as checkboxes
  - Hide console output when no lines are available
  - Create .scanoss directory if it does not exist

### Removed
- Remove "Scan Current Directory" menu option



## [0.5.0] - 2024-01-24
### Added
- Initial open source release
- React-based frontend for diff viewer
- Go backend with scanning capabilities
- Component and license management
- Decision saving and loading from SCANOSS settings file
- File scanning and result filtering
- Keyboard shortcuts support
- Terminal output viewer
- Configuration management through CLI


[0.5.0]: https://github.com/scanoss/scanoss.cc/compare/v0.4.0...v0.5.0
[0.6.0]: https://github.com/scanoss/scanoss.cc/compare/v0.5.0...v0.6.0
[0.6.1]: https://github.com/scanoss/scanoss.cc/compare/v0.6.0...v0.6.1