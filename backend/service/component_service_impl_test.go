//go:build unit

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

package service_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/backend/service"
	internal_test "github.com/scanoss/scanoss.cc/internal"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComponentServiceUndoRedoAfterOpeningFolder reproduces the bug where, after
// opening a folder at runtime whose scanoss.json already has identified
// (completed) files, including a pending file and then undoing it wiped ALL
// completed files back to pending instead of only the file that was just
// included.
//
// The root cause was that ComponentServiceImpl captured its initialFilters
// snapshot once at construction time. Opening a folder at runtime reloads the
// settings repository's BOM but used to leave the service's snapshot stale, so
// Undo's clear-and-replay restored nothing.
func TestComponentServiceUndoRedoAfterOpeningFolder(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	mu := internal_test.NewMockUtils()
	cfg := config.GetInstance()

	// Startup state: a settings file with no BOM entries (no folder opened yet).
	startupPath := filepath.Join(t.TempDir(), "scanoss.json")
	emptySettings, _ := json.Marshal(entities.SettingsFile{})
	mu.On("ReadFile", startupPath).Return(emptySettings, nil)
	cfg.SetScanSettingsFilePath(startupPath)

	settingsRepo := repository.NewScanossSettingsJsonRepository(mu)
	require.NoError(t, settingsRepo.Init())

	// The component service captures its initial filters snapshot here, while
	// the BOM is still empty, and registers a config-change listener.
	svc := service.NewComponentServiceImpl(nil, settingsRepo, nil, nil, nil)

	// Open a folder at runtime: its scanoss.json already has one identified file
	// (an originally-"completed" file living in the Include BOM).
	const originalPath = "vendor/already-identified.go"
	openedSettings, _ := json.Marshal(entities.SettingsFile{
		Bom: entities.Bom{
			Include: []entities.ComponentFilter{
				{Purl: "pkg:github/scanoss/original@1.0.0", Path: originalPath},
			},
		},
	})
	openedPath := filepath.Join(t.TempDir(), "opened", "scanoss.json")
	mu.On("ReadFile", openedPath).Return(openedSettings, nil)

	// Selecting the folder updates the config, which notifies listeners: the
	// settings repo reloads the new BOM and the component service refreshes its
	// initial filters snapshot.
	cfg.SetScanSettingsFilePath(openedPath)

	// Sanity check: the originally-completed file is present after opening.
	require.True(t, bomIncludes(settingsRepo, originalPath),
		"originally-completed file should be in the BOM after opening the folder")

	// Select a pending file and Include it -> it becomes completed.
	const newPath = "src/new-file.go"
	require.NoError(t, svc.FilterComponents([]entities.ComponentFilterDTO{
		{Path: newPath, Purl: "pkg:github/scanoss/new@1.0.0", Action: entities.Include},
	}))
	require.True(t, bomIncludes(settingsRepo, newPath), "newly included file should be completed")

	// Undo the inclusion.
	require.NoError(t, svc.Undo())

	// The newly included file goes back to pending...
	assert.False(t, bomIncludes(settingsRepo, newPath),
		"undone file should move back to pending")
	// ...but the originally-completed file must REMAIN completed (this is the bug).
	assert.True(t, bomIncludes(settingsRepo, originalPath),
		"originally-completed file must stay completed after undoing an unrelated include")

	// Redo restores the included file without disturbing the original.
	require.NoError(t, svc.Redo())
	assert.True(t, bomIncludes(settingsRepo, newPath), "redo should re-include the file")
	assert.True(t, bomIncludes(settingsRepo, originalPath),
		"originally-completed file must stay completed after redo")
}

func bomIncludes(repo repository.ScanossSettingsRepository, path string) bool {
	for _, entry := range repo.GetSettings().Bom.Include {
		if entry.Path == path {
			return true
		}
	}
	return false
}
