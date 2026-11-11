FROM golang:1.26.8-alpine as builder

WORKDIR /app

copy . .
run go build -o backend backend.go

from alpine:latest

workdir /root/

copy --from=builder /app/backend /root/backend
RUN echo "initial todo" > /root/todo.txt

cmd ["./backend"]
