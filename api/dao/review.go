package dao

import (
	"log"
	"crypto/tls"
	"net"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/marsDev31/kuclap-backend/api/models"
)

type ReviewsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "reviews"
)

// Establish a connection to database
func (m *ReviewsDAO) Connect() {

	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(m.Server)
	fmt.Println("CONNECTING: Parseurl finish.")
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	fmt.Println("CONNECTING: TLS configed. 123")
	session, err := mgo.DialWithInfo(dialInfo)
	
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
	fmt.Println("CONNECTED: got session.")
}

// Find list of reviews
func (m *ReviewsDAO) FindAll() ([]models.Review, error) {
	var reviews []models.Review
	err := db.C(COLLECTION).Find(bson.M{}).All(&reviews)
	return reviews, err
}

// Find a review by its id
func (m *ReviewsDAO) FindById(id string) (models.Review, error) {
	var review models.Review
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&review)
	return review, err
}

// Insert a review into database
func (m *ReviewsDAO) Insert(review models.Review) error {
	err := db.C(COLLECTION).Insert(&review)
	return err
}

// Delete an existing review
func (m *ReviewsDAO) Delete(review models.Review) error {
	err := db.C(COLLECTION).Remove(&review)
	return err
}

// Update an existing review
func (m *ReviewsDAO) Update(review models.Review) error {
	err := db.C(COLLECTION).UpdateId(review.ID, &review)
	return err
}
