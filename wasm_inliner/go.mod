module github.com/le0pard/vmail/wasm_inliner

go 1.24.0

toolchain go1.24.1

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20250910110424-f5484e930a46

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/tdewolff/parse/v2 v2.8.3 // indirect
	golang.org/x/net v0.44.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
