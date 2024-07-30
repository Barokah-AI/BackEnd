package controller

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"

	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
)

func LogIn(db *mongo.Database, respwrt http.ResponseWriter, request *http.Request, privatekey string) {
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)

	// error handling
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}
}
