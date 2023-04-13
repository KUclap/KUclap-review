package dao

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION_REVIEWS   = "reviews"
	COLLECTION_CLASSES   = "classes"
	COLLECTION_REPORTS   = "reported"
	COLLECTION_QUESTIONS = "questions"
	COLLECTION_ANSWERS   = "answers"
	COLLECTION_RECAPS    = "recaps"
)

// SessionDAO is struct for allocate info for create connection with mongoDB
type SessionDAO struct {
	Server   string
	Database string

	reviews   *mongo.Collection
	classes   *mongo.Collection
	reports   *mongo.Collection
	questions *mongo.Collection
	answers   *mongo.Collection
	recaps    *mongo.Collection
}

// Connect is Establish a connection to database
func (m *SessionDAO) Connect() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(m.Server))
	if err != nil {
		log.Fatal(errors.Wrap(err, "[SessionDAO.Connect]: unable to create client mongodb"))
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "[SessionDAO.Connect]: unable to connect to mongodb database"))
	}

	db := client.Database(m.Database)

	m.reviews = db.Collection(COLLECTION_REVIEWS)
	m.classes = db.Collection(COLLECTION_CLASSES)
	m.reports = db.Collection(COLLECTION_REPORTS)
	m.questions = db.Collection(COLLECTION_QUESTIONS)
	m.answers = db.Collection(COLLECTION_ANSWERS)
	m.recaps = db.Collection(COLLECTION_RECAPS)

	log.Println("mongodb: Mongo has connected, Server get origin session. 🎉")
}
