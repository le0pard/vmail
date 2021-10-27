package main

import (
	"testing"
)

// ReportFromHTMLSimple calls ReportFromHTML with a simple html
func ReportFromHTMLSimple(t *testing.T) {
	html := "<html><body></body></html>"
	_, err := ReportFromHTML([]byte(html))
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}
}
