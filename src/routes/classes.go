package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"kuclap-review-api/src/helper"
	"kuclap-review-api/src/models"
)

// IndexClassesHandler is index routing for class usecase
func IndexClassesHandler(r *mux.Router) {

	var prefixPath = "/class"

	r.HandleFunc(prefixPath, InsertClassEndpoint).Methods("POST")
	r.HandleFunc(prefixPath+"/{classid}", FindClassByClassIDEndpoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{classid}/stats", UpdateStatsEndPoint).Methods("PUT")
	r.HandleFunc("/classes", AllClassesEndpoint).Methods("GET")

}

// AllClassesEndpoint is GET list of classes
func AllClassesEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	classes, err := repository.FindAllClasses(ctx)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, classes)

}

// InsertClassEndpoint is POST insert class.
func InsertClassEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var class models.Class
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := repository.InsertClass(ctx, class); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, class)

}

// FindClassByClassIDEndpoint is GET class by class_id
func FindClassByClassIDEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	params := mux.Vars(r)
	class, err := repository.FindClassByClassID(ctx, params["classid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, class)

}

// UpdateStatsEndPoint is PUT stats by class_id
func UpdateStatsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var newStats models.StatClass

	defer r.Body.Close()
	params := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&newStats); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := repository.FindClassByClassID(ctx, params["classid"])
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

	if err = repository.UpdateStatsClassByCreated(ctx, params["classid"], newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, newStats)

}
