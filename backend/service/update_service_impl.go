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
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	githubAPIURL    = "https://api.github.com/repos/scanoss/scanoss.cc/releases/latest"
	downloadTimeout = 10 * time.Minute
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
	// For macOS, the update is a .zip containing a .dmg file
	// We need to extract, mount, copy, and install the .app bundle
	log.Info().Msg("Extracting and installing macOS DMG...")

	// Extract the ZIP file to a temporary directory
	extractDir := filepath.Join(os.TempDir(), "scanoss-update-extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Use unzip command to extract
	cmd := exec.Command("unzip", "-o", updatePath, "-d", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract ZIP: %w", err)
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

	// Mount the DMG
	mountPoint := filepath.Join(os.TempDir(), "scanoss-update-mount")
	log.Info().Msgf("Mounting DMG to %s...", mountPoint)
	cmd = exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-mountpoint", mountPoint)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mount DMG: %w", err)
	}
	// Ensure DMG is unmounted on exit
	defer func() {
		log.Info().Msg("Unmounting DMG...")
		exec.Command("hdiutil", "detach", mountPoint, "-force").Run()
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

	// Determine installation path
	installPath := "/Applications/scanoss-cc.app"
	log.Info().Msgf("Installing app to %s...", installPath)

	// Try to use rsync for atomic copy
	if _, err := exec.LookPath("rsync"); err == nil {
		log.Info().Msg("Using rsync for atomic installation...")
		cmd = exec.Command("rsync", "-a", "--delete", appPath+"/", installPath+"/")
		if err := cmd.Run(); err != nil {
			log.Warn().Err(err).Msg("rsync failed, falling back to cp")
			// Fallback to cp if rsync fails
			if err := u.copyAppWithCP(appPath, installPath); err != nil {
				return err
			}
		}
	} else {
		// rsync not available, use cp
		log.Info().Msg("rsync not available, using cp...")
		if err := u.copyAppWithCP(appPath, installPath); err != nil {
			return err
		}
	}

	// Clear macOS quarantine attribute
	log.Info().Msg("Clearing quarantine attributes...")
	exec.Command("xattr", "-r", "-d", "com.apple.quarantine", installPath).Run()

	log.Info().Msg("Installation complete, restarting application...")

	// Restart the application
	cmd = exec.Command("open", installPath)
	if err := cmd.Start(); err != nil {
		log.Warn().Err(err).Msg("failed to restart application")
	}

	// Quit the current instance
	wailsruntime.Quit(u.ctx)
	return nil
}

// copyAppWithCP copies the .app bundle using cp command (fallback method)
func (u *UpdateServiceImpl) copyAppWithCP(sourcePath, destPath string) error {
	// Remove existing installation
	log.Info().Msgf("Removing old installation at %s...", destPath)
	if err := exec.Command("rm", "-rf", destPath).Run(); err != nil {
		log.Warn().Err(err).Msg("failed to remove old installation, continuing anyway")
	}

	// Copy new app bundle
	log.Info().Msgf("Copying new app from %s to %s...", sourcePath, destPath)
	cmd := exec.Command("cp", "-R", sourcePath, destPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy app bundle: %w", err)
	}

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

	// Make the new binary executable
	if err := os.Chmod(newBinaryPath, 0o755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Use staging approach: copy to .next file for atomic swap
	// This avoids issues with running executables on Linux
	nextPath := currentExe + ".next"
	log.Info().Msgf("Staging new binary to %s...", nextPath)

	// Copy the new binary to .next location
	if err := u.copyFile(newBinaryPath, nextPath); err != nil {
		return fmt.Errorf("failed to stage new binary: %w", err)
	}

	// Make sure .next is executable
	if err := os.Chmod(nextPath, 0o755); err != nil {
		os.Remove(nextPath)
		return fmt.Errorf("failed to make staged binary executable: %w", err)
	}

	log.Info().Msg("Update staged successfully, restarting application...")

	// The actual swap will happen on application startup
	// For now, just restart to trigger the swap
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

// copyFile copies a file from src to dst
func (u *UpdateServiceImpl) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}

// CheckPendingUpdate checks if there's a pending update and applies it
// This should be called during application startup (Linux only)
func CheckPendingUpdate() error {
	if runtime.GOOS != "linux" {
		return nil // Only applicable to Linux
	}

	currentExe, err := os.Executable()
	if err != nil {
		return nil // Silently ignore errors during startup check
	}

	nextPath := currentExe + ".next"
	oldPath := currentExe + ".old"

	// Check if there's a pending update
	if _, err := os.Stat(nextPath); os.IsNotExist(err) {
		// No pending update, check if we should cleanup old backup
		if _, err := os.Stat(oldPath); err == nil {
			// Remove old backup from previous successful update
			log.Info().Msg("Removing old backup from previous update...")
			os.Remove(oldPath)
		}
		return nil
	}

	log.Info().Msg("Pending update detected, performing atomic swap...")

	// Backup current executable
	if err := os.Rename(currentExe, oldPath); err != nil {
		log.Error().Err(err).Msg("Failed to backup current executable")
		return err
	}

	// Move .next to current
	if err := os.Rename(nextPath, currentExe); err != nil {
		// Rollback
		log.Error().Err(err).Msg("Failed to swap binaries, rolling back...")
		os.Rename(oldPath, currentExe)
		return err
	}

	// Make sure it's executable
	os.Chmod(currentExe, 0o755)

	log.Info().Msg("Update applied successfully, restarting...")

	// Re-exec into the new binary
	if err := restartApplication(currentExe); err != nil {
		log.Error().Err(err).Msg("Failed to restart after update")
		return err
	}

	// This code should never be reached as we've re-exec'd
	return nil
}

// restartApplication restarts the application by re-executing it
func restartApplication(exePath string) error {
	return exec.Command(exePath, os.Args[1:]...).Start()
}

// GetCurrentVersion returns the current application version
func (u *UpdateServiceImpl) GetCurrentVersion() string {
	return entities.AppVersion
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
