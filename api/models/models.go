package models

import (
	"time"
)

type Review struct {
	// ID			primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text		string		`json:"text" bson:"text"`
	// Author		string		`json:"author" bson:"author"`
	// Grade		string		`json:"grade" bson:"grade"`
	// Auth		string		`json:"auth" bson:"auth"`
	CreatedAt	time.Time	`json:"createdAt" bson:"created_at"`
}
