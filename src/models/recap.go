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
	DeleteReason	string			`json:"deleteReason" bson:"delete_reason"`
	TypeFile		string			`json:"type" bson:"type"`
	Tag				string			`json:"tag" bson:"tag"`
	ClassNameTH		string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN		string			`json:"classNameEN" bson:"class_name_en"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	Reported		bool			`json:"reported" bson:"reported"`
}

type ResRecap struct {
	RecapID			bson.ObjectId	`json:"recapId" bson:"_id"`
	ReviewID		string			`json:"reviewId" bson:"review_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	FileName		string			`json:"filename" bson:"filename"`
	Donwloaded		uint64			`json:"downloaded" bson:"downloaded"`
	Description		string			`json:"description" bson:"description"`
	Author			string			`json:"author" bson:"author"`
	DeleteReason	string			`json:"deleteReason" bson:"delete_reason"`
	TypeFile		string			`json:"type" bson:"type"`
	Tag				string			`json:"tag" bson:"tag"`
	ClassNameTH		string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN		string			`json:"classNameEN" bson:"class_name_en"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	Reported		bool			`json:"reported" bson:"reported"`
}