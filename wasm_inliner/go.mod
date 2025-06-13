module github.com/le0pard/vmail/wasm_inliner

go 1.23.0

toolchain go1.24.1

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20250319212121-46c4e83b7f81

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/tdewolff/parse/v2 v2.7.21 // indirect
	golang.org/x/net v0.38.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
