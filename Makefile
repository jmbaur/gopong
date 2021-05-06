.PHONY: all wasm server

all: wasm server

wasm:
	GOOS=js GOARCH=wasm go build  \
	     -o assets/main.wasm \
	     pong/main.go && \
	     gzip -f assets/main.wasm



server:
	CGO_ENABLED=0 go build \
		    -tags netgo -ldflags '-w' -a \
		    -o bin/server main.go

clean:
	rm -rf bin/ && rm -f assets/main.wasm.gz
