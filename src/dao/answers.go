package dao

import (

	"time"

	"kuclap-review-api/src/models"

	"gopkg.in/mgo.v2/bson"
)

// CreateAnswer is POST create answer.
func (m *SessionDAO) CreateAnswer(answer models.Answer) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).Insert(&answer)
	
	return err
}

// UpdateNumberAnswerByClassID is Update number of review
func (m *SessionDAO) UpdateNumberAnswerByQuestionID(questionId string, number int, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Update(bson.M{"_id": bson.ObjectIdHex(questionId)}, bson.M{ "$inc": bson.M{"number_answers": int32(number)}, "$set": bson.M{"update_at": updateAt} })
	
	return err
}

// FindAllAnswers is Find All of list of answer 
func (m *SessionDAO) FindAllAnswers() ([]models.Answer, error) {
	var answers []models.Answer

	db	:=	session.Copy()
	defer db.Close()
	
	
	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).Find(bson.M{}).All(&answers)
	return answers, err
}

// FindLastAnswersByQuestionID is Find All of list of answer 
func (m *SessionDAO) FindLastAnswersByQuestionID(questionId string) ([]models.ResAnswer, error) {
	var answers []models.ResAnswer

	db	:=	session.Copy()
	defer db.Close()
	
	
	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).Find(bson.M{"question_id": questionId}).Sort("-$natural").All(&answers)
	return answers, err
}

// FindAnswerByID is Find answer by answer_id
func (m *SessionDAO) FindAnswerByID(answerID string) (models.ResAnswer, error) {
	var answer models.ResAnswer

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).Find(bson.M{"_id": bson.ObjectIdHex(answerID)}).One(&answer)
	
	return answer, err
}

// FindAnswerAllPropertyByID is Find answer by _id
func (m *SessionDAO) FindAnswerAllPropertyByID(answerID string) (models.Answer, error) {
	var answer models.Answer

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).Find(bson.M{"_id": bson.ObjectIdHex(answerID)}).One(&answer)
	
	return answer, err
}

// DeleteAnswerByID is Delete an existing review
func (m *SessionDAO) DeleteAnswerByID(answerID string) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).Remove(bson.M{"_id": bson.ObjectIdHex(answerID)})
	
	return err
}

// DeleteAnswersByQuestionID is Delete an existing review
func (m *SessionDAO) DeleteAnswersByQuestionID(questionId string) error {
	db	:=	session.Copy()
	defer db.Close()

	_, err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).RemoveAll(bson.M{"question_id": questionId})
	
	return err
}

// UpdateAnswerReportByID is Update reported
func (m *SessionDAO) UpdateAnswerReportByID(id string, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_ANSWERS).UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"reported": true, "update_at": updateAt}})
	
	return err
}