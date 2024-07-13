package controller

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
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

func callHuggingFaceAPI(prompt string) (string, float64, error) {
	apiUrl := config.GetEnv("HUGGINGFACE_API_URL")
	apiToken := "Bearer " + config.GetEnv("HUGGINGFACE_API_KEY")

	reqBody := model.HFRequest{Inputs: prompt}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error making request to Hugging Face API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("unexpected status code from Hugging Face API: %d | Server HF Response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read and print the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading response body: %v", err)
	}
	responseBody := string(bodyBytes)
	fmt.Println("HF API Response:", responseBody) // Print the raw response

	// Handle the expected nested array structure
	var nestedData [][]map[string]interface{}
	err = json.Unmarshal(bodyBytes, &nestedData)
	if err != nil {
		return "", 0, fmt.Errorf("error decoding response: %v | Server HF Response: %s", err, responseBody)
	}

	// Flatten the nested array structure
	var flatData []map[string]interface{}
	for _, d := range nestedData {
		flatData = append(flatData, d...)
	}

	// Check if the flat data has at least one element
	if len(flatData) == 0 {
		return "", 0, fmt.Errorf("empty response after flattening nested structure: %s", responseBody)
	}

	// Assume the first element contains the best response
	bestResponse := flatData[0]

	// Extract label and score from the best response
	label, ok := bestResponse["label"].(string)
	if !ok {
		return "", 0, fmt.Errorf("missing or invalid label in response: %s", responseBody)
	}
	score, ok := bestResponse["score"].(float64)
	if !ok {
		// Handle the case where the score might be an integer
		if scoreInt, ok := bestResponse["score"].(int); ok {
			score = float64(scoreInt)
		} else {
			return "", 0, fmt.Errorf("missing or invalid score in response: %s", responseBody)
		}
	}

	return label, score, nil
}

func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
	var chat model.AIRequest

	err := json.NewDecoder(req.Body).Decode(&chat)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Permintaan Tidak Valid", "error saat membaca isi permintaan: "+err.Error())
		return
	}

	if chat.Prompt == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Permintaan Tidak Valid", "masukin pertanyaan dulu ya kakak 🤗")
		return
	}

	// Read and use the tokenizer
	vocab, err := readVocab("../helper/vocab.txt")
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca vocab: "+err.Error())
		return
	}

	tokenizerConfig, err := readTokenizerConfig("../helper/tokenizer_config.json")
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca konfigurasi tokenizer: "+err.Error())
		return
	}

	tokens, err := tokenize(chat.Prompt, vocab, tokenizerConfig)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "error saat melakukan tokenisasi: "+err.Error())
		return
	}

	// Convert tokens to string for API call
	tokensStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(tokens)), " "), "[]")

	// Call Hugging Face API with tokenized prompt
	label, score, err := callHuggingFaceAPI(tokensStr)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "model sedang diload: "+err.Error())
		return
	}

	// Load the dataset
	datasetPath := ("../dataset/barokah.csv")
	labelToQA, err := LoadDataset(datasetPath)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "kesalahan server: tidak bisa memuat dataset: "+err.Error())
		return
	}

	record, ok := labelToQA[label]
	if !ok {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "kesalahan server: label tidak ditemukan dalam dataset")
		return
	}

	answer := record[1]

	helper.WriteJSON(respw, http.StatusOK, map[string]string{
		"prompt":   chat.Prompt,
		"response": answer,
		"label":    label,
		"score":    strconv.FormatFloat(score, 'f', -1, 64),
	})
}
