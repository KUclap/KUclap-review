package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Review struct {
	ID			bson.ObjectId	`json:"reviewId" bson:"_id"`
	ClassId		string			`json:"classId" bson:"class_id"`
	Review		string			`json:"text" bson:"text"`
	Author		string			`json:"author" bson:"author"`
	Grade		string			`json:"grade" bson:"grade"`
	Auth		string			`json:"auth" bson:"auth"`
	Clap		uint			`json:"clap" bson:"clap"`
	Boo			uint			`json:"boo" bson:"boo"`
	Grades		Grade			`json:"grades" bson:"grades"`
	CreatedAt	time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt	time.Time		`json:"updateAt" bson:"update_at"`
	Reported	bool			`json:"reported" bson:"reported"`
}


type Grade struct {
	GradeA		uint	`json:"gradeA" bson:"grade_a"`
	GradeBplus	uint	`json:"gradeBPlus" bson:"grade_a_plus"`
	GradeB		uint	`json:"gradeB" bson:"grade_b"`
	GradeCplus	uint	`json:"gradeCPlus" bson:"grade_c_plus"`
	GradeC		uint	`json:"gradeC" bson:"grade_c"`
	GradeDplus	uint	`json:"gradeDPlus" bson:"grade_d_plus"`
	GradeD		uint	`json:"gradeD" bson:"grade_d"`
	GradeF		uint	`json:"gradeF" bson:"grade_f"`
}