// helper/python_helper.go
package helper

// import (
// 	"encoding/json"
// 	"os/exec"

//     model "github.com/Barokah-AI/BackEnd/model"
// )

// func CallPythonScript(message string) (string, error) {
//     cmd := exec.Command("python", "bert-model/chat_model.py", message)
//     output, err := cmd.Output()
//     if err != nil {
//         return "", err
//     }

//     var chatResp model.ChatResponse
//     if err := json.Unmarshal(output, &chatResp); err != nil {
//         return "", err
//     }

//     return chatResp.Response, nil
// }
