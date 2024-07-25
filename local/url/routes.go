package routes

// import package fmt, http, config, helper, controller
import (
	"net/http"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/local/controller"
)

func URL(w http.ResponseWriter, r *http.Request) {
	// Set access control headers for the response writer
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}

	var method, path string = r.Method, r.URL.Path
	// Switch case to handle different routes
	switch {
	// If method is GET and path is /, it will run Home function
	case method == "GET" && path == "/":
		// Call Home function
		Home(w, r)
	// If method is POST and path is /chat, it will run Chat function
	case method == "POST" && path == "/chat":
		// Call Chat function from controller package
		controller.Chat(w, r, config.GetEnv("HUGGINGFACE_API_KEY"))
	// If method is POST and path is /signup, it will run SignUp function
	case method == "POST" && path == "/signup":
		// Call SignUp function from controller package
		controller.SignUp(config.Mongoconn, "users", w, r)
	// If method is POST and path is /login, it will run LogIn function
	case method == "POST" && path == "/login":
		controller.LogIn(config.Mongoconn, w, r, config.GetEnv("PASETOPRIVATEKEY"))
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
