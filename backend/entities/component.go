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

type ComponentLicense interface {
}

type LicenseDTO struct {
	Name             string `json:"name"`
	PatentHints      string `json:"patent_hints"`
	Copyleft         string `json:"copyleft"`
	ChecklistURL     string `json:"checklist_url"`
	OsadlUpdated     string `json:"osadl_updated"`
	Source           string `json:"source"`
	URL              string `json:"url"`
	IncompatibleWith string `json:"incompatible_with,omitempty"`
}

func NewLicenseDTO(name string, patentHints string, copyleft string, checklistURL string, osadlUpdated string, source string, url string, incompatibleWith string) *LicenseDTO {
	return &LicenseDTO{
		Name:             name,
		URL:              url,
		ChecklistURL:     checklistURL,
		Copyleft:         copyleft,
		PatentHints:      patentHints,
		IncompatibleWith: incompatibleWith,
		OsadlUpdated:     osadlUpdated,
		Source:           source,
	}
}

type ComponentDTO struct {
	ID          string             `json:"id"`
	Lines       string             `json:"lines,omitempty"`
	OssLines    string             `json:"oss_lines,omitempty"`
	Matched     string             `json:"matched,omitempty"`
	FileHash    string             `json:"file_hash,omitempty"`
	SourceHash  string             `json:"source_hash,omitempty"`
	FileURL     string             `json:"file_url,omitempty"`
	Purl        []string           `json:"purl"`
	Vendor      string             `json:"vendor,omitempty"`
	Component   string             `json:"component,omitempty"`
	Version     string             `json:"version,omitempty"`
	Latest      string             `json:"latest,omitempty"`
	URL         string             `json:"url,omitempty"`
	Status      string             `json:"status,omitempty"`
	ReleaseDate string             `json:"release_date,omitempty"`
	File        string             `json:"file,omitempty"`
	URLHash     string             `json:"url_hash,omitempty"`
	URLStats    struct{}           `json:"url_stats,omitempty"`
	Provenance  string             `json:"provenance,omitempty"`
	Licenses    []ComponentLicense `json:"licenses,omitempty"`
	Server      struct {
		Version   string `json:"version,omitempty"`
		KbVersion struct {
			Monthly string `json:"monthly,omitempty"`
			Daily   string `json:"daily,omitempty"`
		} `json:"kb_version"`
		Hostname string `json:"hostname,omitempty"`
		Flags    string `json:"flags,omitempty"`
		Elapsed  string `json:"elapsed,omitempty"`
	} `json:"server"`
}

type FilterAction string

const (
	Include FilterAction = "include"
	Remove  FilterAction = "remove"
	Replace FilterAction = "replace"
)

type ComponentFilterDTO struct {
	Path        string       `json:"path,omitempty"`
	Purl        string       `json:"purl" validate:"required"`
	Usage       string       `json:"usage,omitempty"`
	Action      FilterAction `json:"action" validate:"required,eq=include|eq=remove|eq=replace"`
	Comment     string       `json:"comment,omitempty"`
	ReplaceWith string       `json:"replace_with,omitempty" validate:"omitempty,valid-purl"`
	License     string       `json:"license,omitempty"`
}

type Component struct {
	ID          string   `json:"id"`
	Lines       string   `json:"lines,omitempty"`
	OssLines    string   `json:"oss_lines,omitempty"`
	Matched     string   `json:"matched,omitempty"`
	FileHash    string   `json:"file_hash,omitempty"`
	SourceHash  string   `json:"source_hash,omitempty"`
	FileURL     string   `json:"file_url,omitempty"`
	Purl        []string `json:"purl,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Component   string   `json:"component,omitempty"`
	Version     string   `json:"version,omitempty"`
	Latest      string   `json:"latest,omitempty"`
	URL         string   `json:"url,omitempty"`
	Status      string   `json:"status,omitempty"`
	ReleaseDate string   `json:"release_date,omitempty"`
	File        string   `json:"file,omitempty"`
	URLHash     string   `json:"url_hash,omitempty"`
	URLStats    struct{} `json:"url_stats,omitempty"`
	Provenance  string   `json:"provenance,omitempty"`
	Licenses    []struct {
		Name             string `json:"name"`
		PatentHints      string `json:"patent_hints"`
		Copyleft         string `json:"copyleft"`
		ChecklistURL     string `json:"checklist_url"`
		OsadlUpdated     string `json:"osadl_updated"`
		Source           string `json:"source"`
		URL              string `json:"url"`
		IncompatibleWith string `json:"incompatible_with,omitempty"`
	} `json:"licenses,omitempty"`
	Health struct {
		CreationDate string `json:"creation_date"`
		LastUpdate   string `json:"last_update"`
		LastPush     string `json:"last_push"`
		Stars        int    `json:"stars"`
		Issues       int    `json:"issues"`
		Forks        int    `json:"forks"`
	} `json:"health"`
	Dependencies []interface{} `json:"dependencies"`
	Copyrights   []struct {
		Name   string `json:"name"`
		Source string `json:"source"`
	} `json:"copyrights"`
	Vulnerabilities []interface{} `json:"vulnerabilities"`
	Server          struct {
		Version   string `json:"version,omitempty"`
		KbVersion struct {
			Monthly string `json:"monthly,omitempty"`
			Daily   string `json:"daily,omitempty"`
		} `json:"kb_version"`
		Hostname string `json:"hostname,omitempty"`
		Flags    string `json:"flags,omitempty"`
		Elapsed  string `json:"elapsed,omitempty"`
	} `json:"server"`
}

type DeclaredComponent struct {
	Name string `json:"name"`
	Purl string `json:"purl"`
}
