module github.com/le0pard/vmail/wasm_inliner

go 1.24.0

toolchain go1.24.1

require github.com/le0pard/vmail/wasm_inliner/inliner v0.0.0-20251021222452-89de85b36b43

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/tdewolff/parse/v2 v2.8.5 // indirect
	golang.org/x/net v0.46.0 // indirect
)

replace github.com/le0pard/vmail/wasm_inliner/inliner => ./inliner
