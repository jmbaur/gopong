FROM golang:alpine AS builder
RUN apk add --no-cache make
WORKDIR /go/src/github.com/jmbaur/gopong
COPY pong/main.go go.* ./
RUN GOOS=js GOARCH=wasm go build -o main.wasm main.go

FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY assets /usr/share/nginx/html
COPY --from=builder /go/src/github.com/jmbaur/gopong/main.wasm /usr/share/nginx/html/main.wasm
