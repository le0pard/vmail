module github.com/le0pard/vmail/wasm_parser

go 1.17

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20211027174133-e2c0af7319db

require (
	github.com/tdewolff/parse/v2 v2.5.26 // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
