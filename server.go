package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"kuclap-review-api/src/config"
	"kuclap-review-api/src/dao"
	"kuclap-review-api/src/middleware"
	"kuclap-review-api/src/routes"
)

var (
	KIND          string
	PORT          string
	ORIGIN        []string
	serverConfig  config.Config
	configuration config.Configuration
	repository    dao.SessionDAO
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is alive 😎")
}

func init() {
	log.Println("Initial service... 🔧")

	serverConfig.Read()
	configuration = serverConfig.GetConfig()

	KIND = configuration.Kind
	ORIGIN = configuration.OriginAllowed
	repository.Server = configuration.Server
	repository.Database = configuration.Database

	PORT = config.Getenv("PORT", configuration.Port)

	repository.Connect()
	// initialClasses()
}

// Define HTTP request routes
func main() {
	log.Println("Starting server... 🤤")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Authorization", "Content-Type"})
	exposeOk := handlers.ExposedHeaders([]string{""})
	originsOk := handlers.AllowedOrigins(ORIGIN)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()

	routes.InjectAdapterDAO(&repository)

	routes.IndexClassesHandler(r)
	routes.IndexReviewHandler(r)
	routes.IndexQuestionsHandler(r)
	routes.IndexAnswersHandler(r)
	routes.IndexRecapHandler(r)
	routes.IndexAdminHandler(r)

	r.HandleFunc("/healthcheck", Healthcheck).Methods("GET")

	log.Println("Running on " + KIND + " Mode 🌶")
	log.Println("Server listening on port " + PORT + " 🚀")

	if err := http.ListenAndServe(":"+PORT, middleware.LimitMiddleware(handlers.CompressHandler(handlers.CORS(headersOk, exposeOk, originsOk, methodsOk)(r)))); err != nil {
		log.Fatal(err)
	}
}
