package model

import (
	"context"
	"database/sql"
	"fmt"
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
	err := a.DB.QueryRowContext(ctx, query, id).Scan(&anime.Id, &anime.Title, &anime.Episodes, &anime.Studio,
		&anime.Description, &anime.ReleaseYear, &anime.Genre, &anime.Rating)
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

// func (a *AnimeModel) SelectAll() ([]*Anime, error) {
// 	query := `
// 	SELECT * FROM animes`
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	rows, err := a.DB.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var animes []*Anime
// 	for rows.Next() {
// 		var anime Anime
// 		err := rows.Scan(&anime.Id, &anime.Title, &anime.Episodes, &anime.Studio,
// 			&anime.Description, &anime.ReleaseYear, &anime.Genre, &anime.Rating)
// 		if err != nil {
// 			return nil, err
// 		}
// 		animes = append(animes, &anime)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return animes, nil
// }

func (a *AnimeModel) SelectAll(title string, genre string, filter Filters) ([]*Anime, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, title, episodes, studio, description, releaseYear, genre, rating
	FROM animes
	WHERE (STRPOS(LOWER(title), LOWER($1))> 0 or $1 = '')
	AND (STRPOS(LOWER(genre), LOWER($2))> 0 or $2 = '')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4
	`, filter.sortColumn(), filter.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, genre, filter.limit(), filter.offset()}
	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	var animes []*Anime
	totalRecords := 0
	for rows.Next() {
		var anime Anime
		err := rows.Scan(&totalRecords, &anime.Id, &anime.Title, &anime.Episodes, &anime.Studio,
			&anime.Description, &anime.ReleaseYear, &anime.Genre, &anime.Rating)
		if err != nil {
			return nil, Metadata{}, err
		}
		animes = append(animes, &anime)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err

	}
	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return animes, metadata, nil
}

func (a *AnimeModel) Update(anime *Anime) error {
	query := `
	UPDATE animes
	SET title = $1, episodes = $2, studio = $3, description = $4, releaseYear = $5, genre = $6, rating = $7
	WHERE id = $8`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := a.DB.ExecContext(ctx, query, anime.Title, anime.Episodes,
		anime.Studio, anime.Description, anime.ReleaseYear, anime.Genre, anime.Rating, anime.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AnimeModel) Delete(id int) error {
	watchedAnimesQuery := `
    DELETE 
    FROM watched_animes
    WHERE anime_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := a.DB.ExecContext(ctx, watchedAnimesQuery, id)
	if err != nil {
		return err
	}

	animesQuery := `
    DELETE 
    FROM animes
    WHERE id = $1`
	_, err = a.DB.ExecContext(ctx, animesQuery, id)
	if err != nil {
		return err
	}
	return nil
}
