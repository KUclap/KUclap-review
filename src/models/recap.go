package models

import (

	"time"
	"gopkg.in/mgo.v2/bson"

)

type Recap struct {
	RecapID			bson.ObjectId	`json:"recapId" bson:"_id"`
	ReviewID		string			`json:"reviewId" bson:"review_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	FileName		string			`json:"filename" bson:"filename"`
	Donwloaded		uint64			`json:"downloaded" bson:"downloaded"`
	Description		string			`json:"description" bson:"description"`
	Author			string			`json:"author" bson:"author"`
	Auth			string			`json:"auth" bson:"auth"`
	TypeFile		string			`json:"type" bson:"type"`
	Tag				string			`json:"tag" bson:"tag"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
}