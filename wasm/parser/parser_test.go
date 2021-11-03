package parser

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

// ReportFromHTMLSimple calls ReportFromHTML with a simple html
func TestReportFromHTMLSimple(t *testing.T) {
	html := `<html><body>
	<div style="background: red"></div>
</body></html>`
	report, err := ReportFromHTML([]byte(html))
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}
	// log.Printf("report: %v\n", report)
	if !reflect.DeepEqual(report.HtmlTags["div"][""].Lines, map[int]bool{2: true}) {
		t.Errorf("ReportFromHTML HtmlTags div got %v, want %v", report.HtmlTags["div"][""].Lines, map[int]bool{2: true})
	}
	if !reflect.DeepEqual(report.CssProperties["background"][""].Lines, map[int]bool{2: true}) {
		t.Errorf("ReportFromHTML CssProperties background got %v, want %v", report.HtmlTags["background"][""].Lines, map[int]bool{2: true})
	}
}

func TestReportFromHTMLWithStyleTag(t *testing.T) {
	html := `<html><body>
	<style>
	  a {
			color: red;
		}

		a:hover {
			color: yellow;
		}
	</style>
	<a href="https://www.rwpod.com">Link to website</a>
</body></html>`
	report, err := ReportFromHTML([]byte(html))
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}
	var tests = []struct {
		checkType string
		got       map[int]bool
		want      map[int]bool
	}{
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{2: true}},
		{"CssSelectorTypes TYPE_SELECTOR_TYPE", report.CssSelectorTypes["9"].Lines, map[int]bool{3: true, 7: true}},
		{"CssPseudoSelectors hover", report.CssPseudoSelectors["hover"].Lines, map[int]bool{7: true}},
	}

	for _, tt := range tests {
		testname := tt.checkType
		t.Run(testname, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("%s: got %v, want %v", tt.checkType, tt.got, tt.want)
			}
		})
	}
}

func TestReportFromHTMLWithInlineStylesAndStyleTag(t *testing.T) {
	html := `<html><body>
	<style>
	  .button {
			color: red;
			width: 300px;
		}

		@media (max-width: 700px) {
			.button {
				width: 100vh;
			}
		}

		.cover img {
			aspect-ratio: 1/2;
		}
	</style>
	<button class="button">Some button</button>
	<button class="button" style="background: url('img.webp') no-repeat; color: black">Another</button>
	<div class="cover">
		<img src="/some/img.avif" alt="cover" />
	</div>
</body></html>`
	report, err := ReportFromHTML([]byte(html))
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}

	// log.Printf("report: %v\n", report)

	var tests = []struct {
		checkType string
		got       map[int]bool
		want      map[int]bool
	}{
		{"HtmlTags div", report.HtmlTags["div"][""].Lines, map[int]bool{20: true}},
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{2: true}},
		{"CssProperties background", report.CssProperties["background"][""].Lines, map[int]bool{19: true}},
		{"CssProperties width", report.CssProperties["width"][""].Lines, map[int]bool{5: true, 10: true}},
		{"CssProperties aspect-ratio", report.CssProperties["aspect-ratio"][""].Lines, map[int]bool{15: true}},
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{8: true}},
		{"CssSelectorTypes DESCENDANT_COMBINATOR_TYPE", report.CssSelectorTypes["5"].Lines, map[int]bool{14: true}},
		{"CssSelectorTypes TYPE_SELECTOR_TYPE", report.CssSelectorTypes["9"].Lines, map[int]bool{14: true}},
		{"CssDimentions px", report.CssDimentions["px"].Lines, map[int]bool{5: true, 8: true}},
		{"CssDimentions vh", report.CssDimentions["vh"].Lines, map[int]bool{10: true}},
		{"ImgFormats avif", report.ImgFormats["avif"].Lines, map[int]bool{21: true}},
		{"ImgFormats webp", report.ImgFormats["webp"].Lines, map[int]bool{19: true}},
	}

	for _, tt := range tests {
		testname := tt.checkType
		t.Run(testname, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("%s: got %v, want %v", tt.checkType, tt.got, tt.want)
			}
		})
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
