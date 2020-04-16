package models
import "time"

type Class struct {
	ClassId			string			`json:"classId" bson:"class_id"`
	NameTh			string			`json:"nameTh" bson:"name_th"`
	NameEn			string			`json:"nameEn" bson:"name_en"`
	Label			string			`json:"label" bson:"label"`
	Hours			string			`json:"hours" bson:"hours"`
	Unit			uint			`json:"unit" bson:"unit"`
	NumberReviewer	uint			`json:"numberReviewer" bson:"number_reviewer"`
	Stats			StatClass		`json:"stats" bson:"stats"`	
}

type StatClass struct {
	How			float32		`json:"how" bson:"how"`
	Homework	float32		`json:"homework" bson:"homework"`
	Interest	float32		`json:"interest" bson:"interest"`
	UpdateAt	time.Time	`json:"updateAt" bson:"update_at"`
}


type OldClass struct {
	Value	string	`json:"value"`
	Label	string	`json:"label"`
}