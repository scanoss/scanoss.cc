// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/inconshreveable/go-update"
	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	githubAPIURL    = "https://api.github.com/repos/scanoss/scanoss.cc/releases/latest"
	downloadTimeout = 10 * time.Minute

	// updateHelperScript is the shell script that performs atomic .app bundle replacement
	updateHelperScript = `#!/bin/bash
# SCANOSS Update Helper Script
# This script performs atomic .app bundle replacement for macOS updates

set -e  # Exit on error
set -u  # Exit on undefined variable

OLD_APP_PATH="$1"
NEW_APP_PATH="$2"
BACKUP_PATH="$3"
SCAN_ROOT="$4"
OLD_PID="$5"

LOG_FILE="/tmp/scanoss-update.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "$LOG_FILE"
}

log "Update helper started"
log "Old app: $OLD_APP_PATH"
log "New app: $NEW_APP_PATH"
log "Backup: $BACKUP_PATH"
log "Scan root: $SCAN_ROOT"
log "Old PID: $OLD_PID"

# Wait for old process to exit (max 30 seconds)
log "Waiting for old process to exit..."
WAIT_COUNT=0
while kill -0 "$OLD_PID" 2>/dev/null; do
    sleep 1
    WAIT_COUNT=$((WAIT_COUNT + 1))
    if [ $WAIT_COUNT -gt 30 ]; then
        log "ERROR: Old process did not exit in time"
        exit 1
    fi
done
log "Old process exited"

# Additional grace period for file locks to release
sleep 2

# Verify new app signature
log "Verifying new app signature..."
if ! codesign --verify --deep "$NEW_APP_PATH" 2>&1 | tee -a "$LOG_FILE"; then
    log "ERROR: New app signature verification failed"
    exit 1
fi
log "Signature verification passed"

# Create backup of old app
log "Creating backup..."
if [ -d "$BACKUP_PATH" ]; then
    rm -rf "$BACKUP_PATH"
fi
if ! mv "$OLD_APP_PATH" "$BACKUP_PATH" 2>&1 | tee -a "$LOG_FILE"; then
    log "ERROR: Failed to create backup"
    exit 1
fi
log "Backup created successfully"

# Copy new app using ditto (preserves extended attributes)
log "Installing new app..."
if ! ditto --extattr "$NEW_APP_PATH" "$OLD_APP_PATH" 2>&1 | tee -a "$LOG_FILE"; then
    log "ERROR: Failed to install new app, restoring backup..."
    mv "$BACKUP_PATH" "$OLD_APP_PATH"
    exit 1
fi
log "New app installed"

# Remove quarantine attributes from entire bundle
log "Removing quarantine attributes..."
xattr -dr com.apple.quarantine "$OLD_APP_PATH" 2>&1 | tee -a "$LOG_FILE" || true
log "Quarantine removed"

# Launch new app
log "Launching new app..."
if [ -n "$SCAN_ROOT" ]; then
    open "$OLD_APP_PATH" --args --scan-root "$SCAN_ROOT" --check-update-success 2>&1 | tee -a "$LOG_FILE" || {
        log "ERROR: Failed to launch new app, restoring backup..."
        rm -rf "$OLD_APP_PATH"
        mv "$BACKUP_PATH" "$OLD_APP_PATH"
        open "$OLD_APP_PATH" --args --scan-root "$SCAN_ROOT"
        exit 1
    }
else
    open "$OLD_APP_PATH" --args --check-update-success 2>&1 | tee -a "$LOG_FILE" || {
        log "ERROR: Failed to launch new app, restoring backup..."
        rm -rf "$OLD_APP_PATH"
        mv "$BACKUP_PATH" "$OLD_APP_PATH"
        open "$OLD_APP_PATH"
        exit 1
    }
fi

log "Update completed successfully"
exit 0
`
)

type UpdateServiceImpl struct {
	ctx         context.Context
	httpClient  *http.Client
	downloadDir string
}

