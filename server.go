package main

import (
	// "os"
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

var (
	KIND			string
	PORT			string
	ORIGIN			[]string
	serverConfig	config.Config
	configuration	config.Configuration
	mgoDAO			dao.SessionDAO

)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is alive 😎")
}

func init() {
	log.Println("Initial service... 🔧") 
	
	serverConfig.Read()
	configuration	=	serverConfig.GetConfig()

	KIND			=	configuration.Kind
	ORIGIN			=	configuration.OriginAllowed
	mgoDAO.Server	=	configuration.Server
	mgoDAO.Database	=	configuration.Database

	PORT			=	config.Getenv("PORT",	configuration.Port)

	mgoDAO.Connect()
	// initialClasses()
}

// Define HTTP request routes
func main() {
	log.Println("Starting server... 🤤")
	
	headersOk	:=	handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Authorization", "Content-Type"})
	exposeOk	:=	handlers.ExposedHeaders([]string{""})
	originsOk	:=	handlers.AllowedOrigins(ORIGIN)
	methodsOk	:=	handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r			:=	mux.NewRouter()

	routes.InjectAdapterDAO(&mgoDAO)
	
	routes.IndexClassesHandler(r)
	routes.IndexReviewHandler(r)
	routes.IndexQuestionsHandler(r)
	routes.IndexAnswersHandler(r)
	routes.IndexRecapHandler(r)
	routes.IndexAdminHandler(r)

	r.HandleFunc("/healthcheck", Healthcheck).Methods("GET")

	log.Println("Running on " + KIND + " Mode 🌶")	
	log.Println("Server listening on port " + PORT + " 🚀")

	if err	:=	http.ListenAndServe(":" + PORT, middleware.LimitMiddleware(handlers.CompressHandler(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r)))); err != nil {
		log.Fatal(err)
	}	
}