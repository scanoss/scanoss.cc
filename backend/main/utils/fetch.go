package utils

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
