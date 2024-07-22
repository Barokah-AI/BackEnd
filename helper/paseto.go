package helper

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