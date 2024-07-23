package routes

import (
	"net/http"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/local/controller"
)

// func URL(w http.ResponseWriter, r *http.Request) {
// 	if config.SetAccessControlHeaders(w, r) {
// 		return // If it's a preflight request, return early.
// 	}

// 	var method, path string = r.Method, r.URL.Path
// 	switch {
// 	case method == "GET" && path == "/":
// 		Home(w, r)
// 	case method == "POST" && path == "/chat":
// 		controller.Chat(w, r, config.GetEnv("HUGGINGFACE_API_KEY"))
// 	case method == "POST" && path == "/signup":
// 		controller.SignUp(config.Mongoconn, "users", w, r)
// 	case method == "POST" && path == "/login":
// 		controller.LogIn(config.Mongoconn, w, r, config.GetEnv("PASETOPRIVATEKEY"))
// 	default:
// 		helper.ErrorResponse(w, r, http.StatusNotFound, "Not Found", "The requested resource was not found")
// 	}
}

func Home(respw http.ResponseWriter, req *http.Request) {
	resp := map[string]string{
		"github_repo": "https://github.com/barokah-ai/backend",
		"message": "Insyallah Berkah 🤞",
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}