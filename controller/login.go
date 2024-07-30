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
}
