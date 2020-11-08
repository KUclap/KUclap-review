package main

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"

    "github.com/KUclap/KUclap-review/src/config"
	"github.com/KUclap/KUclap-review/src/middleware"
	"github.com/KUclap/KUclap-review/src/dao"
	"github.com/KUclap/KUclap-review/src/models"
)

var limiter = middleware.NewIPRateLimiter(200, 10)
var server_config = config.Config{}
var mdao = dao.SessionDAO{}


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
	class, err := mdao.FindClassByClassID(params["classid"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, class)
}

// PUT update clap by id
func UpdateClapByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7*time.Hour)
	iclap, _ := strconv.ParseUint(params["clap"],10 ,32)
	if err := mdao.UpdateClapById(params["reviewid"], iclap, updateAt); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// PUT update boo by id
func UpdateBooByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7*time.Hour)
	iboo, _ := strconv.ParseUint(params["boo"],10, 32)
	if err := mdao.UpdateBooById(params["reviewid"], iboo, updateAt); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// GET report of reviews by class_id
func UpdateReportByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	updateAt := time.Now().UTC().Add(7*time.Hour)

	if err := mdao.UpdateReportById(params["reviewid"], updateAt); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}
	respondWithJson(w, http.StatusOK,  map[string]string{"result": "success"})
}

func CreateReportEndPoint(w http.ResponseWriter, r *http.Request) { 
	defer r.Body.Close()
	var report models.Report

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7*time.Hour)
	if err := mdao.UpdateReportById(report.ReviewID, updateAt); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	report.CreatedAt = time.Now().UTC().Add(7*time.Hour)
	if err := mdao.InsertReport(report); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, report)
}


// PUT stats by class_id
func UpdateStatsEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var newStats models.StatClass
	var class models.Class	

	if err := json.NewDecoder(r.Body).Decode(&newStats); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := mdao.FindClassByClassID(params["classid"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Process new stats
	var oldStats = class.Stats
	newStats.How = getNewStatsByCreated(class.NumberReviewer, oldStats.How, newStats.How)
	newStats.Homework = getNewStatsByCreated(class.NumberReviewer, oldStats.Homework, newStats.Homework)
	newStats.Interest = getNewStatsByCreated(class.NumberReviewer, oldStats.Interest, newStats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7*time.Hour) 
	
	if err = mdao.UpdateStatsClassByCreated(params["classid"], newStats); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, newStats)
}

// GET list of classes
func AllClassesEndpoint(w http.ResponseWriter, r *http.Request) {
	classes, err := mdao.FindAllClasses()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, classes)
}

// Create class 
func InsertClassEndpoint(w http.ResponseWriter, r * http.Request){
	defer r.Body.Close()
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mdao.InsertClass(class); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, class)
}

// GET list of reviews by class_id
func AllReviewsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := mdao.FindReviewsByClassID(params["classid"], page, offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review classid")
		return
	}
	respondWithJson(w, http.StatusOK, reviews)
}

// GET list of reviews 
// Read param on UrlQuery (eg. /last?offset=5 )
// Paging by query: page={number_page} offset={number_offset}
func LastReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := mdao.LastReviews(page ,offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review offset")
		return
	}
	respondWithJson(w, http.StatusOK, reviews)
}

// GET list of reviews
func AllReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	reviews, err := mdao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, reviews)
}

// POST a new review
func CreateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var class models.Class	
	var review models.Review

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := mdao.FindClassByClassID(review.ClassID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	var oldStats = class.Stats
	var newStats models.StatClass
	newStats.How = getNewStatsByCreated(class.NumberReviewer, oldStats.How, review.Stats.How)
	newStats.Homework = getNewStatsByCreated(class.NumberReviewer, oldStats.Homework, review.Stats.Homework)
	newStats.Interest = getNewStatsByCreated(class.NumberReviewer, oldStats.Interest, review.Stats.Interest)
	newStats.UpdateAt = time.Now().UTC().Add(7*time.Hour) 

	if err = mdao.UpdateStatsClassByCreated(review.ClassID, newStats); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	review.ClassNameTH = class.NameTH
	review.ClassNameEN = class.NameEN
	review.CreatedAt = time.Now().UTC().Add(7*time.Hour) 
	review.UpdateAt = review.CreatedAt
	review.ID = bson.NewObjectId()

	if err := mdao.Insert(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, review)
}

// GET a review by its ID
func FindReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //  param on endpoint
	review, err := mdao.FindById(params["reviewid"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}
	respondWithJson(w, http.StatusOK, review)
}

// DELETE an existing review
func DeleteReviewByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]
	var class models.Class
	
	review, err := mdao.FindById(params["reviewid"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}
	
	class, err = mdao.FindClassByClassID(review.ClassID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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
	
	newStats.UpdateAt = time.Now().UTC().Add(7*time.Hour) 

	if err = mdao.UpdateStatsClassByDeleted(review.ClassID, newStats); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if review.Auth == reqAuth {
		if err := mdao.DeleteById(params["reviewid"]); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	} else {
		respondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}
}

// Parse the server_configuration file 'server_config.toml', and establish a connection to DB
func init() {
	log.Println("Initial service... 🔧")
	// Conection on database
	server_config.Read()
	mdao.Server = goDotEnvVariable("DB_SERVER")
	mdao.Database = server_config.Database
	mdao.Connect() 
	// initialClasses()
}

// Define HTTP request routes
func main() {
	log.Println("Starting server... 🤤")
	port := goDotEnvVariable("PORT")
	origin := goDotEnvVariable("ORIGIN_ALLOWED")
	
	
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

	if err := http.ListenAndServe(":" + port, limitMiddleware(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r))); err != nil {
		log.Fatal(err)
	}
	
}

// Rate Limit base on IP (r = tokens per second, b = maximum burst size of b events)
func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ROOT request
func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi Developers!, Welcome to KUclap services: PRs welcome @https://github.com/KUclap/KUclap-review.")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// func getRemoteAddr(r *http.Request) string {
// 	forwarded := r.Header.Get("X-FORWARDED-FOR")
// 	if forwarded != "" {
// 		log.Println(forwarded)
// 		return forwarded
// 	}
// 	return r.RemoteAddr
// }

//  Getter environment from .env.
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func insetClasstoDatabase(class models.Class) {
	if err := mdao.InsertClass(class); err != nil {
		log.Println("err initial classes : ", err)
	}
}

func initialClasses(){
	// classed.json is old data (KUnit version)
	// classedParsed.json is old data (KUclap version)
	var classes []models.Class
	data, err := ioutil.ReadFile("./classesParsed.json")
    if err != nil {
      fmt.Print(err)
	}
	err = json.Unmarshal(data, &classes)
	if err != nil {
        log.Println("error:", err)
	}
	for _, class := range classes {
		insetClasstoDatabase(class)
	}
}
