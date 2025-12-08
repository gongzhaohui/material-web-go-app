run:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o  web/app.wasm
	go run .