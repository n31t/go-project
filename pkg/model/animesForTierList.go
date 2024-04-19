package model

import (
	"database/sql"
	"log"
)

type AnimesForTierlist struct {
	Id      string `json:"id"`
	AnimeId string `json:"animeId"`
	Tier    string `json:"tier"`
}

type AnimesForTierListModel struct {
	DB       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
