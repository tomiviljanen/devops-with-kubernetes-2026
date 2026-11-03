package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Text string `json:"text"`
}

const imageURL = "https://picsum.photos/400"

var (
	downloadMutex sync.Mutex
	lastUpdate    time.Time
	imagePath     = getImagePath()
)

func getImagePath() string {
	if path := os.Getenv("IMAGE_PATH"); path != "" {
		return path
	}
	return "static/image.jpg"
}

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func refreshImage() {
	downloadMutex.Lock()
	defer downloadMutex.Unlock()

	tmpPath := "static/image_new.jpg"
	if err := downloadFile(tmpPath, imageURL); err == nil {
		_ = os.Remove(imagePath)
		if err := os.Rename(tmpPath, imagePath); err == nil {
			lastUpdate = time.Now()
			fmt.Println("Image refreshed successfully")
		}
	}
}

func startBackgroundRefresh() {
	// Initial download
	go refreshImage()

	// Periodic refresh every 10 minutes
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			refreshImage()
		}
	}()
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appId := uuid.New()

	// Start the background refresh loop
	startBackgroundRefresh()

	fs := http.FileServer(http.Dir("./static/"))
	loggedFs := SimpleLogger(fs)
	http.Handle("/", loggedFs)

	fmt.Printf("Server started in port: %s with %s\n", port, appId.String())
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
