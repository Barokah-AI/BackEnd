package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"

	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
)

func SignUp(database *mongo.Database, col string, respwrt http.ResponseWriter, request *http.Request) {
	var user model.User

	// error handling
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	// check if user data is empty
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// check if email is valid
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// check if email already exists
	userExists, _ := helper.GetUserFromEmail(user.Email, database)
	if user.Email == userExists.Email {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "email sudah terdaftar")
		return
	}

	// check if password and confirm password match
	if strings.Contains(user.Password, " ") {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi")
		return
	}

	// check if password is at least 8 characters
	if len(user.Password) < 8 {
		helper.ErrorResponse(respwrt, request, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter")
		return
	}

	// check if password and confirm password match
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}

	// hash password
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// insert user data to database
	userData := bson.M{
		"namalengkap": user.NamaLengkap,
		"email":       user.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"salt":        hex.EncodeToString(salt),
	}

	// check if user data is empty
	inserted_id, err := helper.InsertOneDoc(database, col, userData)
	if err != nil {
		helper.ErrorResponse(respwrt, request, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : insert data, "+err.Error())
		return
	}

	// response
	response := map[string]any{
		"message":    "berhasil mendaftar",
		"insertedID": inserted_id,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	helper.WriteJSON(respwrt, http.StatusCreated, response)
}
