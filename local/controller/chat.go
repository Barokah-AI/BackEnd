package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
)

func callHuggingFaceAPI(prompt string) (string, float64, error) {
	api_url := config.GetEnv("HUGGINGFACE_API_URL")
	api_token := "Bearer " + config.GetEnv("HUGGINGFACE_API_KEY")

	req_body := model.HFRequest{Inputs: prompt}
	json_req_body, err := json.Marshal(req_body)
	if err != nil {
		return "", 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", api_url, bytes.NewBuffer(json_req_body))
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", api_token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error making request to Hugging Face API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("unexpected status code from Hugging Face API: %d | Server HF Response: %s", resp.StatusCode, string(body_bytes))
	}

	// Read and print the response body
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading response body: %v", err)
	}
	response_body := string(body_bytes)
	fmt.Println("HF API Response:", response_body) // Print the raw response

	// Handle the expected nested array structure
	var nested_data [][]map[string]interface{}
	err = json.Unmarshal(body_bytes, &nested_data)
	if err != nil {
		return "", 0, fmt.Errorf("error decoding response: %v | Server HF Response: %s", err, response_body)
	}

	// Flatten the nested array structure
	var flat_data []map[string]interface{}
	for _, d := range nested_data {
		flat_data = append(flat_data, d...)
	}

	// Check if the flat data has at least one element
	if len(flat_data) == 0 {
		return "", 0, fmt.Errorf("empty response after flattening nested structure: %s", response_body)
	}

	// Assume the first element contains the best response
	best_response := flat_data[0]

	// Extract label and score from the best response
	label, ok := best_response["label"].(string)
	if !ok {
		return "", 0, fmt.Errorf("missing or invalid label in response: %s", response_body)
	}
	score, ok := best_response["score"].(float64)
	if !ok {
		// Handle the case where the score might be an integer
		if score_integer, ok := best_response["score"].(int); ok {
			score = float64(score_integer)
		} else {
			return "", 0, fmt.Errorf("missing or invalid score in response: %s", response_body)
		}
	}

	return label, score, nil
}

func Chat(respwrt http.ResponseWriter, request *http.Request, token_model string) {
	var chat model.AIRequest

	err := json.NewDecoder(request.Body).Decode(&chat)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Permintaan Tidak Valid", "error saat membaca isi permintaan: "+err.Error())
		return
	}

	if chat.Prompt == "" {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Permintaan Tidak Valid", "masukin pertanyaan dulu ya kakak 🤗")
		return
	}

	// Read and use the tokenizer
	vocab, err := readVocab("../helper/vocab.txt")
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca vocab: "+err.Error())
		return
	}

	tokenizerConfig, err := readTokenizerConfig("../helper/tokenizer_config.json")
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca konfigurasi tokenizer: "+err.Error())
		return
	}

	tokens, err := tokenize(chat.Prompt, vocab, tokenizerConfig)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Kesalahan Server Internal", "error saat melakukan tokenisasi: "+err.Error())
		return
	}

	// Convert tokens to string for API call
	tokensStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(tokens)), " "), "[]")

	// Call Hugging Face API with tokenized prompt
	label, score, err := callHuggingFaceAPI(tokensStr)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Kesalahan Server Internal", "model sedang diload: "+err.Error())
		return
	}

	// Load the dataset
	dataset_path := ("../dataset/barokah.csv")
	label_to_qa, err := helper.LoadDatasetLocal(dataset_path)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Kesalahan Server Internal", "kesalahan server: tidak bisa memuat dataset: "+err.Error())
		return
	}

	record, ok := label_to_qa[label]
	if !ok {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Kesalahan Server Internal", "kesalahan server: label tidak ditemukan dalam dataset")
		return
	}

	answer := record[1]

	helper.WriteJSON(respwrt, http.StatusOK, map[string]string{
		"prompt":   chat.Prompt,
		"response": answer,
		"label":    label,
		"score":    strconv.FormatFloat(score, 'f', -1, 64),
	})
}
