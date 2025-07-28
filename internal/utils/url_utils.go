package utils

import "strings"

func RemoveProtocolFromUrl(url string) string {
	if after, ok := strings.CutPrefix(url, "http://"); ok {
		return after
	}
	if after, ok := strings.CutPrefix(url, "https://"); ok {
		return after
	}
	return url
}
