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

	// Hugging Face API URL dan token
	apiUrl := config.GetEnv("HUGGINGFACE_API_URL")
	apiToken := "Bearer " + config.GetEnv("HUGGINGFACE_API_KEY")

	response, err := client.R().
		SetHeader("Authorization", apiToken).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"inputs": "` + chat.Prompt + `"}`).
		Post(apiUrl)

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	// Log response body
	// log.Printf("Response from Hugging Face API: %s", response.String())

	// Periksa jika model sedang dimuat
	if response.StatusCode() == http.StatusServiceUnavailable {
		helper.ErrorResponse(respwd, request, http.StatusServiceUnavailable, "Internal Server Error", "Model sedang dimuat, coba lagi sebentar ya kakak 🙏 | HF Response: "+response.String())
		return
	}

	// Periksa jika model tidak ditemukan
	if response.StatusCode() == http.StatusNotFound {
		helper.ErrorResponse(respwd, request, http.StatusNotFound, "Not Found", "Model tidak ditemukan | HF Response: "+response.String())
		return
	}

	// Periksa jika model mengembalikan status code lain
	if response.StatusCode() != http.StatusOK {
		helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: "+response.String())
		return
	}

	// Handle the expected nested array structure
	var nestedData [][]map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &nestedData)
	if err != nil {
		helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "error decoding response: "+err.Error()+" | Server HF Response: "+response.String())
		return
	}

	// Flatten the nested array structure
	var flatData []map[string]interface{}
	for _, d := range nestedData {
		flatData = append(flatData, d...)
	}

	// Extracting the highest scoring label from the model output
	var bestLabel string
	var highestScore float64
	for _, item := range flatData {
		label, labelOk := item["label"].(string)
		score, scoreOk := item["score"].(float64)
		if labelOk && scoreOk && (bestLabel == "" || score > highestScore) {
			bestLabel = label
			highestScore = score
		}
	}

	if bestLabel != "" {
		// Load the dataset from GCS
		bucketName := config.GetEnv("GCS_BUCKET_NAME")
		objectName := config.GetEnv("GCS_DATASET_FILE")

		labelToQA, err := helper.LoadDatasetGCS(bucketName, objectName)
		if err != nil {
			helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "server error: could not load dataset: "+err.Error())
			return
		}

		// Get the answer corresponding to the best label
		record, ok := labelToQA[bestLabel]
		if !ok {
			helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "server error: label not found in dataset")
			return
		}

		answer := record[1]

		helper.WriteJSON(respwd, http.StatusOK, map[string]string{
			"prompt":   chat.Prompt,
			"response": answer,
			"label":    bestLabel,
			"score":    strconv.FormatFloat(highestScore, 'f', -1, 64),
		})
	} else {
		helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
	}
}
