module github.com/le0pard/vmail/wasm_inliner

go 1.17

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20220807182418-3ab8dd5f03b3

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/tdewolff/parse/v2 v2.6.2 // indirect
	golang.org/x/net v0.0.0-20220812174116-3211cb980234 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