// NewUpdateService creates a new update service instance
func NewUpdateService() UpdateService {
	return &UpdateServiceImpl{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		downloadDir: os.TempDir(),
	}
}

type githubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Name               string  `json:"name"`
		BrowserDownloadURL string  `json:"browser_download_url"`
		Size               int64   `json:"size"`
		Digest             *string `json:"digest"` // SHA256 digest in format "sha256:hash"
	} `json:"assets"`
}

// SetContext sets the context for the service
func (s *UpdateServiceImpl) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// CheckForUpdate checks if a new version is available on GitHub
func (u *UpdateServiceImpl) CheckForUpdate() (*entities.UpdateInfo, error) {
	log.Info().Msg("Checking for updates...")

	currentVersion := u.GetCurrentVersion()
	if currentVersion == "" {
		return nil, fmt.Errorf("current version is not set")
	}

	// Fetch latest release from GitHub
	req, err := http.NewRequestWithContext(u.ctx, "GET", githubAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	// Parse versions
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersionClean := strings.TrimPrefix(currentVersion, "v")

	log.Info().Msgf("Current version: %s, Latest version: %s", currentVersionClean, latestVersion)

	// Compare versions
	isNewer, err := isVersionNewer(currentVersionClean, latestVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to compare versions: %w", err)
	}

	if !isNewer {
		log.Info().Msg("No update available")
		return &entities.UpdateInfo{
			Version:   latestVersion,
			Available: false,
		}, nil
	}

	// Find the appropriate asset for the current platform
	assetName := getAssetNameForPlatform()
	var downloadURL string
	var size int64
	var expectedSHA256 string

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, assetName) {
			downloadURL = asset.BrowserDownloadURL
			size = asset.Size

			// Extract SHA256 digest from GitHub's automatic checksum
			if asset.Digest != nil && *asset.Digest != "" {
				// Digest format is "sha256:hash", extract just the hash
				digest := *asset.Digest
				expectedSHA256, _ = strings.CutPrefix(digest, "sha256:")
				log.Info().Msgf("Found SHA256 digest for asset: %s", expectedSHA256)
			} else {
				log.Warn().Msg("No digest found for release asset - update download will fail verification")
			}
			break
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("no suitable download found for this platform")
	}

	log.Info().Msgf("Update available: %s", latestVersion)

	return &entities.UpdateInfo{
		Version:        latestVersion,
		DownloadURL:    downloadURL,
		ReleaseNotes:   release.Body,
		PublishedAt:    release.PublishedAt,
		Size:           size,
		Available:      true,
		ExpectedSHA256: expectedSHA256,
	}, nil
}

// DownloadUpdate downloads the new version to a temporary location
func (u *UpdateServiceImpl) DownloadUpdate(updateInfo *entities.UpdateInfo) (string, error) {
	log.Info().Msgf("Downloading update from: %s", updateInfo.DownloadURL)

	// Create a temporary file
	tmpFile, err := os.CreateTemp(u.downloadDir, "scanoss-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	tmpPath := tmpFile.Name()

	// Create request with timeout
	ctx, cancel := context.WithTimeout(u.ctx, downloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", updateInfo.DownloadURL, nil)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Download the file
	resp, err := u.httpClient.Do(req)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpPath)
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Copy with progress tracking
	hash := sha256.New()
	multiWriter := io.MultiWriter(tmpFile, hash)

	written, err := io.Copy(multiWriter, resp.Body)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to save update: %w", err)
	}

	// Compute actual checksum as hex string
	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	log.Info().Msgf("Downloaded %d bytes (SHA256: %s)", written, actualChecksum)

	// Verify checksum - this is mandatory for security
	if updateInfo.ExpectedSHA256 == "" {
		os.Remove(tmpPath)
		return "", fmt.Errorf("no expected checksum available for verification - cannot safely install update")
	}

	expectedChecksum := strings.ToLower(strings.TrimSpace(updateInfo.ExpectedSHA256))
	actualChecksumLower := strings.ToLower(actualChecksum)

	if expectedChecksum != actualChecksumLower {
		os.Remove(tmpPath)
		return "", fmt.Errorf("checksum verification failed: expected %s, got %s", expectedChecksum, actualChecksumLower)
	}

	log.Info().Msg("Checksum verification passed successfully")

	return tmpPath, nil
}

