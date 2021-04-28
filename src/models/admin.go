package models

type AdminDeleteByID struct {
	ReviewID		string			`json:"reviewId" bson:"review_id"`
	QuestionID		string			`json:"questionId" bson:"question_id"`
	RecapID			string			`json:"recapId" bson:"recap_id"`
	AnswerID		string			`json:"answerId" bson:"answer_id"`

	DeleteReason	string			`json:"deleteReason" bson:"delete_reason"`
}