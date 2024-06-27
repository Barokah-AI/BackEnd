// // helper/hugface_helper.go
package helper

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/go-resty/resty/v2"
// )

// type HuggingFaceRequest struct {
//     Inputs string `json:"inputs"`
// }

// type HuggingFaceResponse struct {
//     GeneratedText string `json:"generated_text"`
// }

// func CallHuggingFaceAPI(message string) (string, error) {
//     client := resty.New()

//     // Sesuaikan URL dan model API Hugging Face Anda
//     url := "https://api-inference.huggingface.co/models/dimasardnt/barokah-model"
//     apiKey := "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" // Gantilah dengan API key Anda

//     reqBody := HuggingFaceRequest{
//         Inputs: message,
//     }

//     resp, err := client.R().
//         SetHeader("Authorization", "Bearer "+apiKey).
//         SetHeader("Content-Type", "application/json").
//         SetBody(reqBody).
//         Post(url)

//     if err != nil {
//         return "", err
//     }

//     if resp.StatusCode() != http.StatusOK {
//         return "", fmt.Errorf("failed to call Hugging Face API: %s", resp.Status())
//     }

//     var hfResponse []HuggingFaceResponse
//     if err := json.Unmarshal(resp.Body(), &hfResponse); err != nil {
//         return "", err
//     }

//     if len(hfResponse) == 0 {
//         return "", fmt.Errorf("empty response from Hugging Face API")
//     }

//     return hfResponse[0].GeneratedText, nil
// }
