package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Animes AnimeModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Animes: AnimeModel{DB: db, InfoLog: infoLog, ErrorLog: errorLog},
		Users:  UserModel{DB: db, InfoLog: infoLog, ErrorLog: errorLog},
	}
}
