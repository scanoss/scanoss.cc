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
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
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

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, assetName) {
			downloadURL = asset.BrowserDownloadURL
			size = asset.Size
			break
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("no suitable download found for this platform")
	}

	log.Info().Msgf("Update available: %s", latestVersion)

	return &entities.UpdateInfo{
		Version:      latestVersion,
		DownloadURL:  downloadURL,
		ReleaseNotes: release.Body,
		PublishedAt:  release.PublishedAt,
		Size:         size,
		Available:    true,
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

	log.Info().Msgf("Downloaded %d bytes (SHA256: %x)", written, hash.Sum(nil))

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
	// For macOS, the update is typically a .pkg or .dmg
	// We'll open the installer and let the user complete the installation
	log.Info().Msg("Opening macOS installer...")

	cmd := exec.Command("open", updatePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open installer: %w", err)
	}

	// Quit the application to allow installation
	wailsruntime.Quit(u.ctx)
	return nil
}

func (u *UpdateServiceImpl) applyUpdateWindows(updatePath string) error {
	// For Windows, the update is typically an .exe installer
	// We'll run the installer and exit the current application
	log.Info().Msg("Running Windows installer...")

	cmd := exec.Command(updatePath, "/silent") // Adjust flags based on your installer
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run installer: %w", err)
	}

	// Quit the application to allow installation
	wailsruntime.Quit(u.ctx)
	return nil
}

func (u *UpdateServiceImpl) applyUpdateLinux(updatePath string) error {
	// For Linux, we need to determine the update type (AppImage, .deb, etc.)
	ext := filepath.Ext(updatePath)

	switch ext {
	case ".AppImage":
		return u.applyAppImageUpdate(updatePath)
	case ".deb":
		return u.applyDebUpdate(updatePath)
	default:
		return fmt.Errorf("unsupported update format: %s", ext)
	}
}

func (u *UpdateServiceImpl) applyAppImageUpdate(updatePath string) error {
	log.Info().Msg("Applying AppImage update...")

	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}

	// Make the new AppImage executable
	if err := os.Chmod(updatePath, 0o755); err != nil {
		return fmt.Errorf("failed to make AppImage executable: %w", err)
	}

	// Replace the current executable
	backupPath := currentExe + ".backup"
	if err := os.Rename(currentExe, backupPath); err != nil {
		return fmt.Errorf("failed to backup current executable: %w", err)
	}

	if err := os.Rename(updatePath, currentExe); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentExe)
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Restart the application
	cmd := exec.Command(currentExe, os.Args[1:]...)
	cmd.Start()

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
		return "-Installer.pkg" // macOS .pkg installer
	case "windows":
		return "-Setup.exe" // Windows installer
	case "linux":
		return ".AppImage" // Linux AppImage
	default:
		return ""
	}
}

// isVersionNewer compares two semantic versions
func isVersionNewer(current, latest string) (bool, error) {
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// Pad to ensure same length
	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}

	for i := 0; i < maxLen; i++ {
		var currentPart, latestPart int

		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentPart)
		}
		if i < len(latestParts) {
			fmt.Sscanf(latestParts[i], "%d", &latestPart)
		}

		if latestPart > currentPart {
			return true, nil
		} else if latestPart < currentPart {
			return false, nil
		}
	}

	return false, nil
}
