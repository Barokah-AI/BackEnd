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

