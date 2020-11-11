package main

import (
	"os"
	"log"
	"net/http"
	"fmt"
    "kuclap-review-api/src/config"
	"kuclap-review-api/src/middleware"
	"kuclap-review-api/src/routes"
	"kuclap-review-api/src/dao"

	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)


var serverConfig = config.Config{}
var mgoDAO = dao.SessionDAO{}
var kind string
var port string

// ROOT request
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is alive 😎")
}

// Parse the serverConfiguration file 'serverConfig.toml', and establish a connection to DB
func init() {
	log.Println("Initial service... 🔧") 
	
	serverConfig.Read()
	
	if os.Getenv("KIND") == "development" {
		kind = serverConfig.Development.Kind
		port = serverConfig.Application.Port
		mgoDAO.Server = serverConfig.Development.Server
		mgoDAO.Database = serverConfig.Development.Database
	} else if os.Getenv("KIND") == "production" {
		kind = serverConfig.Production.Kind
		port = serverConfig.Application.Port
		mgoDAO.Server = serverConfig.Production.Server
		mgoDAO.Database = serverConfig.Production.Database
	} else {
		kind = serverConfig.Development.Kind + " (staging on heroku)"
		port = os.Getenv("PORT")
		mgoDAO.Server = serverConfig.Development.Server
		mgoDAO.Database = serverConfig.Development.Database
	}

	mgoDAO.Connect()
	// initialClasses()
}

// Define HTTP request routes
func main() {
	log.Println("Starting server... 🤤")
	
	origin := serverConfig.Application.ORIGIN_ALLOWED
	
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Authorization", "Content-Type"})
	exposeOk := handlers.ExposedHeaders([]string{""})
	originsOk := handlers.AllowedOrigins([]string{origin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()
	routes.InjectAdapterDAO(&mgoDAO)
	routes.IndexClassesHandler(r)
	routes.IndexReviewHandler(r)
	routes.IndexQuestionsHandler(r)
	r.HandleFunc("/healthcheck", Healthcheck).Methods("GET")
	log.Println("Running on " + kind + " Mode 🌶")	
	log.Println("Server listening on port " + port + " 🚀")
	if err := http.ListenAndServe(":" + port, middleware.LimitMiddleware(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r))); err != nil {
		log.Fatal(err)
	}	
}