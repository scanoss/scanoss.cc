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

package fetch

import (
	"bytes"
	"io"
	"net/http"
)

type Options struct {
	Method  string
	Headers map[string]string
	Body    []byte
}

// Fetch is a helper function that mimics the fetch API in JavaScript.
func Fetch(url string, options Options) (*http.Response, error) {
	// Set default method to GET if not provided
	if options.Method == "" {
		options.Method = http.MethodGet
	}

	// Create a new HTTP request
	req, err := http.NewRequest(options.Method, url, bytes.NewBuffer(options.Body))
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// Execute the HTTP request
	client := &http.Client{}
	return client.Do(req)
}

// FetchText is a convenience function to get the response body as a string.
func Text(url string, options Options) (string, error) {
	resp, err := Fetch(url, options)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
