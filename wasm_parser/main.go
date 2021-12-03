package main

// Import the package to access the Wasm environment
import (
	"errors"
	"sort"
	"sync"
	"syscall/js"

	"github.com/le0pard/vmail/wasm_parser/parser"
)

type ReportNestedLevelMap struct {
	Data    map[string]map[string]parser.ReportContainer
	JsonKey string
}

type ReportOneLevelMap struct {
	Data    map[string]parser.ReportContainer
	JsonKey string
}

type ReportItemMap struct {
	Data    parser.ReportContainer
	JsonKey string
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

func collectItemReport(item parser.ReportContainer) map[string]interface{} {
	// hash to slice
	lines := make([]int, 0, len(item.Lines))
	for line, _ := range item.Lines {
		lines = append(lines, line)
	}
	// sort slice with positions
	sort.Ints(lines)

	linesObj := make([]interface{}, len(lines))
	for i, line := range lines {
		linesObj[i] = line
	}

	report := map[string]interface{}{
		"rules":      item.Rules,
		"lines":      linesObj,
		"more_lines": item.MoreLines,
	}
	return report
}

func collectOneLevelReport(items map[string]parser.ReportContainer) map[string]interface{} {
	itemsReports := make(map[string]interface{})
	for k1, v1 := range items {
		itemsReports[k1] = collectItemReport(v1)
	}
	return itemsReports
}

func collectNestedLevelReport(items map[string]map[string]parser.ReportContainer) map[string]interface{} {
	itemsReports := make(map[string]interface{})
	for k1, v1 := range items {
		itemsReports[k1] = collectOneLevelReport(v1)
	}
	return itemsReports
}

func normalizeReportForPromise(report *parser.ParseReport) map[string]interface{} {
	var (
		wg sync.WaitGroup
		mx sync.RWMutex
	)

	newReport := make(map[string]interface{})

	var nestedLevelKeys []ReportNestedLevelMap = []ReportNestedLevelMap{
		ReportNestedLevelMap{
			Data:    report.HtmlTags,
			JsonKey: "html_tags",
		},
		ReportNestedLevelMap{
			Data:    report.HtmlAttributes,
			JsonKey: "html_attributes",
		},
		ReportNestedLevelMap{
			Data:    report.CssProperties,
			JsonKey: "css_properties",
		},
		ReportNestedLevelMap{
			Data:    report.AtRuleCssStatements,
			JsonKey: "at_rule_css_statements",
		},
	}

	for _, k := range nestedLevelKeys {
		if len(k.Data) > 0 {
			wg.Add(1)

			go func(items map[string]map[string]parser.ReportContainer, jsonKey string) {
				defer wg.Done()
				attrs := collectNestedLevelReport(items)
				mx.Lock()
				defer mx.Unlock()
				newReport[jsonKey] = attrs
			}(k.Data, k.JsonKey)
		}
	}

	var oneLevelKeys []ReportOneLevelMap = []ReportOneLevelMap{
		ReportOneLevelMap{
			Data:    report.CssSelectorTypes,
			JsonKey: "css_selector_types",
		},
		ReportOneLevelMap{
			Data:    report.CssDimentions,
			JsonKey: "css_dimentions",
		},
		ReportOneLevelMap{
			Data:    report.CssFunctions,
			JsonKey: "css_functions",
		},
		ReportOneLevelMap{
			Data:    report.CssPseudoSelectors,
			JsonKey: "css_pseudo_selectors",
		},
		ReportOneLevelMap{
			Data:    report.ImgFormats,
			JsonKey: "img_formats",
		},
		ReportOneLevelMap{
			Data:    report.LinkTypes,
			JsonKey: "link_types",
		},
	}

	for _, k := range oneLevelKeys {
		if len(k.Data) > 0 {
			wg.Add(1)

			go func(items map[string]parser.ReportContainer, jsonKey string) {
				defer wg.Done()
				attrs := collectOneLevelReport(items)
				mx.Lock()
				defer mx.Unlock()
				newReport[jsonKey] = attrs
			}(k.Data, k.JsonKey)
		}
	}

	var singleItemKeys []ReportItemMap = []ReportItemMap{
		ReportItemMap{
			Data:    report.CssVariables,
			JsonKey: "css_variables",
		},
		ReportItemMap{
			Data:    report.CssImportant,
			JsonKey: "css_important",
		},
		ReportItemMap{
			Data:    report.Html5Doctype,
			JsonKey: "html5_doctype",
		},
	}

	for _, k := range singleItemKeys {
		if len(k.Data.Lines) > 0 {
			wg.Add(1)

			go func(item parser.ReportContainer, jsonKey string) {
				defer wg.Done()
				attrs := collectItemReport(item)
				mx.Lock()
				defer mx.Unlock()
				newReport[jsonKey] = attrs
			}(k.Data, k.JsonKey)
		}
	}

	wg.Wait()

	return newReport
}

// VMailParser returns a JavaScript function
func VMailParser() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Get the HTML as argument
		// args[0] is a js.Value, so we need to get a string out of it
		htmlBody := args[0].String()
		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(promiseThis js.Value, promiseArgs []js.Value) interface{} {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]
			// Now that we have a way to return the response to JS, spawn a goroutine
			// This way, we don't block the event loop and avoid a deadlock
			go func() {

				report, err := parser.ReportFromHTML([]byte(htmlBody))
				if err != nil {
					rejectWithError(reject, err.Error())
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

// Main function: it sets up our Wasm application
func main() {
	// Define the function "VMailParser" in the JavaScript scope
	js.Global().Set("VMailParser", VMailParser())
	// Prevent the function from returning, which is required in a wasm module
	select {}
}
