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
