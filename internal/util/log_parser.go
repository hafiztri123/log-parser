package util

import (
	"strconv"
	"strings"
	"time"
)

func ParseAccessLogTimestamp(timestampStr string) time.Time {
	//magic number
	layout := "22/Jun/2025:12:32:07 +0700"

	if parsed, err := time.Parse(layout, timestampStr); err == nil {
		return parsed
	}

	return time.Now()
}

func ParseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0
	}

	if result, err := strconv.Atoi(s); err == nil {
		return result
	}

	return 0
}

func ParseHTTPRequest(request string) (*string, *string) {
	parts := strings.Fields(request)
	if len(parts) >= 2 {
		return &parts[0], &parts[1]
	}
	return nil, nil
}

func ExtractHTTPVersion(request string) string {
	if idx := strings.LastIndex(request, "HTTP/"); idx != -1 {
		return strings.Fields(request[idx:])[0]
	}
	return "HTTP/1.0"
}