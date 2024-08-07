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
	var nested_data [][]map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &nested_data)
	if err != nil {
		helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "error decoding response: "+err.Error()+" | Server HF Response: "+response.String())
		return
	}

	// Flatten the nested array structure
	var flat_data []map[string]interface{}
	for _, _data := range nested_data {
		flat_data = append(flat_data, _data...)
	}

	// Extracting the highest scoring label from the model output
	var best_label string
	var highest_score float64
	for _, item := range flat_data {
		label, labelOk := item["label"].(string)
		score, scoreOk := item["score"].(float64)
		if labelOk && scoreOk && (best_label == "" || score > highest_score) {
			best_label = label
			highest_score = score
		}
	}

	if best_label != "" {
		// Load the dataset from GCS
		bucket_name := config.GetEnv("GCS_BUCKET_NAME")
		object_name := config.GetEnv("GCS_DATASET_FILE")

		label_to_qa, err := helper.LoadDatasetGCS(bucket_name, object_name)
		if err != nil {
			helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "server error: could not load dataset: "+err.Error())
			return
		}

		// Get the answer corresponding to the best label
		records, ok := label_to_qa[best_label]
		if !ok {
			helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "server error: label not found in dataset")
			return
		}

		answers := records[1]

		helper.WriteJSON(respwd, http.StatusOK, map[string]string{
			"prompt":   chat.Prompt,
			"response": answers,
			"label":    best_label,
			"score":    strconv.FormatFloat(highest_score, 'f', -1, 64),
		})
	} else {
		helper.ErrorResponse(respwd, request, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
	}
}
