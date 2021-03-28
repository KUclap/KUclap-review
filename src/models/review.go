package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Review is model for create review
type Review struct {
	ID			bson.ObjectId	`json:"reviewId" bson:"_id"`
	ClassID		string			`json:"classId" bson:"class_id"`
	Text		string			`json:"text" bson:"text"`
	Author		string			`json:"author" bson:"author"`
	Grade		string			`json:"grade" bson:"grade"`
	Auth		string			`json:"auth,omitempty" bson:"auth"`
	Clap		uint64			`json:"clap" bson:"clap"`
	Boo			uint64			`json:"boo" bson:"boo"`
	Stats		StatReview		`json:"stats" bson:"stats"`	
	ClassNameTH	string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN	string			`json:"classNameEN" bson:"class_name_en"`
	Sec			uint64			`json:"sec" bson:"sec"`
	Semester	uint64			`json:"semester" bson:"semester"`
	Year		uint64			`json:"year" bson:"year"`
	Recap		string			`json:"recap" bson:"recap"`
	CreatedAt	time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt	time.Time		`json:"updateAt" bson:"update_at"`
	Reported	bool			`json:"reported" bson:"reported"`
}

// StatReview is model for storing stat on the review
type StatReview struct {
	How			float64		`json:"how" bson:"how"`
	Homework	float64		`json:"homework" bson:"homework"`
	Interest	float64		`json:"interest" bson:"interest"`
}

// ResReview is struct for body response to client 
type ResReview struct {
	ID			bson.ObjectId	`json:"reviewId" bson:"_id"`
	ClassID		string			`json:"classId" bson:"class_id"`
	Text		string			`json:"text" bson:"text"`
	Author		string			`json:"author" bson:"author"`
	Grade		string			`json:"grade" bson:"grade"`
	Clap		uint64			`json:"clap" bson:"clap"`
	Boo			uint64			`json:"boo" bson:"boo"`
	Stats		StatReview		`json:"stats" bson:"stats"`	
	ClassNameTH	string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN	string			`json:"classNameEN" bson:"class_name_en"`
	Sec			uint64			`json:"sec" bson:"sec"`
	Semester	uint64			`json:"semester" bson:"semester"`
	Year		uint64			`json:"year" bson:"year"`
	Recap		string			`json:"recap" bson:"recap"`
	CreatedAt	time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt	time.Time		`json:"updateAt" bson:"update_at"`
	Reported	bool			`json:"reported" bson:"reported"`
}

// RDeleteReview is require struct for delete the review
type RDeleteReview struct {
	ID			string			`json:"reviewId" bson:"_id"`
	Auth		string			`json:"auth" bson:"auth"`
}

// SetBSON is function for filling default value when load value from mongo
func (review *ResReview) SetBSON(raw bson.Raw) (err error) {
	
	type my ResReview

	if err = raw.Unmarshal((*my)(review)); err != nil {
		return
	}
	review.Stats.How		=	review.Stats.How * 20
	review.Stats.Homework	=	review.Stats.Homework * 20
	review.Stats.Interest	=	review.Stats.Interest * 20

	return
	
}