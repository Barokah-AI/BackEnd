package controller

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"net/http"

// 	"github.com/badoux/checkmail"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"golang.org/x/crypto/argon2"
// )

// // user
// func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
// 	var user model.User
// 	err := json.NewDecoder(req.Body).Decode(&user)
// 	if err != nil {
// 		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body " + err.Error())
// 		return
// 	}
// 	if user.Email == "" || user.Password == "" {
// 		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
// 		return
// 	}
// 	if err = checkmail.ValidateFormat(user.Email); err != nil {
// 		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
// 		return
// 	}
// 	existsDoc, err := helper.GetUserFromEmail(user.Email, db)
// 	if err != nil {
// 		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : get email " + err.Error())
// 		return
// 	}
// 	salt, err := hex.DecodeString(existsDoc.Salt)
// 	if err != nil {
// 		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
// 		return
// 	}
// 	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
// 	if hex.EncodeToString(hash) != existsDoc.Password {
// 		helper.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "password salah")
// 		return
// 	}
// 	tokenstring, err := helper.Encode(user.ID, user.Email, privatekey)
// 	if err != nil {
// 		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : token")
// 		return
// 	}
// 	resp := map[string]any{
// 		"status":  "success",
// 		"message": "login berhasil",
// 		"token":   tokenstring,
// 		"data" : map[string]string{
// 			"email": existsDoc.Email,
// 			"namalengkap": existsDoc.NamaLengkap,
// 		},
// 	}
// 	helper.WriteJSON(respw, http.StatusOK, resp)
// }
