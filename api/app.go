package main

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
	"os"

	"gopkg.in/mgo.v2/bson"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"github.com/marsDev31/kuclap-backend/api/config"
	"github.com/marsDev31/kuclap-backend/api/dao"
	"github.com/marsDev31/kuclap-backend/api/models"
)

var limiter = NewIPRateLimiter(1, 5)
var mcf = config.Config{}
var mdao = dao.SessionDAO{}

func UpdateStatsEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var stats models.StatClass
	// Get payload as json format
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	stats.UpdateAt = time.Now().UTC().Add(7 *time.Hour) 
	if err := mdao.UpdateStatsClass(params["classid"], stats); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// GET list of classes
func AllClassesEndpoint(w http.ResponseWriter, r *http.Request) {
	classes, err := mdao.GetAllClasses()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, classes)
}

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


// GET list of reviews // Read param on UrlQuery (eg. /last?offset=5 )
func LastReviewsEndPoint(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	reviews, err := mdao.LastReviews(offset)
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

// GET a reviews by its ID
func FindReviewEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //  param on endpoint
	review, err := mdao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}
	respondWithJson(w, http.StatusOK, review)
}

// POST a new review
func CreateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	review.CreatedAt = time.Now().UTC().Add(7 *time.Hour) // Parse UTC to GTM+7 Thailand's timezone.
	review.UpdateAt = review.CreatedAt
	review.ID = bson.NewObjectId()
	review.Auth = getRemoteAddr(r)

	if err := mdao.Insert(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, review)
}

// PUT update an existing review
func UpdateReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var review models.Review

	review.UpdateAt = time.Now().UTC().Add(7 *time.Hour)
	review.ID = bson.ObjectIdHex(params["id"])

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mdao.Update(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing review
func DeleteReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mdao.Delete(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
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

// ROOT request
func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi Developers!, Welcome to KUclap services: PRs welcome @https://github.com/marsDev31/kuclap-backend.")
}

func getRemoteAddr(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
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

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	fmt.Println("Initial service..")

	// Conection on database
	mcf.Read()
	mdao.Server = goDotEnvVariable("SERVER")
	mdao.Database = mcf.Database
	mdao.Connect() 
	initialClasses()
	
}

// Define HTTP request routes
func main() {
	
	port := goDotEnvVariable("PORT")
	fmt.Println("Starting services.")
	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/classes", AllClassesEndpoint).Methods("GET")
	r.HandleFunc("/classes/{classid}/stats", UpdateStatsEndPoint).Methods("PUT")
	r.HandleFunc("/last", LastReviewsEndPoint).Methods("GET")
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc("/reviews", CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/{id}", UpdateReviewEndPoint).Methods("PUT")
	r.HandleFunc("/reviews", DeleteReviewEndPoint).Methods("DELETE")
	r.HandleFunc("/reviews/{id}", FindReviewEndpoint).Methods("GET")

	if err := http.ListenAndServe(":" + port, limitMiddleware(r)); err != nil {
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