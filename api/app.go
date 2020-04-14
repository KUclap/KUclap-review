package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/marsDev31/kuclap-backend/api/config"
	"github.com/marsDev31/kuclap-backend/api/dao"
	"github.com/marsDev31/kuclap-backend/api/models"
)

var mcf = config.Config{}
var mdao = dao.ReviewsDAO{}

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
	review.CreatedAt = time.Now().UTC()
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

	review.UpdateAt = time.Now().UTC()
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

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	mcf.Read() // read config
	mdao.Server = mcf.Server
	mdao.Database = mcf.Database
	mdao.Connect() // conecting database
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reviews", AllReviewsEndPoint).Methods("GET")
	r.HandleFunc("/reviews", CreateReviewEndPoint).Methods("POST")
	r.HandleFunc("/reviews/{id}", UpdateReviewEndPoint).Methods("PUT")
	r.HandleFunc("/reviews", DeleteReviewEndPoint).Methods("DELETE")
	r.HandleFunc("/reviews/{id}", FindReviewEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
