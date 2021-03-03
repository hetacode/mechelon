package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env.dev")

	log.Println("api gateway svc is starting")

	router := mux.NewRouter()
	clientsRouter := router.PathPrefix("/clients").Subrouter()
	frontendRouter := router.PathPrefix("/api").Subrouter()

	configureClientHandlers(clientsRouter)
	configureFrontendHandlers(frontendRouter)

	corsRouter := useCORS(router)
	srv := &http.Server{
		Handler:      corsRouter,
		Addr:         "0.0.0.0:4000", // TODO: move PORT to the env file
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 5,
	}

	log.Println("api gateway svc is running")
	log.Fatal(srv.ListenAndServe())
}

func configureClientHandlers(h http.Handler) {

}

func configureFrontendHandlers(h http.Handler) {

}

func useCORS(handler http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"0.0.0.0"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	cred := handlers.AllowCredentials()
	corsHandler := handlers.CORS(headersOk, originsOk, methodsOk, cred)(handler)

	return corsHandler
}
