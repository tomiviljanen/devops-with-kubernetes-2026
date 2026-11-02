package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Message struct {
	Text string `json:"text"`
}

func handler(val *int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("pong " + strconv.Itoa((*val))))
		(*val)++
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	val := 0

	http.Handle("/", handler(&val))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
