package config

import (
	"net/http"
)

// Daftar origins yang diizinkan
var Origins = []string{
	"http://localhost:8080",
	"http://localhost:5173",
	"http://127.0.0.1:5500",
	"http://127.0.0.1:5501",
	"http://127.0.0.1:5503",
	"https://barokah-ai.vercel.app",
}

// Fungsi untuk memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}
