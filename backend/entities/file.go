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
