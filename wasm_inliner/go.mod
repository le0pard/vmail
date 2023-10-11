module github.com/le0pard/vmail/wasm_inliner

go 1.17

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20230918103829-de06841b4c95

require (
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/tdewolff/parse/v2 v2.6.8 // indirect
	golang.org/x/net v0.17.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
