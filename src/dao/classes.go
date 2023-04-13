package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"kuclap-review-api/src/models"
)

// UpdateStatsClassByCreated is Update stats on class by class_id
func (m *SessionDAO) UpdateStatsClassByCreated(ctx context.Context, classID string, newStats models.StatClass) error {
	filter := bson.M{"class_id": classID}
	operation := bson.M{
		"$set": bson.M{"stats": newStats},
		"$inc": bson.M{"number_reviewer": 1},
	}

	r := m.classes.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateStatsClassByCreated]: unable to find class and increase stats")
	}

	return nil
}

// UpdateStatsClassByDeleted is Update stats on class by class_id
func (m *SessionDAO) UpdateStatsClassByDeleted(ctx context.Context, classID string, newStats models.StatClass) error {
	filter := bson.M{"class_id": classID}
	operation := bson.M{
		"$set": bson.M{"stats": newStats},
		"$inc": bson.M{"number_reviewer": -1},
	}

	r := m.classes.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateStatsClassByDeleted]: unable to find class and decrease stats")
	}

	return nil
}

// UpdateNumberReviewByClassID is Update number of review
func (m *SessionDAO) UpdateNumberReviewByClassID(ctx context.Context, classID string, updateAt time.Time) error {
	filter := bson.M{"class_id": classID}
	operation := bson.M{
		"$set": bson.M{"update_at": updateAt},
		"$inc": bson.M{"number_reviewer": 1},
	}

	r := m.classes.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberReviewByClassID]: unable to find class and increase stats")
	}

	return nil
}

// UpdateNumberRecapByClassID is Update number of review
func (m *SessionDAO) UpdateNumberRecapByClassID(ctx context.Context, classID string, number int, updateAt time.Time) error {
	filter := bson.M{"class_id": classID}
	operation := bson.M{
		"$set": bson.M{"update_at": updateAt},
		"$inc": bson.M{"number_recaps": int32(number)},
	}

	r := m.classes.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberRecapByClassID]: unable to find class and increase number of recaps stats")
	}

	return nil
}

// UpdateNumberQuestionByClassID is Update number of review
func (m *SessionDAO) UpdateNumberQuestionByClassID(ctx context.Context, classID string, number int, updateAt time.Time) error {
	filter := bson.M{"class_id": classID}
	operation := bson.M{
		"$set": bson.M{"update_at": updateAt},
		"$inc": bson.M{"number_questions": int32(number)},
	}

	r := m.classes.FindOneAndUpdate(ctx, filter, operation)
	if err := r.Err(); err != nil {
		return errors.Wrap(err, "[SessionDAO.UpdateNumberQuestionByClassID]: unable to find class and increase number of questions stats")
	}

	return nil

}

// FindClassByClassID is Find class by class_id
func (m *SessionDAO) FindClassByClassID(ctx context.Context, classID string) (*models.Class, error) {
	filter := bson.M{"class_id": classID}

	r := m.classes.FindOne(ctx, filter)
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "[SessionDAO.FindClassByClassID]: unable to find one, %s", classID)
	}

	var class *models.Class
	if err := r.Decode(&class); err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindClassByClassID]: unable to decode content")
	}

	class.ToDefault()

	return class, nil
}

// FindAllClasses is Find All of list of classes
func (m *SessionDAO) FindAllClasses(ctx context.Context) ([]models.Class, error) {
	cur, err := m.classes.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "[SessionDAO.FindAllClasses]: unable to find all")
	}

	var cc []models.Class
	for cur.Next(ctx) {
		var c models.Class

		err := cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrap(err, "[SessionDAO.FindAllClasses]: unable to decode content")
		}

		c.ToDefault()

		cc = append(cc, c)
	}

	return cc, err
}

// InsertClass is Insert class to database
func (m *SessionDAO) InsertClass(ctx context.Context, class models.Class) error {
	bb, err := bson.Marshal(class)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.InsertClass]: unable to marshal answer")
	}

	if _, err := m.classes.InsertOne(ctx, bb); err != nil {
		return errors.Wrap(err, "[SessionDAO.InsertClass]: unable to insert class")
	}

	return nil
}
