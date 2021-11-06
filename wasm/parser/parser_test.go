package parser

import (
	"io/ioutil"
	"os"
	"path"
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
	if !reflect.DeepEqual(report.HtmlTags["div"][""].Lines, map[int]bool{2: true}) {
		t.Errorf("ReportFromHTML HtmlTags div got %v, want %v", report.HtmlTags["div"][""].Lines, map[int]bool{2: true})
	}
	if !reflect.DeepEqual(report.CssProperties["background"][""].Lines, map[int]bool{2: true}) {
		t.Errorf("ReportFromHTML CssProperties background got %v, want %v", report.HtmlTags["background"][""].Lines, map[int]bool{2: true})
	}
}

func TestReportFromHTMLPropsNormalization(t *testing.T) {
	html := `<html><body>
	<style>
	  .button {
			margin-top: 10px;
			margin-bottom: 15px!important;
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

func TestReportFromHTMLFirstFixture(t *testing.T) {
	pathDir, _ := os.Getwd()
	html, err := ioutil.ReadFile(path.Join(pathDir, "fixtures", "1.html"))
	if err != nil {
		t.Fatalf(`Not found fixture %s, %v`, "1.html", err)
	}

	report, err := ReportFromHTML(html)
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}

	var tests = []struct {
		checkType string
		got       map[int]bool
		want      map[int]bool
	}{
		{"HtmlTags div", report.HtmlTags["div"][""].Lines, map[int]bool{99: true, 101: true, 109: true, 133: true, 151: true, 158: true, 190: true, 198: true, 206: true, 224: true, 231: true, 257: true, 265: true, 291: true, 299: true, 325: true, 333: true, 359: true, 367: true, 393: true, 401: true, 427: true, 434: true, 466: true, 474: true, 506: true, 513: true}},
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{13: true, 61: true, 74: true, 85: true}},
		{"CssProperties background", report.CssProperties["background"][""].Lines, map[int]bool{101: true, 102: true, 190: true, 191: true, 224: true, 225: true, 291: true, 292: true, 359: true, 360: true, 427: true, 428: true, 445: true, 448: true, 466: true, 467: true, 506: true, 507: true, 541: true, 559: true, 577: true}},
		{"CssProperties width", report.CssProperties["width"][""].Lines, map[int]bool{64: true, 69: true, 76: true, 81: true, 88: true, 92: true, 102: true, 109: true, 120: true, 121: true, 152: true, 158: true, 170: true, 171: true, 191: true, 198: true, 225: true, 231: true, 243: true, 244: true, 257: true, 292: true, 299: true, 311: true, 312: true, 325: true, 360: true, 367: true, 379: true, 380: true}},
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{62: true, 86: true}},
		{"CssSelectorTypes DESCENDANT_COMBINATOR_TYPE", report.CssSelectorTypes["5"].Lines, map[int]bool{14: true, 75: true, 80: true}},
		{"CssSelectorTypes TYPE_SELECTOR_TYPE", report.CssSelectorTypes["9"].Lines, map[int]bool{14: true, 18: true, 25: true, 26: true, 32: true, 41: true}},
		{"CssDimentions %", report.CssDimentions["%"].Lines, map[int]bool{21: true, 22: true, 35: true, 64: true, 65: true, 69: true, 70: true, 76: true, 77: true, 81: true, 82: true, 88: true, 102: true, 109: true, 121: true, 152: true, 158: true, 171: true, 191: true, 198: true, 225: true, 231: true, 244: true, 257: true, 292: true, 299: true, 312: true, 325: true, 360: true, 367: true}},
		{"ImgFormats png", report.ImgFormats["png"].Lines, map[int]bool{121: true, 171: true, 525: true, 545: true, 563: true, 581: true}},
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

func TestReportFromHTMLSecondFixture(t *testing.T) {
	pathDir, _ := os.Getwd()
	html, err := ioutil.ReadFile(path.Join(pathDir, "fixtures", "2.html"))
	if err != nil {
		t.Fatalf(`Not found fixture %s, %v`, "2.html", err)
	}

	report, err := ReportFromHTML(html)
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}

	var tests = []struct {
		checkType string
		got       map[int]bool
		want      map[int]bool
	}{
		{"HtmlTags div", report.HtmlTags["div"][""].Lines, map[int]bool{176: true, 177: true, 179: true, 183: true, 185: true, 187: true, 199: true, 219: true, 221: true, 223: true, 235: true, 261: true, 262: true, 264: true, 269: true, 271: true, 273: true, 285: true, 305: true, 332: true, 353: true, 382: true, 383: true, 385: true, 389: true, 391: true, 393: true, 405: true, 431: true, 432: true}},
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{23: true}},
		{"CssProperties width", report.CssProperties["width"][""].Lines, map[int]bool{66: true, 74: true, 78: true, 82: true, 101: true, 105: true, 167: true, 179: true, 185: true, 221: true, 264: true, 271: true, 355: true, 385: true, 391: true, 434: true, 440: true, 520: true, 543: true, 561: true, 592: true, 598: true, 616: true, 640: true, 720: true, 750: true, 756: true, 836: true, 859: true, 877: true}},
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{34: true, 64: true, 87: true, 143: true}},
		{"CssSelectorTypes CHILD_COMBINATOR_TYPE", report.CssSelectorTypes["3"].Lines, map[int]bool{108: true}},
		{"CssDimentions %", report.CssDimentions["%"].Lines, map[int]bool{36: true, 40: true, 44: true, 89: true, 96: true, 101: true, 105: true, 163: true, 167: true, 179: true, 185: true, 199: true, 201: true, 221: true, 235: true, 237: true, 264: true, 271: true, 285: true, 287: true, 305: true, 307: true, 332: true, 334: true, 355: true, 357: true, 385: true, 391: true, 405: true, 407: true}},
		{"ImgFormats png", report.ImgFormats["png"].Lines, map[int]bool{264: true, 561: true, 616: true, 877: true, 953: true, 1211: true, 1316: true, 1335: true, 1354: true, 1373: true}},
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

func TestReportFromHTMLThirdFixture(t *testing.T) {
	pathDir, _ := os.Getwd()
	html, err := ioutil.ReadFile(path.Join(pathDir, "fixtures", "3.html"))
	if err != nil {
		t.Fatalf(`Not found fixture %s, %v`, "3.html", err)
	}

	report, err := ReportFromHTML(html)
	if err != nil {
		t.Fatalf(`ReportFromHTML("%s"), %v`, html, err)
	}

	var tests = []struct {
		checkType string
		got       map[int]bool
		want      map[int]bool
	}{
		{"HtmlTags style", report.HtmlTags["style"][""].Lines, map[int]bool{23: true}},
		{"CssProperties display: none", report.CssProperties["display"]["none"].Lines, map[int]bool{141: true}},
		{"AtRuleCssStatements @media", report.AtRuleCssStatements["@media"][""].Lines, map[int]bool{34: true, 64: true, 83: true, 139: true}},
		{"ImgFormats png", report.ImgFormats["png"].Lines, map[int]bool{223: true, 356: true, 602: true, 940: true, 1048: true, 1367: true, 1403: true, 1422: true, 1441: true, 1460: true}},
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
