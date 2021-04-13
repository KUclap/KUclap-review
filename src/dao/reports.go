package dao

import (
	"kuclap-review-api/src/models"
)

// InsertReport is Insert report to database
func (m *SessionDAO) InsertReport(report models.Report) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_REPORTS).Insert(&report)

	return err	
}