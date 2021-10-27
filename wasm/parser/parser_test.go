package parser

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// ReportFromHTMLSimple calls ReportFromHTML with a simple html
func TestReportFromHTMLSimple(t *testing.T) {
	html := "<html><body></body></html>"
	_, err := ReportFromHTML([]byte(html))
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}
}

type htmlFixture struct {
	File string
}

var htmlFixtures = []htmlFixture{
	{"1.html"},
	{"2.html"},
	{"3.html"},
}

func TestReportFromHTMLFixtures(t *testing.T) {
	for _, fixture := range htmlFixtures {
		pathDir, _ := os.Getwd()
		html, err := ioutil.ReadFile(path.Join(pathDir, "fixtures", fixture.File))
		if err != nil {
			t.Fatalf(`Not found fixture %s, %v`, fixture.File, err)
			continue
		}

		_, err = ReportFromHTML(html)
		if err != nil {
			t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
		}
	}

}
