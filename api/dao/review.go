package dao

import (
	"log"


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
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of users
func (m *ReviewsDAO) FindAll() ([]models.Review, error) {
	var users []models.Review
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

// Find a user by its id
func (m *ReviewsDAO) FindById(id string) (models.Review, error) {
	var user models.Review
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Insert a user into database
func (m *ReviewsDAO) Insert(user models.Review) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// Delete an existing user
func (m *ReviewsDAO) Delete(user models.Review) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

// Update an existing user
func (m *ReviewsDAO) Update(user models.Review) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}
