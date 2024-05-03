package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	ErrRecordNotFound = errors.New("record not found")

	ErrEditConflict = errors.New("edit conflict: resource has been modified by another user")
)

type Models struct {
	Animes      AnimeModel
	Permissions PermissionModel
	Tokens      TokenModel
	Users       UserModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Animes:      AnimeModel{DB: db, InfoLog: infoLog, ErrorLog: errorLog},
		Permissions: PermissionModel{DB: db, InfoLog: infoLog, ErrorLog: errorLog},
		Tokens:      TokenModel{DB: db, InfoLog: infoLog, ErrorLog: errorLog},
		Users:       UserModel{DB: db, InfoLog: infoLog, ErrorLog: errorLog},
	}
}
