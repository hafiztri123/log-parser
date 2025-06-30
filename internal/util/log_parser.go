package util

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseLogTimestamp(timestampStr string) time.Time {
	accessLogTimestampRegex := regexp.MustCompile(`^\d{1,2}/\[a-zA-Z]+/\d{4}:\d{1,2}:\d{1,2}:\d{1,2}\s+\+\d+`)
	applicationLogTimestampRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s+\d{1,2}:\d{1,2}:\d{1,2}\.\d{1,3}`)

	accessLoglayout := "22/Jun/2025:12:32:07 -0700"
	applicationLogLayout := "2006-01-02 15:04:05.000"

	if accessLogTimestampRegex.MatchString(timestampStr) {
		parsed, err := time.Parse(accessLoglayout, timestampStr)
		if err == nil {
			return parsed
		}

	} else if applicationLogTimestampRegex.MatchString(timestampStr) {
		parsed, err := time.Parse(applicationLogLayout, timestampStr) 
		if err == nil {
			return parsed
		}
	}

	//Fallback
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
