package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(os.Getenv("INPUT_DIR") + "output.txt")
		if err != nil {
			log.Fatalf("Failed to read file: %s", err)
		}

		w.Write(content)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/", handler())

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
