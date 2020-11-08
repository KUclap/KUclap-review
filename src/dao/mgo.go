package dao

import (
	"log"
	"crypto/tls"
	"net"
	"time"
	"strconv"
	"kuclap-review-api/src/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION_REVIEWS = "reviews"
	COLLECTION_CLASSES = "classes"
	COLLECTION_REPORTS = "reported"
)

type SessionDAO struct {
	Server   string
	Database string
}

var session *mgo.Session

// Establish a connection to database
func (m *SessionDAO) Connect() {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(m.Server)
	log.Println("MGO: Parseurl finish, Connecting... ✅")
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	log.Println("MGO: TLS configed, Connecting... ✅")
	session, err = mgo.DialWithInfo(dialInfo)
	
	if err != nil {
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)

	// session = mgoSession.DB(m.Database)
	log.Println("MGO: Mongo has connected, Server get origin session. 🎉")
}

// Insert report to database
func (m *SessionDAO) InsertReport(report models.Report) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_REPORTS).Insert(&report)
	return err	
}

// Update clap by id
func (m *SessionDAO) UpdateClapById(id string, newClap uint64, updateAt time.Time) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"clap": newClap}, "$set": bson.M{"update_at": updateAt}})
	return err
}

// Update boo by id
func (m *SessionDAO) UpdateBooById(id string, newBoo uint64, updateAt time.Time) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"boo": newBoo}, "$set": bson.M{"update_at": updateAt}})
	return err
}

// Update reported
func (m *SessionDAO) UpdateReportById(id string, updateAt time.Time) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"reported": true, "update_at": updateAt}})
	return err
}

// Update stats on class by class_id 
func (m *SessionDAO) UpdateStatsClassByCreated(classId string, newStats models.StatClass) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classId}, bson.M{"$set": bson.M{"stats": newStats}, "$inc": bson.M{"number_reviewer": 1}})
	return err
}

// Update stats on class by class_id 
func (m *SessionDAO) UpdateStatsClassByDeleted(classId string, newStats models.StatClass) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classId}, bson.M{"$set": bson.M{"stats": newStats}, "$inc": bson.M{"number_reviewer": -1}})
	return err
}

// Update number of review
func (m *SessionDAO) UpdateNuberReviewByClassID(classId string, updateAt time.Time) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classId}, bson.M{"$inc": bson.M{"number_reviewer": 1}})
	return err
}

// Find class by class_id
func (m *SessionDAO) FindClassByClassID(classId string) (models.Class, error) {
	db := session.Copy()
	defer db.Close()
	var class models.Class
	err := db.DB(m.Database).C(COLLECTION_CLASSES).Find(bson.M{"class_id": classId}).One(&class)
	return class, err
}

// Find All of list of classes 
func (m *SessionDAO) FindAllClasses() ([]models.Class, error) {
	db := session.Copy()
	defer db.Close()
	var classes []models.Class
	err := db.DB(m.Database).C(COLLECTION_CLASSES).Find(bson.M{}).All(&classes)
	return classes, err
}

// Insert class to database
func (m *SessionDAO) InsertClass(class models.Class) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_CLASSES).Insert(&class)
	return err	
}

// Find reviews by class_id
func (m *SessionDAO) FindReviewsByClassID(classId string, page string, offset string) ([]models.ResReview, error) {
	db := session.Copy()
	defer db.Close()
	var reviews []models.ResReview
	iPage, err := strconv.Atoi(page)
	iOffset, err := strconv.Atoi(offset)
	err = db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{"class_id": classId}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&reviews)
	return reviews, err
}

// Find last reviews range with offset
func (m *SessionDAO) LastReviews(page string,offset string) ([]models.ResReview, error) {
	db := session.Copy()
	defer db.Close()
	var reviews []models.ResReview
	iPage, err := strconv.Atoi(page)
	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		log.Println("err : atoi.", err)
	}
	// err = db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{}).Sort("-$natural").Limit(iOffset).All(&reviews)
	err = db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&reviews)
	return reviews, err
}

// Find list of reviews: All of reviews
func (m *SessionDAO) FindAll() ([]models.ResReview, error) {
	db := session.Copy()
	defer db.Close()
	var reviews []models.ResReview
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{}).All(&reviews)
	return reviews, err
}

// Find a review by its id
func (m *SessionDAO) FindById(id string) (models.Review, error) {
	db := session.Copy()
	defer db.Close()
	var review models.Review
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).FindId(bson.ObjectIdHex(id)).One(&review)
	return review, err
}

// Insert a review into database
func (m *SessionDAO) Insert(review models.Review) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).Insert(&review)
	return err
}

// Delete an existing review
func (m *SessionDAO) DeleteById(id string) error {
	db := session.Copy()
	defer db.Close()
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).RemoveId(bson.ObjectIdHex(id))
	return err
}

// Find reviews reported
// func (m *SessionDAO) FindReviewsReported() ([]models.Review, error) {
// 	var reviews []models.Review
// 	err := db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{"reported": bson.M{"$ne": false}}).All(&reviews)
// 	return reviews, err
// }

// Update an existing review
// func (m *SessionDAO) Update(review models.Review) error {
// 	err := db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(review.ID, &review)
// 	return err
// }
