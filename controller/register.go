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

	"github.com/Barokah-AI/BackEnd/model"
	"github.com/Barokah-AI/BackEnd/helper"
)

func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	// error handling
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	// check if user data is empty
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data ya kak")
		return
	}

	// check if email is valid
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "sprry email tidak valid kak")
		return
	}

	// check if email already exists
	userExists, _ := helper.GetUserFromEmail(user.Email, db)
	if user.Email == userExists.Email {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "waduh email sudah terdaftar kak")
		return
	}

	// check if password and confirm password match
	if strings.Contains(user.Password, " ") {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi kak")
		return
	}

	// check if password is at least 8 characters
	// if len(user.Password) < 8 {
	// 	helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter ya kak")
	// 	return
	// }

	// check if password and confirm password match
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
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
	insertedID, err := helper.InsertOneDoc(db, col, userData)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : insert data, "+err.Error())
		return
	}

	// response
	resp := map[string]any{
		"message":    "berhasil mendaftar kak",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	helper.WriteJSON(respw, http.StatusCreated, resp)
}
