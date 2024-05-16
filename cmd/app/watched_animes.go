package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/n31t/go-project/pkg/model"
	"github.com/n31t/go-project/pkg/validator"
)

type output struct {
	Anime     *model.Anime `json:"anime"`
	WasViewed time.Time    `json:"wasViewed"`
	Tier      string       `json:"tier"`
}

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
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "anime added"})
}

func (app *application) watchedAnimeRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime id")
		return
	}
	user := app.contextGetUser(r)
	watchedAnimeExists, err := app.models.WatchedAnimes.IsWatched(id, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	if !watchedAnimeExists {
		app.respondWithError(w, http.StatusBadRequest, "Anime isn't in watched list")
		return
	}
	watchedAnime, err := app.models.WatchedAnimes.Select(id, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	anime, err := app.models.Animes.Select(watchedAnime.AnimeId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return

	}
	res := output{
		Anime:     anime,
		WasViewed: watchedAnime.WasViewed,
		Tier:      watchedAnime.Tier,
	}
	app.respondWithJSON(w, http.StatusOK, res)
}

func (app *application) watchedAnimeList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Genre string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Genre = app.readString(qs, "genre", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "episodes", "studio", "releaseYear", "genre", "rating", "-id", "-title", "-episodes", "-studio", "-releaseYear", "-genre", "-rating"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := app.contextGetUser(r)
	animes, metadata, err := app.models.WatchedAnimes.SelectAll(int(user.Id), input.Genre, input.Filters)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	var response []output
	for _, anime := range animes {
		animeInfo, err := app.models.Animes.Select(anime.AnimeId)
		if err != nil {
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
			return

		}
		response = append(response, output{
			Anime:     animeInfo,
			WasViewed: anime.WasViewed,
			Tier:      anime.Tier,
		})
	}

	app.respondWithJSONMetadata(w, http.StatusOK, response, metadata)

}

func (app *application) watchedAnimeListByTier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["tier"]
	param = strings.ToUpper(param)
	if !validator.Tier(param, validator.Tiers...) {
		app.respondWithError(w, http.StatusBadRequest, "Invalid tier")
		return
	}
	user := app.contextGetUser(r)
	animes, err := app.models.WatchedAnimes.SelectAllWithTier(param, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	var response []output
	for _, anime := range animes {
		animeInfo, err := app.models.Animes.Select(anime.AnimeId)
		if err != nil {
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
			return

		}
		response = append(response, output{
			Anime:     animeInfo,
			WasViewed: anime.WasViewed,
			Tier:      anime.Tier,
		})
	}
	app.respondWithJSON(w, http.StatusOK, response)
}
func (app *application) watchedAnimeUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime id")
		return
	}
	user := app.contextGetUser(r)
	watchedAnimeExists, err := app.models.WatchedAnimes.IsWatched(id, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	if !watchedAnimeExists {
		app.respondWithError(w, http.StatusBadRequest, "Anime isn't in watched list")
		return
	}

	anime, err := app.models.WatchedAnimes.Select(id, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	var input struct {
		Tier      string    `json:"tier"`
		WasViewed time.Time `json:"wasViewed"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Tier != "" {
		if !validator.Tier(input.Tier, validator.Tiers...) {
			app.respondWithError(w, http.StatusBadRequest, "Invalid tier")
			return
		}
		anime.Tier = input.Tier
	}
	if !input.WasViewed.IsZero() {
		anime.WasViewed = input.WasViewed
	}

	err = app.models.WatchedAnimes.Update(anime, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) watchedAnimeDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime id")
		return
	}
	user := app.contextGetUser(r)
	watchedAnimeExists, err := app.models.WatchedAnimes.IsWatched(id, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	if !watchedAnimeExists {
		app.respondWithError(w, http.StatusBadRequest, "Anime isn't in watched list")
		return
	}
	err = app.models.WatchedAnimes.Delete(id, int(user.Id))
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
