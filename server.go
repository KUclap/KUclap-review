package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
	"strconv"
	"strings"
	"kuclap-review-api/src/helper"
    "kuclap-review-api/src/config"
	"kuclap-review-api/src/middleware"
	"kuclap-review-api/src/dao"
	"kuclap-review-api/src/models"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)


var serverConfig = config.Config{}
var reviewDAO = dao.SessionDAO{}


func getNewStatsByCreated(oldN float64, oldstat float64, newStats float64) float64 {
	return ((newStats / 5 * 100) + (oldstat * oldN)) / (oldN + 1)
}


func getNewStatsByDeleted(oldN float64, oldstat float64, newStats float64) float64 {
	if oldN - 1 <= 0 {
		return oldstat
	} 
		return ((oldstat * oldN) - (newStats / 5 * 100) ) / (oldN - 1)
	
}

// GET class by class_id
func FindClassByClassIDEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	class, err := reviewDAO.FindClassByClassID(params["classid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, class)
}

// PUT update clap by id
func UpdateClapByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)
	iclap, _ := strconv.ParseUint(params["clap"],10 ,32)
	if err := reviewDAO.UpdateClapById(params["reviewid"], iclap, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// PUT update boo by id
func UpdateBooByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)
	iboo, _ := strconv.ParseUint(params["boo"],10, 32)
	if err := reviewDAO.UpdateBooById(params["reviewid"], iboo, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// GET report of reviews by class_id
func UpdateReportByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err := reviewDAO.UpdateReportById(params["reviewid"], updateAt); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}
	helper.RespondWithJson(w, http.StatusOK,  map[string]string{"result": "success"})
}

func CreateReportEndPoint(w http.ResponseWriter, r *http.Request) { 
	defer r.Body.Close()
	var report models.Report

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err := reviewDAO.UpdateReportById(report.ReviewID, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	report.CreatedAt = time.Now().UTC().Add(7 * time.Hour)
	if err := reviewDAO.InsertReport(report); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, report)
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

	class, err := reviewDAO.FindClassByClassID(params["classid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Process new stats
	var oldStats = class.Stats
	newStats.How = getNewStatsByCreated(class.NumberReviewer, oldStats.How, newStats.How)
	newStats.Homework = getNewStatsByCreated(class.NumberReviewer, oldStats.Homework, newStats.Homework)
	newStats.Interest = getNewStatsByCreated(class.NumberReviewer, oldStats.Interest, newStats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour) 
	
	if err = reviewDAO.UpdateStatsClassByCreated(params["classid"], newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, newStats)
}

// GET list of classes
func AllClassesEndpoint(w http.ResponseWriter, r *http.Request) {
	classes, err := reviewDAO.FindAllClasses()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, classes)
}

// Create class 
func InsertClassEndpoint(w http.ResponseWriter, r * http.Request){
	defer r.Body.Close()
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := reviewDAO.InsertClass(class); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusCreated, class)
}

// GET list of reviews by class_id
func AllReviewsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := reviewDAO.FindReviewsByClassID(params["classid"], page, offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review classid")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, reviews)
}

// GET list of reviews 
// Read param on UrlQuery (eg. /last?offset=5 )
// Paging by query: page={number_page} offset={number_offset}
func LastReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := reviewDAO.LastReviews(page ,offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review offset")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, reviews)
}

// GET list of reviews
func AllReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	reviews, err := reviewDAO.FindAll()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusOK, reviews)
}

