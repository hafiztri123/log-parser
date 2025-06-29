package main

import (
	"fmt"
	"hafiztri123/log-pipeline/internal/parser"
)

func main() {
	detector := parser.NewFormatDetector()

	lines := []string{
		`172.16.31.45 - jane [27/Jun/2025:20:00:02 +0700] "POST /contact HTTP/1.1" 400 6795 "-" "python-requests/2.25.1"`,
		`2025-06-28 06:12:45.107 - ERROR - `,
		`  File "database.go", line 74, in main`,
		`    render_template()`,
	}

	for _, line := range lines {
		format := detector.DetectFormat(line)
		fmt.Printf("Format: %s | Line: %s\n", format, line)
	}
}
