// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2025 SCANOSS.COM
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

// ComponentSearchRequest represents the request for searching components online
type ComponentSearchRequest struct {
	Search    string `json:"search" validate:"required,min=1"`
	Vendor    string `json:"vendor,omitempty"`
	Component string `json:"component,omitempty"`
	Package   string `json:"package,omitempty"`
	Limit     int32  `json:"limit,omitempty"`
	Offset    int32  `json:"offset,omitempty"`
}

// ComponentSearchResponse represents the response from component search
type ComponentSearchResponse struct {
	Components []SearchedComponent `json:"components"`
	Status     StatusResponse      `json:"status"`
}

// SearchedComponent represents a component found in the search results
type SearchedComponent struct {
	Component string `json:"component"`
	Purl      string `json:"purl"`
	URL       string `json:"url"`
}

// StatusResponse represents the status information from the API
type StatusResponse struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
