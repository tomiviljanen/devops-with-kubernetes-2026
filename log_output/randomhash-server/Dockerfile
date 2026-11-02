FROM golang:1.26.8-alpine as builder

WORKDIR /app

copy go.mod go.sum ./
run go mod download

copy . .

run go build -o randomhash-server randomhash-server.go

from alpine:latest

workdir /root/

copy --from=builder /app/randomhash-server .
RUN mkdir -p /usr/src/app/files
cmd ["./randomhash-server"]
