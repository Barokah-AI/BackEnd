package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struct User
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap     string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	Confirmpassword string             `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
	Salt            string             `bson:"salt,omitempty" json:"salt,omitempty"`
}

// Struct Password
type Password struct {
	Password        string `bson:"password,omitempty" json:"password,omitempty"`
	Newpassword     string `bson:"newpass,omitempty" json:"newpass,omitempty"`
	Confirmpassword string `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
}

// Struct for AI Request
type AIRequest struct {
	Prompt    string    `bson:"prompt,omitempty" json:"prompt,omitempty"`
	AIResp    string    `bson:"airesp,omitempty" json:"airesp,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// Struct for AI Response from Hugging Face API
type HFRequest struct {
	Inputs string `json:"inputs"`
}

// Struct for AI Response from Hugging Face API
type HFResponse struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

// Struct untuk membaca tokenizer_config.json
type TokenizerConfig struct {
	DoLowerCase bool   `json:"do_lower_case"`
	ClsToken    string `json:"cls_token"`
	PadToken    string `json:"pad_token"`
	SepToken    string `json:"sep_token"`
	MaskToken   string `json:"mask_token"`
	UnkToken    string `json:"unk_token"`
}

// Struct for Credential
type Credential struct {
	Status  int    `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// Struct for Response
type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// Struct for Payload
type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}
