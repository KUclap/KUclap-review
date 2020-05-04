package models

import "time"


type Report struct {
	ReviewID		string			`json:"reviewId" bson:"review_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	Text			string			`json:"text" bson:"text"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
}