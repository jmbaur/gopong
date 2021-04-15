.PHONY: all

all:
	GOOS=js GOARCH=wasm go build -o assets/main.wasm
