package dao

import (
	"log"
	"time"
	"strconv"

	"kuclap-review-api/src/models"

	"gopkg.in/mgo.v2/bson"
)

// CreateQuestion is POST create question.
func (m *SessionDAO) CreateQuestion(question models.Question) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Insert(&question)
	
	return err
}

// UpdateNumberQuestionByClassID is Update number of review
func (m *SessionDAO) UpdateNumberQuestionByClassID(classID string, number int, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classID}, bson.M{ "$inc": bson.M{"number_questions": int32(number)}, "$set": bson.M{"update_at": updateAt} })
	
	return err
}

// FindAllQuestions is Find All of list of question 
func (m *SessionDAO) FindAllQuestions() ([]models.Question, error) {
	var questions []models.Question

	db	:=	session.Copy()
	defer db.Close()
	
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{}).All(&questions)
	return questions, err
}

// FindQuestionsByClassID is Find question by class_id
func (m *SessionDAO) FindQuestionsByClassID(classID string, page string, offset string) ([]models.ResQuestion, error) {
	var questions []models.ResQuestion

	db := session.Copy()
	defer db.Close()

	iPage, err		:=	strconv.Atoi(page)
	if err != nil {
		log.Println("err : atoi.", err)
	}

	iOffset, err	:=	strconv.Atoi(offset)
	if err != nil {
		log.Println("err : atoi.", err)
	}

	err = db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{"class_id": classID}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&questions)
	
	return questions, err
}

// LastQuestions is Find last questions range with offset
func (m *SessionDAO) LastQuestions(page string, offset string) ([]models.ResQuestion, error) {
	var questions []models.ResQuestion

	db	:=	session.Copy()
	defer db.Close()
	
	iPage, err		:=	strconv.Atoi(page)
	if err != nil {
		log.Println("err : atoi.", err)
	}

	iOffset, err	:=	strconv.Atoi(offset)
	if err != nil {
		log.Println("err : atoi.", err)
	}

	err	=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&questions)
	
	return questions, err
}

// FindQuestionByID is Find question by _id
func (m *SessionDAO) FindQuestionByID(questionID string) (models.ResQuestion, error) {
	var question models.ResQuestion

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{"_id": bson.ObjectIdHex(questionID) }).One(&question)
	
	return question, err
}

// FindQuestionAllPropertyByID is Find question by _id
func (m *SessionDAO) FindQuestionAllPropertyByID(questionID string) (models.Question, error) {
	var question models.Question

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{"_id": bson.ObjectIdHex(questionID) }).One(&question)
	
	return question, err
}

// DeleteQuestionByID is Delete an existing review
func (m *SessionDAO) DeleteQuestionByID(questionID string) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Remove(bson.M{"_id": bson.ObjectIdHex(questionID)})
	
	return err
}

// UpdateQuestionReportByID is Update reported
func (m *SessionDAO) UpdateQuestionReportByID(id string, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"reported": true, "update_at": updateAt}})
	
	return err
}