package helper

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	return token.V4Sign(secretKey, nil), err
}

func Decode(publicKey string, tokenstring string) (payload Payload, err error) {
	var token *paseto.Token
	var pubKey paseto.V4AsymmetricPublicKey
	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publicKey) 
	if err != nil {
		return payload, fmt.Errorf("Decode NewV4AsymmetricPublicKeyFromHex : %v", err)
	}
	parser := paseto.NewParser()                                
	token, err = parser.ParseV4Public(pubKey, tokenstring, nil) 
	if err != nil {
		return payload, fmt.Errorf("Decode ParseV4Public : %v", err)
	}
	err = json.Unmarshal(token.ClaimsJSON(), &payload)
	return payload, err
}

func GenerateKey() (privateKey, publicKey string) {
	secretKey := paseto.NewV4AsymmetricSecretKey() 
	publicKey = secretKey.Public().ExportHex()    
	privateKey = secretKey.ExportHex()
	return privateKey, publicKey
}
