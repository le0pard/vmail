module github.com/le0pard/vmail/wasm_inliner

go 1.17

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20211027174133-e2c0af7319db

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/tdewolff/parse/v2 v2.5.22 // indirect
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
