package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Text string `json:"text"`
}

func handler(v4 uuid.UUID) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(os.Getenv("INPUT_DIR") + "output.txt")
		if err != nil {
			log.Fatalf("Failed to read file: %s", err)
		}

		data := time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + ": " + v4.String() + "\nPing / Pongs: " + string(content)

		w.Write([]byte(data))
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	v4 := uuid.New()

	http.Handle("/", handler(v4))

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
