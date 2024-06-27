package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
	"github.com/go-resty/resty/v2"
)

// Struct untuk hasil prediksi
type Prediction struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

func Chat(respw http.ResponseWriter, req *http.Request) {
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
	apiUrl := "https://api-inference.huggingface.co/models/dimasardnt/barokah-model"
	apiToken := "Bearer " + config.GetEnv("TOKENMODEL")

	response, err := client.R().
		SetHeader("Authorization", apiToken).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"inputs": "` + chat.Promt + `"}`).
		Post(apiUrl)

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	var predictions [][]Prediction

	err = json.Unmarshal(response.Body(), &predictions)
	if err != nil {
		fmt.Println("Response:", response.String())
		fmt.Println("Error decoding JSON:", err)
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
		return
	}

	if len(predictions) > 0 && len(predictions[0]) > 0 {
		// Sort predictions by score in descending order
		sort.Slice(predictions[0], func(i, j int) bool {
			return predictions[0][i].Score > predictions[0][j].Score
		})

		// Send the top prediction
		helper.WriteJSON(respw, http.StatusOK, predictions[0][0])
	} else {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
	}
}
