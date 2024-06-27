package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
	"github.com/go-resty/resty/v2"
)

func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
	var chat model.AIRequest

	err := json.NewDecoder(req.Body).Decode(&chat)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if chat.Promt == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	client := resty.New()

    // Hugging Face API URL dan token
    apiUrl := "anu"
    apiToken := "Bearer " + "ini"

	response, err := client.R().
            SetHeader("Authorization", apiToken).
            SetHeader("Content-Type", "application/json").
            SetBody(`{"inputs": "`+chat.Promt+`"}`).
            Post(apiUrl)

        if err != nil {
            log.Fatalf("Error making request: %v", err)
        }

	var data []map[string]string

	err = json.Unmarshal([]byte(response.String()), &data)
	if err != nil {
        fmt.Println("Response:", response.String())
        fmt.Println("token", tokenmodel)
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if len(data) > 0 {
		helper.WriteJSON(respw, http.StatusOK, data[0])
	} else {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
	}
}