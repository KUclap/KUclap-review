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

	"gopkg.in/mgo.v2/bson"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
    "github.com/marsDev31/kuclap-backend/api/config"
	"github.com/marsDev31/kuclap-backend/api/middleware"
	"github.com/marsDev31/kuclap-backend/api/dao"
	"github.com/marsDev31/kuclap-backend/api/models"
)

var limiter = middleware.NewIPRateLimiter(200, 10)
var mcf = config.Config{}
var mdao = dao.SessionDAO{}

func getNewStats(oldN float64, oldstat float64, newStats float64) float64 {
	return ((newStats / 5 * 100) + (oldstat * oldN)) / (oldN + 1)
}

// GET class by class_id
func FindClassByClassIdEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	class, err := mdao.FindClassByClassId(params["classid"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, class)
}

// PUT update clap by id
func UpdateClapByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateAt := time.Now().UTC()
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
	updateAt := time.Now().UTC()
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
	updateAt := time.Now().UTC()

	if err := mdao.UpdateReportById(params["reviewid"], updateAt); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}
	respondWithJson(w, http.StatusOK,  map[string]string{"result": "success"})
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

	class, err := mdao.FindClassByClassId(params["classid"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Process new stats
	var oldStats = class.Stats
	newStats.How = getNewStats(class.NumberReviewer, oldStats.How, newStats.How)
	newStats.Homework = getNewStats(class.NumberReviewer, oldStats.Homework, newStats.Homework)
	newStats.Interest = getNewStats(class.NumberReviewer, oldStats.Interest, newStats.Interest)
	newStats.UpdateAt = time.Now().UTC() 
	
	if err = mdao.UpdateStatsClass(params["classid"], newStats); err != nil {
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
func AllReviewsByClassIdEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	reviews, err := mdao.FindReviewsByClassId(params["classid"], page, offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review classid")
		return
	}
	respondWithJson(w, http.StatusOK, reviews)
}

// GET list of reviews // Read param on UrlQuery (eg. /last?offset=5 )
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
	var review models.Review

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	review.CreatedAt = time.Now().UTC() // Parse UTC to GTM+7 Thailand's timezone.
	review.UpdateAt = review.CreatedAt
	review.ID = bson.NewObjectId()
	review.Auth = getRemoteAddr(r)

	if err := mdao.Insert(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, review)
}

// GET a reviews by its ID
func FindReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //  param on endpoint
	review, err := mdao.FindById(params["reviewid"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review Id or haven't your id in DB")
		return
	}
	respondWithJson(w, http.StatusOK, review)
}

// DELETE an existing review
func DeleteReviewByIdEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	if err := mdao.DeleteById(params["reviewid"]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	fmt.Println("Initial service..")
	// Conection on database
	mcf.Read()
	mdao.Server = goDotEnvVariable("SERVER")
	mdao.Database = mcf.Database
	mdao.Connect() 
	// initialClasses()
}

// Define HTTP request routes
func main() {
	
	port := goDotEnvVariable("PORT")
	fmt.Println("Starting services.")
	headersOk := handlers.AllowedHeaders([]string{"Origin", "Authorization", "Content-Type"})
	exposeOk := handlers.ExposedHeaders([]string{""})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/classes", AllClassesEndpoint).Methods("GET")
	r.HandleFunc("/class", InsertClassEndpoint).Methods("POST")
	r.HandleFunc("/class/{classid}", FindClassByClassIdEndpoint).Methods("GET")
	r.HandleFunc("/class/{classid}/stats", UpdateStatsEndPoint).Methods("PUT")
	r.HandleFunc("/reviews/last", LastReviewsEndPoint).Methods("GET")
	r.HandleFunc("/review", CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/{classid}", AllReviewsByClassIdEndPoint).Methods("GET")
	r.HandleFunc("/review/{reviewid}", FindReviewEndpoint).Methods("GET")
	r.HandleFunc("/review/report/{reviewid}", UpdateReportByIdEndPoint).Methods("PUT")
	r.HandleFunc("/review/clap/{reviewid}/{clap}", UpdateClapByIdEndPoint).Methods("PUT")
	r.HandleFunc("/review/boo/{reviewid}/{boo}", UpdateBooByIdEndPoint).Methods("PUT")
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc("/review/{reviewid}", DeleteReviewByIdEndPoint).Methods("DELETE")
	// r.HandleFunc("/reviews/reported", FindReviewReportedEndpoint).Methods("GET")
	// r.HandleFunc("/reviews/{reviewid}", UpdateReviewEndPoint).Methods("PUT")
	
	

	if err := http.ListenAndServe(":" + port, limitMiddleware(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r))); err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port " + port)
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
	fmt.Fprint(w, "Hi Developers!, Welcome to KUclap services: PRs welcome @https://github.com/marsDev31/kuclap-backend.")
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

func getRemoteAddr(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		fmt.Println(forwarded)
		return forwarded
	}
	return r.RemoteAddr
}

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
		fmt.Println("err initial classes : ", err)
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
        fmt.Println("error:", err)
	}
	for _, class := range classes {
		insetClasstoDatabase(class)
	}
}

// GET class by class_id
// func FindReviewReportedEndpoint(w http.ResponseWriter, r *http.Request) {
// 	reviews, err := mdao.FindReviewsReported()
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, reviews)
// }

// PUT update an existing review
// func UpdateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	params := mux.Vars(r)
// 	var review models.Review

// 	review.UpdateAt = time.Now().UTC()
// 	review.ID = bson.ObjectIdHex(params["reviewid"])

// 	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	if err := mdao.Update(review); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
// }