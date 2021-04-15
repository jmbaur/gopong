FROM golang:alpine AS builder
RUN apk add --no-cache make
WORKDIR /go/src/github.com/jmbaur/gopong
COPY main.go Makefile go.mod go.sum ./
RUN make

FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY assets /usr/share/nginx/html
COPY --from=builder /go/src/github.com/jmbaur/gopong/assets/main.wasm /usr/share/nginx/html/main.wasm
