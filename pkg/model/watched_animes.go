package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type WatchedAnime struct {
	Id        string    `json:"id"`
	AnimeId   string    `json:"animeId"`
	UserId    string    `json:"userId"`
	WasViewed time.Time `json:"wasViewed"`
	Tier      string    `json:"tier"`
}

type WatchedAnimeModel struct {
	DB       *sql.DB
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (w *WatchedAnimeModel) Insert(watchedAnime *WatchedAnime) error {
	query := `INSERT INTO watched_animes (anime_id, user_id) 
	VALUES ($1, $2)
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := w.DB.QueryRowContext(ctx, query, watchedAnime.AnimeId, watchedAnime.UserId).Scan(&watchedAnime.Id)
	if err != nil {
		return err
	}
	return nil
}

func (w *WatchedAnimeModel) Select(id int, userId int) (*WatchedAnime, error) {
	query := `
    SELECT id, anime_id, user_id, was_viewed, tier
    FROM watched_animes
    WHERE id = $1 AND user_id = $2`
	var watchedAnime WatchedAnime
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := w.DB.QueryRowContext(ctx, query, id, userId).Scan(&watchedAnime.Id, &watchedAnime.AnimeId, &watchedAnime.UserId, &watchedAnime.WasViewed, &watchedAnime.Tier)
	if err != nil {
		return nil, err
	}
	return &watchedAnime, nil
}

func (w *WatchedAnimeModel) Update(watchedAnime *WatchedAnime, userId int) error {
	query := `UPDATE watched_animes SET was_viewed = $1, tier = $2 WHERE id = $3 AND user_id = $4`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := w.DB.ExecContext(ctx, query, watchedAnime.WasViewed, watchedAnime.Tier, userId, watchedAnime.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (w *WatchedAnimeModel) Delete(id int, userId int) error {
	query := `DELETE FROM watched_animes WHERE id = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := w.DB.ExecContext(ctx, query, id, userId)
	if err != nil {
		return err
	}
	return nil
}

func (w *WatchedAnimeModel) SelectAll(userId int, genre string, filter Filters) ([]*WatchedAnime, Metadata, error) {
	query := fmt.Sprintf(`
    SELECT count(*) OVER(), wa.id, wa.anime_id, wa.user_id, wa.was_viewed, wa.tier
    FROM watched_animes wa
    INNER JOIN animes a ON wa.anime_id = a.id
    WHERE wa.user_id = $1 
	AND (STRPOS(LOWER(a.genre), LOWER($2))> 0 or $2 = '')
    ORDER BY %s %s, wa.id ASC
    LIMIT $3 OFFSET $4
    `, filter.sortColumn(), filter.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{userId, genre, filter.limit(), filter.offset()}
	rows, err := w.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	var watchedAnimes []*WatchedAnime
	totalRecords := 0
	for rows.Next() {
		var watchedAnime WatchedAnime
		err := rows.Scan(&totalRecords, &watchedAnime.Id, &watchedAnime.AnimeId, &watchedAnime.UserId, &watchedAnime.WasViewed, &watchedAnime.Tier)
		if err != nil {
			return nil, Metadata{}, err
		}
		watchedAnimes = append(watchedAnimes, &watchedAnime)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return watchedAnimes, metadata, nil
}

func (w *WatchedAnimeModel) AddTier(tier string, watchedAnime *WatchedAnime, userId int) error {
	query := `UPDATE watched_animes SET tier = $1 WHERE id = $2 AND user_id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := w.DB.ExecContext(ctx, query, tier, watchedAnime.Id, watchedAnime.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (w *WatchedAnimeModel) SelectAllWithTier(tier string, userId int) ([]*WatchedAnime, error) {
	query := `
	SELECT id, anime_id, user_id, was_viewed, tier
	FROM watched_animes
	WHERE user_id = $1 AND tier = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := w.DB.QueryContext(ctx, query, userId, tier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var watchedAnimes []*WatchedAnime
	for rows.Next() {
		var watchedAnime WatchedAnime
		err := rows.Scan(&watchedAnime.Id, &watchedAnime.AnimeId, &watchedAnime.UserId, &watchedAnime.WasViewed, &watchedAnime.Tier)
		if err != nil {
			return nil, err
		}
		watchedAnimes = append(watchedAnimes, &watchedAnime)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return watchedAnimes, nil
}
