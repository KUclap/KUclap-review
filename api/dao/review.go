package dao

import (
	"log"
	"crypto/tls"
	"net"
	"fmt"
	"time"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/marsDev31/kuclap-backend/api/models"
)

type SessionDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	CREVIEWS = "reviews"
	CCLASSES = "classes"
)

// Establish a connection to database
func (m *SessionDAO) Connect() {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(m.Server)
	fmt.Println("CONNECTING: Parseurl finish.")
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	fmt.Println("CONNECTING: TLS configed.")
	session, err := mgo.DialWithInfo(dialInfo)
	
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
	fmt.Println("CONNECTED: got session.")
}

// Update clap by id
func (m *SessionDAO) UpdateClapById(id string, newClap uint64, updateAt time.Time) error {
	err := db.C(CREVIEWS).Update(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"clap": newClap}, "$set": bson.M{"update_at": updateAt}})
	return err
}

// Update boo by id
func (m *SessionDAO) UpdateBooById(id string, newBoo uint64, updateAt time.Time) error {
	err := db.C(CREVIEWS).Update(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"clap": newBoo}, "$set": bson.M{"update_at": updateAt}})
	return err
}

// Update reported
func (m *SessionDAO) UpdateReportById(id string, updateAt time.Time) error {
	err := db.C(CREVIEWS).Update(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"reported": true, "update_at": updateAt}})
	return err
}

// Update stats on class by class_id 
func (m *SessionDAO) UpdateStatsClass(classId string, newStats models.StatClass) error {
	err := db.C(CCLASSES).Update(bson.M{"class_id": classId}, bson.M{"$set": bson.M{"stats": newStats}})
	return err
}

// Find class by class_id
func (m *SessionDAO) FindClassByClassId(classId string) (models.Class, error) {
	var class models.Class
	err := db.C(CCLASSES).Find(bson.M{"class_id": classId}).One(&class)
	return class, err
}

// Find All of list of classes 
func (m *SessionDAO) FindAllClasses() ([]models.Class, error) {
	var classes []models.Class
	err := db.C(CCLASSES).Find(bson.M{}).All(&classes)
	return classes, err
}

// Insert class to database
func (m *SessionDAO) InsertClass(class models.Class) error {
	err := db.C(CCLASSES).Insert(&class)
	return err	
}

// Find reviews by class_id
func (m *SessionDAO) FindReviewsByClassId(classId string) ([]models.Class, error) {
	var reviews []models.Class
	err := db.C(CREVIEWS).Find(bson.M{"class_id": classId}).All(&reviews)
	return reviews, err
}

// Find last reviews range with offset
func (m *SessionDAO) LastReviews(offset string) ([]models.Review, error) {
	var reviews []models.Review
	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Println("err : atoi.", err)
	}
	err = db.C(CREVIEWS).Find(bson.M{}).Sort("-$natural").Limit(iOffset).All(&reviews)
	return reviews, err
}

// Find list of reviews: All of reviews
func (m *SessionDAO) FindAll() ([]models.Review, error) {
	var reviews []models.Review
	err := db.C(CREVIEWS).Find(bson.M{}).All(&reviews)
	return reviews, err
}

// Find a review by its id
func (m *SessionDAO) FindById(id string) (models.Review, error) {
	var review models.Review
	err := db.C(CREVIEWS).FindId(bson.ObjectIdHex(id)).One(&review)
	return review, err
}

// Insert a review into database
func (m *SessionDAO) Insert(review models.Review) error {
	err := db.C(CREVIEWS).Insert(&review)
	return err
}

// Delete an existing review
func (m *SessionDAO) Delete(review models.Review) error {
	err := db.C(CREVIEWS).Remove(&review)
	return err
}

// Update an existing review
func (m *SessionDAO) Update(review models.Review) error {
	err := db.C(CREVIEWS).UpdateId(review.ID, &review)
	return err
}
