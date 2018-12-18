package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	title := "Jenkins X SonarQube App Integration Test"

	from := ""
	if r.URL != nil {
		from = r.URL.String()
	}
	if from != "/favicon.ico" {
		log.Printf("title: %s\n", title)
	}

	fmt.Fprintf(w, "Hello from this title :  "+title+"\n")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
