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
	"encoding/json"
	"fmt"
	"slices"
)

type ComponentFilterUsage string

var (
	ScanossSettingsJson *ScanossSettings
)

type ScanossSettings struct {
	SettingsFile *SettingsFile
}

type SettingsFile struct {
	Settings ScanossSettingsSchema `json:"settings,omitempty"`
	Bom      Bom                   `json:"bom,omitempty"`
}

type ScanossSettingsSchema struct {
	Skip SkipSettings `json:"skip,omitempty"`
}

type SkipSettings struct {
	Patterns SkipPatterns `json:"patterns,omitempty"`
	Sizes    Sizes        `json:"sizes,omitempty"`
}

type SkipPatterns struct {
	Scanning       []string `json:"scanning,omitempty"`
	Fingerprinting []string `json:"fingerprinting,omitempty"`
}

type Sizes struct {
	Scanning       []SizesSkipSettings `json:"scanning,omitempty"`
	Fingerprinting []SizesSkipSettings `json:"fingerprinting,omitempty"`
}

type SizesSkipSettings struct {
	Patterns []string `json:"patterns,omitempty"`
	Min      int      `json:"min,omitempty"`
	Max      int      `json:"max,omitempty"`
}

type Bom struct {
	Include []ComponentFilter `json:"include,omitempty"`
	Remove  []ComponentFilter `json:"remove,omitempty"`
	Replace []ComponentFilter `json:"replace,omitempty"`
}

type ComponentFilter struct {
	Path        string               `json:"path,omitempty"`
	Purl        string               `json:"purl"`
	Usage       ComponentFilterUsage `json:"usage,omitempty"`
	Comment     string               `json:"comment,omitempty"`
	ReplaceWith string               `json:"replace_with,omitempty"`
	License     string               `json:"license,omitempty"`
}

type InitialFilters struct {
	Include []ComponentFilter
	Remove  []ComponentFilter
	Replace []ComponentFilter
}

func (sf *SettingsFile) Equal(other *SettingsFile) (bool, error) {
	originalSettingsFileJson, err := json.Marshal(sf)
	if err != nil {
		return false, fmt.Errorf("error parsing json: %v", err)
	}

	currentContentJSON, err := json.Marshal(other)
	if err != nil {
		return false, fmt.Errorf("error parsing json: %v", err)
	}

	hasChanges := string(originalSettingsFileJson) != string(currentContentJSON)

	return hasChanges, nil
}

func (sf *SettingsFile) GetResultWorkflowState(result Result) WorkflowState {
	included, _ := sf.IsResultIncluded(result)
	removed, _ := sf.IsResultRemoved(result)
	replaced, _ := sf.IsResultReplaced(result)

	if included || removed || replaced {
		return Completed
	}

	return Pending
}

func (sf *SettingsFile) IsResultIncluded(result Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Include)
}

func (sf *SettingsFile) IsResultRemoved(result Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Remove)
}

func (sf *SettingsFile) IsResultReplaced(result Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Replace)
}

func (sf *SettingsFile) IsResultInList(result Result, list []ComponentFilter) (bool, int) {
	i := slices.IndexFunc(list, func(cf ComponentFilter) bool {
		if cf.Path != "" && cf.Purl != "" {
			return cf.Path == result.Path && slices.Contains(*result.Purl, cf.Purl)
		}
		return slices.Contains(*result.Purl, cf.Purl)
	})

	return i != -1, i
}

func (sf *SettingsFile) GetResultFilterConfig(result Result) FilterConfig {
	var filterAction FilterAction
	var filterType FilterType

	if included, i := sf.IsResultIncluded(result); included {
		filterAction = Include
		filterType = getResultFilterType(sf.Bom.Include[i])
	} else if removed, i := sf.IsResultRemoved(result); removed {
		filterAction = Remove
		filterType = getResultFilterType(sf.Bom.Remove[i])
	} else if replaced, i := sf.IsResultReplaced(result); replaced {
		filterAction = Replace
		filterType = getResultFilterType(sf.Bom.Replace[i])
	}

	return FilterConfig{
		Action: filterAction,
		Type:   filterType,
	}
}

