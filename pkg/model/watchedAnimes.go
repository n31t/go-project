package model

import (
	"database/sql"
	"log"
)

type WatchedAnimes struct {
	Id        string `json:"id"`
	AnimeId   string `json:"animeId"`
	UserId    string `json:"userId"`
	WasViewed string `json:"wasViewed"`
}

type WatchedAnimesModel struct {
	DB       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
