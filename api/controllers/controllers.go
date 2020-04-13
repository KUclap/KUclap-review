package controllers
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
	"fmt"
	"time"
	"os"

	"gopkg.in/mgo.v2"
	"github.com/marsDev31/kuclap-backend/api/models"
)

var reviews *mgo.Collection

func ConnectDB() *mgo.Collection{
	// Connect to mongo
	session, err := mgo.Dial("mongo:27017")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		fmt.Println("mongo err")
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	
	reviews := session.DB("ku-clap").C("reviews")
	fmt.Println("Conected DB")
	return reviews
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	review := &models.Review{}
	err = json.Unmarshal(data, review)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	review.CreatedAt = time.Now().UTC()
	
	reviews := ConnectDB()
	if err := reviews.Insert(review); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSON(w, review)
}

func GetReviews(w http.ResponseWriter, r *http.Request) {
	result := []models.Review{}
	reviews := ConnectDB()
	
	if err := reviews.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}