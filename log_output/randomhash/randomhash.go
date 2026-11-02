package main

import (
	"encoding/json"
	"fmt"
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
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			time.Now().UTC().Format("2006-01-02T15:04:05.000Z"): v4.String(),
		}

		json.NewEncoder(w).Encode(response)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	v4 := uuid.New()

	http.Handle("/", handler(v4))

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}()
	for true {
		fmt.Println(time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + ": " + v4.String())
		time.Sleep(5 * time.Second)
	}
}