func getResultFilterType(cf ComponentFilter) FilterType {
	if cf.Path != "" && cf.Purl != "" {
		return ByFile
	}
	return ByPurl
}

func (sf *SettingsFile) GetBomEntryFromResult(result Result) ComponentFilter {
	if included, i := sf.IsResultIncluded(result); included {
		return sf.Bom.Include[i]
	}

	if removed, i := sf.IsResultRemoved(result); removed {
		return sf.Bom.Remove[i]
	}

	if replaced, i := sf.IsResultReplaced(result); replaced {
		return sf.Bom.Replace[i]
	}

	return ComponentFilter{}
}

var DefaultSkippedDirs = []string{
	"nbproject",
	"nbbuild",
	"nbdist",
	"__pycache__",
	"venv",
	"_yardoc",
	"eggs",
	"wheels",
	"htmlcov",
	"__pypackages__",
	"example",
	"examples",
	"docs",
	"doc",
}

var DefaultSkippedFiles = []string{
	"gradlew",
	"gradlew.bat",
	"mvnw",
	"mvnw.cmd",
	"gradle-wrapper.jar",
	"maven-wrapper.jar",
	"thumbs.db",
	"babel.config.js",
	"license.txt",
	"license.md",
	"copying.lib",
	"makefile",
}

var DefaultSkippedDirExtensions = []string{".egg-info"}

var DefaultSkippedExtensions = []string{
	".1",
	".2",
	".3",
	".4",
	".5",
	".6",
	".7",
	".8",
	".9",
	".ac",
	".adoc",
	".am",
	".asciidoc",
	".bmp",
	".build",
	".cfg",
	".chm",
	".class",
	".cmake",
	".cnf",
	".conf",
	".config",
	".contributors",
	".copying",
	".crt",
	".csproj",
	".css",
	".csv",
	".dat",
	".data",
	".doc",
	".docx",
	".dtd",
	".dts",
	".iws",
	".c9",
	".c9revisions",
	".dtsi",
	".dump",
	".eot",
	".eps",
	".geojson",
	".gdoc",
	".gif",
	".glif",
	".gmo",
	".gradle",
	".guess",
	".hex",
	".htm",
	".html",
	".ico",
	".iml",
	".in",
	".inc",
	".info",
	".ini",
	".ipynb",
	".jpeg",
	".jpg",
	".json",
	".jsonld",
	".lock",
	".log",
	".m4",
	".map",
	".markdown",
	".md",
	".md5",
	".meta",
	".mk",
	".mxml",
	".o",
	".otf",
	".out",
	".pbtxt",
	".pdf",
	".pem",
	".phtml",
	".plist",
	".png",
	".po",
	".ppt",
	".prefs",
	".properties",
	".pyc",
	".qdoc",
	".result",
	".rgb",
	".rst",
	".scss",
	".sha",
	".sha1",
	".sha2",
	".sha256",
	".sln",
	".spec",
	".sql",
	".sub",
	".svg",
	".svn-base",
	".tab",
	".template",
	".test",
	".tex",
	".tiff",
	".toml",
	".ttf",
	".txt",
	".utf-8",
	".vim",
	".wav",
	".woff",
	".woff2",
	".xht",
	".xhtml",
	".xls",
	".xlsx",
	".xml",
	".xpm",
	".xsd",
	".xul",
	".yaml",
	".yml",
	".wfp",
	".editorconfig",
	".dotcover",
	".pid",
	".lcov",
	".egg",
	".manifest",
	".cache",
	".coverage",
	".cover",
	".gem",
	".lst",
	".pickle",
	".pdb",
	".gml",
	".pot",
	".plt",
	".whml",
	".pom",
	".smtml",
	".min.js",
	".mf",
	".base64",
	".s",
	".diff",
	".patch",
	".rules",
	//  File endings
	"-doc",
	"changelog",
	"config",
	"copying",
	"license",
	"authors",
	"news",
	"licenses",
	"notice",
	"readme",
	"swiftdoc",
	"texidoc",
	"todo",
	"version",
	"ignore",
	"manifest",
	"sqlite",
	"sqlite3",
}
