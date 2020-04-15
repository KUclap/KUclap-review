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


var mcf = config.Config{}
var mdao = dao.ReviewsDAO{}
var classes []models.Classes

// ROOT request
func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi Developers!, Welcome to KUclap services: PRs welcome @https://github.com/marsDev31/kuclap-backend.")
}

// GET list of classes
func AllClassesEndpoint(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, classes)
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

//  Getter environment from .env.
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
  }

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	fmt.Println("Initial service..")

	// Conection on database
	mcf.Read()
	mdao.Server = goDotEnvVariable("SERVER")
	mdao.Database = mcf.Database
	mdao.Connect() 
	
	// Read json file
	// classed.json is old data (KUnit version)
	// classedParsed.json is old data (KUclap version)
	data, err := ioutil.ReadFile("./classesParsed.json")
    if err != nil {
      fmt.Print(err)
	}
	
	err = json.Unmarshal(data, &classes)
	if err != nil {
        fmt.Println("error:", err)
    }
	
}

// Define HTTP request routes
func main() {
	
	port := goDotEnvVariable("PORT")
	fmt.Println("Starting services.")
	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/classes", AllClassesEndpoint).Methods("GET")
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc("/reviews", CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/{id}", UpdateReviewEndPoint).Methods("PUT")
	r.HandleFunc("/reviews", DeleteReviewEndPoint).Methods("DELETE")
	r.HandleFunc("/reviews/{id}", FindReviewEndpoint).Methods("GET")

	if err := http.ListenAndServe(":" + port, r); err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port " + port)
	
	
}
