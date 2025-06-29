//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// --- Configuration ---
const (
	logFileName       = "generated_logs.log"
	totalLogEntries   = 50000
	timeframeSeconds  = 3600 * 24 * 7 // Simulate logs over the last 7 days
)

// --- Log Format Definitions ---

// ApacheLogFormat defines the structure for an Apache combined log entry.
// Example: 127.0.0.1 - - [28/Jun/2025:14:22:36 +0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I)"
const apacheLogFormat = `%s - %s [%s] "%s %s %s" %d %d "%s" "%s"`

// AppLogFormat defines the structure for a standard application log.
// Example: 2025-06-28 14:22:36.123 - INFO - User 'admin' logged in successfully.
const appLogFormat = `%s - %s - %s`

// ErrorLogFormat defines the structure for an application error log.
const errorLogFormat = `%s - ERROR - %s`
const tracebackTemplate = `Traceback (most recent call last):
  File "%s", line %d, in %s
    %s()
  File "%s", line %d, in %s
    raise %s("%s")
%s: %s`

// --- Data for Randomization ---

var (
	ipAddresses = []string{
		"192.168.1.101", "10.0.0.5", "172.16.31.45", "203.0.113.19",
		"198.51.100.87", "127.0.0.1", "209.123.12.34", "64.233.160.1",
	}
	remoteUsers = []string{"-", "john", "jane", "admin_user", "api_key_user"}
	httpMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	requestPaths = []string{
		"/index.html", "/api/v1/users", "/api/v1/products", "/images/logo.png",
		"/about.html", "/contact", "/admin/dashboard", "/login", "/logout",
		"/assets/style.css", "/assets/main.js", "/data/../boot.ini",
	}
	httpVersions = []string{"HTTP/1.1", "HTTP/2.0"}
	statusCodes  = []int{200, 201, 204, 400, 401, 403, 404, 500, 503}
	referers     = []string{
		"http://www.google.com", "http://www.bing.com", "http://www.yoursite.com",
		"http://localhost:3000/dashboard", "-", "https://t.co/xyz",
	}
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"python-requests/2.25.1", "curl/7.68.0",
	}
	appMessages = map[string][]string{
		"INFO": {
			"User '%s' logged in successfully.",
			"Processing payment for order #%d.",
			"Data export job '%s' started.",
			"Cache cleared for key '%s'.",
		},
		"WARN": {
			"API rate limit exceeded for IP %s.",
			"Disk space is running low (%d%% remaining).",
			"Deprecated function 'old_function' was called.",
			"Failed login attempt for user '%s'.",
		},
	}
	errorMessages = []string{
		"Database connection failed: Timeout expired",
		"NullPointerException: Attempt to invoke method 'toString' on a null object reference",
		"Failed to write to file: /var/log/app.log",
		"Uncaught TypeError: Cannot read properties of undefined (reading 'id')",
	}
	fileNames = []string{"app.py", "database.go", "main.go", "routes.js"}
	funcNames = []string{"connect_db", "process_request", "handle_payment", "render_template"}
	errorTypes = []string{"ValueError", "ConnectionError", "FileNotFoundError", "TypeError"}
)


// --- Generator Functions ---

