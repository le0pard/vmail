module github.com/le0pard/vmail/wasm_inliner

go 1.22.0

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20241228222635-552335c7442d

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/tdewolff/parse/v2 v2.7.19 // indirect
	golang.org/x/net v0.33.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
