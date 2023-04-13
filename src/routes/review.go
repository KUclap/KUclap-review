package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"kuclap-review-api/src/helper"
	"kuclap-review-api/src/models"
)

// IndexReviewHandler is index routing for review usecase
func IndexReviewHandler(r *mux.Router) {

	var prefixPath = "/review"

	r.HandleFunc(prefixPath, CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/last", LastReviewsEndPoint).Methods("GET")
	r.HandleFunc("/reviews/{classid}", AllReviewsByClassIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{reviewid}", FindReviewEndpoint).Methods("GET")
	r.HandleFunc(prefixPath+"/report/{reviewid}", UpdateReportByIDEndPoint).Methods("PUT")
	r.HandleFunc(prefixPath+"/clap/{reviewid}/{clap}", UpdateClapByIDEndPoint).Methods("PUT")
	r.HandleFunc(prefixPath+"/boo/{reviewid}/{boo}", UpdateBooByIDEndPoint).Methods("PUT")
	r.HandleFunc(prefixPath+"/{reviewid}", DeleteReviewByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/report", CreateReviewReportEndPoint).Methods("POST")

	// r.HandleFunc("/reviews/reported", FindReviewReportedEndpoint).Methods("GET")
	// r.HandleFunc("/reviews/{reviewid}", UpdateReviewEndPoint).Methods("PUT")

}

// LastReviewsEndPoint is GET list of reviews
// Read param on UrlQuery (eg. /last?offset=5 )
// Paging by query: page={number_page} offset={number_offset}
func LastReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")

	filter := new(models.ReviewFilterField)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(filter, r.URL.Query()); err != nil {
		log.Println("Error in decoding Query on request: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Query")
		return
	}

	querys := helper.CreateQueryFiltering(filter)

	reviews, err := repository.LastReviews(ctx, page, offset, querys)

	if err != nil {
		log.Println("Error in LastReviews DAO : ", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review offset")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, reviews)

}

// CreateReviewEndPoint is POST a new review
func CreateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var review models.Review

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := repository.FindClassByClassID(ctx, review.ClassID)
	if err != nil {
		log.Println("Error in FindClassByClassID DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var newStats models.StatClass

	var oldStats = class.Stats
	newStats.How = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.How, review.Stats.How)
	newStats.Homework = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.Homework, review.Stats.Homework)
	newStats.Interest = helper.GetNewStatsByCreated(class.NumberReviewer, oldStats.Interest, review.Stats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour)

	if err = repository.UpdateStatsClassByCreated(ctx, review.ClassID, newStats); err != nil {
		log.Println("Error in UpdateStatsClassByCreated DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	review.ClassNameTH = class.NameTH
	review.ClassNameEN = class.NameEN
	review.CreatedAt = time.Now().UTC().Add(7 * time.Hour)
	review.UpdateAt = review.CreatedAt
	review.ID = primitive.NewObjectID()

	if err := repository.Insert(ctx, review); err != nil {
		log.Println("Error in Insert DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, review)

}

// AllReviewsByClassIDEndPoint is GET list of reviews by class_id
func AllReviewsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")

	reviews, err := repository.FindReviewsByClassID(ctx, params["classid"], page, offset)
	if err != nil {
		log.Println("Error in FindReviewsByClassID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, reviews)

}

// FindReviewEndpoint is GET a review by its ID
func FindReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	review, err := repository.FindByID(ctx, params["reviewid"])
	if err != nil {
		log.Println("Error in FindByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, review)

}

// UpdateReportByIDEndPoint is GET report of reviews by class_id
func UpdateReportByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	defer r.Body.Close()

	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err := repository.UpdateReportByID(ctx, params["reviewid"], updateAt); err != nil {
		log.Println("Error in UpdateReportByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

}

// UpdateClapByIDEndPoint is PUT update clap by id
func UpdateClapByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)
	iclap, _ := strconv.ParseUint(params["clap"], 10, 32)

	if iclap > 25 {
		iclap = 25
	}

	if err := repository.UpdateClapByID(ctx, params["reviewid"], iclap, updateAt); err != nil {
		log.Println("Error in UpdateClapByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

}

// UpdateBooByIDEndPoint is PUT update boo by id
func UpdateBooByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	iboo, _ := strconv.ParseUint(params["boo"], 10, 32)

	if iboo > 25 {
		iboo = 25
	}

	if err := repository.UpdateBooByID(ctx, params["reviewid"], iboo, updateAt); err != nil {
		log.Println("Error in UpdateBooByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

}

// AllReviewsEndPoint is GET list of reviews
func AllReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	reviews, err := repository.FindAll(ctx)
	if err != nil {
		log.Println("Error in FindAll DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, reviews)

}

// DeleteReviewByIDEndPoint is DELETE an existing review
func DeleteReviewByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]

	review, err := repository.FindReviewAllPropertyByID(ctx, params["reviewid"])
	if err != nil {
		log.Println("Error in FindReviewAllPropertyByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}

	if review.Auth == reqAuth {

		// Delete the review.
		class, err := repository.FindClassByClassID(ctx, review.ClassID)
		if err != nil {
			log.Println("Error in FindClassByClassID DAO", err.Error())
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var newStats models.StatClass
		var oldStats = class.Stats

		if class.NumberReviewer == 1 {
			// Ignore NaN when we divide with zero
			newStats.How = 0
			newStats.Homework = 0
			newStats.Interest = 0
		} else {
			newStats.How = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.How, review.Stats.How)
			newStats.Homework = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Homework, review.Stats.Homework)
			newStats.Interest = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Interest, review.Stats.Interest)
		}

		newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour)

		if err := repository.DeleteByID(ctx, params["reviewid"]); err != nil {
			log.Println("Error in DeleteByID DAO", err.Error())
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err = repository.UpdateStatsClassByDeleted(ctx, review.ClassID, newStats); err != nil {
			log.Println("Error in UpdateStatsClassByDeleted DAO", err.Error())
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if review.RecapID != "" {
			// Delete recap on the review.
			if err := repository.DeleteRecapByID(ctx, review.RecapID); err != nil {
				log.Println("Error in DeleteRecapByID DAO", err.Error())
				helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			updateAt := time.Now().UTC().Add(7 * time.Hour)
			if err = repository.UpdateNumberRecapByClassID(ctx, review.ClassID, -1, updateAt); err != nil {
				log.Println("Error in UpdateNumberRecapByClassID DAO", err.Error())
				helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

	} else {
		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}

}

// CreateReportEndPoint is POST create report for the review
func CreateReviewReportEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var report models.Report
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err := repository.UpdateReportByID(ctx, report.ReviewID, updateAt); err != nil {
		log.Println("Error in UpdateReportByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	report.CreatedAt = time.Now().UTC().Add(7 * time.Hour)

	if err := repository.InsertReport(ctx, report); err != nil {
		log.Println("Error in InsertReport DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, report)

}
