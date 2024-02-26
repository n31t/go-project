package model

import (
	"database/sql"
	"log"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserModel struct {
	DB       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
