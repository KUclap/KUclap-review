package dao

import (
	"time"
	
	"kuclap-review-api/src/models"
	"gopkg.in/mgo.v2/bson"
)

// UpdateStatsClassByCreated is Update stats on class by class_id 
func (m *SessionDAO) UpdateStatsClassByCreated(classID string, newStats models.StatClass) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classID}, bson.M{"$set": bson.M{"stats": newStats}, "$inc": bson.M{"number_reviewer": 1}})

	return err
}

// UpdateStatsClassByDeleted is Update stats on class by class_id 
func (m *SessionDAO) UpdateStatsClassByDeleted(classID string, newStats models.StatClass) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classID}, bson.M{"$set": bson.M{"stats": newStats}, "$inc": bson.M{"number_reviewer": -1}})

	return err
}

// UpdateNuberReviewByClassID is Update number of review
func (m *SessionDAO) UpdateNuberReviewByClassID(classID string, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classID}, bson.M{"$inc": bson.M{"number_reviewer": 1}})
	
	return err
}

// FindClassByClassID is Find class by class_id
func (m *SessionDAO) FindClassByClassID(classID string) (models.Class, error) {
	var class models.Class

	db	:=	session.Copy()
	defer db.Close()

	err := db.DB(m.Database).C(COLLECTION_CLASSES).Find(bson.M{"class_id": classID}).One(&class)

	return class, err
}

// FindAllClasses is Find All of list of classes 
func (m *SessionDAO) FindAllClasses() ([]models.Class, error) {
	var classes []models.Class

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Find(bson.M{}).All(&classes)
	
	return classes, err
}

// InsertClass is Insert class to database
func (m *SessionDAO) InsertClass(class models.Class) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Insert(&class)

	return err	
}