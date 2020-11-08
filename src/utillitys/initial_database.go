package utillitys

import (
	"encoding/json"
	"io/ioutil"
	"kuclap-review-api/src/models"
	"log"
	"kuclap-review-api/src/dao"
)

var reviewDAO = dao.SessionDAO{}

func insetClasstoDatabase(class models.Class) {
	if err := reviewDAO.InsertClass(class); err != nil {
		log.Println("err initial classes : ", err)
	}
}

func InitialClasses(){
	
	// classed.json is old data (KUnit version)
	// classedParsed.json is old data (KUclap version)

	var classes []models.Class
	data, err := ioutil.ReadFile("./classesParsed.json")
    if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(data, &classes)
	if err != nil {
        log.Println("error:", err)
	}
	for _, class := range classes {
		insetClasstoDatabase(class)
	}
}