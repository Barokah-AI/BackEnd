package model

import (
	"time"
)

// Struct untuk membaca request dari user
type AIRequest struct {
	Prompt   	    string             `bson:"prompt,omitempty" json:"prompt,omitempty"`
	AIResp          string             `bson:"airesp,omitempty" json:"airesp,omitempty"`
	CreatedAt       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// Struct untuk membaca request dari Hugging Face API
type HFRequest struct {
    Inputs string `json:"inputs"`
}

type HFResponse struct {
    Label string  `json:"label"`
    Score float64 `json:"score"`
}

// Struct untuk membaca tokenizer_config.json
type TokenizerConfig struct {
	DoLowerCase   bool   `json:"do_lower_case"`
	ClsToken      string `json:"cls_token"`
	PadToken      string `json:"pad_token"`
	SepToken      string `json:"sep_token"`
	MaskToken     string `json:"mask_token"`
	UnkToken      string `json:"unk_token"`
}

type Credential struct {
	Status  int    `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}