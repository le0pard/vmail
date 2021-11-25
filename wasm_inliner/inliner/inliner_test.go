package inliner

import (
	"testing"
)

func TestInlineCssInHTMLSimple(t *testing.T) {
	htmlDoc := `<html>
<head>
	<title>Email</title>
	<link href="https://getbootstrap.com/docs/5.1/dist/css/bootstrap.css" rel="stylesheet" />
	<link href="https://getbootstrap.com/docs/5.1/assets/css/docs.css" rel="stylesheet" media="screen and (max-width: 600px)" />
	<link href="https://getbootstrap.com/docs/5.1/assets/css/docs2.css" rel="stylesheet" media="print" />
	<link href="https://getbootstrap.com/docs/5.1/assets/css/docs3.css" rel="stylesheet" media="(min-width: 400px)" />
</head>
<body>
	<div>
		<h2>Title</h2>
	</div>
</body>
</html>`
	_, err := InlineCssInHTML([]byte(htmlDoc))
	if err != nil {
		t.Fatalf(`InlineCssInHTML("%s"), %v`, htmlDoc, err)
	}
	// log.Printf("report: %v\n", report)
}
