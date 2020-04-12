package main
import (
	"net/http"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
)


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

	// Get posts collection
	reviews = session.DB("ku-clap").C("reviews")
	
	r := mux.NewRouter()
	// r.Host("www.example.com") // set origin 
	r.HandleFunc("/create_review", createPost).Methods("POST")
	r.HandleFunc("/posts", readPosts).Methods("GET")
	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
	
}