// ApplyUpdate applies the downloaded update and restarts the application
func (u *UpdateServiceImpl) ApplyUpdate(updatePath string) error {
	log.Info().Msgf("Applying update from: %s", updatePath)

	switch runtime.GOOS {
	case "darwin":
		return u.applyUpdateMacOS(updatePath)
	case "windows":
		return u.applyUpdateWindows(updatePath)
	case "linux":
		return u.applyUpdateLinux(updatePath)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (u *UpdateServiceImpl) applyUpdateMacOS(updatePath string) error {
	// For macOS, the update is a .zip containing a .dmg file with a signed .app bundle
	// We use a helper script to atomically replace the entire .app bundle
	log.Info().Msg("Extracting and installing macOS update...")

	// Extract the ZIP file to a temporary directory
	extractDir := filepath.Join(os.TempDir(), "scanoss-update-extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Use unzip command to extract
	log.Info().Msg("Extracting ZIP...")
	cmd := exec.Command("unzip", "-o", updatePath, "-d", extractDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to extract ZIP: %w (output: %s)", err, string(output))
	}

	// Find the DMG file in the extracted directory
	var dmgPath string
	err := filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".dmg") {
			dmgPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find DMG in ZIP: %w", err)
	}

	if dmgPath == "" {
		return fmt.Errorf("no DMG file found in update ZIP")
	}

	log.Info().Msgf("Found DMG: %s", dmgPath)

	// Verify DMG signature
	log.Info().Msg("Verifying DMG signature...")
	cmd = exec.Command("codesign", "--verify", "--deep", dmgPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("DMG signature verification failed: %w (output: %s)", err, string(output))
	}
	log.Info().Msg("DMG signature verification passed")

	// Mount the DMG
	mountPoint := filepath.Join(os.TempDir(), "scanoss-update-mount")
	if err := os.MkdirAll(mountPoint, 0o755); err != nil {
		return fmt.Errorf("failed to prepare mount point: %w", err)
	}
	log.Info().Msgf("Mounting DMG to %s...", mountPoint)
	cmd = exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-mountpoint", mountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to mount DMG: %w (output: %s)", err, string(output))
	}
	// Ensure DMG is unmounted on exit
	defer func() {
		log.Info().Msg("Unmounting DMG...")
		cmd := exec.Command("hdiutil", "detach", mountPoint, "-force")
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Warn().Err(err).Msgf("Failed to unmount DMG (output: %s)", string(output))
		}
		os.RemoveAll(mountPoint)
	}()

	// Find the .app bundle in the mounted DMG
	var appPath string
	err = filepath.Walk(mountPoint, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".app") {
			appPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find .app bundle in DMG: %w", err)
	}

	if appPath == "" {
		return fmt.Errorf("no .app bundle found in DMG")
	}

	log.Info().Msgf("Found .app bundle: %s", appPath)

	// Verify the .app bundle signature
	log.Info().Msg("Verifying .app bundle signature...")
	cmd = exec.Command("codesign", "--verify", "--deep", appPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf(".app bundle signature verification failed: %w (output: %s)", err, string(output))
	}
	log.Info().Msg(".app bundle signature verification passed")

	// Copy the entire .app bundle to a staging location
	stagingDir := filepath.Join(os.TempDir(), "scanoss-update-staging")
	if err := os.MkdirAll(stagingDir, 0o755); err != nil {
		return fmt.Errorf("failed to create staging directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	newAppPath := filepath.Join(stagingDir, filepath.Base(appPath))
	log.Info().Msgf("Copying .app bundle to staging location: %s", newAppPath)
	cmd = exec.Command("ditto", "--extattr", appPath, newAppPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to copy .app bundle: %w (output: %s)", err, string(output))
	}

	// Current app bundle path (assuming /Applications installation)
	currentAppPath := "/Applications/scanoss-cc.app"
	log.Info().Msgf("Current app bundle: %s", currentAppPath)

	// Verify the current app exists
	if _, err := os.Stat(currentAppPath); err != nil {
		return fmt.Errorf("current app not found at %s: %w", currentAppPath, err)
	}

	// Backup path
	backupPath := "/Applications/.scanoss-cc.app.backup"

	// Get current process PID
	pid := os.Getpid()

	// Get current scan root
	currentScanRoot := config.GetInstance().GetScanRoot()

	// Write the helper script to a temporary file
	helperScriptPath := filepath.Join(os.TempDir(), "scanoss-update-helper.sh")
	if err := os.WriteFile(helperScriptPath, []byte(updateHelperScript), 0o755); err != nil {
		return fmt.Errorf("failed to write helper script: %w", err)
	}
	defer os.Remove(helperScriptPath)

	log.Info().Msg("Launching update helper script...")

	// Launch the helper script in the background
	cmd = exec.Command(helperScriptPath, currentAppPath, newAppPath, backupPath, currentScanRoot, fmt.Sprintf("%d", pid))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch helper script: %w", err)
	}

	log.Info().Msg("Helper script launched, quitting application...")

	// Quit the current instance - the helper will wait for us to exit
	wailsruntime.Quit(u.ctx)
	return nil
}

