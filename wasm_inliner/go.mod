module github.com/le0pard/vmail/wasm_inliner

go 1.26.0

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20261001214527-23030f2024f8

require (
	github.com/andybalholm/cascadia v1.3.5 // indirect
	github.com/tdewolff/parse/v2 v2.8.16 // indirect
	golang.org/x/net v0.59.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
