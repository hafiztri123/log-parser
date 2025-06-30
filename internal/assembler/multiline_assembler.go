package assembler

import (
	"hafiztri123/log-pipeline/internal/parser"
	"hafiztri123/log-pipeline/internal/util"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	LogLevel *string
	SourceType string
	RawMessage string
	IPAddress *string
	UserID *string
	HTTPMethod *string
	HTTPPath *string
	HTTPStatus *int
	ServiceName *string
	FileSource string
	ParsedData map[string]interface{}
}

type MultiLineAssembler struct {
	detector *parser.FormatDetector
	pendingEntry *LogEntry
	pendingLines []string
	inMultilineBlock bool
}

func NewMultiLineAssembler() *MultiLineAssembler {
	return &MultiLineAssembler{
		detector: parser.NewFormatDetector(),
		pendingLines: make([]string, 0),
	}
}

func (mla *MultiLineAssembler) ProcessLine(line string, filesource string) *LogEntry {
	format := mla.detector.DetectFormat(line)

	switch format {
	case parser.AccessLog:
		completed := mla.completePendingEntry()

		entry := mla.processAccessLog(line, filesource)
		mla.pendingEntry = entry
		mla.inMultilineBlock = false

		return completed
	case parser.ApplicationLog:
		completed := mla.completePendingEntry()

	}
}

func (mla *MultiLineAssembler) completePendingEntry() *LogEntry {
	if mla.pendingEntry == nil {
		return nil
	}

	completed := mla.pendingEntry
	mla.pendingEntry = nil
	mla.pendingLines = nil
	mla.inMultilineBlock = false

	return completed
}


func (mla *MultiLineAssembler) processAccessLog(line, filesource string) *LogEntry {
	parsed := mla.detector.ParseAccessLog(line)
	if parsed == nil {
		mla.processUnknownLog(line, filesource)
	}

	timestamp := util.ParseLogTimestamp(parsed["timestamp"])
	status :=  util.ParseInt(parsed["status"])
	responseSize := util.ParseInt(parsed["size"])
	request := parsed["request"]
	ipValue := parsed["ip"]
	method, path := util.ParseHTTPRequest(request)

	entry := &LogEntry{
		Timestamp: timestamp,
		LogLevel: nil,
		SourceType: "access",
		RawMessage: line,
		IPAddress: &ipValue,
		HTTPMethod: method,
		HTTPPath: path,
		HTTPStatus: &status,
		FileSource: filesource,
		ParsedData: map[string]interface{}{
			"user_agent": parsed["user_agent"],
			"referer": parsed["referer"],
			"response_size": responseSize,
			"http_version" : util.ExtractHTTPVersion(request),
		},
	}

	if parsed["user"] != "-" {
		userID := parsed["user"]
		entry.UserID = &userID
	}

	return entry
}

func (mla *MultiLineAssembler) processUnknownLog(line,  fileSource string) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now(),
		SourceType: "unknown",
		RawMessage: line,
		FileSource: fileSource,
		ParsedData: make(map[string]interface{}),
	}
}

func (mla *MultiLineAssembler) processApplicationLog(line, filesource string) *LogEntry {
	parsed := mla.detector.ParseApplicationLog(line)
	if parsed == nil {
		return mla.processUnknownLog(line, filesource)
	}

	timestamp := util.ParseLogTimestamp(parsed["timestamp"])
	logLevel := parsed["level"]

	entry := &LogEntry{
		Timestamp: timestamp,
		LogLevel: &logLevel,
		SourceType: "access",
		RawMessage: parsed["message"],
		IPAddress : nil,
		UserID : nil,
		HTTPMethod : nil,
		HTTPPath : nil,
		HTTPStatus : nil,
		ServiceName : nil,
		FileSource : "",
		ParsedData : map[string]interface{}{
			"user_agent": "test",
		},
	}
}