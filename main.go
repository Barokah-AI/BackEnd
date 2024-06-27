package main

import (
	"net/http"

	routes "github.com/Barokah-AI/BackEnd/url"
)

func main() {
	http.HandleFunc("/", routes.URL)
	http.ListenAndServe(":8080", nil)
}