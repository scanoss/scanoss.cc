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

type ScanResponse struct {
	Output    string `json:"output,omitempty"`
	ErrOutput string `json:"error_output,omitempty"`
	Error     error  `json:"error,omitempty"`
}

type ScanArgDef struct {
	Name           string
	Shorthand      string
	Default        interface{}
	Usage          string
	Tooltip        string
	Type           string
	IsCore         bool
	IsFileSelector bool
}

var (
	ScanArguments = []ScanArgDef{
		{"wfp", "w", "", "Scan a WFP File instead of a folder (optional)", "", "string", false, true},
		{"dep", "p", "", "Use a dependency file instead of a folder (optional)", "", "string", false, true},
		{"stdin", "s", "", "Scan the file contents supplied via STDIN (optional)", "", "string", false, false},
		{"files", "e", []string{}, "List of files to scan.", "", "stringSlice", false, false},
		{"identify", "i", "", "Scan and identify components in SBOM file", "", "string", false, true},
		{"ignore", "n", "", "Ignore components specified in the SBOM file", "", "string", false, true},
		{"output", "o", "", "Output result file name (optional - default stdout).", "Location where the scan results will be saved", "string", true, true},
		{"format", "f", "plain", "Result output format (optional - default: plain)", "", "string", false, false},
		{"threads", "T", 5, "Number of threads to use while scanning (optional - default 5)", "", "int", false, false},
		{"flags", "F", 0, "Scanning engine flags", "Advanced scanning engine configuration flags", "int", false, false},
		{"post-size", "P", 32, "Number of kilobytes to limit the post to while scanning (optional - default 32)", "Limits the size of each scan request to improve performance and reliability", "int", false, false},
		{"timeout", "M", 180, "Timeout (in seconds) for API communication (optional - default 180)", "", "int", false, false},
		{"retry", "R", 5, "Retry limit for API communication (optional - default 5)", "", "int", false, false},
		{"no-wfp-output", "", false, "Skip WFP file generation", "", "bool", false, false},
		{"dependencies", "D", false, "Add Dependency scanning", "", "bool", false, false},
		{"dependencies-only", "", false, "Run Dependency scanning only", "", "bool", false, false},
		{"sc-command", "", "", "Scancode command and path if required (optional - default scancode).", "", "string", false, false},
		{"sc-timeout", "", 600, "Timeout (in seconds) for scancode to complete (optional - default 600)", "", "int", false, false},
		{"dep-scope", "", "", "Filter dependencies by scope - default all (options: dev/prod)", "", "string", false, false},
		{"dep-scope-inc", "", "", "Include dependencies with declared scopes", "", "string", false, false},
		{"dep-scope-exc", "", "", "Exclude dependencies with declared scopes", "", "string", false, false},
		{"settings", "", "", "Settings file to use for scanning (optional - default scanoss.json)", "Configuration file that defines scanning behavior and rules", "string", true, true},
		{"skip-settings-file", "", false, "Skip default settings file (scanoss.json) if it exists", "", "", false, false},
		{"debug", "d", false, "Enable debug messages", "Show detailed diagnostic information during scanning", "bool", true, false},
		{"trace", "t", false, "Enable trace messages, including API posts", "Display all API communication and detailed execution steps", "bool", true, false},
		{"quiet", "q", true, "Enable quiet mode", "Suppress non-essential output messages", "bool", true, false},
	}
)
