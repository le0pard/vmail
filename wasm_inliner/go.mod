module github.com/le0pard/vmail/wasm_inliner

go 1.17

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20220919190428-08b7a50c0642

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/tdewolff/parse/v2 v2.6.3 // indirect
	golang.org/x/net v0.0.0-20220930213112-107f3e3c3b0b // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
