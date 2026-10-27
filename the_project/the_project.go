package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func helloworld() http.Handler {
	type response struct {
		Message string `json:"message"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := response{Message: "Hello world"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", helloworld())
	fmt.Println("Server started in port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
