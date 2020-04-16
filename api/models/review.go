package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Review struct {
	ID			bson.ObjectId	`json:"reviewId" bson:"_id"`
	ClassId		string			`json:"classId" bson:"class_id"`
	Text		string			`json:"text" bson:"text"`
	Author		string			`json:"author" bson:"author"`
	Grade		string			`json:"grade" bson:"grade"`
	Auth		string			`json:"auth" bson:"auth"`
	Clap		uint64			`json:"clap" bson:"clap"`
	Boo			uint64			`json:"boo" bson:"boo"`
	CreatedAt	time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt	time.Time		`json:"updateAt" bson:"update_at"`
	Reported	bool			`json:"reported" bson:"reported"`
}
