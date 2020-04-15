package models

type Class struct {
	ClassId		string			`json:"classId" bson:"class_id"`
	NameTh		string			`json:"nameTh" bson:"name_th"`
	NameEn		string			`json:"nameEn" bson:"name_en"`
	Label		string			`json:"label" bson:"label"`
	Stats		StatClass		`json:"stats" bson:"stats"`	
}

type StatClass struct {
	How			float32		`json:"how" bson:"how"`
	Homework	float32		`json:"homework" bson:"homework"`
	Interest	float32		`json"interest" bson:"interest"`
}


type OldClass struct {
	Value	string	`json:"value"`
	Label	string	`json:"label"`
}