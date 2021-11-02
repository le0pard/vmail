package main

// Import the package to access the Wasm environment
import (
	"errors"
	"sort"
	"sync"
	"syscall/js"

	"github.com/le0pard/vmail/wasm/parser"
)

// Main function: it sets up our Wasm application
func main() {
	// Define the function "VMail" in the JavaScript scope
	js.Global().Set("VMail", VMail())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}

func rejectWithError(reject js.Value, message string) {
	err := errors.New(message)

	// Create a JS Error object and pass it to the reject function
	// The constructor for Error accepts a string,
	// so we need to get the error message as string from "err"
	errorConstructor := js.Global().Get("Error")
	errorObject := errorConstructor.New(err.Error())
	reject.Invoke(errorObject)
}

func normalizeReportForPromise(report *parser.ParseReport) map[string]interface{} {
	var (
		wg sync.WaitGroup
		mx sync.RWMutex
	)

	newReport := make(map[string]interface{})

	if len(report.HtmlTags) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			tagReports := make(map[string]interface{})
			for k1, v1 := range report.HtmlTags {
				tagAttributeReports := make(map[string]interface{})
				for k2, v2 := range v1 {
					// hash to slice
					lines := make([]int, 0, len(v2.Lines))
					for line, _ := range v2.Lines {
						lines = append(lines, line)
					}
					// sort slice with positions
					sort.Ints(lines)

					linesObj := make([]interface{}, len(lines))
					for i, line := range lines {
						linesObj[i] = line
					}

					tagAttributeReports[k2] = map[string]interface{}{
						"rules":      v2.Rules,
						"lines":      linesObj,
						"more_lines": v2.MoreLines,
					}
				}
				tagReports[k1] = tagAttributeReports
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["html_tags"] = tagReports
		}()
	}

	if len(report.CssProperties) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			cssPropertyReports := make(map[string]interface{})
			for k1, v1 := range report.CssProperties {
				cssValReports := make(map[string]interface{})
				for k2, v2 := range v1 {
					// hash to slice
					lines := make([]int, 0, len(v2.Lines))
					for line, _ := range v2.Lines {
						lines = append(lines, line)
					}
					// sort slice with positions
					sort.Ints(lines)

					linesObj := make([]interface{}, len(lines))
					for i, line := range lines {
						linesObj[i] = line
					}

					cssValReports[k2] = map[string]interface{}{
						"rules":      v2.Rules,
						"lines":      linesObj,
						"more_lines": v2.MoreLines,
					}
				}
				cssPropertyReports[k1] = cssValReports
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["css_properties"] = cssPropertyReports
		}()
	}

	if len(report.CssSelectorTypes) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			cssSelectorsTypeReports := make(map[string]interface{})
			for k1, v1 := range report.CssSelectorTypes {
				// hash to slice
				lines := make([]int, 0, len(v1.Lines))
				for line, _ := range v1.Lines {
					lines = append(lines, line)
				}
				// sort slice with positions
				sort.Ints(lines)

				linesObj := make([]interface{}, len(lines))
				for i, line := range lines {
					linesObj[i] = line
				}

				cssSelectorsTypeReports[k1] = map[string]interface{}{
					"rules":      v1.Rules,
					"lines":      linesObj,
					"more_lines": v1.MoreLines,
				}
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["css_selector_types"] = cssSelectorsTypeReports
		}()
	}

	if len(report.CssDimentions) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			cssDimentionsReports := make(map[string]interface{})
			for k1, v1 := range report.CssDimentions {
				// hash to slice
				lines := make([]int, 0, len(v1.Lines))
				for line, _ := range v1.Lines {
					lines = append(lines, line)
				}
				// sort slice with positions
				sort.Ints(lines)

				linesObj := make([]interface{}, len(lines))
				for i, line := range lines {
					linesObj[i] = line
				}

				cssDimentionsReports[k1] = map[string]interface{}{
					"rules":      v1.Rules,
					"lines":      linesObj,
					"more_lines": v1.MoreLines,
				}
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["css_dimentions"] = cssDimentionsReports
		}()
	}

	if len(report.CssFunctions) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			cssFunctionsReports := make(map[string]interface{})
			for k1, v1 := range report.CssFunctions {
				// hash to slice
				lines := make([]int, 0, len(v1.Lines))
				for line, _ := range v1.Lines {
					lines = append(lines, line)
				}
				// sort slice with positions
				sort.Ints(lines)

				linesObj := make([]interface{}, len(lines))
				for i, line := range lines {
					linesObj[i] = line
				}

				cssFunctionsReports[k1] = map[string]interface{}{
					"rules":      v1.Rules,
					"lines":      linesObj,
					"more_lines": v1.MoreLines,
				}
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["css_functions"] = cssFunctionsReports
		}()
	}

	if len(report.CssPseudoSelectors) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			cssPseudoSelectorsReports := make(map[string]interface{})
			for k1, v1 := range report.CssPseudoSelectors {
				// hash to slice
				lines := make([]int, 0, len(v1.Lines))
				for line, _ := range v1.Lines {
					lines = append(lines, line)
				}
				// sort slice with positions
				sort.Ints(lines)

				linesObj := make([]interface{}, len(lines))
				for i, line := range lines {
					linesObj[i] = line
				}

				cssPseudoSelectorsReports[k1] = map[string]interface{}{
					"rules":      v1.Rules,
					"lines":      linesObj,
					"more_lines": v1.MoreLines,
				}
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["css_pseudo_selectors"] = cssPseudoSelectorsReports
		}()
	}

	if len(report.ImgFormats) > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			imgFormatsReports := make(map[string]interface{})
			for k1, v1 := range report.ImgFormats {
				// hash to slice
				lines := make([]int, 0, len(v1.Lines))
				for line, _ := range v1.Lines {
					lines = append(lines, line)
				}
				// sort slice with positions
				sort.Ints(lines)

				linesObj := make([]interface{}, len(lines))
				for i, line := range lines {
					linesObj[i] = line
				}

				imgFormatsReports[k1] = map[string]interface{}{
					"rules":      v1.Rules,
					"lines":      linesObj,
					"more_lines": v1.MoreLines,
				}
			}
			mx.Lock()
			defer mx.Unlock()
			newReport["img_formats"] = imgFormatsReports
		}()
	}

	wg.Wait()

	return newReport
}

// VMail returns a JavaScript function
func VMail() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Get the URL as argument
		// args[0] is a js.Value, so we need to get a string out of it
		htmlBody := args[0].String()
		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			// Now that we have a way to return the response to JS, spawn a goroutine
			// This way, we don't block the event loop and avoid a deadlock
			go func() {
				report, err := parser.ReportFromHTML([]byte(htmlBody))
				if err != nil {
					rejectWithError(reject, "Error to parser HTML")
					return
				}

				// Resolve the Promise, passing anything back to JavaScript
				// This is done by invoking the "resolve" function passed to the handler
				resolve.Invoke(normalizeReportForPromise(report))
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
