//go:build unit

// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2026 SCANOSS.COM
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

package config_test

import (
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initConfig(t *testing.T, scanRoot, inputFile, settingsFile, workDir string) *config.Config {
	t.Helper()
	config.ResetInstance()
	cfg := config.GetInstance()
	err := cfg.InitializeConfig("", scanRoot, "", "", inputFile, settingsFile, workDir, false)
	require.NoError(t, err)
	return cfg
}

func TestInitializePathConfig(t *testing.T) {
	workDir := "/tmp/workdir"

	t.Run("--input without --scan-root is not overridden", func(t *testing.T) {
		customInput := "/custom/path/results.json"
		cfg := initConfig(t, "", customInput, "", workDir)

		assert.Equal(t, customInput, cfg.GetResultFilePath(),
			"--input value should not be overridden when --scan-root is omitted")
		assert.Equal(t, workDir, cfg.GetScanRoot(),
			"scan root should default to workDir when --scan-root is omitted")
	})

	t.Run("--settings without --scan-root is not overridden", func(t *testing.T) {
		customSettings := "/custom/path/scanoss.json"
		cfg := initConfig(t, "", "", customSettings, workDir)

		assert.Equal(t, customSettings, cfg.GetScanSettingsFilePath(),
			"--settings value should not be overridden when --scan-root is omitted")
	})

	t.Run("--input with --scan-root is respected", func(t *testing.T) {
		customInput := "/custom/path/results.json"
		customRoot := "/custom/root"
		cfg := initConfig(t, customRoot, customInput, "", workDir)

		assert.Equal(t, customInput, cfg.GetResultFilePath(),
			"--input value should be respected when --scan-root is also provided")
		assert.Equal(t, customRoot, cfg.GetScanRoot())
	})

	t.Run("--settings with --scan-root is respected", func(t *testing.T) {
		customSettings := "/custom/path/scanoss.json"
		customRoot := "/custom/root"
		cfg := initConfig(t, customRoot, "", customSettings, workDir)

		assert.Equal(t, customSettings, cfg.GetScanSettingsFilePath(),
			"--settings value should be respected when --scan-root is also provided")
	})

	t.Run("defaults are used when no flags provided", func(t *testing.T) {
		cfg := initConfig(t, "", "", "", workDir)

		expectedResults := filepath.Join(workDir, ".scanoss", "results.json")
		expectedSettings := filepath.Join(workDir, "scanoss.json")

		assert.Equal(t, workDir, cfg.GetScanRoot())
		assert.Equal(t, expectedResults, cfg.GetResultFilePath())
		assert.Equal(t, expectedSettings, cfg.GetScanSettingsFilePath())
	})

	t.Run("--scan-root defaults results and settings to that root", func(t *testing.T) {
		customRoot := "/custom/root"
		cfg := initConfig(t, customRoot, "", "", workDir)

		expectedResults := filepath.Join(customRoot, ".scanoss", "results.json")
		expectedSettings := filepath.Join(customRoot, "scanoss.json")

		assert.Equal(t, customRoot, cfg.GetScanRoot())
		assert.Equal(t, expectedResults, cfg.GetResultFilePath())
		assert.Equal(t, expectedSettings, cfg.GetScanSettingsFilePath())
	})
}
