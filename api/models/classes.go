package models
import "time"

type Class struct {
	ClassID			string			`json:"classId" bson:"class_id"`
	NameTh			string			`json:"nameTh" bson:"name_th"`
	NameEn			string			`json:"nameEn" bson:"name_en"`
	Label			string			`json:"label" bson:"label"`
	Hours			string			`json:"hours" bson:"hours"`
	Unit			uint64			`json:"unit" bson:"unit"`
	NumberReviewer	float64			`json:"numberReviewer" bson:"number_reviewer"`
	Stats			StatClass		`json:"stats" bson:"stats"`	
}

type StatClass struct {
	How			float64		`json:"how" bson:"how"`
	Homework	float64		`json:"homework" bson:"homework"`
	Interest	float64		`json:"interest" bson:"interest"`
	UpdateAt	time.Time	`json:"updateAt" bson:"update_at"`
}


type OldClass struct {
	Value	string	`json:"value"`
	Label	string	`json:"label"`
}