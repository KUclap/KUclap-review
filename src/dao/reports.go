package dao

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"kuclap-review-api/src/models"
)

// InsertReport is Insert report to database
func (m *SessionDAO) InsertReport(ctx context.Context, report models.Report) error {
	bb, err := bson.Marshal(report)
	if err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateRecap]: unable to marshal report")
	}

	if _, err := m.reports.InsertOne(ctx, bb); err != nil {
		return errors.Wrap(err, "[SessionDAO.CreateRecap]: unable to insert report")
	}

	return nil
}
