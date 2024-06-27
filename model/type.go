package model

import (
	"time"
)

// type ChatRequest struct {
// 	Message string `json:"message"`
// }

// type ChatResponse struct {
// 	Response string `json:"response"`
// }

// type HuggingFaceRequest struct {
// 	Inputs string `json:"inputs"`
// }

// type HuggingFaceResponse struct {
// 	GeneratedText string `json:"generated_text"`
// }

type AIRequest struct {
	Promt   	    string             `bson:"promt,omitempty" json:"promt,omitempty"`
	CreatedAt       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}