// POST a new review
func CreateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var class models.Class	
	var review models.Review

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := reviewDAO.FindClassByClassID(review.ClassID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	var oldStats = class.Stats
	var newStats models.StatClass
	newStats.How = getNewStatsByCreated(class.NumberReviewer, oldStats.How, review.Stats.How)
	newStats.Homework = getNewStatsByCreated(class.NumberReviewer, oldStats.Homework, review.Stats.Homework)
	newStats.Interest = getNewStatsByCreated(class.NumberReviewer, oldStats.Interest, review.Stats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour) 

	if err = reviewDAO.UpdateStatsClassByCreated(review.ClassID, newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	review.ClassNameTH = class.NameTH
	review.ClassNameEN = class.NameEN
	review.CreatedAt = time.Now().UTC().Add(7 * time.Hour) 
	review.UpdateAt = review.CreatedAt
	review.ID = bson.NewObjectId()

	if err := reviewDAO.Insert(review); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJson(w, http.StatusCreated, review)
}

// GET a review by its ID
func FindReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //  param on endpoint
	review, err := reviewDAO.FindById(params["reviewid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}
	helper.RespondWithJson(w, http.StatusOK, review)
}

// DELETE an existing review
func DeleteReviewByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]
	var class models.Class
	
	review, err := reviewDAO.FindById(params["reviewid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}
	
	class, err = reviewDAO.FindClassByClassID(review.ClassID)
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
		newStats.How = getNewStatsByDeleted(class.NumberReviewer, oldStats.How, review.Stats.How)
		newStats.Homework = getNewStatsByDeleted(class.NumberReviewer, oldStats.Homework,  review.Stats.Homework)
		newStats.Interest = getNewStatsByDeleted(class.NumberReviewer, oldStats.Interest,  review.Stats.Interest)
	}
	
	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour) 

	if err = reviewDAO.UpdateStatsClassByDeleted(review.ClassID, newStats); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if review.Auth == reqAuth {
		if err := reviewDAO.DeleteById(params["reviewid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	} else {
		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}
}

// ROOT request
func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "https://github.com/kuclap 😎")
}

// Parse the serverConfiguration file 'serverConfig.toml', and establish a connection to DB
func init() {
	log.Println("Initial service... 🔧")
	// Conection on database
	serverConfig.Read()
	reviewDAO.Server = helper.GetENV("DB_SERVER")
	reviewDAO.Database = serverConfig.Database
	reviewDAO.Connect() 
	// initialClasses()
}

// Define HTTP request routes
func main() {
	log.Println("Starting server... 🤤")
	port := helper.GetENV("PORT")
	origin := helper.GetENV("ORIGIN_ALLOWED")
	
	
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Authorization", "Content-Type"})
	exposeOk := handlers.ExposedHeaders([]string{""})
	originsOk := handlers.AllowedOrigins([]string{origin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/classes", AllClassesEndpoint).Methods("GET")
	r.HandleFunc("/class", InsertClassEndpoint).Methods("POST")
	r.HandleFunc("/class/{classid}", FindClassByClassIDEndpoint).Methods("GET")
	r.HandleFunc("/class/{classid}/stats", UpdateStatsEndPoint).Methods("PUT")
	r.HandleFunc("/reviews/last", LastReviewsEndPoint).Methods("GET")
	r.HandleFunc("/review", CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/{classid}", AllReviewsByClassIDEndPoint).Methods("GET")
	r.HandleFunc("/review/{reviewid}", FindReviewEndpoint).Methods("GET")
	r.HandleFunc("/review/report/{reviewid}", UpdateReportByIdEndPoint).Methods("PUT")
	r.HandleFunc("/review/clap/{reviewid}/{clap}", UpdateClapByIdEndPoint).Methods("PUT")
	r.HandleFunc("/review/boo/{reviewid}/{boo}", UpdateBooByIdEndPoint).Methods("PUT")
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc("/review/{reviewid}", DeleteReviewByIdEndPoint).Methods("DELETE")
	r.HandleFunc("/report", CreateReportEndPoint).Methods("POST")
	// r.HandleFunc("/reviews/reported", FindReviewReportedEndpoint).Methods("GET")
	// r.HandleFunc("/reviews/{reviewid}", UpdateReviewEndPoint).Methods("PUT")
	
	log.Println("Server listening on port " + port + " 🚀")

	if err := http.ListenAndServe(":" + port, middleware.LimitMiddleware(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r))); err != nil {
		log.Fatal(err)
	}
	
}