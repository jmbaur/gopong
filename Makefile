.PHONY: all wasm server

all: wasm server

wasm:
	GOOS=js GOARCH=wasm go build \
	     -o assets/main.wasm \
	     pong/main.go && \
	     gzip -f assets/main.wasm

server:
	go build -o bin/server main.go

clean:
	rm -rf bin/
