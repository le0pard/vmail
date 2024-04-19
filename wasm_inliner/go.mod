module github.com/le0pard/vmail/wasm_inliner

go 1.17

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20240403175846-8a4d1f8710e5

require (
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/tdewolff/parse/v2 v2.7.12 // indirect
	golang.org/x/net v0.24.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
