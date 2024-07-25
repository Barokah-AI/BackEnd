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

// user
func LogIn(db *mongo.Database, respwrt http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)

	// error handling
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	// check if email and password is empty
	if user.Email == "" || user.Password == "" {
		helper.ErrorResponse(respwrt, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// check if email is valid
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// check if email exists
	exists_doc, err := helper.GetUserFromEmail(user.Email, db)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : get email "+err.Error())
		return
	}

	// check if password is correct
	salt, err := hex.DecodeString(exists_doc.Salt)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}

	// compare password
	hashpwd := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hashpwd) != exists_doc.Password {
		helper.ErrorResponse(respwrt, req, http.StatusUnauthorized, "Unauthorized", "password salah")
		return
	}

	// generate token
	token_string, err := helper.Encode(user.ID, user.Email, privatekey)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : token")
		return
	}

	// response
	resp := map[string]any{
		"status":  "success",
		"message": "login berhasil",
		"token":   token_string,
		"data": map[string]string{
			"email":       exists_doc.Email,
			"namalengkap": exists_doc.NamaLengkap,
		},
	}
	helper.WriteJSON(respwrt, http.StatusOK, resp)
}
