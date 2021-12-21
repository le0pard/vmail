package parser

import (
	"os"
	"reflect"
	"testing"
)

func TestReportFromHTMLSimple(t *testing.T) {
	html := `<html><body>
	<div style="background: red"></div>
</body></html>`
	report, err := ReportFromHTML([]byte(html))
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}
	// log.Printf("report: %v\n", report)
	if !reflect.DeepEqual(report.CssProperties["background"][""].Lines, map[int]bool{2: true}) {
		t.Errorf("ReportFromHTML CssProperties background got %v, want %v", report.HtmlTags["background"][""].Lines, map[int]bool{2: true})
	}
}

func TestReportFromHTMLPropsNormalization(t *testing.T) {
	html := `<html><body>
	<style>
	  .button {
			margin-top: 10px;
			margin-bottom: 15px!imporTant;
			padding: 10px;
		}

		.not-button {
			margin-left: 10px;
			margin-right: 15px;
			padding-top: 30px !important;
		}
	</style>
	<button class="button not-button">To be or not to be</button>
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
		{"CssProperties margin", report.CssProperties["margin"][""].Lines, map[int]bool{4: true, 5: true, 10: true, 11: true}},
		{"CssProperties padding", report.CssProperties["padding"][""].Lines, map[int]bool{6: true, 12: true}},
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

func TestReportFromHTMLImages(t *testing.T) {
	html := `<html><body>
	<style>
	  .button {
			background: url(img.svg) no-repeat;
		}
	</style>
	<img src="elva-fairy-800w.svg" alt="Elva dressed as a fairy" />
	<img srcset="elva-fairy-480w.svg 480w,
             elva-fairy-800w.svg 800w"
     sizes="(max-width: 600px) 480px,
            800px"
     src="elva-fairy-800w.webp"
     alt="Elva dressed as a fairy" />
	<img srcset="elva-fairy-320w.svg,
             elva-fairy-480w.svg 1.5x,
             elva-fairy-640w.svg 2x"
     src="elva-fairy-640w.avif"
     alt="Elva dressed as a fairy" />

	<picture>
		<source media="(max-width: 799px)" srcset="elva-480w-close-portrait.svg">
		<source media="(min-width: 800px)" srcset="elva-800w.svg">
		<img src="elva-800w.svg" alt="Chris standing up holding his daughter Elva">
	</picture>
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
		{"HtmlTags picture", report.HtmlTags["picture"][""].Lines, map[int]bool{20: true}},
		{"ImgFormats svg", report.ImgFormats["svg"].Lines, map[int]bool{4: true, 7: true, 8: true, 14: true, 21: true, 22: true, 23: true}},
		{"ImgFormats webp", report.ImgFormats["webp"].Lines, map[int]bool{8: true}},
		{"ImgFormats avif", report.ImgFormats["avif"].Lines, map[int]bool{14: true}},
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

func TestReportFromHTMLDifferentSelectors(t *testing.T) {
	html := `<html><body>
	<style>
		* {
			padding: 0;
			margin: 0;
		}

	  button {
			color: red;
		}

		.button {
			color: red;
		}

		.button[disabled="disabled"] {
			color: yellow;
		}

		.cover .img {
			color: red;
		}

		.cover > .img {
			color: red;
		}

		.cover ~ .img {
			color: red;
		}

		h1 + p {
			color: red;
		}

		.grid, .notgrid {
			color: red;
		}

		.grid.flatten {
			color: red;
		}

		#someID {
			color: red;
		}
	</style>
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
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{2: true}},
		{"CssProperties margin", report.CssProperties["margin"][""].Lines, map[int]bool{5: true}},
		{"CssProperties padding", report.CssProperties["padding"][""].Lines, map[int]bool{4: true}},
		{"CssSelectorTypes ADJACENT_SIBLING_COMBINATOR_TYPE", report.CssSelectorTypes["0"].Lines, map[int]bool{32: true}},
		{"CssSelectorTypes ATTRIBUTE_SELECTOR_TYPE", report.CssSelectorTypes["1"].Lines, map[int]bool{16: true}},
		{"CssSelectorTypes CHAINING_SELECTORS_TYPE", report.CssSelectorTypes["2"].Lines, map[int]bool{40: true}},
		{"CssSelectorTypes CHILD_COMBINATOR_TYPE", report.CssSelectorTypes["3"].Lines, map[int]bool{24: true}},
		{"CssSelectorTypes CLASS_SELECTOR_TYPE", report.CssSelectorTypes["4"].Lines, map[int]bool{12: true, 16: true, 20: true, 24: true, 28: true, 36: true, 40: true}},
		{"CssSelectorTypes DESCENDANT_COMBINATOR_TYPE", report.CssSelectorTypes["5"].Lines, map[int]bool{20: true}},
		{"CssSelectorTypes GENERAL_SIBLING_COMBINATOR_TYPE", report.CssSelectorTypes["6"].Lines, map[int]bool{28: true}},
		{"CssSelectorTypes GROUPING_SELECTORS_TYPE", report.CssSelectorTypes["7"].Lines, map[int]bool{36: true}},
		{"CssSelectorTypes ID_SELECTOR_TYPE", report.CssSelectorTypes["8"].Lines, map[int]bool{44: true}},
		{"CssSelectorTypes TYPE_SELECTOR_TYPE", report.CssSelectorTypes["9"].Lines, map[int]bool{8: true}},
		{"CssSelectorTypes UNIVERSAL_SELECTOR_STAR_TYPE", report.CssSelectorTypes["10"].Lines, map[int]bool{3: true}},
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
		<img src="data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUA
    AAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO
        9TXL0Y4OHwAAAABJRU5ErkJggg==" alt="cover" />
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
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{2: true}},
		{"CssProperties background", report.CssProperties["background"][""].Lines, map[int]bool{19: true}},
		{"CssProperties width", report.CssProperties["width"][""].Lines, map[int]bool{5: true, 10: true}},
		{"CssProperties aspect-ratio", report.CssProperties["aspect-ratio"][""].Lines, map[int]bool{15: true}},
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{8: true}},
		{"CssSelectorTypes DESCENDANT_COMBINATOR_TYPE", report.CssSelectorTypes["5"].Lines, map[int]bool{14: true}},
		{"CssSelectorTypes TYPE_SELECTOR_TYPE", report.CssSelectorTypes["9"].Lines, map[int]bool{14: true}},
		{"CssDimentions vh", report.CssDimentions["vh"].Lines, map[int]bool{10: true}},
		{"ImgFormats avif", report.ImgFormats["avif"].Lines, map[int]bool{21: true}},
		{"ImgFormats webp", report.ImgFormats["webp"].Lines, map[int]bool{19: true}},
		{"ImgFormats base64", report.ImgFormats["base64"].Lines, map[int]bool{22: true}},
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

func TestReportFromHTMLWithComplexStuff(t *testing.T) {
	html := `<html><body>
	<style>
	  .flex {
			display: flex;
		}

		.grid {
			display: grid;
		}

		.hidden {
			display: none;
		}

		.day   { background: #eee; color: black; }
		.night { background: #333; color: white; }

		@media (prefers-color-scheme: dark) {
			.day.dark-scheme   { background:  #333; color: white; }
			.night.dark-scheme { background: black; color:  #ddd; }
		}

		@media (prefers-color-scheme: light) {
			.day.light-scheme   { background: white; color:  #555; }
			.night.light-scheme { background:  #eee; color: black; }
		}
	</style>
	<button type="submit">Submit</button>
	<button type="reset">Reset</button>
	<input type="hidden" />
	<input type="checkbox" />
	<input type="radio" />
	<input type="submit" />
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
		{"HtmlTags button type=submit", report.HtmlTags["button"]["type||submit"].Lines, map[int]bool{28: true}},
		{"HtmlTags button type=reset", report.HtmlTags["button"]["type||reset"].Lines, map[int]bool{29: true}},
		{"HtmlTags input type=hidden", report.HtmlTags["input"]["type||hidden"].Lines, map[int]bool{30: true}},
		{"HtmlTags input type=checkbox", report.HtmlTags["input"]["type||checkbox"].Lines, map[int]bool{31: true}},
		{"HtmlTags input type=radio", report.HtmlTags["input"]["type||radio"].Lines, map[int]bool{32: true}},
		{"HtmlTags input type=submit", report.HtmlTags["input"]["type||submit"].Lines, map[int]bool{33: true}},
		{"CssProperties display flex", report.CssProperties["display"]["flex"].Lines, map[int]bool{4: true}},
		{"CssProperties display grid", report.CssProperties["display"]["grid"].Lines, map[int]bool{8: true}},
		{"CssProperties display none", report.CssProperties["display"]["none"].Lines, map[int]bool{12: true}},
		{"CssProperties background", report.CssProperties["background"][""].Lines, map[int]bool{15: true, 16: true, 19: true, 20: true, 24: true, 25: true}},
		{"AtRuleCssStatements @media prefers-color-scheme", report.AtRuleCssStatements["@media"]["prefers-color-scheme"].Lines, map[int]bool{18: true, 23: true}},
		{"CssSelectorTypes CLASS_SELECTOR_TYPE", report.CssSelectorTypes["4"].Lines, map[int]bool{3: true, 7: true, 11: true, 15: true, 16: true, 19: true, 20: true, 24: true, 25: true}},
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

func TestReportFromHTMLNestedMedia(t *testing.T) {
	html := `<html><body>
	<style>
	  @media screen {
			@media (min-width: 1px) {
				@media (min-height: 1px) {
					@media (max-width: 9999px) {
						@media (max-height: 9999px) {
							@media (prefers-reduced-motion: reduce) {
								body {
									background: red;
								}
							}
						}
					}
				}
			}
		}
	</style>
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
		{"CssProperties background", report.CssProperties["background"][""].Lines, map[int]bool{10: true}},
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{3: true, 4: true, 5: true, 6: true, 7: true, 8: true}},
		{"AtRuleCssStatements @media prefers-reduced-motion", report.AtRuleCssStatements["@media"]["prefers-reduced-motion"].Lines, map[int]bool{8: true}},
		{"CssSelectorTypes TYPE_SELECTOR_TYPE", report.CssSelectorTypes["9"].Lines, map[int]bool{9: true}},
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

func TestReportFromHTMLDifferentDimentionsFormat(t *testing.T) {
	html := `<html><body>
	<style>
	  .one-class {
			padding: 1rem;
		}
		.two-class {
			padding: 1.5rem;
		}
		.another-class {
			padding: .56rem;
		}
		.another-class {
			padding: 100.5343546456rem;
		}
		.another2-class {
			padding: -100rem;
		}
		.another3-class {
			padding: +123.123rem;
		}
	</style>
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
		{"CssDimentions rem", report.CssDimentions["rem"].Lines, map[int]bool{4: true, 7: true, 10: true, 13: true, 16: true, 19: true}},
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

func TestReportFromHTML5Doctype(t *testing.T) {
	html := `<!DOCTYPE html>
<html><body>
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
		{"Html5Doctype", report.Html5Doctype.Lines, map[int]bool{1: true}},
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

func TestReportFromHTMLCssImportant(t *testing.T) {
	html := `<html><body>
	<style>
		.button {
			font-size: 14px !important;
			margin: 10px;
		}
	</style>
	<button class="button" style="padding:10px!important">Test</button>
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
		{"CssImportant", report.CssImportant.Lines, map[int]bool{4: true, 8: true}},
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

func TestReportFromHTMLLinkTypes(t *testing.T) {
	html := `<html><body>
	<ul>
		<li><a href="#test1">Anchor to #test1</a></li>
		<li><a href="#test2">Anchor to #test2</a></li>
		<li><a href="#test3">Anchor to #test3</a></li>
		<li><a href="#test4">Anchor to #test4</a></li>
		<li><a href="#">Not Anchor</a></li>
	</ul>
	<a href="mailto:to@example.com">Test 1</a>
	<a href="mailto:to@example.com?bcc=bcc@example.com">Test 3</a>
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
		{"LinkTypes anchor", report.LinkTypes["anchor"].Lines, map[int]bool{3: true, 4: true, 5: true, 6: true}},
		{"LinkTypes mailto", report.LinkTypes["mailto"].Lines, map[int]bool{9: true, 10: true}},
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

func TestReportFromHTMLCSSWithNoSemicolumns(t *testing.T) {
	html := `<html><body>
<style>
  @media (min-width: 576px) {
    .container, .container-sm {
        max-width: 540px
        }
    }
  @media (min-width: 100px) {
    .container, .container-sm {





        max-width: 1000px
        }
    }
  .button {max-width: 1000px



	}
	.button2 {
				max-width: 1000px



	}
	.button3 {

		max-width: 3000px

	}

	.button3 {

		max-width: 3000px}
</style>
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
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{3: true, 8: true}},
		{"CssProperties max-width", report.CssProperties["max-width"][""].Lines, map[int]bool{5: true, 15: true, 18: true, 24: true, 31: true, 37: true}},
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

func BenchmarkReportFromHTML(b *testing.B) {
	html, err := os.ReadFile("./bench.html")
	if err != nil {
		b.Fatalf(`Error to read bench.html, %v`, err)
	}

	for i := 0; i < b.N; i++ {
		_, err = ReportFromHTML(html)
		if err != nil {
			b.Fatalf(`ReportFromHTML, %v`, err)
		}
	}

}
