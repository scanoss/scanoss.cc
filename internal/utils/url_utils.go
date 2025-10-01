package utils

import "strings"

// RemoveProtocolFromURL removes the HTTP or HTTPS protocol prefix from a URL string.
// If no protocol is found, returns the original string unchanged.
func RemoveProtocolFromURL(url string) string {
	if url == "" {
		return url
	}

	if after, ok := strings.CutPrefix(url, "http://"); ok {
		return after
	}
	if after, ok := strings.CutPrefix(url, "https://"); ok {
		return after
	}
	return url
}
