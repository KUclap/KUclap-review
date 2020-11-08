package routes 

import (
	"encoding/json"
	"net/http"
	"time"
	"kuclap-review-api/src/models"
	"kuclap-review-api/src/helper"

	"github.com/gorilla/mux"
)

func IndexClassesHandler(r *mux.Router) {
	r.HandleFunc("/classes", AllClassesEndpoint).Methods("GET")
	r.HandleFunc("/class", InsertClassEndpoint).Methods("POST")
	r.HandleFunc("/class/{classid}", FindClassByClassIDEndpoint).Methods("GET")
	r.HandleFunc("/class/{classid}/stats", UpdateStatsEndPoint).Methods("PUT")
}

// GET list of classes
func AllClassesEndpoint(w http.ResponseWriter, r *http.Request) {
	classes, err := mgoDAO.FindAllClasses()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, classes)
}

func InsertClassEndpoint(w http.ResponseWriter, r * http.Request){
	defer r.Body.Close()
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mgoDAO.InsertClass(class); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusCreated, class)
}

// GET class by class_id
func FindClassByClassIDEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	class, err := mgoDAO.FindClassByClassID(params["classid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, class)
}

// PUT stats by class_id
func UpdateStatsEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var newStats models.StatClass
	var class models.Class	

	if err := json.NewDecoder(r.Body).Decode(&newStats); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := mgoDAO.FindClassByClassID(params["classid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Process new stats
	var oldStats = class.Stats
	newStats.How = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.How, newStats.How)
	newStats.Homework = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.Homework, newStats.Homework)
	newStats.Interest = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.Interest, newStats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour) 
	
	if err = mgoDAO.UpdateStatsClassByCreated(params["classid"], newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, newStats)
}