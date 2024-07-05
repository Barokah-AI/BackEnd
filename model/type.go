package model

import (
	"time"
)


type AIRequest struct {
	Prompt   	    string             `bson:"prompt,omitempty" json:"prompt,omitempty"`
	AIResp          string             `bson:"airesp,omitempty" json:"airesp,omitempty"`
	CreatedAt       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}