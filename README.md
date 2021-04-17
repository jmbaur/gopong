# GoPong

## Running

```
$ docker build -t gopong .
$ docker run -p 8080:80 gopong
```

## Building the WASM

```
$ GOOS=js GOARCH=wasm go build -o assets/main.wasm pong/main.go
```

## Running the HTTP server

This will build the WASM, place it in the ./assets directory, and serve up files in that directory.

```
$ go run main.go
```
