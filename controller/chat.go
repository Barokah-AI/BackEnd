package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
	"github.com/go-resty/resty/v2"
)

func Chat(respwd http.ResponseWriter, request *http.Request, tokenmodel string) {
	var chat model.AIRequest

	err := json.NewDecoder(request.Body).Decode(&chat)
	if err != nil {
		helper.ErrorResponse(respwd, request, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if chat.Prompt == "" {
		helper.ErrorResponse(respwd, request, http.StatusBadRequest, "Bad Request", "masukin pertanyaan dulu ya kakak 🤗")
		return
	}

	client := resty.New()

}
