package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question is model for create question
type Question struct {
	ID           primitive.ObjectID `json:"questionId" bson:"_id"`
	ClassID      string             `json:"classId" bson:"class_id"`
	Question     string             `json:"question" bson:"question"`
	Author       string             `json:"author,omitempty" bson:"author"`
	Auth         string             `json:"auth,omitempty" bson:"auth"`
	DeleteReason string             `json:"deleteReason" bson:"delete_reason"`
	ClassNameTH  string             `json:"classNameTH" bson:"class_name_th"`
	ClassNameEN  string             `json:"classNameEN" bson:"class_name_en"`
	NumberAnswer uint64             `json:"numberAnswer" bson:"number_answers"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdateAt     time.Time          `json:"updateAt" bson:"update_at"`
	Reported     bool               `json:"reported" bson:"reported"`
}

// ResQuestion is struct for body response to client
type ResQuestion struct {
	ID           primitive.ObjectID `json:"questionId" bson:"_id"`
	ClassID      string             `json:"classId" bson:"class_id"`
	Question     string             `json:"question" bson:"question"`
	Author       string             `json:"author,omitempty" bson:"author"`
	DeleteReason string             `json:"deleteReason" bson:"delete_reason"`
	ClassNameTH  string             `json:"classNameTH" bson:"class_name_th"`
	ClassNameEN  string             `json:"classNameEN" bson:"class_name_en"`
	NumberAnswer uint64             `json:"numberAnswer" bson:"number_answers"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdateAt     time.Time          `json:"updateAt" bson:"update_at"`
	Reported     bool               `json:"reported" bson:"reported"`
}

// // GetBSON is function for filling default value when save value on mongo
// func (quest *Question) GetBSON() (interface{}, error) {
//     if quest.Answer == "" {
// 		quest.Answer = "ยังไม่มีข้อมูล"
//     }
//     type my *Question
//     return my(quest), nil
// }
