# VMail

### Build wasm file

```bash
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm
```
