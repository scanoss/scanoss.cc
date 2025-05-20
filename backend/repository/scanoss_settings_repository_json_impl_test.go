// go:build unit

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

package repository_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/internal/config"
	utilsMocks "github.com/scanoss/scanoss.cc/internal/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var SCANOSS_SETTINGS_PATH = "scanoss.json"

// setupTest initializes the testing environment and returns a cleanup function
func setupTest(t *testing.T) (*utilsMocks.MockFileReader, func()) {
	t.Helper()

	// Create a mock file reader
	mockFileReader := new(utilsMocks.MockFileReader)

	// Initialize the test environment
	cleanup := func() {
		// Reset global state
		entities.ScanossSettingsJson = nil
	}

	return mockFileReader, cleanup
}

// setupConfig creates a temporary configuration for testing
func setupConfig(t *testing.T, settingsPath string) {
	t.Helper()

	cfg := config.GetInstance()
	cfg.SetScanRoot(t.TempDir())
	cfg.SetScanSettingsFilePath(settingsPath)
}

// createMockSettingsFile creates a mock settings file content
func createMockSettingsFile() entities.SettingsFile {
	return entities.SettingsFile{
		Bom: entities.Bom{
			Include: []entities.ComponentFilter{
				{Purl: "pkg:npm/test1@1.0.0", Path: "test1/path"},
			},
			Remove: []entities.ComponentFilter{
				{Purl: "pkg:npm/test2@1.0.0", Path: "test2/path"},
			},
			Replace: []entities.ComponentFilter{
				{Purl: "pkg:npm/test3@1.0.0", Path: "test3/path", ReplaceWith: "pkg:npm/test3-replacement@1.0.0"},
			},
		},
		Settings: entities.ScanossSettingsSchema{
			Skip: entities.SkipSettings{
				Patterns: entities.SkipPatterns{
					Scanning: []string{"node_modules/", "*.min.js", "custom/pattern"},
				},
			},
		},
	}
}

