package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/n31t/go-project/pkg/model"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) animeCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string  `json:"title"`
		Episodes    int     `json:"episodes"`
		Studio      string  `json:"studio"`
		Description string  `json:"description"`
		ReleaseYear int     `json:"releaseYear"`
		Genre       string  `json:"genre"`
		Rating      float64 `json:"rating"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	anime := &model.Anime{
		Title:       input.Title,
		Episodes:    input.Episodes,
		Studio:      input.Studio,
		Description: input.Description,
		ReleaseYear: input.ReleaseYear,
		Genre:       input.Genre,
		Rating:      input.Rating,
	}

	err = app.models.Animes.Insert(anime)
	if err != nil {
		// app.respondWithError(w, http.StatusInternalServerError, "TEST")
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusCreated, anime)
}

func (app *application) animesList(w http.ResponseWriter, r *http.Request) {
	animes, err := app.models.Animes.SelectAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, animes)
}

func (app *application) animeRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	anime, err := app.models.Animes.Select(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Anime not found")
		return
	}
	app.respondWithJSON(w, http.StatusOK, anime)
}

func (app *application) animeUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}
	anime, err := app.models.Animes.Select(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Anime not found")
		return
	}

	var input struct {
		Title       *string  `json:"title"`
		Episodes    *int     `json:"episodes"`
		Studio      *string  `json:"studio"`
		Description *string  `json:"description"`
		ReleaseYear *int     `json:"releaseYear"`
		Genre       *string  `json:"genre"`
		Rating      *float64 `json:"rating"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		anime.Title = *input.Title
	}
	if input.Episodes != nil {
		anime.Episodes = *input.Episodes
	}
	if input.Studio != nil {
		anime.Studio = *input.Studio
	}
	if input.Description != nil {
		anime.Description = *input.Description
	}
	if input.ReleaseYear != nil {
		anime.ReleaseYear = *input.ReleaseYear
	}
	if input.Genre != nil {
		anime.Genre = *input.Genre
	}
	if input.Rating != nil {
		anime.Rating = *input.Rating
	}
	err = app.models.Animes.Update(anime)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, anime)
}

func (app *application) animeDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}
	err = app.models.Animes.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return err
	}

	return nil
}
