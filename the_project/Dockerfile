FROM golang:1.26.8-alpine as builder

WORKDIR /app

copy . .

run go build -o the_project the_project.go

from alpine:latest

workdir /root/

copy --from=builder /app/the_project .

cmd ["./the_project"]
