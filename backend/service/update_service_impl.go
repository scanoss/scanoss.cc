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
	// We need to extract the ZIP and open the DMG
	log.Info().Msg("Extracting and opening macOS DMG...")

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

	// Open the DMG file
	cmd = exec.Command("open", dmgPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open DMG: %w", err)
	}

	// Quit the application to allow installation
	wailsruntime.Quit(u.ctx)
	return nil
}

func (u *UpdateServiceImpl) applyUpdateWindows(updatePath string) error {
	// For Windows, the update is a ZIP file containing the new executable
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

	// Replace the current executable
	backupPath := currentExe + ".old"
	if err := os.Rename(currentExe, backupPath); err != nil {
		return fmt.Errorf("failed to backup current executable: %w", err)
	}

	if err := os.Rename(newBinaryPath, currentExe); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentExe)
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Restart the application
	cmd = exec.Command(currentExe, os.Args[1:]...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to restart application after update: %w", err)
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
	case ".deb":
		return u.applyDebUpdate(updatePath)
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
	cmd := exec.Command("unzip", "-o", updatePath, "-d", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract ZIP: %w", err)
	}

	// Find the binary in the extracted directory
	var newBinaryPath string
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Look for executable files that are not directories
		if !info.IsDir() && (info.Mode()&0o111 != 0 || strings.Contains(info.Name(), "scanoss-cc")) {
			newBinaryPath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to find binary in ZIP: %w", err)
	}

	if newBinaryPath == "" {
		return fmt.Errorf("no executable found in update ZIP")
	}

	// Make the new binary executable
	if err := os.Chmod(newBinaryPath, 0o755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Replace the current executable
	backupPath := currentExe + ".backup"
	if err := os.Rename(currentExe, backupPath); err != nil {
		return fmt.Errorf("failed to backup current executable: %w", err)
	}

	if err := os.Rename(newBinaryPath, currentExe); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentExe)
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Restart the application
	cmd = exec.Command(currentExe, os.Args[1:]...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to restart application after update: %w", err)
	}

	// Quit the current instance
	wailsruntime.Quit(u.ctx)
	return nil
}

func (u *UpdateServiceImpl) applyDebUpdate(updatePath string) error {
	log.Info().Msg("Opening .deb package for installation...")

	// Open with the default package manager
	cmd := exec.Command("xdg-open", updatePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open .deb package: %w", err)
	}

	// Quit the application to allow installation
	wailsruntime.Quit(u.ctx)
	return nil
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
		return "-linux.zip" // Linux binary in a ZIP file
	default:
		return ""
	}
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
