package main

import (
	"log"
	"net/http"

	routes "github.com/Barokah-AI/BackEnd/local/url"
)

func main() {
	http.HandleFunc("/", routes.URL)
	port := ":8080"
	log.Printf("Starting server on port %s", port + " or click this link http://localhost:8080")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
