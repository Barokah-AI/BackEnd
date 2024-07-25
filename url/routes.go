package routes

import (
	"net/http"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/controller"
	"github.com/Barokah-AI/BackEnd/helper"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}

	var method, path string = r.Method, r.URL.Path
	switch {
	// jika method GET dan path / maka akan menjalankan fungsi Home
	case method == "GET" && path == "/":
		Home(w, r)
	// jika method POST dan path /chat maka akan menjalankan fungsi Chat
	case method == "POST" && path == "/chat":
		controller.Chat(w, r, config.GetEnv("HUGGINGFACE_API_KEY"))
	// jika method POST dan path /ngobrol maka akan menjalankan fungsi Ngobrol
	case method == "POST" && path == "/ngobrol":
		controller.Ngobrol(w, r, config.GetEnv("HUGGINGFACE_API_KEY"))
	// jika method POST dan path /signup maka akan menjalankan fungsi SignUp
	case method == "POST" && path == "/signup":
		controller.SignUp(config.Mongoconn, "users", w, r)
	// jika method POST dan path /login maka akan menjalankan fungsi LogIn
	case method == "POST" && path == "/login":
		controller.LogIn(config.Mongoconn, w, r, config.GetEnv("PASETOPRIVATEKEY"))
	// jika path tidak ditemukan maka akan merespon dengan status 404
	default:
		helper.ErrorResponse(w, r, http.StatusNotFound, "Not Found", "The requested resource was not found")
	}
}

func Home(respw http.ResponseWriter, req *http.Request) {
	resp := map[string]string{
		"github_repo": "https://github.com/barokah-ai/backend",
		"message":     "Insyallah Berkah 🤞",
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}