func (u *UpdateServiceImpl) applyUpdateWindows(updatePath string) error {
	// For Windows, the update is a ZIP file containing the new executable
	// We use go-update library which handles Windows file locking correctly
	log.Info().Msg("Applying Windows update from ZIP...")

	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}

	// Extract the ZIP file to a temporary directory
	extractDir := filepath.Join(os.TempDir(), "scanoss-update-extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Extract using PowerShell Expand-Archive (available on all modern Windows)
	log.Info().Msg("Extracting update ZIP...")
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force", updatePath, extractDir))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract ZIP: %w", err)
	}

	// Find the .exe file in the extracted directory
	var newBinaryPath string
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Look for .exe files
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".exe") {
			newBinaryPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find executable in ZIP: %w", err)
	}

	if newBinaryPath == "" {
		return fmt.Errorf("no executable found in update ZIP")
	}

	log.Info().Msgf("Found new binary: %s", newBinaryPath)

	// Open the new binary file
	newBinary, err := os.Open(newBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to open new binary: %w", err)
	}
	defer newBinary.Close()

	// Apply the update using go-update
	log.Info().Msg("Applying update...")
	err = update.Apply(newBinary, update.Options{
		TargetPath: currentExe,
	})
	if err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}

	log.Info().Msg("Update applied successfully, restarting application...")

	// Restart the application
	cmd = exec.Command(currentExe, os.Args[1:]...)
	if err := cmd.Start(); err != nil {
		log.Warn().Err(err).Msg("failed to restart application")
	}

	// Quit the current instance
	wailsruntime.Quit(u.ctx)
	return nil
}

func (u *UpdateServiceImpl) applyUpdateLinux(updatePath string) error {
	// For Linux, we support .zip and .deb packages
	ext := filepath.Ext(updatePath)

	switch ext {
	case ".zip":
		return u.applyZipUpdate(updatePath)
	default:
		return fmt.Errorf("unsupported update format: %s", ext)
	}
}

