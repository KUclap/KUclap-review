package models

import (
	"time"
)

type Review struct {
	// ID			primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Review		string		`json:"review" bson:"review"`
	Author		string		`json:"author" bson:"author"`
	Grade		string		`json:"grade" bson:"grade"`
	Auth		string		`json:"auth" bson:"auth"`
	CreatedAt	time.Time	`json:"createdAt" bson:"created_at"`
}
