package utils

import (
	"encoding/json"
	"os"
	"testing"

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

	data, err := ReadFile(tmpfile.Name())
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

	parsedData, err := JSONParse[map[string]string](jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result = parsedData
	if result["key"] != "value" {
		t.Fatalf("Expected value to be 'value', got %s", result["key"])
	}

	invalidJsonData := []byte(`{"key":"value}`)
	_, err = JSONParse[map[string]string](invalidJsonData)
	require.Error(t, err)
}

func TestWriteJsonFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	data := map[string]string{"key": "value"}
	err = WriteJsonFile(tmpfile.Name(), data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	readData, err := ReadFile(tmpfile.Name())
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

	err = FileExist(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = FileExist("nonexistentfile.json")
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

	if !IsWritableFile(tmpfile.Name()) {
		t.Fatalf("Expected file to be writable")
	}

	// Create a read-only file
	readonlyFile, err := os.CreateTemp("", "readonly.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(readonlyFile.Name())

	if err := readonlyFile.Chmod(0444); err != nil {
		t.Fatal(err)
	}

	if IsWritableFile(readonlyFile.Name()) {
		t.Fatalf("Expected file to be read-only")
	}
}
