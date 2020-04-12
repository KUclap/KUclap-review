package main
import (
	"net/http"
	"log"

	"github.com/marsDev31/kuclap-backend/api/controllers"
	"github.com/gorilla/mux"	
	"github.com/rs/cors"
)


// var reviews *mgo.Collection

func main(){
	// // Connect to mongo
	// session, err := mgo.Dial("mongo:27017")
	// if err != nil {
	// 	log.Fatalln(err)
	// 	log.Fatalln("mongo err")
	// 	os.Exit(1)
	// }
	// defer session.Close()
	// session.SetMode(mgo.Monotonic, true)

	// // Get posts collection
	// reviews = session.DB("ku-clap").C("reviews")
	
	r := mux.NewRouter()
	// r.Host("www.example.com") // set origin 
	r.HandleFunc("/create_review", controllers.CreateReview).Methods("POST")
	r.HandleFunc("/get_reviews", controllers.GetReviews).Methods("GET")
	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
	
}