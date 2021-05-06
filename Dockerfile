FROM golang:alpine as builder

WORKDIR /project
RUN apk add make
COPY . .
RUN make

FROM alpine
COPY --from=builder /project/bin/server /server
COPY --from=builder /project/assets /assets
CMD ["/server"]
