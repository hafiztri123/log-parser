package parser

import (
	"testing"
)

func TestFormatDetector(t *testing.T) {
	detector := NewFormatDetector()

	testCases := []struct {
		line     string
		expected LogFormat
	}{
		// Access log
		{
			`172.16.31.45 - jane [27/Jun/2025:20:00:02 +0700] "POST /contact HTTP/1.1" 400 6795 "-" "python-requests/2.25.1"`,
			AccessLog,
		},
		// Application log
		{
			`2025-06-28 06:12:45.107 - ERROR - `,
			ApplicationLog,
		},
		{
			`2025-06-23 01:29:17.107 - INFO - Cache cleared for key 'jane'.`,
			ApplicationLog,
		},
		// Stack trace continuation
		{
			`  File "database.go", line 74, in main`,
			StackTraceContinuation,
		},
		{
			`    render_template()`,
			StackTraceContinuation,
		},
		// Unknown
		{
			`some random log format`,
			Unknown,
		},
	}

	for i, tc := range testCases {
		result := detector.DetectFormat(tc.line)
		if result != tc.expected {
			t.Errorf("Test %d failed: expected %s, got %s for line: %s",
				i, tc.expected, result, tc.line,
			)
		}
	}
}

func TestAccessLogParsing(t *testing.T) {
	detector := NewFormatDetector()
	line := `172.16.31.45 - jane [27/Jun/2025:20:00:02 +0700] "POST /contact HTTP/1.1" 400 6795 "-" "python-requests/2.25.1"`

	parsed := detector.ParseAccessLog(line)

	expected := map[string]string{
		"ip":         "172.16.31.45",
		"user":       "jane",
		"timestamp":  "27/Jun/2025:20:00:02 +0700",
		"request":    "POST /contact HTTP/1.1",
		"status":     "400",
		"size":       "6795",
		"referer":    "-",
		"user_agent": "python-requests/2.25.1",
	}

	for key, expectedValue := range expected {
		if parsed[key] != expectedValue {
			t.Errorf("Expected %s=%s got %s", key, expectedValue, parsed[key])
		}
	}
}

func TestApplicationLogParsing(t *testing.T) {
	detector := NewFormatDetector()
	line := "2025-06-27 05:46:31.182 - INFO - User 'admin_user' logged in successfully."

	parsed := detector.ParseApplicationLog(line)

	expected := map[string]string{
		"timestamp": "2025-06-27 05:46:31.182",
		"level":     "INFO",
		"message":   "User 'admin_user' logged in successfully.",
	}

	for key, expectedValue := range expected {
		if parsed[key] != expectedValue {
			t.Errorf("Expected %s=%s got %s", key, expectedValue, parsed[key])
		}
	}
}

func TestLogTypeDetector(t *testing.T) {

	formatDetector := NewFormatDetector()

	testcases := []struct {
		line         string
		expectedType string
	}{
		{
			`172.16.31.45 - jane [27/Jun/2025:20:00:02 +0700] "POST /contact HTTP/1.1" 400 6795 "-" "python-requests/2.25.1"`,
			"access",
		},
		// Application log
		{
			`2025-06-28 06:12:45.107 - ERROR - `,
			"application",
		},
		{
			`  File "database.go", line 74, in main`,
			"continuation",
		},
		{
			`some random log format`,
			"unknown",
		},
		{
			"",
			"unknown",
		},
	}

	for _, tc := range testcases {
		logType := formatDetector.DetectFormat(tc.line)
		if tc.expectedType != logType.String() {
			t.Errorf("[ERROR] logtype: %s is not the same as expected type %s", logType.String(), tc.expectedType)
		}
	}

}
