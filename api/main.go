
package main
import (
	"net/http"
	"log"
	// "encoding/json"
	// "io/ioutil"
	// "time"
	// "os"

	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "github.com/marsDev31/kuclap-backend/api/models"
	"github.com/marsDev31/kuclap-backend/api/controllers"
	"github.com/gorilla/mux"	
	"github.com/rs/cors"
)

// var reviews *mgo.Collection

func main(){
	// Connect to mongo
	// session, err := mgo.Dial("mongo:27017")
	// if err != nil {
	// 	log.Fatalln(err)
	// 	log.Fatalln("mongo err")
	// 	os.Exit(1)
	// }
	// defer session.Close()
	// session.SetMode(mgo.Monotonic, true)

	// Get reviews collection
	// reviews = session.DB("ku-clap").C("reviews")

	r := mux.NewRouter()
	// r.Host("www.example.com") // set origin 
	r.HandleFunc("/create_review", controllers.CreateReview).Methods("POST")
	r.HandleFunc("/get_reviews", controllers.GetReviews).Methods("GET")
	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
	
}

// func createReview(w http.ResponseWriter, r *http.Request) {
// 	data, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		responseError(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	review := &models.Review{}
// 	err = json.Unmarshal(data, review)
// 	if err != nil {
// 		responseError(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

	
// 	review.CreatedAt = time.Now().UTC()
// 	review.ReviewID = bson.NewObjectId()

// 	if err := reviews.Insert(review); err != nil {
// 		responseError(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	responseJSON(w, review)
// }

// func getReviews(w http.ResponseWriter, r *http.Request) {
// 	result := []models.Review{}
	
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