package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"kuclap-review-api/src/models"
)

// CreateAnswer is POST create answer.
func (m *SessionDAO) CreateAnswer(ctx context.Context, answer models.Answer) error {
	bb, err := bson.Marshal(answer)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateAnswer]: unable to marshal answer")
	}

	if _, err := m.answers.InsertOne(ctx, bb); err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateAnswer]: unable to insert answer")
	}

	return nil
}

// FindAllAnswers is Find All of list of answer
func (m *SessionDAO) FindAllAnswers(ctx context.Context) ([]models.Answer, error) {
	cur, err := m.answers.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAllAnswers]: unable to find all")
	}

	var aa []models.Answer
	for cur.Next(ctx) {
		var c models.Answer

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindAllAnswers]: unable to decode content")
		}

		aa = append(aa, c)
	}

	return aa, err
}

// FindLastAnswersByQuestionID is Find All of list of answer
func (m *SessionDAO) FindLastAnswersByQuestionID(ctx context.Context, questionId string) ([]models.ResAnswer, error) {
	filter := bson.M{"question_id": questionId}

	cur, err := m.answers.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var aa []models.ResAnswer
	for cur.Next(ctx) {
		var c models.ResAnswer

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindLastAnswersByQuestionID]: unable to decode ResAnswer")
		}

		aa = append(aa, c)
	}

	return aa, err
}

// FindAnswerByID is Find answer by answer_id
func (m *SessionDAO) FindAnswerByID(ctx context.Context, answerID string) (*models.ResAnswer, error) {
	oid, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAnswerByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.answers.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindAnswerByID]: unable to find one, %s", answerID)
	}

	var answer *models.ResAnswer
	if err := r.Decode(&answer); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAnswerByID]: unable to decode content")
	}

	return answer, nil
}

// FindAnswerAllPropertyByID is Find answer by _id
func (m *SessionDAO) FindAnswerAllPropertyByID(ctx context.Context, answerID string) (*models.Answer, error) {
	oid, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAnswerAllPropertyByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.answers.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindAnswerAllPropertyByID]: unable to find one, %s", answerID)
	}

	var answer *models.Answer
	if err := r.Decode(&answer); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAnswerAllPropertyByID]: unable to decode content")
	}

	return answer, nil
}

// DeleteAnswerByID is Delete an existing review
func (m *SessionDAO) DeleteAnswerByID(ctx context.Context, answerID string) error {
	oid, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.DeleteAnswerByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.answers.FindOneAndDelete(ctx, filter)
	if err := r.Err(); err != nil {
		return errors.Wrapf(err, "[SessionDAO.DeleteAnswerByID]: unable to find one and delete, %s", answerID)
	}

	return nil
}

// DeleteAnswersByQuestionID is Delete an existing review
func (m *SessionDAO) DeleteAnswersByQuestionID(ctx context.Context, questionId string) error {
	filter := bson.M{"question_id": questionId}

	_, err := m.answers.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrapf(err, "[SessionDAO.DeleteAnswersByQuestionID]: unable to find one and delete all answer, %s", questionId)
	}

	return nil
}

// UpdateAnswerReportByID is Update reported
func (m *SessionDAO) UpdateAnswerReportByID(ctx context.Context, answerID string, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateAnswerReportByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{
			"reported":  true,
			"update_at": updateAt,
		},
	}

	r := m.answers.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateAnswerReportByID]: unable to find answer and report it")
	}

	return nil
}