// randomChoice selects a random element from a string slice.
func randomChoice(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

// generateTimestamp creates a random timestamp within a given timeframe.
func generateTimestamp(start time.Time, secondsAgo int) time.Time {
	randomSeconds := rand.Intn(secondsAgo)
	return start.Add(-time.Duration(randomSeconds) * time.Second)
}

// generateApacheLog creates a single Apache combined log entry.
func generateApacheLog(ts time.Time) string {
	ip := randomChoice(ipAddresses)
	user := randomChoice(remoteUsers)
	// Apache format: 02/Jan/2006:15:04:05 -0700
	timestampStr := ts.Format("02/Jan/2006:15:04:05 -0700")
	method := randomChoice(httpMethods)
	path := randomChoice(requestPaths)
	version := randomChoice(httpVersions)
	status := statusCodes[rand.Intn(len(statusCodes))]

	// Make status code more realistic based on method and path
	if method == "POST" {
		status = []int{200, 201, 400}[rand.Intn(3)]
	} else if strings.HasPrefix(path, "/admin") {
		status = []int{200, 401, 403}[rand.Intn(3)]
	} else if strings.Contains(path, "..") { // Simple check for path traversal
		status = 400
	}

	size := rand.Intn(29901) + 100 // 100 to 30000
	referer := randomChoice(referers)
	userAgent := randomChoice(userAgents)

	return fmt.Sprintf(apacheLogFormat, ip, user, timestampStr, method, path, version, status, size, referer, userAgent)
}

// generateAppLog creates a single application log entry (INFO or WARN).
func generateAppLog(ts time.Time) string {
	// App format: 2006-01-02 15:04:05.000
	timestampStr := ts.Format("2006-01-02 15:04:05.000")
	level := "INFO"
	if rand.Intn(4) == 0 { // 25% chance of a WARN message
		level = "WARN"
	}
	
	msgTemplate := randomChoice(appMessages[level])
	var message string
	// Fill in templates with random data
	switch {
	case strings.Contains(msgTemplate, "'%s'"):
		message = fmt.Sprintf(msgTemplate, randomChoice(remoteUsers))
	case strings.Contains(msgTemplate, "#%d"):
		message = fmt.Sprintf(msgTemplate, rand.Intn(99999)+1000)
	case strings.Contains(msgTemplate, "%s"):
		message = fmt.Sprintf(msgTemplate, randomChoice(ipAddresses))
	case strings.Contains(msgTemplate, "%%"):
		message = fmt.Sprintf(msgTemplate, rand.Intn(20)) // 0-19% disk space
	default:
		message = msgTemplate
	}

	return fmt.Sprintf(appLogFormat, timestampStr, level, message)
}

// generateErrorLog creates a multi-line error log with a traceback.
func generateErrorLog(ts time.Time) string {
	timestampStr := ts.Format("2006-01-02 15:04:05.000")
	baseError := randomChoice(errorMessages)

	// Create a fake traceback
	fileName1 := randomChoice(fileNames)
	line1 := rand.Intn(200) + 50
	funcName1 := randomChoice(funcNames)
	fileName2 := randomChoice(fileNames)
	line2 := rand.Intn(100) + 1
	funcName2 := "main"
	errorType := randomChoice(errorTypes)

	traceback := fmt.Sprintf(tracebackTemplate,
		fileName2, line2, funcName2,
		funcName1,
		fileName1, line1, funcName1,
		errorType, baseError,
		errorType, baseError,
	)

	// The first line is the standard error message
	firstLine := fmt.Sprintf(errorLogFormat, timestampStr, "")
	return firstLine + "\n" + traceback
}

func main() {
	// Initialize random seed
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create and open the log file
	file, err := os.Create(logFileName)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	startTime := time.Now()
	
	fmt.Printf("Generating %d log entries into %s...\n", totalLogEntries, logFileName)

	for i := 0; i < totalLogEntries; i++ {
		// Generate a timestamp for the log entry
		ts := generateTimestamp(startTime, timeframeSeconds)

		// Choose a random log type to generate
		logType := rand.Intn(100)
		var logEntry string

		switch {
		case logType < 70: // 70% chance for Apache logs
			logEntry = generateApacheLog(ts)
		case logType < 95: // 25% chance for App logs (INFO/WARN)
			logEntry = generateAppLog(ts)
		default: // 5% chance for Error logs
			logEntry = generateErrorLog(ts)
		}
		
		// Write the entry to the file buffer
		_, err := writer.WriteString(logEntry + "\n")
		if err != nil {
			log.Fatalf("Failed to write to buffer: %v", err)
		}

		// Periodically flush the buffer to disk
		if i%1000 == 0 {
			writer.Flush()
		}
	}

	// Final flush to ensure everything is written
	writer.Flush()
	
	fmt.Println("Log generation complete.")
}
