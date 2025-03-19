module github.com/le0pard/vmail/wasm_parser

go 1.22.0
toolchain go1.24.1

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20241228224109-0ee97081d19f

require (
	github.com/tdewolff/parse/v2 v2.7.19 // indirect
	golang.org/x/net v0.36.0 // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
