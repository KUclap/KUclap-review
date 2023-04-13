package models

import (
	"time"
)

// Class is model for create class
type Class struct {
	ClassID        string    `json:"classId" bson:"class_id"`
	NameTH         string    `json:"nameTh" bson:"name_th"`
	NameEN         string    `json:"nameEn" bson:"name_en"`
	Label          string    `json:"label" bson:"label"`
	Category       string    `json:"category" bson:"category"`
	Hours          string    `json:"hours" bson:"hours"`
	Unit           uint64    `json:"unit" bson:"unit"`
	NumberQuestion uint64    `json:"numberQuestion" bson:"number_questions"`
	NumberRecap    uint64    `json:"numberRecap" bson:"number_recaps"`
	NumberReviewer float64   `json:"numberReviewer" bson:"number_reviewer"`
	Stats          StatClass `json:"stats" bson:"stats"`
}

// StatClass is model for storing stat on the class
type StatClass struct {
	How      float64   `json:"how" bson:"how"`
	Homework float64   `json:"homework" bson:"homework"`
	Interest float64   `json:"interest" bson:"interest"`
	UpdateAt time.Time `json:"updateAt" bson:"update_at"`
}

// OldClass is struct for mapping old class structure
type OldClass struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ToDefault setting default value for class
func (class *Class) ToDefault() {
	if class.Category == "" {
		class.Category = "ยังไม่มีข้อมูล"
	}
}
