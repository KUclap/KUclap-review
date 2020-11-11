package models
import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Class is model for create class
type Class struct {
	ClassID			string			`json:"classId" bson:"class_id"`
	NameTH			string			`json:"nameTh" bson:"name_th"`
	NameEN			string			`json:"nameEn" bson:"name_en"`
	Label			string			`json:"label" bson:"label"`
	Category		string			`json:"category" bson:"category"`
	Hours			string			`json:"hours" bson:"hours"`
	Unit			uint64			`json:"unit" bson:"unit"`
	NumberQuestion	uint64			`json:"numberQuestion" bson:"number_questions"`
	NumberReviewer	float64			`json:"numberReviewer" bson:"number_reviewer"`
	Stats			StatClass		`json:"stats" bson:"stats"`	
}

// StatClass is model for storing stat on the class
type StatClass struct {
	How			float64		`json:"how" bson:"how"`
	Homework	float64		`json:"homework" bson:"homework"`
	Interest	float64		`json:"interest" bson:"interest"`
	UpdateAt	time.Time	`json:"updateAt" bson:"update_at"`
}

// OldClass is struct for mapping old class structure
type OldClass struct {
	Value	string	`json:"value"`
	Label	string	`json:"label"`
}

// GetBSON is function for filling default value when save value on mongo
func (class *Class) GetBSON() (interface{}, error) {
    if class.Category == "" {
		class.Category = "ยังไม่มีข้อมูล" 
    }
    type my *Class
    return my(class), nil
}

// SetBSON is function for filling default value when load value from mongo
func (class *Class) SetBSON(raw bson.Raw) (err error) {
	type my Class
    if err = raw.Unmarshal((*my)(class)); err != nil {
        return
	}
	if class.Category == "" {
		class.Category = "ยังไม่มีข้อมูล" 
	}
    return
}