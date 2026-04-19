module github.com/le0pard/vmail/wasm_parser

go 1.25.0

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20260418123752-128b1209f74d

require (
	github.com/tdewolff/parse/v2 v2.8.11 // indirect
	golang.org/x/net v0.53.0 // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
