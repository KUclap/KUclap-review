package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Review struct {
	ID			bson.ObjectId	`json:"reviewId" bson:"_id"`
	Review		string			`json:"text" bson:"text"`
	Author		string			`json:"author" bson:"author"`
	Grade		string			`json:"grade" bson:"grade"`
	Auth		string			`json:"auth" bson:"auth"`
	CreatedAt	time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt	time.Time		`json:"updateAt" bson:"update_at"`
}
