package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Anime struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Episodes    int     `json:"episodes"`
	Studio      string  `json:"studio"`
	Description string  `json:"description"`
	ReleaseYear int     `json:"releaseYear"`
	Genre       string  `json:"genre"`
	Rating      float64 `json:"rating"`
}

type AnimeModel struct {
	DB       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (a *AnimeModel) Insert(anime *Anime) error {
	query := `INSERT INTO animes (title, episodes, studio, description, releaseYear, genre, rating) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := a.DB.QueryRowContext(ctx, query, anime.Title, anime.Episodes,
		anime.Studio, anime.Description, anime.ReleaseYear,
		anime.Genre, anime.Rating).Scan(&anime.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AnimeModel) Select(id int) (*Anime, error) {
	query := `
	SELECT id, title, episodes, studio, description, releaseYear, genre, rating
	FROM animes
	WHERE id = $1`
	var anime Anime
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := a.DB.QueryRowContext(ctx, query, anime.Id).Scan(&anime.Id, &anime.Title, &anime.Episodes, &anime.Studio,
		&anime.Description, &anime.ReleaseYear, &anime.Genre, &anime.Rating)
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

func (a *AnimeModel) SelectAll() ([]*Anime, error) {
	query := `
	SELECT * FROM animes`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var animes []*Anime
	for rows.Next() {
		var anime Anime
		err := rows.Scan(&anime.Id, &anime.Title, &anime.Episodes, &anime.Studio,
			&anime.Description, &anime.ReleaseYear, &anime.Genre, &anime.Rating)
		if err != nil {
			return nil, err
		}
		animes = append(animes, &anime)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeModel) Update(anime *Anime) error {
	query := `
	UPDATE animes
	SET title = $1, episodes = $2, studio = $3, description = $4, releaseYear = $5, genre = $6, rating = $7
	WHERE id = $8`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := a.DB.QueryRowContext(ctx, query, anime.Title, anime.Episodes,
		anime.Studio, anime.Description, anime.ReleaseYear, anime.Genre, anime.Rating, anime.Id).Scan()
	if err != nil {
		return err
	}
	return nil
}
