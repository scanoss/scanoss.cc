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
package utils_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/scanoss/scanoss.cc/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte(`{"key": "value"}`)
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	fr := &utils.DefaultFileReader{}
	data, err := fr.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(data) != string(content) {
		t.Fatalf("Expected %s, got %s", string(content), string(data))
	}
}

func TestJSONParse(t *testing.T) {
	jsonData := []byte(`{"key": "value"}`)
	var result map[string]string

	parsedData, err := utils.JSONParse[map[string]string](jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result = parsedData
	if result["key"] != "value" {
		t.Fatalf("Expected value to be 'value', got %s", result["key"])
	}

	invalidJsonData := []byte(`{"key":"value}`)
	_, err = utils.JSONParse[map[string]string](invalidJsonData)
	require.Error(t, err)
}

func TestWriteJsonFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	data := map[string]string{"key": "value"}
	err = utils.WriteJsonFile(tmpfile.Name(), data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	fr := &utils.DefaultFileReader{}
	readData, err := fr.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]string
	err = json.Unmarshal(readData, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["key"] != "value" {
		t.Fatalf("Expected value to be 'value', got %s", result["key"])
	}
}

func TestFileExist(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = utils.FileExist(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = utils.FileExist("nonexistentfile.json")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestIsWritableFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if !utils.IsWritableFile(tmpfile.Name()) {
		t.Fatalf("Expected file to be writable")
	}

	readonlyFile, err := os.CreateTemp("", "readonly.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(readonlyFile.Name())

	if err := readonlyFile.Chmod(0444); err != nil {
		t.Fatal(err)
	}

	if utils.IsWritableFile(readonlyFile.Name()) {
		t.Fatalf("Expected file to be read-only")
	}
}
