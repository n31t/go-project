package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/n31t/go-project/pkg/model"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}
type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://adilovamir:password@localhost:5435/adilovamir?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	} else {
		log.Println("Connected to the database")
	}

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()

}

func (app *application) run() {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Animes
	v1.HandleFunc("/animes", app.animesList).Methods("GET")
	v1.HandleFunc("/animes", app.animeCreate).Methods("POST")
	v1.HandleFunc("/animes/{id}", app.animeRetrieve).Methods("GET")
	v1.HandleFunc("/animes/{id}", app.animeUpdate).Methods("PUT")
	// v1.HandleFunc("/animes/{id}", app.animeDelete).Methods("DELETE")

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
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
