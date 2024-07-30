package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorResponse(respwd http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respwd, statusCode, resp)
}

func WriteJSON(respwd http.ResponseWriter, statusCode int, content any) {
	respwd.Header().Set("Content-Type", "application/json")
	respwd.WriteHeader(statusCode)
	respwd.Write([]byte(Jsonstr(content)))
}

func Jsonstr(strc any) string {
	json_data, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(json_data)
}
