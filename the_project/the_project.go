package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Message struct {
	Text string `json:"text"`
}

func handler(appId uuid.UUID) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		newString := uuid.New()

		response := map[string]string{
			appId.String(): newString.String(),
		}

		fmt.Printf("Responding with %s\n", newString.String())
		json.NewEncoder(w).Encode(response)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appId := uuid.New()

	http.Handle("/", handler(appId))
	fmt.Printf("Server started in port: %s with %s\n", port, appId.String())
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
