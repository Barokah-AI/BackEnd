package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
	"github.com/go-resty/resty/v2"
)

// LoadDataset loads the dataset from the given CSV file and returns a map of label to question-answer pair
func LoadDataset(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dataset file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|' // Set the delimiter to pipe
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read dataset file: %v", err)
	}

	labelToQA := make(map[string][]string)
	for i, record := range records {
		if len(record) != 2 {
			log.Printf("Skipping invalid record at line %d: %v\n", i+1, record)
			continue
		}
		label := "LABEL_" + strconv.Itoa(i)
		labelToQA[label] = record
	}
	return labelToQA, nil
}

// NormalizeText normalizes the text by converting to lowercase and removing punctuation
func NormalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	// Remove punctuation
	text = strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789 ", r) {
			return r
		}
		return -1
	}, text)
	return text
}

// Chat is the controller for handling chat requests
func Chat(respw http.ResponseWriter, req *http.Request, apiKey string) {
	var chat model.AIRequest

	err := json.NewDecoder(req.Body).Decode(&chat)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if chat.Prompt == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "masukin pertanyaan dulu ya kak 🤗")
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

	// Handle the expected nested array structure
	var nestedData [][]map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &nestedData)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error decoding response: "+err.Error()+" | Server HF Response: "+response.String())
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
		// Path relatif ke dataset
		datasetPath := ("../dataset/barokah.csv")
		// datasetPath := config.GetEnv("DATASET_PATH")
		// datasetPath := "./dataset/questions.csv"

		// Log path dataset untuk debugging
        log.Printf("Dataset path : %s", datasetPath)
		
		// Load the dataset
		labelToQA, err := LoadDataset(datasetPath)
		if err != nil {
			helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "server error: could not load dataset: "+err.Error())
			return
		}

		// Get the answer corresponding to the best label
		record, ok := labelToQA[bestLabel]
		if !ok {
			helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "server error: label not found in dataset")
			return
		}

		answer := record[1]

		helper.WriteJSON(respw, http.StatusOK, map[string]string{
			"prompt":   chat.Prompt,
			"response": answer,
			"label":    bestLabel,
			"score":    strconv.FormatFloat(highestScore, 'f', -1, 64),
		})
	} else {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
	}
}