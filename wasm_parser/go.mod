module github.com/le0pard/vmail/wasm_parser

go 1.23.0

toolchain go1.24.1

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20250319212121-46c4e83b7f81

require (
	github.com/tdewolff/parse/v2 v2.7.21 // indirect
	golang.org/x/net v0.37.0 // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
