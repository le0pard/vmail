package inliner

import (
	"testing"
)

// htmlDoc := `<html>
// <head>
// 	<title>Email</title>
// 	<link href="https://getbootstrap.com/docs/5.1/dist/css/bootstrap.css" rel="stylesheet" />
// 	<link href="https://getbootstrap.com/docs/5.1/assets/css/docs.css" rel="stylesheet" media="screen and (max-width: 600px)" />
// 	<link href="https://getbootstrap.com/docs/5.1/assets/css/docs2.css" rel="stylesheet" media="print" />
// 	<link href="https://getbootstrap.com/docs/5.1/assets/css/docs3.css" rel="stylesheet" media="(min-width: 400px)" />
// 	<style type="text/css">
// 		/* Client-specific Styles */
// 		#outlook a{padding:0;} /* Force Outlook to provide a "view in browser" button. */
// 		body{width:100% !important;} .ReadMsgBody{width:100%;} .ExternalClass{width:100%;} /* Force Hotmail to display emails at full width */
// 		body{-webkit-text-size-adjust:none;} /* Prevent Webkit platforms from changing default text sizes. */

// 		/* Reset Styles */
// 		body{margin:0; padding:0;}
// 		img{border:0; height:auto; line-height:100%; outline:none; text-decoration:none;}
// 		table td{border-collapse:collapse;}
// 		#backgroundTable{height:100% !important; margin:0; padding:0; width:100% !important;}

// 		body, #backgroundTable{
// 			/*@editable*/ background-color:#FAFAFA;
// 		}

// 		#templateContainer{
// 			/*@editable*/ border: 1px solid #DDDDDD;
// 		}

// 		h1, .h1, .h1-header {
// 			color:#202020;
// 			display:block;
// 			font-family:Arial;
// 			font-size:34px;
// 			font-weight:bold;
// 			line-height:100%;
// 			margin-top:0 !important;
// 			margin-right:0 !important;
// 			margin-bottom:10px !important;
// 			margin-left:0 !important;
// 			text-align:left;
// 		}
// 	</style>
// </head>
// <body>
// 	<div>
// 		<h2>Title</h2>
// 	</div>
// </body>
// </html>`

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
	_, err := InlineCssInHTML([]byte(htmlDoc))
	if err != nil {
		t.Fatalf(`InlineCssInHTML("%s"), %v`, htmlDoc, err)
	}
	// log.Printf("report: %v\n", report)
}
