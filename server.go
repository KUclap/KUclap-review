package main

import (
	"log"
	"net/http"
	"fmt"
	"kuclap-review-api/src/helper"
    "kuclap-review-api/src/config"
	"kuclap-review-api/src/middleware"
	"kuclap-review-api/src/routes"
	"kuclap-review-api/src/dao"

	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

var mgoDAO = dao.SessionDAO{}

// ROOT request
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is alive 😎")
}

// Parse the serverConfiguration file 'serverConfig.toml', and establish a connection to DB
func init() {
	log.Println("Initial service... 🔧")
	serverConfig := config.Config{}
	serverConfig.Read()
	mgoDAO.Database = serverConfig.Database
	mgoDAO.Server = helper.GetENV("DB_SERVER")
	mgoDAO.Connect()
	// initialClasses()
}

// Define HTTP request routes
func main() {
	log.Println("Starting server... 🤤")
	port := helper.GetENV("PORT")
	origin := helper.GetENV("ORIGIN_ALLOWED")
	
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Authorization", "Content-Type"})
	exposeOk := handlers.ExposedHeaders([]string{""})
	originsOk := handlers.AllowedOrigins([]string{origin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()
	routes.InjectAdapterDAO(&mgoDAO)
	routes.IndexClassesHandler(r)
	routes.IndexReviewHandler(r)
	r.HandleFunc("/healthcheck", Healthcheck).Methods("GET")
		
	log.Println("Server listening on port " + port + " 🚀")
	if err := http.ListenAndServe(":" + port, middleware.LimitMiddleware(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r))); err != nil {
		log.Fatal(err)
	}
	
}