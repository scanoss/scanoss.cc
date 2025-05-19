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

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
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

// fullySplitPath splits a path into ALL its components.
// The gitignore library expects to be handed a slice of everything busted ALL the way apart, but
// the core filepath.Split ONLY breaks off the final path component.
//
// We theoretically could split on our platform-native path separator BUT it's not quite that simple: on windows,
// both forward AND back slashes as path separators, plus there's special handling of potential volume specifiers.
//
// We COULD mimic filepath.Split's implementation stepping through and checking IsPathSeparator(), but to keep
// it simple, we'll just repeatedly call that directly. These are all short enough it isn't that expensive.
func FullySplitPath(path string) (split []string) {
	for path != "" {
		dir, file := filepath.Split(filepath.Clean(path))
		if file == "" {
			break
		}
		split = append(split, file)
		path = dir
	}

	slices.Reverse(split)

	return
}

// We use this function to normalize paths on windows/unix systems, converting all separators to forward slashes and removing "." and ".."
func NormalizePathToSlash(p string) string {
	if p == "" {
		return p
	}

	return path.Clean(filepath.ToSlash(p))
}
