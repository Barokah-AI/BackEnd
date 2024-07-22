package helper

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"aidanwoods.dev/go-paseto"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}

func Encode(id primitive.ObjectID, email, privateKey string) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.Set("id", id)
	token.SetString("email", email)
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKey)
