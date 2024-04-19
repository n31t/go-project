package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/n31t/go-project/pkg/model"

	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres:password@localhost:5435/adilovamir?sslmode=disable", "PostgreSQL DSN")
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
	// migrationDown(db)
	migrationUp(db)

	// err = m.Force(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = m.Down()
	// if err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal(err)
	// }
	// err = m.Up()
	// if err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal(err)
	// }

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
