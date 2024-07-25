package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, error, massage string) {
	resp := map[string]string{
		"error":   error,
		"message": massage,
	}
	WriteJSON(respw, statusCode, resp)
}

func WriteJSON(respw http.ResponseWriter, statusCode int, content any) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

func Jsonstr(strc any) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}
