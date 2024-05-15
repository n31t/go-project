package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// func (app *application) run() {
// 	r := mux.NewRouter()
// 	v1 := r.PathPrefix("/api/v1").Subrouter()

// 	// v1.NotFoundHandler = http.HandlerFunc(app.notFoundHandler)
// 	// v1.MethodNotAllowedHandler = hhtp.HandlerFunc(app.methodNotAllowedResponse)
// 	// Animes
// 	v1.HandleFunc("/animes", app.animesList).Methods("GET")
// 	v1.HandleFunc("/animes", app.animeCreate).Methods("POST")
// 	v1.HandleFunc("/animes/{id:[0-9]+}", app.animeRetrieve).Methods("GET")
// 	v1.HandleFunc("/animes/{id:[0-9]+}", app.animeUpdate).Methods("PUT")
// 	v1.HandleFunc("/animes/{id}", app.animeDelete).Methods("DELETE")

// 	// Healthcheck
// 	v1.HandleFunc("/healthcheck", app.healthCheck).Methods("GET")

// 	// Users
// 	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
// 	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")

// 	// Tokens
// 	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

// 	log.Printf("Starting server on %s\n", app.config.port)
// 	// err := http.ListenAndServe(app.config.port, r)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// }

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Animes
	v1.HandleFunc("/animes", app.animesList).Methods("GET")
	v1.HandleFunc("/animes", app.requirePermission("animes:create", app.animeCreate)).Methods("POST")
	v1.HandleFunc("/animes/{id:[0-9]+}", app.animeRetrieve).Methods("GET")
	v1.HandleFunc("/animes/{id:[0-9]+}", app.requirePermission("animes:update", app.animeUpdate)).Methods("PUT")
	v1.HandleFunc("/animes/{id}", app.requirePermission("animes:delete", app.animeDelete)).Methods("DELETE")

	// Watched Animes
	v1.HandleFunc("/watched-animes", app.requireActivatedUser(app.watchedAnimeCreate)).Methods("POST")
	// Healthcheck
	v1.HandleFunc("/healthcheck", app.healthCheck).Methods("GET")

	// Users
	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")

	// Tokens
	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}
