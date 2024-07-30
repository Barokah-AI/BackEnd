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
