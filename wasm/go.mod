module github.com/le0pard/vmail/wasm

go 1.17

require github.com/le0pard/vmail/wasm/parser v0.0.0-20211027174133-e2c0af7319db

require (
	github.com/gorilla/css v1.0.0 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
)

replace github.com/le0pard/vmail/wasm/parser => ./parser
