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

// CreateRecap is POST create recap.
func (m *SessionDAO) CreateRecap(ctx context.Context, recap models.Recap) error {
	bb, err := bson.Marshal(recap)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateRecap]: unable to marshal recap")
	}

	if _, err := m.recaps.InsertOne(ctx, bb); err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateRecap]: unable to insert recap")
	}

	return nil
}

// UpdateNumberDownloadedByRecapID is Update number of review
func (m *SessionDAO) UpdateNumberDownloadedByRecapID(ctx context.Context, recapID string, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(recapID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberDownloadedByRecapID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{"update_at": updateAt},
		"$inc": bson.M{"downloaded": 1},
	}

	r := m.recaps.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberDownloadedByRecapID]: unable to find recap and increase number of download")
	}

	return nil
}

// FindAllRecaps is Find All of list of recap
func (m *SessionDAO) FindAllRecaps(ctx context.Context) ([]models.ResRecap, error) {
	cur, err := m.recaps.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAllRecaps]: unable to find all")
	}

	var rr []models.ResRecap
	for cur.Next(ctx) {
		var c models.ResRecap

		err := cur.Decode(&c)

		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindAllRecaps]: unable to decode content")
		}

		rr = append(rr, c)
	}

	return rr, err
}

// FindRecapsByClassID is Find recap by class_id
func (m *SessionDAO) FindRecapsByClassID(ctx context.Context, classID string, page string, offset string) ([]models.ResRecap, error) {
	iPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindRecapsByClassID]: unable to convert page to number")
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindRecapsByClassID]: unable to convert offset to number")
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

	cur, err := m.recaps.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindRecapsByClassID]: unable to find all with class id %s", classID)
	}

	var rr []models.ResRecap
	for cur.Next(ctx) {
		var c models.ResRecap

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindRecapsByClassID]: unable to decode content")
		}

		rr = append(rr, c)
	}

	return rr, err
}

// LastRecaps is Find last recaps range with offset
func (m *SessionDAO) LastRecaps(ctx context.Context, page string, offset string) ([]models.ResRecap, error) {
	iPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastRecaps]: unable to convert page to number")
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastRecaps]: unable to convert offset to number")
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

	cur, err := m.recaps.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastRecaps]: unable to find all")
	}

	var rr []models.ResRecap
	for cur.Next(ctx) {
		var c models.ResRecap

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.LastRecaps]: unable to decode content")
		}

		rr = append(rr, c)
	}

	return rr, err
}

// FindRecapByID is Find recap by recap_id
func (m *SessionDAO) FindRecapByID(ctx context.Context, recapID string) (*models.ResRecap, error) {
	oid, err := primitive.ObjectIDFromHex(recapID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindRecapByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.recaps.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindQuestionByID]: unable to find one with recap id %s", recapID)
	}

	var recap *models.ResRecap
	if err := r.Decode(&recap); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindQuestionByID]: unable to decode content")
	}

	return recap, nil
}

// FindRecapAllPropertyByID is Find recap by recap_id
func (m *SessionDAO) FindRecapAllPropertyByID(ctx context.Context, recapID string) (*models.Recap, error) {
	oid, err := primitive.ObjectIDFromHex(recapID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindRecapAllPropertyByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.recaps.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindRecapAllPropertyByID]: unable to find one with recap id %s", recapID)
	}

	var recap *models.Recap
	if err := r.Decode(&recap); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindRecapAllPropertyByID]: unable to decode content")
	}

	return recap, nil
}

// DeleteRecapByID is Delete an existing review
func (m *SessionDAO) DeleteRecapByID(ctx context.Context, recapID string) error {
	oid, err := primitive.ObjectIDFromHex(recapID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.DeleteRecapByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	_, err = m.recaps.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrapf(err, "[SessionDAO.DeleteRecapByID]: unable to find one and delete with recap id %s", recapID)
	}

	return nil
}

// UpdateReportByID is Update reported
func (m *SessionDAO) UpdateRecapReportByID(ctx context.Context, recapID string, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(recapID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateRecapReportByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{
			"reported":  true,
			"update_at": updateAt,
		},
	}

	r := m.recaps.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateRecapReportByID]: unable to find recap and report it")
	}

	return nil
}
