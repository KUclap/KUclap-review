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
	COLLECTION_REVIEWS	=	"reviews"
	COLLECTION_CLASSES	=	"classes"
	COLLECTION_REPORTS	=	"reported"
	COLLECTION_QUESTION	=	"questions"
)

// SessionDAO is struct for allocate info for create connection with mongoDB
type SessionDAO struct {
	Server   string
	Database string
}

var session *mgo.Session

// ###################################
// ######## MGO APDAPTER #############
// ###################################

// Connect is Establish a connection to database
func (m *SessionDAO) Connect() {

	tlsConfig		:=	&tls.Config{}

	dialInfo, err	:=	mgo.ParseURL(m.Server)
	if err != nil {
		log.Fatal(err)
	}

	dialInfo.DialServer	=	func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err	:=	tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err	=	mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)

	log.Println("MGO: Mongo has connected, Server get origin session. 🎉")
}


// ###################################
// ######## REPORT APDAPTER ##########
// ###################################

// InsertReport is Insert report to database
func (m *SessionDAO) InsertReport(report models.Report) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_REPORTS).Insert(&report)

	return err	
}

// ###################################
// ######## CLASS APDAPTER ###########
// ###################################

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

// ###################################
// ######## REVIEW APDAPTER ##########
// ###################################

// UpdateClapByID is Update clap by id
func (m *SessionDAO) UpdateClapByID(id string, newClap uint64, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"clap": newClap}, "$set": bson.M{"update_at": updateAt}})
	
	return err
}

// UpdateBooByID is Update boo by id
func (m *SessionDAO) UpdateBooByID(id string, newBoo uint64, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"boo": newBoo}, "$set": bson.M{"update_at": updateAt}})
	
	return err
}

// UpdateReportByID is Update reported
func (m *SessionDAO) UpdateReportByID(id string, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"reported": true, "update_at": updateAt}})
	
	return err
}

// FindReviewsByClassID is Find reviews by class_id
func (m *SessionDAO) FindReviewsByClassID(classID string, page string, offset string) ([]models.ResReview, error) {
	var reviews []models.ResReview

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

	err	=	db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{"class_id": classID}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&reviews)

	return reviews, err
}

// LastReviews is Find last reviews range with offset
func (m *SessionDAO) LastReviews(page string,offset string) ([]models.ResReview, error) {
	var reviews []models.ResReview

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
	
	err = db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&reviews)

	return reviews, err
}

// FindAll is Find list of reviews: All of reviews
func (m *SessionDAO) FindAll() ([]models.ResReview, error) {
	db := session.Copy()
	defer db.Close()
	var reviews []models.ResReview
	err := db.DB(m.Database).C(COLLECTION_REVIEWS).Find(bson.M{}).All(&reviews)
	return reviews, err
}

// FindByID is Find a review by its id
func (m *SessionDAO) FindByID(id string) (models.ResReview, error) {
	var review models.ResReview

	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).FindId(bson.ObjectIdHex(id)).One(&review)
	
	return review, err
}

// FindReviewAllPropertyByID is Find a review by its id
func (m *SessionDAO) FindReviewAllPropertyByID(id string) (models.Review, error) {
	var review models.Review

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).FindId(bson.ObjectIdHex(id)).One(&review)
	
	return review, err
}

// Insert a review into database
func (m *SessionDAO) Insert(review models.Review) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).Insert(&review)
	
	return err
}

// DeleteByID is Delete an existing review
func (m *SessionDAO) DeleteByID(id string) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_REVIEWS).RemoveId(bson.ObjectIdHex(id))
	
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

// ###################################
// ######## QUESTION APDAPTER ########
// ###################################

// CreateQuestion is POST create question.
func (m *SessionDAO) CreateQuestion(question models.Question) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Insert(&question)
	
	return err
}

// UpdateNumberQuestionByClassID is Update number of review
func (m *SessionDAO) UpdateNumberQuestionByClassID(classID string, number int, updateAt time.Time) error {
	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_CLASSES).Update(bson.M{"class_id": classID}, bson.M{ "$inc": bson.M{"number_questions": int32(number)}, "$set": bson.M{"update_at": updateAt} })
	
	return err
}

// FindAllQuestions is Find All of list of question 
func (m *SessionDAO) FindAllQuestions() ([]models.Question, error) {
	var questions []models.Question

	db	:=	session.Copy()
	defer db.Close()
	
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{}).All(&questions)
	return questions, err
}

// FindQuestionsByClassID is Find question by class_id
func (m *SessionDAO) FindQuestionsByClassID(classID string, page string, offset string) ([]models.ResQuestion, error) {
	var questions []models.ResQuestion

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

	err = db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{"class_id": classID}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&questions)
	
	return questions, err
}

// LastQuestions is Find last questions range with offset
func (m *SessionDAO) LastQuestions(page string, offset string) ([]models.ResQuestion, error) {
	var questions []models.ResQuestion

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

	err	=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{}).Sort("-$natural").Skip(iPage * iOffset).Limit(iOffset).All(&questions)
	
	return questions, err
}

// FindQuestionByID is Find question by question_id
func (m *SessionDAO) FindQuestionByID(questionID string) (models.ResQuestion, error) {
	var question models.ResQuestion

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{"question_id": bson.ObjectIdHex(questionID) }).One(&question)
	
	return question, err
}

// FindQuestionAllPropertyByID is Find question by question_id
func (m *SessionDAO) FindQuestionAllPropertyByID(questionID string) (models.Question, error) {
	var question models.Question

	db	:=	session.Copy()
	defer db.Close()
	
	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Find(bson.M{"question_id": bson.ObjectIdHex(questionID) }).One(&question)
	
	return question, err
}

// DeleteQuestionByID is Delete an existing review
func (m *SessionDAO) DeleteQuestionByID(questionID string) error {
	db	:=	session.Copy()
	defer db.Close()

	err	:=	db.DB(m.Database).C(COLLECTION_QUESTION).Remove(bson.M{"question_id": bson.ObjectIdHex(questionID)})
	
	return err
}