package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // To register the driver.
)

var (
	Db *sql.DB
)

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
		rows, err := Db.Query("SELECT todo from todo")
		if err != nil {
			log.Printf("Error querying todos: %v", err)
			http.Error(w, "Could not query todos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []string
		for rows.Next() {
			var t string
			if err := rows.Scan(&t); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			todos = append(todos, t)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error during rows iteration: %v", err)
			http.Error(w, "Could not process todos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		for _, t := range todos {
			w.Write([]byte(t + "\n"))
		}
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

		todoText := string(body)
		_, err = Db.Exec("INSERT INTO todo (todo) VALUES ($1)", todoText)
		if err != nil {
			log.Printf("Error inserting into database: %v", err)
			http.Error(w, "Could not save todo", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Todo added successfully\n"))
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

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s connect_timeout=5 sslmode=disable",
		dbHost, dbUser, dbPassword, dbName,
	)
	Db = connectWithRetry(connStr)
	defer Db.Close()
	if _, err := Db.Exec("CREATE TABLE IF NOT EXISTS todo (id SERIAL PRIMARY KEY, todo TEXT)"); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	if _, err := Db.Exec("INSERT INTO todo (id, todo) VALUES (1, 'initial todo') ON CONFLICT (id) DO NOTHING"); err != nil {
		log.Fatalf("Failed to seed todo: %v", err)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /todo", handleGet())
	mux.Handle("POST /todo", handlePost())

	fmt.Printf("Server started in port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
