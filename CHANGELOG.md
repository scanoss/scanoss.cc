# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Upcoming changes...

## [0.7.5] 2025-05-20
### Fixed
- Handle malformed JSON in results file

## [0.7.4] 2025-05-20
### Fixed
- Normalize paths when adding/removing skip patterns

## [0.7.3] 2025-05-01
### Fixed
- Fix splitting paths on Windows

## [0.7.1] 2025-04-28
### Fixed
- Disable pop-up warning message that appeared when launching the application from a graphical interface on Windows.


## [0.7.0] 2025-03-24

### Added
- Add abort scan button

## [0.6.5] 2025-03-21

### Added
- Add UI for modifying scanoss.json skip settings
- Add Settings menu

### Modified
- Improve scanning performance


## [0.6.4] 2025-03-18

### Fixed
- Fix search bar not resetting results
- Fix scan terminal output not showing output immediately
- Fix scan dialog size

## [0.6.3] 2025-03-12

### Fixed
- Fix windows build

## [0.6.2] 2025-03-12

### Fixed
- Fix windows build

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
[0.6.2]: https://github.com/scanoss/scanoss.cc/compare/v0.6.1...v0.6.2
[0.6.3]: https://github.com/scanoss/scanoss.cc/compare/v0.6.2...v0.6.3
[0.6.4]: https://github.com/scanoss/scanoss.cc/compare/v0.6.3...v0.6.4
[0.6.5]: https://github.com/scanoss/scanoss.cc/compare/v0.6.4...v0.6.5
[0.7.0]: https://github.com/scanoss/scanoss.cc/compare/v0.6.5...v0.7.0
[0.7.1]: https://github.com/scanoss/scanoss.cc/compare/v0.7.0...v0.7.1
[0.7.2]: https://github.com/scanoss/scanoss.cc/compare/v0.7.0...v0.7.2
[0.7.3]: https://github.com/scanoss/scanoss.cc/compare/v0.7.2...v0.7.3
[0.7.4]: https://github.com/scanoss/scanoss.cc/compare/v0.7.3...v0.7.4
[0.7.5]: https://github.com/scanoss/scanoss.cc/compare/v0.7.4...v0.7.5

