package models

type Class struct {
	ClassId		string			`json:"classId"`
	NameTh		string			`json:"nameTh"`
	NameEn		string			`json:"nameEn"`
	Label		string			`json:"label"`
	
}


type OldClass struct {
	Value	string	`json:"value"`
	Label	string	`json:"label"`
}