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

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type DefaultFileReader struct{}

func NewDefaultFileReader() *DefaultFileReader {
	return &DefaultFileReader{}
}

func (d *DefaultFileReader) ReadFile(filePath string) ([]byte, error) {
	// Open and read the JSON file
	expandedPath := ExpandPath(filePath)
	file, err := os.Open(expandedPath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return byteValue, nil
}

func JSONParse[T any](file []byte) (T, error) {
	var intermediateMap T

	if err := json.Unmarshal(file, &intermediateMap); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling JSON")
		return intermediateMap, err
	}

	return intermediateMap, nil
}

func WriteJsonFile(path string, in any) error {
	j, err := JSONSerialize(in)
	if err != nil {
		return err
	}

	err = WriteFile(path, j)
	if err != nil {
		return err
	}

	return nil
}

func JSONSerialize(in any) ([]byte, error) {
	out, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

func WriteFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func FileExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func IsWritableFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fileInfo.Mode()

	return mode&0200 != 0
}

func ExpandPath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if path == "~" {
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path
}

func GetRelativePath(absolutePath string) (string, error) {
	if !filepath.IsAbs(absolutePath) {
		return absolutePath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	relPath, err := filepath.Rel(cwd, absolutePath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	return relPath, nil
}
