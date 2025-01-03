//go:build integration

// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/stretchr/testify/assert"
)

func TestScanServicePythonImpl_Integration(t *testing.T) {
	scanService := service.NewScanServicePythonImpl()

	t.Run("CheckDependencies", func(t *testing.T) {
		err := scanService.CheckDependencies()
		assert.NoError(t, err)
	})

	t.Run("Scan", func(t *testing.T) {
		content := []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}")
		tmpfile, err := os.CreateTemp("", "example.*.go")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		args := []string{"--files", tmpfile.Name(), "--no-wfp-output"}
		err = scanService.Scan(args)
		assert.NoError(t, err)
	})

	t.Run("ScanStream", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "scanoss_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		testFile := filepath.Join(tmpDir, "test.go")
		content := []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatal(err)
		}

		args := []string{"--files", testFile, "--no-wfp-output"}
		err = scanService.ScanStream(args)
		assert.NoError(t, err)
	})
}
