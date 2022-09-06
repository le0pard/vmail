module github.com/le0pard/vmail/wasm_parser

go 1.17

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20220901104053-41146bd2af63

require (
	github.com/tdewolff/parse/v2 v2.6.3 // indirect
	golang.org/x/net v0.0.0-20220906165146-f3363e06e74c // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