func (u *UpdateServiceImpl) applyZipUpdate(updatePath string) error {
	log.Info().Msg("Applying ZIP update...")

	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}

	// Extract the ZIP file to a temporary directory
	extractDir := filepath.Join(os.TempDir(), "scanoss-update-extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Use unzip command to extract
	log.Info().Msg("Extracting update ZIP...")
	cmd := exec.Command("unzip", "-o", updatePath, "-d", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract ZIP: %w", err)
	}

	// Find the binary in the extracted directory with better detection
	var newBinaryPath string
	var candidates []string
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Look for executable files that are not directories
		if !info.IsDir() && info.Mode()&0o111 != 0 {
			// Prefer files with "scanoss" in the name
			if strings.Contains(strings.ToLower(info.Name()), "scanoss") {
				newBinaryPath = path
				return filepath.SkipAll
			}
			// Otherwise collect as candidates
			candidates = append(candidates, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find binary in ZIP: %w", err)
	}

	// If no scanoss binary found, use first candidate
	if newBinaryPath == "" && len(candidates) > 0 {
		newBinaryPath = candidates[0]
	}

	if newBinaryPath == "" {
		return fmt.Errorf("no executable found in update ZIP")
	}

	log.Info().Msgf("Found new binary: %s", newBinaryPath)

	// Verify the binary is valid
	if err := u.verifyLinuxBinary(newBinaryPath); err != nil {
		return fmt.Errorf("binary verification failed: %w", err)
	}

	// Open the new binary file
	newBinary, err := os.Open(newBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to open new binary: %w", err)
	}
	defer newBinary.Close()

	// Apply the update using go-update library
	// Linux allows replacing running executables, so this works directly
	log.Info().Msg("Applying update to binary...")
	if err := update.Apply(newBinary, update.Options{
		TargetPath: currentExe,
	}); err != nil {
		if rollbackErr := update.RollbackError(err); rollbackErr != nil {
			log.Error().Err(rollbackErr).Msg("failed to rollback after Linux update error")
		}
		return fmt.Errorf("failed to apply update: %w", err)
	}

	log.Info().Msg("Update applied successfully, restarting application...")

	// Restart the application
	cmd = exec.Command(currentExe, os.Args[1:]...)
	if err := cmd.Start(); err != nil {
		log.Warn().Err(err).Msg("failed to restart application")
	}

	// Quit the current instance
	wailsruntime.Quit(u.ctx)
	return nil
}

// verifyLinuxBinary performs basic verification on a Linux binary
func (u *UpdateServiceImpl) verifyLinuxBinary(path string) error {
	// Check file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file does not exist: %w", err)
	}

	// Check minimum size (should be at least a few MB for a real binary)
	if info.Size() < 1024*1024 {
		return fmt.Errorf("binary too small: %d bytes", info.Size())
	}

	// Check if it's a valid ELF binary by reading magic bytes
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	magic := make([]byte, 4)
	if _, err := file.Read(magic); err != nil {
		return fmt.Errorf("cannot read file header: %w", err)
	}

	// ELF magic: 0x7f 'E' 'L' 'F'
	if magic[0] != 0x7f || magic[1] != 'E' || magic[2] != 'L' || magic[3] != 'F' {
		return fmt.Errorf("not a valid ELF binary")
	}

	return nil
}

// GetCurrentVersion returns the current application version
func (u *UpdateServiceImpl) GetCurrentVersion() string {
	return entities.AppVersion
}

// VerifyUpdateSuccess checks if an update completed successfully and cleans up backup
// This should be called on app startup when --check-update-success flag is present
func (u *UpdateServiceImpl) VerifyUpdateSuccess() error {
	if runtime.GOOS != "darwin" {
		return nil // Only implemented for macOS currently
	}

	backupPath := "/Applications/.scanoss-cc.app.backup"

	// Check if backup exists
	if _, err := os.Stat(backupPath); err != nil {
		if os.IsNotExist(err) {
			// No backup exists, nothing to clean up
			return nil
		}
		return fmt.Errorf("failed to check backup status: %w", err)
	}

	// Backup exists - update was successful, clean it up
	log.Info().Msg("Update completed successfully, removing backup...")
	if err := os.RemoveAll(backupPath); err != nil {
		log.Warn().Err(err).Msg("Failed to remove backup (non-fatal)")
		return nil // Non-fatal error
	}

	log.Info().Msg("Update verification complete")
	return nil
}

