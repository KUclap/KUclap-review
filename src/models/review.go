package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Review is model for create review
type Review struct {
	ID				bson.ObjectId	`json:"reviewId" bson:"_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	Text			string			`json:"text" bson:"text"`
	Author			string			`json:"author" bson:"author"`
	Grade			string			`json:"grade" bson:"grade"`
	Auth			string			`json:"auth,omitempty" bson:"auth"`
	Clap			uint64			`json:"clap" bson:"clap"`
	Boo				uint64			`json:"boo" bson:"boo"`
	Stats			StatReview		`json:"stats" bson:"stats"`	
	ClassNameTH		string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN		string			`json:"classNameEN" bson:"class_name_en"`
	Sec				uint64			`json:"sec" bson:"sec"`
	Semester		uint64			`json:"semester" bson:"semester"`
	Year			uint64			`json:"year" bson:"year"`
	RecapID			string			`json:"recapId" bson:"recap_id"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt		time.Time		`json:"updateAt" bson:"update_at"`
	Reported		bool			`json:"reported" bson:"reported"`
}

// StatReview is model for storing stat on the review
type StatReview struct {
	How				float64			`json:"how" bson:"how"`
	Homework		float64			`json:"homework" bson:"homework"`
	Interest		float64			`json:"interest" bson:"interest"`
}

// ResReview is struct for body response to client 
type ResReview struct {
	ID				bson.ObjectId	`json:"reviewId" bson:"_id"`
	ClassID			string			`json:"classId" bson:"class_id"`
	Text			string			`json:"text" bson:"text"`
	Author			string			`json:"author" bson:"author"`
	Grade			string			`json:"grade" bson:"grade"`
	Clap			uint64			`json:"clap" bson:"clap"`
	Boo				uint64			`json:"boo" bson:"boo"`
	Stats			StatReview		`json:"stats" bson:"stats"`	
	ClassNameTH		string			`json:"classNameTH" bson:"class_name_th"`
	ClassNameEN		string			`json:"classNameEN" bson:"class_name_en"`
	Sec				uint64			`json:"sec" bson:"sec"`
	Semester		uint64			`json:"semester" bson:"semester"`
	Year			uint64			`json:"year" bson:"year"`
	RecapID			string			`json:"recapId" bson:"recap_id"`
	CreatedAt		time.Time		`json:"createdAt" bson:"created_at"`
	UpdateAt		time.Time		`json:"updateAt" bson:"update_at"`
	Reported		bool			`json:"reported" bson:"reported"`
}

type ReviewFilterField struct {
	// specific filtering
	ClassID			*string			`schema:"class_id" bson:"class_id" type:"match"`
	RecapID			*string			`schema:"recap_id" bson:"recap_id" type:"match"`
	Author			*string			`schema:"author" bson:"author" type:"match"`
	ClassNameTH		*string			`schema:"class_name_th" bson:"class_name_th" type:"match"`
	ClassNameEN		*string			`schema:"class_name_en" bson:"class_name_en" type:"match"`
	Grade			*string			`schema:"grade" bson:grade" type:"match"`
	Sec				*uint64			`schema:"sec" bson:"sec" type:"match"`
	Semester		*uint64			`schema:"semester" bson:"semester" type:"match"`
	Year			*uint64			`schema:"year" bson:"year" type:"match"`
	Reported		*bool			`schema:"reported" bson:"reported" type:"match"`

	// substring filtering
	Text			*string			`schema:"text" bson:"text" type:"text"`

	// length greater than filtering
	ClapGte			*uint64			`schema:"clap_gte" bson:"clap" type:"length" operation:"$gte"`
	BooGte			*uint64			`schema:"boo_gte" bson:"boo" type:"length" operation:"$gte"`
	HowGte			*float64		`schema:"stat_how_gte" bson:"how" type:"length" operation:"$gte"`
	HomeworkGte		*float64		`schema:"stat_homework_gte" bson:"homework" type:"length" operation:"$gte"`
	InterestGte		*float64		`schema:"stat_interest_gte" bson:"interest" type:"length" operation:"$gte"`

	// length less than filtering
	ClapLte			*uint64			`schema:"clap_lte" bson:"clap" type:"length" operation:"$lte"`
	BooLte			*uint64			`schema:"boo_lte" bson:"boo" type:"length" operation:"$lte"`
	HowLte			*float64		`schema:"stat_how_lte" bson:"how" type:"length" operation:"$lte"`
	HomeworkLte		*float64		`schema:"stat_homework_lte" bson:"homework" type:"length" operation:"$lte"`
	InterestLte		*float64		`schema:"stat_interest_lte" bson:"interest" type:"length" operation:"$lte"`

	// date greater than filtering
	CreatedAtGte	*string			`schema:"created_at_gte" bson:"created_at" type:"date" operation:"$gte"`
	UpdateAtGte		*string			`schema:"update_at_gte" bson:"update_at" type:"date" operation:"$gte"`

	// date less than filtering
	CreatedAtLte	*string			`schema:"created_at_lte" bson:"created_at" type:"date" operation:"$lte"`
	UpdateAtLte		*string			`schema:"update_at_lte" bson:"update_at" type:"date" operation:"$lte"`
}

// RDeleteReview is require struct for delete the review
type RDeleteReview struct {
	ID				string			`json:"reviewId" bson:"_id"`
	Auth			string			`json:"auth" bson:"auth"`
}

// SetBSON is function for filling default value when load value from mongo
func (review *ResReview) SetBSON(raw bson.Raw) (err error) {
	
	type my ResReview

	if err = raw.Unmarshal((*my)(review)); err != nil {
		return
	}
	review.Stats.How		=	review.Stats.How * 20
	review.Stats.Homework	=	review.Stats.Homework * 20
	review.Stats.Interest	=	review.Stats.Interest * 20

	return
	
}