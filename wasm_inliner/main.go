package main

// Import the package to access the Wasm environment
import (
	"errors"
	"syscall/js"

	"github.com/le0pard/vmail/wasm_inliner/inliner"
)

// Main function: it sets up our Wasm application
func main() {
	// Define the function "VMailInliner" in the JavaScript scope
	js.Global().Set("VMailInliner", VMailInliner())
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
func VMailInliner() js.Func {
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

				htmlResult, err := inliner.InlineCssInHTML([]byte(htmlBody))
				if err != nil {
					rejectWithError(reject, err.Error())
					return
				}

				// Resolve the Promise, passing anything back to JavaScript
				// This is done by invoking the "resolve" function passed to the handler
				resolve.Invoke(string(htmlResult))
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