// CheckForFailedUpdate checks if the previous update failed and performs rollback if needed
// This should be called on app startup before checking for update success
func (u *UpdateServiceImpl) CheckForFailedUpdate() error {
	if runtime.GOOS != "darwin" {
		return nil // Only implemented for macOS currently
	}

	backupPath := "/Applications/.scanoss-cc.app.backup"
	currentAppPath := "/Applications/scanoss-cc.app"

	// Check if backup exists
	if _, err := os.Stat(backupPath); err != nil {
		if os.IsNotExist(err) {
			// No backup exists, no failed update
			return nil
		}
		return fmt.Errorf("failed to check backup status: %w", err)
	}

	// Backup exists - this means either:
	// 1. Update is in progress (helper script is running)
	// 2. Update failed and we need to rollback

	// If we're here with a backup and no --check-update-success flag,
	// it means the update failed (new version crashed before startup)

	log.Warn().Msg("Detected failed update, attempting rollback...")

	// Remove the current (failed) app
	if err := os.RemoveAll(currentAppPath); err != nil {
		return fmt.Errorf("failed to remove failed app during rollback: %w", err)
	}

	// Restore from backup
	if err := os.Rename(backupPath, currentAppPath); err != nil {
		return fmt.Errorf("failed to restore backup during rollback: %w", err)
	}

	log.Info().Msg("Rollback successful - previous version restored")

	// Restart the app
	cmd := exec.Command("open", currentAppPath)
	if err := cmd.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to restart after rollback")
	}

	// Exit this (failed) instance
	os.Exit(1)
	return nil
}

// Helper function to get the appropriate asset name for the current platform
func getAssetNameForPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		return "-mac.zip" // macOS DMG in a ZIP file
	case "windows":
		return "-win.zip" // Windows executable in a ZIP file
	case "linux":
		// Detect webkit version for Linux. The resulting zip file will be named
		// -linux-amd64-webkit40.zip or -linux-amd64-webkit41.zip
		webkitVersion := detectLinuxWebkitVersion()
		return fmt.Sprintf("-linux-amd64-%s.zip", webkitVersion)
	default:
		return ""
	}
}

// detectLinuxWebkitVersion detects which webkit version to use
// Returns "webkit40" or "webkit41" based on distro version
func detectLinuxWebkitVersion() string {
	// Read /etc/os-release to detect distribution
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		log.Warn().Err(err).Msg("Could not read /etc/os-release, defaulting to webkit40")
		return "webkit40"
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var distroID string
	var versionID string

	for _, line := range lines {
		if token, ok := strings.CutPrefix(line, "ID="); ok {
			distroID = strings.Trim(token, "\"")
		}
		if token, ok := strings.CutPrefix(line, "VERSION_ID="); ok {
			versionID = strings.Trim(token, "\"")
		}
	}

	// Ubuntu >= 24.04 or Debian >= 13 uses webkit41
	switch distroID {
	case "ubuntu":
		// Parse version like "24.04"
		parts := strings.Split(versionID, ".")
		if len(parts) >= 1 {
			var ver int
			if _, err := fmt.Sscanf(parts[0], "%d", &ver); err == nil && ver >= 24 {
				return "webkit41"
			}
		}
	case "debian":
		// Parse version like "13" or "13.1"
		parts := strings.Split(versionID, ".")
		if len(parts) >= 1 {
			var ver int
			if _, err := fmt.Sscanf(parts[0], "%d", &ver); err == nil && ver >= 13 {
				return "webkit41"
			}
		}
	}

	// Default to webkit40 for older versions or other distros
	return "webkit40"
}

// isVersionNewer compares two semantic versions
func isVersionNewer(current, latest string) (bool, error) {
	currentVersion, err := semver.NewVersion(current)
	if err != nil {
		return false, fmt.Errorf("failed to parse current version: %w", err)
	}
	latestVersion, err := semver.NewVersion(latest)
	if err != nil {
		return false, fmt.Errorf("failed to parse latest version: %w", err)
	}
	return latestVersion.GreaterThan(currentVersion), nil
}
