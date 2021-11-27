package inliner

import (
	"regexp"
	"strings"
	"testing"
)

func TestInlineCssInHTMLSimple(t *testing.T) {
	htmlDoc := `<html>
<head>
	<title>Email</title>
	<style type="text/css">
		@media (prefers-reduced-motion: no-preference) {
			:root {
				scroll-behavior: smooth;
			}
		}

		#outlook a {
      padding: 0;
    }

    body {
      margin: 0;
      padding: 0;
      -webkit-text-size-adjust: 100%;
      -ms-text-size-adjust: 100%;
    }

    table,
    td {
      border-collapse: collapse;
      mso-table-lspace: 0pt;
      mso-table-rspace: 0pt;
    }

    img {
      border: 0;
      height: auto;
      line-height: 100%;
      outline: none;
      text-decoration: none;
      -ms-interpolation-mode: bicubic;
    }

		a:focus {
			text-decoration: none;
		}

		h1 {
			color:#202020;
			display:block !important;
			font-family:Arial;
			font-size:34px;
			font-weight:bold;
		}

		.h1 {
			color: #000;
			line-height:100%;
			margin-top:0 !important;
			margin-right:0 !important;
			margin-bottom:10px !important;
			margin-left:0 !important;
			text-align:left;
		}
	</style>
</head>
<body>
	<div>
		<h1 class="h1">Title</h1>
	</div>
</body>
</html>`
	htmlResult, err := InlineCssInHTML([]byte(htmlDoc))
	if err != nil {
		t.Fatalf(`InlineCssInHTML("%s"), %v`, htmlDoc, err)
	}

	htmlResultStr := string(htmlResult)

	// log.Printf("report: %v\n", string(htmlResult))
	if !strings.Contains(htmlResultStr, "<style type=\"text/css\">@media(prefers-reduced-motion:no-preference){:root{scroll-behavior:smooth;}") {
		t.Errorf("InlineCssInHTML not found inlined style tag in %v", htmlResultStr)
	}

	colorRe := regexp.MustCompile(`(?i)<h1.*style=".*color:#000;.*"`)
	if !colorRe.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined color property in h1 tag in %v", htmlResultStr)
	}
	lineHeightRe := regexp.MustCompile(`(?i)<h1.*style=".*line-height:100%;.*"`)
	if !lineHeightRe.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined line-height property in h1 tag in %v", htmlResultStr)
	}
	marginRightRe := regexp.MustCompile(`(?i)<h1.*style=".*margin-right:0;.*"`)
	if !marginRightRe.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined margin-right property in h1 tag in %v", htmlResultStr)
	}
}

func TestInlineCssInUrls(t *testing.T) {
	htmlDoc := `<html>
<head>
	<title>Email</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" />
	<link href="docs/5.1/assets/css/docs.css" rel="stylesheet" media="screen and (max-width: 600px)" />
	<link href="https://getbootstrap.com/not-exists-css-file.css" rel="stylesheet" media="all" />
	<link href="https://example.com/skipped.css" rel="stylesheet" media="print" />
	<link href="https://example.com/skipped.css" rel="stylesheet" media="(min-width: 400px)" />
	<link href="/url-not-works.css" rel="stylesheet" media="all" />
</head>
<body>
	<div class="container">
		<h2>
			Title
			<small class="text-muted">With faded secondary text</small>
		</h2>
	</div>
</body>
</html>`
	htmlResult, err := InlineCssInHTML([]byte(htmlDoc))
	if err != nil {
		t.Fatalf(`InlineCssInHTML("%s"), %v`, htmlDoc, err)
	}

	htmlResultStr := string(htmlResult)

	// log.Printf("report: %v\n", string(htmlResult))
	stylesRe := regexp.MustCompile(`(?i)<style.*>.*@charset "UTF-8";:root{.*"`)
	if !stylesRe.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined style tag in %v", htmlResultStr)
	}

	divRe := regexp.MustCompile(`(?i)<div.*style=".*box-sizing:border-box;.*"`)
	if !divRe.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined box-sizing property in div tag in %v", htmlResultStr)
	}
	h2Re := regexp.MustCompile(`(?i)<h2.*style=".*line-height:1.2;.*"`)
	if !h2Re.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined line-height property in h2 tag in %v", htmlResultStr)
	}
	smallRe := regexp.MustCompile(`(?i)<small.*style=".*font-size:.875em;.*"`)
	if !smallRe.MatchString(htmlResultStr) {
		t.Errorf("InlineCssInHTML not found inlined font-size property in small tag in %v", htmlResultStr)
	}
}
