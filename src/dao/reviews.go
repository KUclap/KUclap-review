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

// UpdateClapByID is Update clap by id
func (m *SessionDAO) UpdateClapByID(ctx context.Context, reviewID string, newClap uint64, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateClapByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{"update_at": updateAt},
		"$inc": bson.M{"clap": newClap},
	}

	r := m.reviews.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrapf(err, "[SessionDAO.UpdateClapByID]: unable to find review and set clap of review with review id, %s", reviewID)
	}

	return nil
}

// UpdateBooByID is Update boo by id
func (m *SessionDAO) UpdateBooByID(ctx context.Context, reviewID string, newBoo uint64, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateBooByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{"update_at": updateAt},
		"$inc": bson.M{"boo": newBoo},
	}

	r := m.reviews.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrapf(err, "[SessionDAO.UpdateBooByID]: unable to find review and set boo of review with review id, %s", reviewID)
	}

	return nil
}

// UpdateReportByID is Update reported
func (m *SessionDAO) UpdateReportByID(ctx context.Context, reviewID string, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateReportByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{
			"reported":  true,
			"update_at": updateAt,
		},
	}

	r := m.reviews.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrapf(err, "[SessionDAO.UpdateReportByID]: unable to find review and report with review id, %s", reviewID)
	}

	return nil
}

// UpdateAdminDeleteByID is Update reported
func (m *SessionDAO) UpdateAdminDeleteByID(ctx context.Context, reviewID string, deleteReason string, updateAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateAdminDeleteByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}
	operation := bson.M{
		"$set": bson.M{
			"delete_reason": deleteReason,
			"update_at":     updateAt,
		},
	}

	r := m.reviews.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrapf(err, "[SessionDAO.UpdateAdminDeleteByID]: unable to find review and mark delete with review id, %s", reviewID)
	}

	return nil
}

// FindReviewsByClassID is Find reviews by class_id
func (m *SessionDAO) FindReviewsByClassID(ctx context.Context, classID string, page string, offset string) ([]models.ResReview, error) {
	iPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindReviewsByClassID]: unable to convert page to number")
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindReviewsByClassID]: unable to convert offset to number")
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

	cur, err := m.reviews.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindReviewsByClassID]: unable to find all with class id %s", classID)
	}

	var rr []models.ResReview
	for cur.Next(ctx) {
		var c models.ResReview

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindReviewsByClassID]: unable to decode content")
		}

		c.ToDefault()

		rr = append(rr, c)
	}

	return rr, err
}

// LastReviews is Find last reviews range with offset
func (m *SessionDAO) LastReviews(ctx context.Context, page string, offset string, query bson.M) ([]models.ResReview, error) {
	iPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastReviews]: unable to convert page to number")
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastReviews]: unable to convert offset to number")
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

	cur, err := m.reviews.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.LastReviews]: unable to find all")
	}

	var rr []models.ResReview
	for cur.Next(ctx) {
		var c models.ResReview

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.LastReviews]: unable to decode content")
		}

		c.ToDefault()

		rr = append(rr, c)
	}

	return rr, err
}

// FindAll is Find list of reviews: All of reviews
func (m *SessionDAO) FindAll(ctx context.Context) ([]models.ResReview, error) {
	cur, err := m.reviews.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAll]: unable to find all")
	}

	var rr []models.ResReview
	for cur.Next(ctx) {
		var c models.ResReview

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindAll]: unable to decode content")
		}

		c.ToDefault()

		rr = append(rr, c)
	}

	return rr, err
}

// FindByID is Find a review by its id
func (m *SessionDAO) FindByID(ctx context.Context, reviewID string) (*models.ResReview, error) {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.reviews.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindByID]: unable to find one, %s", reviewID)
	}

	var review *models.ResReview
	if err := r.Decode(&review); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindByID]: unable to decode content")
	}

	review.ToDefault()

	return review, nil
}

// FindReviewAllPropertyByID is Find a review by its id
func (m *SessionDAO) FindReviewAllPropertyByID(ctx context.Context, reviewID string) (*models.Review, error) {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindReviewAllPropertyByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	r := m.reviews.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindReviewAllPropertyByID]: unable to find one, %s", reviewID)
	}

	var review *models.Review
	if err := r.Decode(&review); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindReviewAllPropertyByID]: unable to decode content")
	}

	review.ToDefault()

	return review, nil
}

// Insert a review into database
func (m *SessionDAO) Insert(ctx context.Context, review models.Review) error {
	bb, err := bson.Marshal(review)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.Insert]: unable to marshal review")
	}

	if _, err := m.reviews.InsertOne(ctx, bb); err != nil {
		return errors.Wrap(err, "[SessionDAO.Insert]: unable to insert review")
	}

	return nil
}

// DeleteByID is Delete an existing review
func (m *SessionDAO) DeleteByID(ctx context.Context, reviewID string) error {
	oid, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.DeleteByID]: unable to parse OID from hex")
	}

	filter := bson.M{"_id": oid}

	_, err = m.reviews.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrapf(err, "[SessionDAO.DeleteByID]: unable to find one and delete review with review id %s", reviewID)
	}

	return nil
}

// Update an existing review
// func (m *SessionDAO) Update(review models.Review) error {
// 	db	:=	session.Copy()
// 	defer db.Close()

// 	err := db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(review.ID, &review)

// 	return err
// }

// Find reviews reported
// func (m *SessionDAO) FindReviewsReported() ([]models.Review, error) {
// 	var reviews []models.Review
// 	err := db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{"reported": bson.M{"$ne": false}}).All(&reviews)
// 	return reviews, err
// }
