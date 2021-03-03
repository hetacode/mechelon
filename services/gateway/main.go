package main

import (
	"log"
	"net/http"
	"os"
	"time"

	gtwhandlers "github.com/hetacode/mechelon/services/gateway/handlers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gtwcontainer "github.com/hetacode/mechelon/services/gateway/container"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env.dev")

	log.Println("api gateway svc is starting")

	c := gtwcontainer.NewContainer()

	router := mux.NewRouter()
	clientsRouter := router.PathPrefix("/clients").Subrouter()
	frontendRouter := router.PathPrefix("/api").Subrouter()

	gtwhandlers.NewClientsHandlers(c, clientsRouter)
	gtwhandlers.NewFrontendHandlers(c, frontendRouter)

	corsRouter := useCORS(router)
	srv := &http.Server{
		Handler:      corsRouter,
		Addr:         "0.0.0.0:" + os.Getenv("SVC_API_GATEWAY_PORT"),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 5,
	}

	log.Println("api gateway svc is running")
	log.Fatal(srv.ListenAndServe())
}

func useCORS(handler http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"0.0.0.0"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	cred := handlers.AllowCredentials()
	corsHandler := handlers.CORS(headersOk, originsOk, methodsOk, cred)(handler)

	return corsHandler
}
