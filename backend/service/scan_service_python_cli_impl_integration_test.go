//go:build integration

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
