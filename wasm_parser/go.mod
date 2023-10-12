module github.com/le0pard/vmail/wasm_parser

go 1.17

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20231009104354-3f512a0fe065

require (
	github.com/tdewolff/parse/v2 v2.6.8 // indirect
	golang.org/x/net v0.17.0 // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
