package main

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	v4 := uuid.New()

	for true {
		data := time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + ": " + v4.String() + "\n"

		f, err := os.OpenFile(os.Getenv("OUTPUT_DIR")+"output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if _, err := f.WriteString(data); err != nil {
			log.Fatal(err)
		}
		time.Sleep(5 * time.Second)
	}
}
