package routes 

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"kuclap-review-api/src/models"
	"kuclap-review-api/src/helper"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// IndexReviewHandler is index routing for review usecase
func IndexReviewHandler(r *mux.Router) {
	r.HandleFunc("/reviews/last", LastReviewsEndPoint).Methods("GET")
	r.HandleFunc("/review", CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/{classid}", AllReviewsByClassIDEndPoint).Methods("GET")
	r.HandleFunc("/review/{reviewid}", FindReviewEndpoint).Methods("GET")
	r.HandleFunc("/review/report/{reviewid}", UpdateReportByIDEndPoint).Methods("PUT")
	r.HandleFunc("/review/clap/{reviewid}/{clap}", UpdateClapByIDEndPoint).Methods("PUT")
	r.HandleFunc("/review/boo/{reviewid}/{boo}", UpdateBooByIDEndPoint).Methods("PUT")
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc("/review/{reviewid}", DeleteReviewByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/report", CreateReportEndPoint).Methods("POST")
	// r.HandleFunc("/reviews/reported", FindReviewReportedEndpoint).Methods("GET")
	// r.HandleFunc("/reviews/{reviewid}", UpdateReviewEndPoint).Methods("PUT")
}

// LastReviewsEndPoint is GET list of reviews 
// Read param on UrlQuery (eg. /last?offset=5 )
// Paging by query: page={number_page} offset={number_offset}
func LastReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := mgoDAO.LastReviews(page ,offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review offset")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, reviews)
}

// CreateReviewEndPoint is POST a new review
func CreateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var class models.Class	
	var review models.Review

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := mgoDAO.FindClassByClassID(review.ClassID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	var oldStats = class.Stats
	var newStats models.StatClass
	newStats.How = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.How, review.Stats.How)
	newStats.Homework = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.Homework, review.Stats.Homework)
	newStats.Interest = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.Interest, review.Stats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour) 

	if err = mgoDAO.UpdateStatsClassByCreated(review.ClassID, newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	review.ClassNameTH = class.NameTH
	review.ClassNameEN = class.NameEN
	review.CreatedAt = time.Now().UTC().Add(7 * time.Hour) 
	review.UpdateAt = review.CreatedAt
	review.ID = bson.NewObjectId()

	if err := mgoDAO.Insert(review); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusCreated, review)
}

// AllReviewsByClassIDEndPoint is GET list of reviews by class_id
func AllReviewsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := mgoDAO.FindReviewsByClassID(params["classid"], page, offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review classid")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, reviews)
}


// FindReviewEndpoint is GET a review by its ID
func FindReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //  param on endpoint
	review, err := mgoDAO.FindById(params["reviewid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, review)
}

// UpdateReportByIDEndPoint is GET report of reviews by class_id
func UpdateReportByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err := mgoDAO.UpdateReportById(params["reviewid"], updateAt); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}
	helper.RespondWithJson(w, http.StatusOK,  map[string]string{"result": "success"})
}

// UpdateClapByIDEndPoint is PUT update clap by id
func UpdateClapByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)
	iclap, _ := strconv.ParseUint(params["clap"],10 ,32)
	if err := mgoDAO.UpdateClapById(params["reviewid"], iclap, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// UpdateBooByIDEndPoint is PUT update boo by id
func UpdateBooByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)
	iboo, _ := strconv.ParseUint(params["boo"],10, 32)
	if err := mgoDAO.UpdateBooById(params["reviewid"], iboo, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// AllReviewsEndPoint is GET list of reviews
func AllReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	reviews, err := mgoDAO.FindAll()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, reviews)
}

// DeleteReviewByIDEndPoint is DELETE an existing review
func DeleteReviewByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]
	var class models.Class
	
	review, err := mgoDAO.FindById(params["reviewid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}
	
	class, err = mgoDAO.FindClassByClassID(review.ClassID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var newStats models.StatClass
	var oldStats = class.Stats
	if class.NumberReviewer ==  1 {
		// Ignore NaN when we divide with zero
		newStats.How = 0
		newStats.Homework = 0
		newStats.Interest = 0
	} else {
		newStats.How = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.How, review.Stats.How)
		newStats.Homework = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Homework,  review.Stats.Homework)
		newStats.Interest = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Interest,  review.Stats.Interest)
	}
	
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour) 

	if err = mgoDAO.UpdateStatsClassByDeleted(review.ClassID, newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if review.Auth == reqAuth {
		if err := mgoDAO.DeleteById(params["reviewid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	} else {
		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}
}

// CreateReportEndPoint is POST create report for the review
func CreateReportEndPoint(w http.ResponseWriter, r *http.Request) { 
	defer r.Body.Close()
	var report models.Report

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err := mgoDAO.UpdateReportById(report.ReviewID, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	report.CreatedAt = time.Now().UTC().Add(7 * time.Hour)
	if err := mgoDAO.InsertReport(report); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, report)
}