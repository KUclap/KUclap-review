package models

import "time"

type Report struct {
	ReviewID   string    `json:"reviewId" bson:"review_id"`
	QuestionID string    `json:"questionId" bson:"question_id"`
	RecapID    string    `json:"recapId" bson:"recap_id"`
	AnswerID   string    `json:"answerId" bson:"answer_id"`
	ClassID    string    `json:"classId" bson:"class_id"`
	Text       string    `json:"text" bson:"text"`
	CreatedAt  time.Time `json:"createdAt" bson:"created_at"`
}
