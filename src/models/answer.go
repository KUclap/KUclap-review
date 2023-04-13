package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Answer is model for answer on question
type Answer struct {
	ID           primitive.ObjectID `json:"answerId" bson:"_id"`
	QuestionID   string             `json:"questionId" bson:"question_id"`
	Answer       string             `json:"answer" bson:"answer"`
	Author       string             `json:"author,omitempty" bson:"author"`
	Auth         string             `json:"auth,omitempty" bson:"auth"`
	DeleteReason string             `json:"deleteReason" bson:"delete_reason"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdateAt     time.Time          `json:"updateAt" bson:"update_at"`
	Reported     bool               `json:"reported" bson:"reported"`
}

// ResAnswer is struct for body response to client
type ResAnswer struct {
	ID           primitive.ObjectID `json:"answerId" bson:"_id"`
	QuestionID   string             `json:"questionId" bson:"question_id"`
	Answer       string             `json:"answer" bson:"answer"`
	Author       string             `json:"author,omitempty" bson:"author"`
	DeleteReason string             `json:"deleteReason" bson:"delete_reason"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdateAt     time.Time          `json:"updateAt" bson:"update_at"`
	Reported     bool               `json:"reported" bson:"reported"`
}
