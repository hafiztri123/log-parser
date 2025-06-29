package parser

import (
	"regexp"
	"strings"
)

type LogFormat int

const (
	AccessLog LogFormat = iota
	ApplicationLog
	StackTraceContinuation
	Unknown
)

type FormatDetector struct {
	accessLogRegex  *regexp.Regexp
	appLogRegex     *regexp.Regexp
	stackTraceRegex *regexp.Regexp
}

func NewFormatDetector() *FormatDetector {
	return &FormatDetector{
		accessLogRegex:  regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+-\s+(.+?)\s\[(.+?)\]\s+"(.+?)"\s+(\d{3})\s+(\d+)\s+"(.+?)"\s+"(.+?)"$`),
		appLogRegex:     regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\.\d{1,3})\s+-\s+(ERROR|WARN|INFO|DEBUG)\s+-\s+(.*)$`),
		stackTraceRegex: regexp.MustCompile(`^(\s+(.+))`),
	}
}

func (fd *FormatDetector) DetectFormat(line string) LogFormat {
	if strings.TrimSpace(line) == "" {
		return Unknown
	}

	if fd.stackTraceRegex.MatchString(line) {
		return StackTraceContinuation
	}

	if fd.appLogRegex.MatchString(line) {
		return ApplicationLog
	}

	if fd.accessLogRegex.MatchString(line) {
		return AccessLog
	}

	return Unknown
}

func (fd *FormatDetector) ParseAccessLog(line string) map[string]string {
	matches := fd.accessLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	return map[string]string{
		"ip":         matches[1],
		"user":       matches[2],
		"timestamp":  matches[3],
		"request":    matches[4],
		"status":     matches[5],
		"size":       matches[6],
		"referer":    matches[7],
		"user_agent": matches[8],
	}
}

func (fd *FormatDetector) ParseApplicationLog(line string) map[string]string {
	matches := fd.appLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	return map[string]string{
		"timestamp": matches[1],
		"level":     matches[2],
		"message":   matches[3],
	}
}

func (f LogFormat) String() string {
	switch f {
	case AccessLog:
		return "access"
	case ApplicationLog:
		return "application"
	case StackTraceContinuation:
		return "continuation"
	default:
		return "unknown"
	}
}
