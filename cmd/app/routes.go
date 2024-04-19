package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) run() {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// v1.NotFoundHandler = http.HandlerFunc(app.notFoundHandler)
	// v1.MethodNotAllowedHandler = hhtp.HandlerFunc(app.methodNotAllowedResponse)
	// Animes
	v1.HandleFunc("/animes", app.animesList).Methods("GET")
	v1.HandleFunc("/animes", app.animeCreate).Methods("POST")
	v1.HandleFunc("/animes/{id:[0-9]+}", app.animeRetrieve).Methods("GET")
	v1.HandleFunc("/animes/{id:[0-9]+}", app.animeUpdate).Methods("PUT")
	v1.HandleFunc("/animes/{id}", app.animeDelete).Methods("DELETE")

	// Healthcheck
	v1.HandleFunc("/healthcheck", app.healthCheck).Methods("GET")

	// Users
	// v1.HandleFunc("/users", app.usersList).Methods("GET")
	// v1.HandleFunc("/users/{id}", app.userRetrieve).Methods("GET")
	// v1.HandleFunc("/users/{id}", app.userCreate).Methods("POST")
	// v1.HandleFunc("/users/{id}", app.userUpdate).Methods("PUT")
	// v1.HandleFunc("/users/{id}", app.userDelete).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	if err != nil {
		log.Fatal(err)
	}
}
