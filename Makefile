.PHONY: all

all:
	GOOS=js GOARCH=wasm go build \
	     -o assets/main.wasm \
	     pong/main.go && \
	     gzip -f assets/main.wasm
