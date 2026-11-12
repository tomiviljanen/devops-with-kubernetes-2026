package main

import (
	"bytes"
	"cmp"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var templates *template.Template

type Message struct {
	Text string `json:"text"`
}

var (
	todoUrl       = getTodoUrl()
	imageURL      = getImageUrl()
	downloadMutex sync.Mutex
	lastUpdate    time.Time
	imagePath     = getImagePath()
)

func getTodoUrl() string {
	todoUrl := os.Getenv("TODO_URL")
	if todoUrl == "" {
		log.Fatal("TODO_URL missing")
	}
	return todoUrl
}

func getImageUrl() string {
	pictureUrl := os.Getenv("PICTURE_URL")
	if pictureUrl == "" {
		log.Fatal("PICTURE_URL missing")
	}
	return pictureUrl
}

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

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func getTodos() []string {
	resp, err := http.Get(todoUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	todos := strings.Split(string(bytes), "\n")
	return filter(todos, func(s string) bool { return strings.TrimSpace(s) != "" })
}

func initTemplates() {
	var err error
	templateLocation := os.Getenv("TEMPLATE_LOCATION")
	if templateLocation == "" {
		log.Fatal("TEMPLATE_LOCATION missing")
	}
	templates, err = template.ParseFiles(templateLocation)
	if err != nil {
		log.Fatal(err)
	}
}

type PageData struct {
	Port  string
	Todos []string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Port:  cmp.Or(os.Getenv("PORT"), "8080"),
		Todos: getTodos(),
	}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func handleProject() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				http.Error(w, "Error reading body", http.StatusInternalServerError)
				return
			}
			log.Printf("Forwarding POST request to backend with body: %s", string(body))
			resp, err := http.Post(todoUrl, "text/plain", bytes.NewBuffer(body))
			if err != nil {
				log.Printf("Error forwarding to backend: %v", err)
				http.Error(w, "Error forwarding to backend", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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
	initTemplates()

	http.Handle("/", SimpleLogger(http.HandlerFunc(indexHandler)))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/todo", handleProject())

	fmt.Printf("Server started in port: %s with %s\n", port, appId.String())
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
