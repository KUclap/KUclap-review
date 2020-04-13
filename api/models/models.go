package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Review struct {
	ReviewID	bson.ObjectId	`json:"review_id" bson:"_id"`
	Text		string			`json:"text" bson:"text"`
	Author		string			`json:"author" bson:"author"`
	Grade		string			`json:"grade" bson:"grade"`
	Auth		string			`json:"auth" bson:"auth"`
	CreatedAt	time.Time		`json:"createdAt" bson:"created_at"`
}
