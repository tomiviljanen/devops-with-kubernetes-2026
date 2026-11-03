package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Text string `json:"text"`
}

var (
	countUrl = getUrl()
)

func getUrl() string {
	if path := os.Getenv("COUNT_URL"); path != "" {
		return path
	}
	return "http://localhost:8080/count"
}

func handler(v4 uuid.UUID) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp, err := http.Get(countUrl)
		if err != nil {
			log.Fatalf("Failed to read url: %s", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		data := time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + ": " + v4.String() + "\nPing / Pongs: " + string(body)

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
