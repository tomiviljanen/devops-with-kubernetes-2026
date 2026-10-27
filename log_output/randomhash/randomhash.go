package main

import ("fmt"
	"time"

        "github.com/google/uuid"
)

func main() {
    v4 := uuid.New()
    for true {
       fmt.Println(time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + ": " + v4.String())
       time.Sleep(5 * time.Second)
    }
}
