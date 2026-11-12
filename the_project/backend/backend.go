package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	filename = getFilename()
)

func getFilename() string {
	filename := os.Getenv("FILENAME")
	if filename == "" {
		log.Fatal("FILENAME empty")
	}
	return filename
}

func SimpleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve the file
		next.ServeHTTP(w, r)

		// Log the result
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func handleGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Error reading todo.txt: %v", err)
			http.Error(w, "Could not read todos file", http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully read todo.txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	})
}

func handlePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		log.Println("writing to todos.txt with", body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			log.Fatal(err)
		}
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening todo.txt for writing: %v", err)
			log.Fatal(err)
		}
		defer f.Close()

		if _, err := f.Write(body); err != nil {
			log.Printf("Error writing body to todo.txt: %v", err)
			log.Fatal(err)
		}
		if _, err := f.WriteString("\n"); err != nil {
			log.Printf("Error appending newline to todo.txt: %v", err)
			log.Fatal(err)
		}
		log.Printf("Successfully wrote to todo.txt")
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()

	mux.Handle("GET /todo", handleGet())
	mux.Handle("POST /todo", handlePost())

	fmt.Printf("Server started in port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
