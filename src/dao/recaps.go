package dao

import (
	"log"
	"strconv"

	"kuclap-review-api/src/models"

	"gopkg.in/mgo.v2/bson"
)

// CreateRecap is POST create recap.
func (m *SessionDAO) CreateRecap(recap models.Recap) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).Insert(&recap)
	
	return err
}

// UpdateNumberRecapByClassID is Update number of review
func (m *SessionDAO) UpdateNumberDonwloadedByRecapID(recapID string) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).Update(bson.M{"_id": bson.ObjectIdHex(recapID)}, bson.M{ "$inc": bson.M{"downloaded": int32(1)} })
	
	return err
}

// UpdateNumberRecapByClassID is Update number of review
func (m *SessionDAO) UpdateNumberRecapByClassID(classID string, number int) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classID}, bson.M{ "$inc": bson.M{"number_recaps": int32(number)} })
	
	return err
}

// FindAllRecaps is Find All of list of recap 
func (m *SessionDAO) FindAllRecaps() ([]models.ResRecap, error) {
	var recaps []models.ResRecap

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).Find(bson.M{}).All(&recaps)
	return recaps, err
}

// FindRecapsByClassID is Find recap by class_id
func (m *SessionDAO) FindRecapsByClassID(classID string, page string, offset string) ([]models.ResRecap, error) {
	var recaps []models.ResRecap

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

	err = db.DB(m.Database).C(COLLECTION_RECAPS).Find(bson.M{"class_id": classID}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&recaps)
	
	return recaps, err
}

// LastRecaps is Find last recaps range with offset
func (m *SessionDAO) LastRecaps(page string, offset string) ([]models.ResRecap, error) {
	var recaps []models.ResRecap

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

	err	=	db.DB(m.Database).C(COLLECTION_RECAPS).Find(bson.M{}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&recaps)
	
	return recaps, err
}

// FindRecapByID is Find recap by recap_id
func (m *SessionDAO) FindRecapByID(recapID string) (models.ResRecap, error) {
	var recap models.ResRecap

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).Find(bson.M{"_id": bson.ObjectIdHex(recapID) }).One(&recap)
	
	return recap, err
}

// FindRecapAllPropertyByID is Find recap by recap_id
func (m *SessionDAO) FindRecapAllPropertyByID(recapID string) (models.Recap, error) {
	var recap models.Recap

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).Find(bson.M{"_id": bson.ObjectIdHex(recapID) }).One(&recap)
	
	return recap, err
}

// DeleteRecapByID is Delete an existing review
func (m *SessionDAO) DeleteRecapByID(recapID string) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).Remove(bson.M{"_id": bson.ObjectIdHex(recapID)})
	
	return err
}

// UpdateReportByID is Update reported
func (m *SessionDAO) UpdateRecapReportByID(id string) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_RECAPS).UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"reported": true}})
	
	return err
}