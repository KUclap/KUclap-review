// package main

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/rs/cors"

// 	"github.com/gorilla/mux"
// 	"gopkg.in/mgo.v2"
// )

// type Review struct {
// 	Text      string    `json:"text" bson:"text"`
// 	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
// }

// var reviews *mgo.Collection

// func main() {
// 	// Connect to mongo
// 	session, err := mgo.Dial("mongo:27017")
// 	if err != nil {
// 		log.Fatalln(err)
// 		log.Fatalln("mongo err")
// 		os.Exit(1)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)

// 	// Get reviews collection
// 	reviews = session.DB("KUCLAP").C("reviews")

// 	// Set up routes
// 	r := mux.NewRouter()
// 	r.HandleFunc("/create_review", createReview).
// 		Methods("POST")
// 	r.HandleFunc("/get_reviews", readReviews).
// 		Methods("GET")

// 	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
// 	log.Println("Listening on port 8080...")
// }

// func createReview(w http.ResponseWriter, r *http.Request) {
// 	// Read body
// 	data, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		responseError(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Read review
// 	review := &Review{}
// 	err = json.Unmarshal(data, review)
// 	if err != nil {
// 		responseError(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	review.CreatedAt = time.Now().UTC()

// 	// Insert new review
// 	if err := reviews.Insert(review); err != nil {
// 		responseError(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	responseJSON(w, review)
// }

// func readReviews(w http.ResponseWriter, r *http.Request) {
// 	result := []Review{}
// 	if err := reviews.Find(nil).Sort("-created_at").All(&result); err != nil {
// 		responseError(w, err.Error(), http.StatusInternalServerError)
// 	} else {
// 		responseJSON(w, result)
// 	}
// }

// func responseError(w http.ResponseWriter, message string, code int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(map[string]string{"error": message})
// }

// func responseJSON(w http.ResponseWriter, data interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(data)
// }


package main
// import  (
// 	"net/http"
// 	"log"

// 	"github.com/marsDev31/kuclap-backend/api/controllers"
// 	"github.com/gorilla/mux"	
// 	"github.com/rs/cors"
// )
import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	// "fmt"
	"time"
	"os"

	"gopkg.in/mgo.v2"
	"github.com/marsDev31/kuclap-backend/api/models"
	"github.com/gorilla/mux"	
	"github.com/rs/cors"
		// "github.com/marsDev31/kuclap-backend/api/controllers"
)

// type Review struct {
// 	Text		string		`json:"text" bson:"text"`
// 	CreatedAt	time.Time	`json:"createdAt" bson:"created_at"`
// }
// type Review struct {
// 	// ID			primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	Text		string		`json:"text" bson:"text"`
// 	// Author		string		`json:"author" bson:"author"`
// 	// Grade		string		`json:"grade" bson:"grade"`
// 	// Auth		string		`json:"auth" bson:"auth"`
// 	CreatedAt	time.Time	`json:"createdAt" bson:"created_at"`
// }

var reviews *mgo.Collection

func main(){
	// Connect to mongo
	session, err := mgo.Dial("mongo:27017")
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Get reviews collection
	reviews = session.DB("ku-clap").C("reviews")
	
	r := mux.NewRouter()
	// r.Host("www.example.com") // set origin 
	// r.HandleFunc("/create_review", controllers.CreateReview).Methods("POST")
	// r.HandleFunc("/get_reviews", controllers.GetReviews).Methods("GET")
	r.HandleFunc("/create_review", createReview).Methods("POST")
	r.HandleFunc("/get_reviews", getReviews).Methods("GET")
	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
	
}

// func connectDB() *mgo.Collection{
// 	// Connect to mongo
// 	session, err := mgo.Dial("mongo:27017")
// 	if err != nil {
// 		log.Fatalln(err)
// 		log.Fatalln("mongo err")
// 		// fmt.Println("mongo err")
// 		os.Exit(1)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
	
// 	reviews := session.DB("ku-clap").C("reviews")
// 	// fmt.Println("Conected DB")
// 	return reviews
// }

func createReview(w http.ResponseWriter, r *http.Request) {
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
	
	// reviews := connectDB()
	if err := reviews.Insert(review); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSON(w, review)
}

func getReviews(w http.ResponseWriter, r *http.Request) {
	result := []models.Review{}
	// reviews := connectDB()
	
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