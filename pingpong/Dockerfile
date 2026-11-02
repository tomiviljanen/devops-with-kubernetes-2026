FROM golang:1.26.8-alpine as builder

WORKDIR /app

copy . .

run go build -o pingpong pingpong.go

from alpine:latest

workdir /root/

copy --from=builder /app/pingpong .
RUN mkdir -p /usr/src/app/files
cmd ["./pingpong"]
