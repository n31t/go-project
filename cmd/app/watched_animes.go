package main

import (
	"net/http"

	"github.com/n31t/go-project/pkg/model"
)

func (app *application) watchedAnimeCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AnimeId int `json:"animeId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user := app.contextGetUser(r)

	watchedAnime := &model.WatchedAnime{
		AnimeId: input.AnimeId,
		UserId:  int(user.Id),
	}
	watchedAnimeExists, err := app.models.WatchedAnimes.IsWatched(watchedAnime.AnimeId, watchedAnime.UserId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	if watchedAnimeExists {
		app.respondWithError(w, http.StatusBadRequest, "Anime already in watched list")
		return
	}
	err = app.models.WatchedAnimes.Insert(watchedAnime)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusCreated, "Anime added to watched list")
}
