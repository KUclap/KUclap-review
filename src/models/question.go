package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Question is model for create question
type Question struct {
	ID				bson.ObjectId	`json:"questionId" bson:"question_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	Question		string			`json:"question" bson:"question"`
	Answer			[]Answer		`json:"answer" bson:"author"`
	AuthorQuestion	string			`json:"authorQuestion,omitempty" bson:"author_question"`
	Auth			string			`json:"auth,omitempty" bson:"auth"`
	ClassNameTH		string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN		string			`json:"classNameEN" bson:"class_name_en"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt		time.Time		`json:"updateAt" bson:"update_at"`
}

// Answer is model for answer on question
type Answer struct {
	AnswerID		bson.ObjectId	`json:"answerId" bson:"answer_id"`
	QuestionID		string			`json:"questionId" bson:"question_id"`
	Answer			string			`json:"answer" bson:"author"`
	AuthorAnswer	string			`json:"authorAnswer,omitempty" bson:"author_answer"`
	Auth			string			`json:"auth,omitempty" bson:"auth"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt		time.Time		`json:"updateAt" bson:"update_at"`
}

// ResQuestion is struct for body response to client
type ResQuestion struct {
	ID				bson.ObjectId	`json:"questionId" bson:"question_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	Question		string			`json:"question" bson:"question"`
	Answer			[]Answer		`json:"answer" bson:"author"`
	AuthorQuestion	string			`json:"authorQuestion,omitempty" bson:"author_question"`
	ClassNameTH		string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN		string			`json:"classNameEN" bson:"class_name_en"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt		time.Time		`json:"updateAt" bson:"update_at"`
}

// ResAnswer is struct for body response to client
type ResAnswer struct {
	AnswerID		bson.ObjectId	`json:"answerId" bson:"answer_id"`
	QuestionID		string			`json:"questionId" bson:"question_id"`
	Answer			string			`json:"answer" bson:"author"`
	AuthorAnswer	string			`json:"authorAnswer,omitempty" bson:"author_answer"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt		time.Time		`json:"updateAt" bson:"update_at"`
}

// // GetBSON is function for filling default value when save value on mongo
// func (quest *Question) GetBSON() (interface{}, error) {
//     if quest.Answer == "" {
// 		quest.Answer = "ยังไม่มีข้อมูล" 
//     }
//     type my *Question
//     return my(quest), nil
// }