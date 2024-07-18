package controller

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"net/http"

// 	"github.com/badoux/checkmail"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"golang.org/x/crypto/argon2"

// 	"github.com/Barokah-AI/BackEnd/helper"
// 	"github.com/Barokah-AI/BackEnd/model"
// )

// user
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
