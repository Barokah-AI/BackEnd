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
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)

	// error handling
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body " + err.Error())
		return
	}

	// check if email and password is empty
	// if user.Email == "" || user.Password == "" {
	// 	helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
	// 	return
	// }

	// check if email is valid
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// check if email exists
	existsDoc, err := helper.GetUserFromEmail(user.Email, db)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : get email " + err.Error())
		return
	}

	// check if password is correct
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}

	// compare password
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		helper.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "password salah")
		return
	}

	// generate token
	tokenstring, err := helper.Encode(user.ID, user.Email, privatekey)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : token")
		return
	}

	// response
	resp := map[string]any{
		"status":  "success",
		"message": "login berhasil",
		"token":   tokenstring,
		"data" : map[string]string{
			"email": existsDoc.Email,
			"namalengkap": existsDoc.NamaLengkap,
		},
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}
