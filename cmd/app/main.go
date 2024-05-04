package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/n31t/go-project/pkg/model"
	"github.com/n31t/go-project/pkg/model/filler"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type config struct {
	port int
	env  string
	fill bool
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
	// flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.BoolVar(&cfg.fill, "fill", false, "Fill the database with initial data")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres:password@postgres:5432/adilovamir?sslmode=disable", "PostgreSQL DSN")
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

	// Migrations
	migrationDown(db)
	migrationUp(db)

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	if cfg.fill {
		err := filler.FillDatabase(app.models)
		if err != nil {
			log.Fatal(err)
		}
	}

	// app.run()
	if err := app.serve(); err != nil {
		fmt.Sprintf("error starting server: %v\n", err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
