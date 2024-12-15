module github.com/le0pard/vmail/wasm_inliner

go 1.22.0

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20241209100531-e0b5efff49a5

require (
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/tdewolff/parse/v2 v2.7.19 // indirect
	golang.org/x/net v0.32.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
