package main

// Import the package to access the Wasm environment
import (
	"encoding/json"
	"errors"
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

				reportStr, err := json.Marshal(report)
				if err != nil {
					rejectWithError(reject, "Error to dump JSON")
					return
				}

				// Resolve the Promise, passing anything back to JavaScript
				// This is done by invoking the "resolve" function passed to the handler
				resolve.Invoke(string(reportStr))
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
