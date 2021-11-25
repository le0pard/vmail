module github.com/le0pard/vmail/wasm_parser

go 1.17

require github.com/le0pard/vmail/wasm_parser/parser v0.0.0-20211027174133-e2c0af7319db

require (
	github.com/tdewolff/parse/v2 v2.5.23-0.20211101212351-646f46fcfe51 // indirect
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3 // indirect
)

replace github.com/le0pard/vmail/wasm_parser/parser => ./parser
