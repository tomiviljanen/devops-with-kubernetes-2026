FROM golang:1.26.8-alpine as builder

WORKDIR /app

copy . .

run go build -o randomhash-server randomhash-server.go

from alpine:latest

workdir /root/

copy --from=builder /app/randomhash-server .

cmd ["./randomhash-server"]
