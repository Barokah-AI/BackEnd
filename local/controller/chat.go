package controller

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	io "io/ioutil"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/Barokah-AI/BackEnd/config"
// 	"github.com/Barokah-AI/BackEnd/helper"
// 	"github.com/Barokah-AI/BackEnd/model"
// )

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

