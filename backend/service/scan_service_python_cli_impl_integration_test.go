//go:build integration

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

		args := []string{"--files", tmpfile.Name()}
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

		args := []string{"--files", testFile}
		err = scanService.ScanStream(args)
		assert.NoError(t, err)
	})
}
