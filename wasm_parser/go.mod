module github.com/le0pard/vmail/wasm_parser

go 1.17

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20230105173102-b3796e5d77cf

require (
	github.com/tdewolff/parse/v2 v2.6.5 // indirect
	golang.org/x/net v0.5.0 // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