func TestScanossSettingsRepositoryJsonImpl(t *testing.T) {
	// TestNewScanossSettingsJsonRepository tests creation of the repository
	t.Run("TestNewScanossSettingsJsonRepository", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		assert.NotNil(t, repo, "Repository should not be nil")
	})

	t.Run("TestInit", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		// Setup mock for Read operation
		settingsContent, _ := json.Marshal(createMockSettingsFile())
		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		// Verify that the settings were loaded
		assert.NotNil(t, entities.ScanossSettingsJson, "ScanossSettingsJson should not be nil after initialization")
		mockFileReader.AssertExpectations(t)
	})

	// TestRead tests reading settings from a file
	t.Run("TestRead", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		// Create test settings
		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		// Setup mock
		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		// Create repository
		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		settings, err := repo.Read()

		// Verify results
		assert.NoError(t, err, "Read should not return an error")
		assert.Equal(t, mockSettings.Bom.Include[0].Purl, settings.Bom.Include[0].Purl, "Include purls should match")
		assert.Equal(t, mockSettings.Bom.Remove[0].Purl, settings.Bom.Remove[0].Purl, "Remove purls should match")
		assert.Equal(t, mockSettings.Bom.Replace[0].Purl, settings.Bom.Replace[0].Purl, "Replace purls should match")
		assert.Equal(t, mockSettings.Settings.Skip.Patterns.Scanning, settings.Settings.Skip.Patterns.Scanning, "Scanning patterns should match")

		mockFileReader.AssertExpectations(t)
	})

	// TestReadFileNotFound tests reading when file doesn't exist
	t.Run("TestReadFileNotFound", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), "nonexistent.json")
		setupConfig(t, settingsPath)

		// Setup mock to return a file not found error
		mockFileReader.On("ReadFile", settingsPath).Return([]byte{}, os.ErrNotExist)

		// Create repository
		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		settings, err := repo.Read()

		// Verify results - should return empty settings without error
		assert.NoError(t, err, "Read should not return an error when file doesn't exist")
		assert.Empty(t, settings.Bom.Include, "Include list should be empty")
		assert.Empty(t, settings.Bom.Remove, "Remove list should be empty")
		assert.Empty(t, settings.Bom.Replace, "Replace list should be empty")

		mockFileReader.AssertExpectations(t)
	})

	// TestReadInvalidJSON tests reading invalid JSON content
	t.Run("TestReadInvalidJSON", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), "invalid.json")
		setupConfig(t, settingsPath)

		// Setup mock to return invalid JSON
		mockFileReader.On("ReadFile", settingsPath).Return([]byte("{invalid json}"), nil)

		// Create repository
		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		_, err := repo.Read()

		// Verify results
		assert.Error(t, err, "Read should return an error for invalid JSON")

		mockFileReader.AssertExpectations(t)
	})

	// TestSave tests saving settings to a file
	t.Run("TestSave", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		assert.NotPanics(t, func() { repo.Save() }, "Save should not panic")

		mockFileReader.AssertExpectations(t)
	})

	// TestGetDeclaredPurls tests getting all declared PURLs
	t.Run("TestGetDeclaredPurls", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		purls := repo.GetDeclaredPurls()

		assert.Contains(t, purls, "pkg:npm/test1@1.0.0", "Should contain include PURL")
		assert.Contains(t, purls, "pkg:npm/test2@1.0.0", "Should contain remove PURL")
		assert.Contains(t, purls, "pkg:npm/test3@1.0.0", "Should contain replace PURL")
		assert.Len(t, purls, 3, "Should have 3 PURLs total")

		mockFileReader.AssertExpectations(t)
	})

	// TestAddBomEntry tests adding a BOM entry
	t.Run("TestAddBomEntry", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")
		newEntry := entities.ComponentFilter{
			Purl: "pkg:npm/new-component@1.0.0",
			Path: "new-component/path",
		}
		err = repo.AddBomEntry(newEntry, "include")
		assert.NoError(t, err, "AddBomEntry should not return an error")

		// Add an entry
		settings := repo.GetSettings()
		assert.Len(t, settings.Bom.Include, 2, "Include list should have 2 entries")
		assert.Equal(t, "pkg:npm/new-component@1.0.0", settings.Bom.Include[1].Purl, "New PURL should be added")

		// Add an entry to remove list
		newEntry = entities.ComponentFilter{
			Purl: "pkg:npm/new-remove@1.0.0",
			Path: "new-remove/path",
		}
		err = repo.AddBomEntry(newEntry, "remove")
		assert.NoError(t, err, "AddBomEntry should not return an error")

		// Verify the entry was added
		settings = repo.GetSettings()
		assert.Len(t, settings.Bom.Remove, 2, "Remove list should have 2 entries")
		assert.Equal(t, "pkg:npm/new-remove@1.0.0", settings.Bom.Remove[1].Purl, "New PURL should be added")

		// Add an entry to replace list
		newEntry = entities.ComponentFilter{
			Purl:        "pkg:npm/new-replace@1.0.0",
			Path:        "new-replace/path",
			ReplaceWith: "pkg:npm/replacement@1.0.0",
		}
		err = repo.AddBomEntry(newEntry, "replace")
		assert.NoError(t, err, "AddBomEntry should not return an error")

		// Verify the entry was added
		settings = repo.GetSettings()
		assert.Len(t, settings.Bom.Replace, 2, "Replace list should have 2 entries")
		assert.Equal(t, "pkg:npm/new-replace@1.0.0", settings.Bom.Replace[1].Purl, "New PURL should be added")

		// Test invalid filter action
		err = repo.AddBomEntry(newEntry, "invalid")
		assert.Error(t, err, "AddBomEntry should return an error for invalid action")

		mockFileReader.AssertExpectations(t)
	})

	// TestClearAllFilters tests clearing all BOM filters
	t.Run("TestClearAllFilters", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		err = repo.ClearAllFilters()
		assert.NoError(t, err, "ClearAllFilters should not return an error")

		settings := repo.GetSettings()
		assert.Empty(t, settings.Bom.Include, "Include list should be empty")
		assert.Empty(t, settings.Bom.Remove, "Remove list should be empty")
		assert.Empty(t, settings.Bom.Replace, "Replace list should be empty")

		mockFileReader.AssertExpectations(t)
	})

	// TestAddStagedScanningSkipPattern tests adding a staged scanning skip pattern
	t.Run("TestAddStagedScanningSkipPattern", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		testCases := []struct {
			name            string
			inputPattern    string
			expectedPattern string
		}{
			{
				name:            "Unix path",
				inputPattern:    "unix/style/path/",
				expectedPattern: "unix/style/path/",
			},
			{
				name:            "Windows path",
				inputPattern:    "windows\\style\\path\\",
				expectedPattern: "windows/style/path/",
			},
			{
				name:            "Mixed path",
				inputPattern:    "mixed/style\\path/",
				expectedPattern: "mixed/style/path/",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err = repo.AddStagedScanningSkipPattern(tc.inputPattern)
				assert.NoError(t, err, "AddStagedScanningSkipPattern should not return an error")

				assert.True(t, repo.HasStagedScanningSkipPatternChanges(), "Should have staged changes")

				effectivePatterns := repo.GetEffectiveScanningSkipPatterns()
				assert.Contains(t, effectivePatterns, tc.expectedPattern, "Effective patterns should contain our normalized pattern")
			})
		}

		mockFileReader.AssertExpectations(t)
	})

	// TestRemoveStagedScanningSkipPattern tests removing a staged scanning skip pattern
	t.Run("TestRemoveStagedScanningSkipPattern", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		testCases := []struct {
			name            string
			patternToRemove string
		}{
			{
				name:            "Remove Unix pattern",
				patternToRemove: "custom/pattern",
			},
			{
				name:            "Remove Windows path format",
				patternToRemove: "node_modules\\",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err = repo.RemoveStagedScanningSkipPattern(tc.patternToRemove)
				assert.NoError(t, err, "RemoveStagedScanningSkipPattern should not return an error")

				assert.True(t, repo.HasStagedScanningSkipPatternChanges(), "Should have staged changes")
			})
		}

		mockFileReader.AssertExpectations(t)
	})

	// TestCommitStagedScanningSkipPatterns tests committing staged patterns
	t.Run("TestCommitStagedScanningSkipPatterns", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		newPattern := "new/pattern/"
		err = repo.AddStagedScanningSkipPattern(newPattern)
		assert.NoError(t, err, "AddStagedScanningSkipPattern should not return an error")
		assert.True(t, repo.HasStagedScanningSkipPatternChanges(), "Should have staged changes")

		err = repo.CommitStagedScanningSkipPatterns()
		assert.NoError(t, err, "CommitStagedScanningSkipPatterns should not return an error")

		assert.False(t, repo.HasStagedScanningSkipPatternChanges(), "Should not have staged changes after commit")

		mockFileReader.AssertExpectations(t)
	})

	// TestDiscardStagedScanningSkipPatterns tests discarding staged patterns
	t.Run("TestDiscardStagedScanningSkipPatterns", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		newPattern := "pattern/to/discard/"
		err = repo.AddStagedScanningSkipPattern(newPattern)
		assert.NoError(t, err, "AddStagedScanningSkipPattern should not return an error")
		assert.True(t, repo.HasStagedScanningSkipPatternChanges(), "Should have staged changes")
		err = repo.DiscardStagedScanningSkipPatterns()
		assert.NoError(t, err, "DiscardStagedScanningSkipPatterns should not return an error")

		assert.False(t, repo.HasStagedScanningSkipPatternChanges(), "Should not have staged changes after discard")

		mockFileReader.AssertExpectations(t)
	})

	// TestMatchesEffectiveScanningSkipPattern tests pattern matching with different path formats
	t.Run("TestMatchesEffectiveScanningSkipPattern", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), "settings.json")
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		mockSettings.Settings.Skip.Patterns.Scanning = append(
			mockSettings.Settings.Skip.Patterns.Scanning,
			"build/",
			"*.log",
			"custom/",
		)
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		// Create repository and initialize
		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		tempDir := t.TempDir()
		nodeModulesDir := filepath.Join(tempDir, "node_modules")
		buildDir := filepath.Join(tempDir, "build")
		customDir := filepath.Join(tempDir, "custom")

		for _, dir := range []string{nodeModulesDir, buildDir, customDir} {
			err = os.MkdirAll(dir, 0755)
			require.NoError(t, err, "Failed to create test directory structure")
		}

		minJsPath := filepath.Join(tempDir, "file.min.js")
		regularJsPath := filepath.Join(tempDir, "file.js")
		logFilePath := filepath.Join(tempDir, "output.log")
		nodeModulesFile := filepath.Join(nodeModulesDir, "package.json")

		for _, path := range []string{minJsPath, regularJsPath, logFilePath, nodeModulesFile} {
			err = os.WriteFile(path, []byte("test"), 0644)
			require.NoError(t, err, "Failed to create test file")
		}

		err = os.MkdirAll(customDir, 0755)
		require.NoError(t, err)

		customPatternPath := filepath.Join(customDir, "pattern.go")
		err = os.WriteFile(customPatternPath, []byte("test"), 0644)
		require.NoError(t, err)

		// Test cases with different path formats that should work on both Windows and Unix
		testCases := []struct {
			name        string
			path        string
			shouldMatch bool
			description string
		}{
			// Test node_modules pattern
			{
				name:        "node_modules directory with native separators",
				path:        nodeModulesDir,
				shouldMatch: true,
				description: "Should match node_modules directory pattern with native separators",
			},
			{
				name:        "node_modules directory with slash separators",
				path:        filepath.ToSlash(nodeModulesDir),
				shouldMatch: true,
				description: "Should match node_modules directory pattern with slash separators",
			},
			{
				name:        "node_modules file with native separators",
				path:        nodeModulesFile,
				shouldMatch: true,
				description: "Should match file inside node_modules directory",
			},

			// Test *.min.js pattern
			{
				name:        "minified JS file with native separators",
				path:        minJsPath,
				shouldMatch: true,
				description: "Should match *.min.js pattern with native separators",
			},
			{
				name:        "minified JS file with slash separators",
				path:        filepath.ToSlash(minJsPath),
				shouldMatch: true,
				description: "Should match *.min.js pattern with slash separators",
			},

			// Test that non-matching files don't match
			{
				name:        "Regular JS file (should not match)",
				path:        regularJsPath,
				shouldMatch: false,
				description: "Should not match regular JS file",
			},

			// Test custom pattern
			{
				name:        "Custom pattern with native separators",
				path:        customPatternPath,
				shouldMatch: true,
				description: "Should match custom pattern with native separators",
			},
			{
				name:        "Custom pattern with slash separators",
				path:        filepath.ToSlash(customPatternPath),
				shouldMatch: true,
				description: "Should match custom pattern with slash separators",
			},

			// Test additional patterns
			{
				name:        "Build directory",
				path:        buildDir,
				shouldMatch: true,
				description: "Should match build/ pattern",
			},
			{
				name:        "Log file",
				path:        logFilePath,
				shouldMatch: true,
				description: "Should match *.log pattern",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test if path matches any of the effective scanning skip patterns
				matches := repo.MatchesEffectiveScanningSkipPattern(tc.path)

				log.Info().Msgf("Path %s matches: %v", tc.path, matches)

				// On Windows, we might need to adjust the expected result if the test is affected
				// by Windows-specific path handling behavior
				expectedMatch := tc.shouldMatch

				if runtime.GOOS == "windows" {
					// If there are any Windows-specific adjustments needed, handle them here
					// For example, some patterns might not match correctly on Windows if they
					// use forward slashes in a way Windows doesn't recognize

					// This is a placeholder for any platform-specific adjustments
					// In most cases, the repository implementation should handle this properly
				}

				assert.Equal(t, expectedMatch, matches, tc.description)
			})
		}

		// Test path conversion specifically
		t.Run("Path conversion handling", func(t *testing.T) {
			// Test with backslashes on any platform
			windowsStylePath := strings.ReplaceAll(filepath.Join(tempDir, "node_modules", "some", "package"), "/", "\\")
			matches := repo.MatchesEffectiveScanningSkipPattern(windowsStylePath)
			assert.True(t, matches, "Windows-style path with backslashes should match node_modules pattern")

			// Test with forward slashes on any platform
			unixStylePath := filepath.ToSlash(filepath.Join(tempDir, "node_modules", "some", "package"))
			matches = repo.MatchesEffectiveScanningSkipPattern(unixStylePath)
			assert.True(t, matches, "Unix-style path with forward slashes should match node_modules pattern")

			// Test mixed slashes on any platform
			mixedPath := filepath.Join(tempDir, "node_modules") + "\\some/package"
			matches = repo.MatchesEffectiveScanningSkipPattern(mixedPath)
			assert.True(t, matches, "Mixed slash path should match node_modules pattern")
		})

		mockFileReader.AssertExpectations(t)
	})

	// TestCrossOSPathHandling tests specific cross-OS path handling
	t.Run("TestCrossOSPathHandling", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil)

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		// Test adding patterns with different OS path formats
		t.Run("AddPatternWithDifferentOSFormats", func(t *testing.T) {
			// Windows-style path
			windowsPath := "windows\\style\\path\\file.txt"
			err = repo.AddStagedScanningSkipPattern(windowsPath)
			assert.NoError(t, err)

			// Unix-style path
			unixPath := "unix/style/path/file.txt"
			err = repo.AddStagedScanningSkipPattern(unixPath)
			assert.NoError(t, err)

			// Mixed-style path
			mixedPath := "mixed/style\\path/file.txt"
			err = repo.AddStagedScanningSkipPattern(mixedPath)
			assert.NoError(t, err)

			// Get effective patterns and check normalization
			effectivePatterns := repo.GetEffectiveScanningSkipPatterns()
			assert.Contains(t, effectivePatterns, "windows/style/path/file.txt", "Windows path should be normalized to forward slashes")
			assert.Contains(t, effectivePatterns, "unix/style/path/file.txt", "Unix path should remain as is")
			assert.Contains(t, effectivePatterns, "mixed/style/path/file.txt", "Mixed path should be normalized to forward slashes")
		})

		// Test matching patterns with different OS path formats
		t.Run("MatchPathWithDifferentOSFormats", func(t *testing.T) {
			// Add a specific pattern for testing
			pattern := "path/subpath/*.go"
			err = repo.AddStagedScanningSkipPattern(pattern)
			assert.NoError(t, err)

			// Test with different path formats
			testPaths := []struct {
				path        string
				shouldMatch bool
			}{
				{"path/subpath/file.go", true},         // Unix format, should match
				{"path\\subpath\\file.go", true},       // Windows format, should match when normalized
				{"path/subpath/subdir/file.go", false}, // Shouldn't match subdirectory
				{"different/path/file.go", false},      // Different path, shouldn't match
			}

			for _, tp := range testPaths {
				// Since there might be OS-specific behavior, we'll check this in a platform-neutral way
				matches := repo.MatchesEffectiveScanningSkipPattern(tp.path)
				assert.Equal(t, tp.shouldMatch, matches, fmt.Sprintf("Path %s should %s", tp.path, map[bool]string{true: "match", false: "not match"}[tp.shouldMatch]))
			}
		})

		mockFileReader.AssertExpectations(t)
	})

	// TestHasUnsavedChanges tests checking for unsaved changes
	t.Run("TestHasUnsavedChanges", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		mockSettings := createMockSettingsFile()
		settingsContent, _ := json.Marshal(mockSettings)

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil).Once()

		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		err := repo.Init()
		assert.NoError(t, err, "Init should not return an error")

		newEntry := entities.ComponentFilter{
			Purl: "pkg:npm/new-component@1.0.0",
			Path: "new-component/path",
		}
		err = repo.AddBomEntry(newEntry, "include")
		assert.NoError(t, err, "AddBomEntry should not return an error")

		mockFileReader.On("ReadFile", settingsPath).Return(settingsContent, nil).Once()

		_, err = repo.HasUnsavedChanges()
		assert.NoError(t, err, "HasUnsavedChanges should not return an error")

		mockFileReader.AssertExpectations(t)
	})

	// TestReadError tests handling read errors
	t.Run("TestReadError", func(t *testing.T) {
		mockFileReader, cleanup := setupTest(t)
		defer cleanup()

		// Setup config
		settingsPath := filepath.Join(t.TempDir(), SCANOSS_SETTINGS_PATH)
		setupConfig(t, settingsPath)

		// Setup mock to return an error
		mockErr := errors.New("read error")
		mockFileReader.On("ReadFile", settingsPath).Return([]byte{}, mockErr)

		// Create repository
		repo := repository.NewScanossSettingsJsonRepository(mockFileReader)
		_, err := repo.Read()

		// Verify results
		assert.Error(t, err, "Read should return an error when ReadFile fails")
		assert.Equal(t, mockErr, err, "Error should be the same as the mock error")

		mockFileReader.AssertExpectations(t)
	})

}
