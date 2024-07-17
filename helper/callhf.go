package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/model"
)

func CallHuggingFaceAPI(prompt string) (string, float64, error) {
    apiUrl := config.GetEnv("HUGGINGFACE_API_URL")
    apiToken := "Bearer " + config.GetEnv("HUGGINGFACE_API_KEY")

    reqBody := model.HFRequest{Inputs: prompt}
    jsonReqBody, err := json.Marshal(reqBody)
    if err != nil {
        return "", 0, fmt.Errorf("error marshalling in request body: %v", err)
    }

    req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonReqBody))
    if err != nil {
        return "", 0, fmt.Errorf("error creating in request: %v", err)
    }
    req.Header.Set("Authorization", apiToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", 0, fmt.Errorf("error making request in to Hugging Face API: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        return "", 0, fmt.Errorf("unexpected status code from in Hugging Face API: %d | Server HF Response: %s", resp.StatusCode, string(bodyBytes))
    }

    var hfResponse []model.HFResponse
    err = json.NewDecoder(resp.Body).Decode(&hfResponse)
    if err != nil {
        return "", 0, fmt.Errorf("error decoding in response: %v", err)
    }

    if len(hfResponse) == 0 {
        return "", 0, fmt.Errorf("empty response in from Hugging Face API")
    }

    bestResponse := hfResponse[0]
    return bestResponse.Label, bestResponse.Score, nil
}
