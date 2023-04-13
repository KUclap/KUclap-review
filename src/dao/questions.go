package dao

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"kuclap-review-api/src/models"
)

// CreateQuestion is POST create question.
func (m *SessionDAO) CreateQuestion(ctx context.Context, question models.Question) error {
	bb, err := bson.Marshal(question)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateQuestion]: unable to marshal question")
	}

	if _, err := m.questions.InsertOne(ctx, bb); err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateQuestion]: unable to insert question")
	}

	return nil
}

// UpdateNumberAnswerByClassID is Update number of review
func (m *SessionDAO) UpdateNumberAnswerByQuestionID(ctx context.Context, questionId string, number int, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(questionId)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberAnswerByQuestionID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$inc": bson.M{"number_answers": int32(number)},
		"$set": bson.M{"update_at": updateAt},
	}

	r := m.questions.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberAnswerByQuestionID]: unable to find question and update it")
	}

	return nil
}

// FindAllQuestions is Find All of list of question
func (m *SessionDAO) FindAllQuestions(ctx context.Context) ([]models.Question, error) {
	cur, err := m.questions.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAllQuestions]: unable to find all")
	}

	var qq []models.Question
	for cur.Next(ctx) {
		var c models.Question

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindAllQuestions]: unable to decode content")
		}

		qq = append(qq, c)
	}

	return qq, err
}

// FindQuestionsByClassID is Find question by class_id
func (m *SessionDAO) FindQuestionsByClassID(ctx context.Context, classID string, page string, offset string) ([]models.ResQuestion, error) {
	iPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionsByClassID]: unable to convert page to number")
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionsByClassID]: unable to convert offset to number")
	}

	limit := int64(iOffset)
	skip := int64(iPage) * int64(iOffset)
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.M{
			"created_at": -1,
		})

	filter := bson.M{"class_id": classID}

	cur, err := m.questions.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindQuestionsByClassID]: unable to find all with class id %s", classID)
	}

	var qq []models.ResQuestion
	for cur.Next(ctx) {
		var c models.ResQuestion

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindQuestionsByClassID]: unable to decode content")
		}

		qq = append(qq, c)
	}

	return qq, err
}

// LastQuestions is Find last questions range with offset
func (m *SessionDAO) LastQuestions(ctx context.Context, page string, offset string) ([]models.ResQuestion, error) {
	iPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastQuestions]: unable to convert page to number")
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastQuestions]: unable to convert offset to number")
	}

	limit := int64(iOffset)
	skip := int64(iPage) * int64(iOffset)
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.M{
			"created_at": -1,
		})

	filter := bson.M{}

	cur, err := m.questions.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastQuestions]: unable to find all")
	}

	var qq []models.ResQuestion
	for cur.Next(ctx) {
		var c models.ResQuestion

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.LastQuestions]: unable to decode content")
		}

		qq = append(qq, c)
	}

	return qq, err
}

// FindQuestionByID is Find question by _id
func (m *SessionDAO) FindQuestionByID(ctx context.Context, questionID string) (*models.ResQuestion, error) {
	oid, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.questions.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindQuestionByID]: unable to find one with question id %s", questionID)
	}

	var question *models.ResQuestion
	if err := r.Decode(&question); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionByID]: unable to decode content")
	}

	return question, nil
}

// FindQuestionAllPropertyByID is Find question by _id
func (m *SessionDAO) FindQuestionAllPropertyByID(ctx context.Context, questionID string) (*models.Question, error) {
	oid, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionAllPropertyByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.questions.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindQuestionAllPropertyByID]: unable to find one with question id %s", questionID)
	}

	var question *models.Question
	if err := r.Decode(&question); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionAllPropertyByID]: unable to decode content")
	}

	return question, nil
}

// DeleteQuestionByID is Delete an existing review
func (m *SessionDAO) DeleteQuestionByID(ctx context.Context, questionID string) error {
	oid, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.DeleteQuestionByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	_, err = m.questions.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrapf(err, "[SessionDAO.DeleteQuestionByID]: unable to find one and delete question with question id %s", questionID)
	}

	return nil
}

// UpdateQuestionReportByID is Update reported
func (m *SessionDAO) UpdateQuestionReportByID(ctx context.Context, questionID string, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateQuestionReportByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{
			"reported":  true,
			"update_at": updateAt,
		},
	}

	r := m.questions.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrapf(err, "[SessionDAO.UpdateQuestionReportByID]: unable to find question and report with question id, %s", questionID)
	}

	return nil
}
