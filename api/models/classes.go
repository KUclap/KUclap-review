package models

type Classes struct {
	ClassID		string			`json:"classID"`
	NameTH		string			`json:"nameTH"`
	NameEN		string			`json:"nameEN"`
	Label		string			`json:"label"`
	
}


type OldClasses struct {
	Value	string	`json:"value"`
	Label	string	`json:"label"`
}