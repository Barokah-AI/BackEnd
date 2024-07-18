package controller

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/Barokah-AI/BackEnd/config"
// 	"github.com/Barokah-AI/BackEnd/helper"
// 	"github.com/Barokah-AI/BackEnd/model"
// 	"github.com/go-resty/resty/v2"
// )

func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
    var chat model.AIRequest

	err := json.NewDecoder(req.Body).Decode(&chat)
    if err != nil {
        helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
        return
    }

    if chat.Prompt == "" {
        helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "masukin pertanyaan dulu ya kakak 🤗")
        return
    }

    client := resty.New()

	  // Hugging Face API URL dan token
	  apiUrl := config.GetEnv("HUGGINGFACE_API_URL")
	  apiToken := "Bearer " + config.GetEnv("HUGGINGFACE_API_KEY")
  
	