package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // To register the driver.
)

type Message struct {
	Text string `json:"text"`
}

var Db *sql.DB

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Incoming request: method=%s uri=%q host=%q remote=%s",
			r.Method, r.URL.RequestURI(), r.Host, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed request: method=%s uri=%q duration=%s",
			r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var counter int
		err := Db.QueryRow("UPDATE counters SET val = val + 1 WHERE id = $1 RETURNING val", 1).Scan(&counter)
		if err != nil {
			log.Fatal(err)
		}
		numString := strconv.Itoa((counter))
		w.Write([]byte("pong " + numString))
	})
}

func countHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var counter int
		err := Db.QueryRow("SELECT val from counters where id = $1", 1).Scan(&counter)
		if err != nil {
			log.Fatal(err)
		}
		numString := strconv.Itoa((counter))
		w.Write([]byte(numString))
	})
}

func connectWithRetry(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Invalid DB config: %v", err)
	}

	for {
		err = db.Ping()
		if err == nil {
			log.Println("Connected to PostgreSQL!")
			return db
		}
		log.Printf("Waiting for PostgreSQL to be ready (%v)... retrying in 2s", err)
		time.Sleep(2 * time.Second)

	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	connStr := "host=postgres-service user=postgres password=password dbname=pingpong_db connect_timeout=5 sslmode=disable"
	Db = connectWithRetry(connStr)
	defer Db.Close()
	if _, err := Db.Exec("CREATE TABLE IF NOT EXISTS counters (id SERIAL PRIMARY KEY, val INT)"); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	if _, err := Db.Exec("INSERT INTO counters (id, val) VALUES (1, 0) ON CONFLICT (id) DO NOTHING"); err != nil {
		log.Fatalf("Failed to seed counter: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler())
	mux.Handle("/count", countHandler())

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
