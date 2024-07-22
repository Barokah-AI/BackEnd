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
