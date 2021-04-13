package dao

import (
	"log"
	"time"
	"strconv"

	"kuclap-review-api/src/models"

	"gopkg.in/mgo.v2/bson"
)

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