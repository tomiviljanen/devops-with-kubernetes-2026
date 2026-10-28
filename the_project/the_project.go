package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dir := http.Dir("./static")
	handler := http.FileServer(dir)
	http.Handle("/", handler)
	fmt.Println("Server started in port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
