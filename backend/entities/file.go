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

package entities

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrReadingFile     = errors.New("error reading file")
	ErrFetchingContent = errors.New("error fetching remote file")
)

var languagesMap = map[string]string{
	"sol":  "solidity",
	"js":   "javascript",
	"ts":   "typescript",
	"tsx":  "typescript",
	"py":   "python",
	"rb":   "ruby",
	"sh":   "bash",
	"go":   "go",
	"java": "java",
	"c":    "c",
	"cpp":  "cpp",
	"h":    "c",
	"hpp":  "cpp",
	"cs":   "csharp",
	"css":  "css",
	"html": "htmlbars",
	"xml":  "xml",
	"json": "json",
	"md":   "markdown",
	"yml":  "yaml",
	"scss": "scss",
	"less": "less",
	"sass": "sass",
	"sql":  "sql",
	"txt":  "text",
}

type File struct {
	path     string
	basePath string
	content  []byte
}

func NewFile(basePath string, path string, content []byte) *File {
	return &File{
		basePath: basePath,
		path:     path,
		content:  content,
	}
}

func (f *File) GetName() string {
	return filepath.Base(f.path)
}

func (f *File) GetRelativePath() string {
	return f.path
}

func (f *File) GetAbsolutePath() string {
	return filepath.Join(f.basePath, f.path)
}

func (f *File) GetContent() []byte {
	return f.content
}

func (f *File) GetMD5Sum() string {
	hash := md5.Sum(f.content)
	return hex.EncodeToString(hash[:])
}

func (f *File) GetLanguage() string {
	fileExtension := ""
	if dotIndex := strings.LastIndex(f.path, "."); dotIndex != -1 {
		fileExtension = f.path[dotIndex+1:]
	}
	return languagesMap[fileExtension]
}

type FileDTO struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	Language string `json:"language"`
}
